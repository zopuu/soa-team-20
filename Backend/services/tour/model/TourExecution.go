package model

import (
	"time"

	"github.com/google/uuid"
)

type TourExecutionStatus string

const (
	TourExecActive    TourExecutionStatus = "Active"
	TourExecCompleted TourExecutionStatus = "Completed"
	TourExecAbandoned TourExecutionStatus = "Abandoned"
)

type KeyPointRef struct {
	Id          uuid.UUID   `json:"id" bson:"id"`
	Title       string      `json:"title" bson:"title"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
	Order       int         `json:"order" bson:"order"`
}

type VisitedKeyPoint struct {
	KeyPointRef
	VisitedAt time.Time `json:"visitedAt" bson:"visitedAt"`
}

type TourExecution struct {
    Id                     uuid.UUID         `json:"id" bson:"_id"`
    TourId                 uuid.UUID         `json:"tourId" bson:"tourId"`
    UserId                 string            `json:"userId" bson:"userId"`
    Status                 TourExecutionStatus `json:"status" bson:"status"`
    CurrentTouristPosition Coordinates       `json:"currentTouristPosition" bson:"currentTouristPosition"`
    StartedAt              time.Time         `json:"startedAt" bson:"startedAt"`
    LastActivityAt         time.Time         `json:"lastActivityAt" bson:"lastActivityAt"`
    EndedAt                *time.Time        `json:"endedAt,omitempty" bson:"endedAt,omitempty"`
    KeyPointsRemaining     []KeyPointRef     `json:"keyPointsRemaining" bson:"keyPointsRemaining"`
    KeyPointsVisited       []VisitedKeyPoint `json:"keyPointsVisited" bson:"keyPointsVisited"`

    // optional new fields
    TotalKeyPoints    int        `json:"totalKeyPoints" bson:"totalKeyPoints"`
    NextKeyPointIndex int        `json:"nextKeyPointIndex" bson:"nextKeyPointIndex"`
    LastKnownCoords   Coordinates `json:"lastKnownCoords" bson:"lastKnownCoords"`
}

