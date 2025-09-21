package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"tour.xws.com/model"
	"tour.xws.com/repository"
)

const proximityRadiusMeters = 10.0

type TourExecutionService struct {
	TourRepo         *repository.TourRepository
	KeyPointRepo     *repository.KeyPointRepository
	CurrentLocRepo   *repository.CurrentLocationRepository
	TourExecRepo     *repository.TourExecutionRepository
}

func (s *TourExecutionService) Start(userId string, tourId uuid.UUID) (*model.TourExecution, error) {
	
	if existing, err := s.TourExecRepo.GetActiveByUserAndTour(userId, tourId); err != nil {
		return nil, err
	} else if existing != nil {
		return existing, nil
	}
	// 0) Prevent multiple actives for this user
	// if active, err := s.TourExecRepo.GetActiveByUserId(userId); err != nil {
	// 	return nil, err
	// } else if active != nil {
	// 	return nil, errors.New("user already has an active tour execution")
	// }

	// 1) Validate tour status (Published/Archived)
	tour, err := s.TourRepo.GetById(tourId)
	if err != nil {
		return nil, err
	}
	if !(tour.Status == model.Published || tour.Status == model.Archived) {
		return nil, errors.New("tour not published or archived")
	}

	// 2) Load ordered keypoints
	kps, err := s.KeyPointRepo.GetAllByTourSortedByCreatedAt(tourId)
	if err != nil {
		return nil, err
	}
	if len(kps) == 0 {
		return nil, errors.New("tour has no keypoints")
	}

	// 3) Build remaining list
	remaining := make([]model.KeyPointRef, 0, len(kps))
	for i, kp := range kps {
		remaining = append(remaining, model.KeyPointRef{
			Id:          kp.Id,
			Title:       kp.Title,
			Coordinates: kp.Coordinates,
			Order:       i + 1,
		})
	}

	// 4) Get current user location to use as start
	loc, err := s.CurrentLocRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, errors.New("no current location for user")
	}

	// 5) Create execution
	now := time.Now()
	te := &model.TourExecution{
		Id:                     uuid.New(),
		TourId:                 tourId,
		UserId:                 userId,
		Status:                 model.TourExecActive, // e.g., "Active"
		CurrentTouristPosition: loc.Coordinates,
		StartedAt:              now,
		LastActivityAt:         now,
		KeyPointsRemaining:     remaining,
		KeyPointsVisited:       []model.VisitedKeyPoint{},
		TotalKeyPoints:         len(kps),
		NextKeyPointIndex:      0,
		LastKnownCoords:        loc.Coordinates, // if you keep this field too
		KeyPointsCompletitionTimes: []time.Time{},
	}

	// 6) Persist execution
	if err := s.TourExecRepo.Insert(te); err != nil {
		return nil, err
	}

	// 7) (Optional) Upsert location again as “last known”
	if s.CurrentLocRepo != nil {
		_ = s.CurrentLocRepo.Upsert(userId, loc.Coordinates)
	}

	return te, nil
}

type ProximityCheckResult struct {
	Reached             bool                 `json:"reached"`
	DistanceMeters      float64              `json:"distanceMeters"`
	NextKeyPoint        *model.KeyPointRef   `json:"nextKeyPoint,omitempty"`
	RemainingCount      int                  `json:"remainingCount"`
	JustCompletedPoint  *model.KeyPointRef   `json:"justCompletedPoint,omitempty"`
	CompletedSession    bool                 `json:"completedSession"`
}

func (s *TourExecutionService) GetActiveByUserAndTour(userId string, tourId uuid.UUID) (*model.TourExecution, error) {
	return s.TourExecRepo.GetActiveByUserAndTour(userId, tourId)
}

func (s *TourExecutionService) CheckProximity(userId string, coords model.Coordinates) (*ProximityCheckResult, error) {
	te, err := s.TourExecRepo.GetActiveByUser(userId)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, errors.New("no active tour execution")
	}

	now := time.Now()
	te.CurrentTouristPosition = coords
	te.LastActivityAt = now

	res := &ProximityCheckResult{
		Reached:        false,
		DistanceMeters: 0,
		RemainingCount: len(te.KeyPointsRemaining),
	}

	if len(te.KeyPointsRemaining) == 0 {
		// već završena po tačkama => kompletiraj
		te.Status = model.TourExecCompleted
		te.EndedAt = &now
		res.CompletedSession = true
		if err := s.TourExecRepo.Update(te); err != nil {
			return nil, err
		}
		return res, nil
	}

	next := te.KeyPointsRemaining[0]
	dist := HaversineMeters(coords.Latitude, coords.Longitude, next.Coordinates.Latitude, next.Coordinates.Longitude)
	res.DistanceMeters = dist
	res.NextKeyPoint = &next

	if dist <= proximityRadiusMeters {
		// pomeri “head” iz remaining u visited
		te.KeyPointsRemaining = te.KeyPointsRemaining[1:]
		te.KeyPointsVisited = append(te.KeyPointsVisited, model.VisitedKeyPoint{
			KeyPointRef: next,
			VisitedAt:   now,
		})

        te.KeyPointsCompletitionTimes = append(te.KeyPointsCompletitionTimes, now)

		res.Reached = true
		res.JustCompletedPoint = &next
		res.RemainingCount = len(te.KeyPointsRemaining)

		if len(te.KeyPointsRemaining) == 0 {
			te.Status = model.TourExecCompleted
			te.EndedAt = &now
			res.CompletedSession = true
		}
	}

	if err := s.TourExecRepo.Update(te); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TourExecutionService) Abandon(userId string) error {
	te, err := s.TourExecRepo.GetActiveByUser(userId)
	if err != nil {
		return err
	}
	if te == nil {
		return errors.New("no active tour execution")
	}
	now := time.Now()
	te.Status = model.TourExecAbandoned
	te.EndedAt = &now
	te.LastActivityAt = now
	return s.TourExecRepo.Update(te)
}

func (s *TourExecutionService) GetActive(userId string) (*model.TourExecution, error) {
	return s.TourExecRepo.GetActiveByUser(userId)
}

func HaversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0 // meters
	rad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := rad(lat2 - lat1)
	dLon := rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rad(lat1))*math.Cos(rad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
