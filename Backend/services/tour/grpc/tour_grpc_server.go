package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"tour.xws.com/model"
	tourpb "tour.xws.com/proto"
	"tour.xws.com/service"
)

type TourGRPCServer struct {
	tourpb.UnimplementedTourServiceServer
	tourService *service.TourService
}

func NewTourGRPCServer(tourService *service.TourService) *TourGRPCServer {
	return &TourGRPCServer{
		tourService: tourService,
	}
}

func (s *TourGRPCServer) CreateTour(ctx context.Context, req *tourpb.CreateTourRequest) (*tourpb.CreateTourResponse, error) {
	// Convert protobuf message to domain model
	tour := model.BeforeCreateTour(
		req.AuthorId,
		req.Title,
		req.Description,
		req.Tags,
		model.TourDifficulty(req.Difficulty),
	)
	
	// Set transport type from request
	tour.TransportType = model.TransportType(req.TransportType)

	// Call the service
	if err := s.tourService.Create(tour); err != nil {
		log.Printf("gRPC CreateTour ERROR err=%v", err)
		return nil, err
	}

	// Convert domain model back to protobuf
	pbTour := &tourpb.Tour{
		Id:          tour.Id.String(),
		AuthorId:    tour.AuthorId,
		Title:       tour.Title,
		Description: tour.Description,
		Difficulty:  tourpb.TourDifficulty(tour.Difficulty),
		Tags:        tour.Tags,
		Status:      tourpb.TourStatus(tour.Status),
		Price:       tour.Price,
		Distance:    tour.Distance,
		Duration:    tour.Duration,
		TransportType: tourpb.TransportType(tour.TransportType),
	}

	// Set timestamps if they're not zero
	if !tour.PublishedAt.IsZero() {
		pbTour.PublishedAt = timestamppb.New(tour.PublishedAt)
	}
	if !tour.ArchivedAt.IsZero() {
		pbTour.ArchivedAt = timestamppb.New(tour.ArchivedAt)
	}

	log.Printf("gRPC CreateTour OK id=%s", tour.Id)
	return &tourpb.CreateTourResponse{
		Tour:    pbTour,
		Message: "Tour created successfully",
	}, nil
}

func (s *TourGRPCServer) DeleteTour(ctx context.Context, req *tourpb.DeleteTourRequest) (*tourpb.DeleteTourResponse, error) {
	
	// Parse UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("gRPC DeleteTour BAD_ID raw=%s", req.Id)
		return &tourpb.DeleteTourResponse{
			Message: "Invalid tour ID",
			Success: false,
		}, nil
	}
	// Call the service
	if err = s.tourService.Delete(id); err != nil {
		log.Printf("gRPC DeleteTour ERROR id=%s err=%v", id, err)
		return &tourpb.DeleteTourResponse{
			Message: err.Error(),
			Success: false,
		}, nil
	}

	log.Printf("gRPC DeleteTour OK id=%s", id)
	return &tourpb.DeleteTourResponse{
		Message: "Tour deleted successfully",
		Success: true,
	}, nil
}