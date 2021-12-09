package services

import (
	"Dp218Go/models"
	"Dp218Go/protos"
	"Dp218Go/repositories"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

const (
	step = 0.001
	interval = 1
)

type GrpcScooterService struct {
	repositories.ScooterRepo

}

type GrpcScooterClient struct {
	Id uint64
	coordinate models.Coordinate
	stream protos.ScooterService_ReceiveClient
	repositories.ScooterRepo
}

func NewGrpcScooterService(repo repositories.ScooterRepo) *GrpcScooterService {
	return &GrpcScooterService{
		repo,
	}
}

func NewGrpcScooterClient(id uint64, coordinate models.Coordinate,
	stream protos.ScooterService_ReceiveClient) *GrpcScooterClient {
	return &GrpcScooterClient{
		Id: id,
		coordinate: coordinate,
		stream: stream,
	}
}

func (gss *GrpcScooterService) InitAndRun(scooterID int,
	coordinate models.Coordinate) error{
	scooter,err := gss.GetScooterById(scooterID)
	if err!= nil {
		fmt.Println(err)
		return err
	}
	// TODO по айди станции нахожу координаты

	scooterStatus, err := gss.GetScooterStatus(scooter.ID)
	if err!= nil {
		fmt.Println(err)
	}

	conn, err := grpc.DialContext(context.Background(), ":8000", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	sClient := protos.NewScooterServiceClient(conn)
	stream, err := sClient.Receive(context.Background())
	if err != nil {
		panic(err)
	}

	client := NewGrpcScooterClient(uint64(scooterID),
		scooterStatus.Location, stream)
	err = client.Run(coordinate)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (s *GrpcScooterClient) grpcScooterMessage()  {
	intPol := time.Duration(interval) * time.Second

	fmt.Println("executing run in client")
	msg := &protos.ClientMessage{
		Id: s.Id,
		Latitude:  s.coordinate.Latitude,
		Longitude:  s.coordinate.Longitude,
	}
	err := s.stream.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(intPol)
}

func (s *GrpcScooterClient) Run(station models.Coordinate) error {

	switch {
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude; s.coordinate.Latitude,
		s.coordinate.Longitude = s.coordinate.Latitude+step,s.coordinate.Longitude+step {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude; s.coordinate.Latitude,
			s.coordinate.Longitude = s.coordinate.Latitude-step,s.coordinate.Longitude+step {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude; s.coordinate.Latitude,
			s.coordinate.Longitude = s.coordinate.Latitude-step,s.coordinate.Longitude-step  {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude; s.coordinate.Latitude,
			s.coordinate.Longitude = s.coordinate.Latitude+step,s.coordinate.Longitude-step {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
	default:
		return fmt.Errorf("you are at this point now")
	}
	return nil
}