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
	Latitude float64
	Longitude float64
	stream protos.ScooterService_ReceiveClient
}

func NewGrpcScooterService(repo repositories.ScooterRepo) *GrpcScooterService {
	return &GrpcScooterService{
		repo,
	}
}

func NewGrpcScooterClient(id uint64, lat, lon float64,
	stream protos.ScooterService_ReceiveClient) *GrpcScooterClient {
	return &GrpcScooterClient{
		Id: id,
		Latitude: lat,
		Longitude: lon,
		stream: stream,
	}
}

func (s *GrpcScooterClient) grpcScooterMessage()  {
	intPol := time.Duration(interval) * time.Second

	fmt.Println("executing run in client")
	msg := &protos.ClientMessage{
		Id: s.Id,
		Latitude:  s.Latitude,
		Longitude:  s.Longitude,
	}
	err := s.stream.Send(msg)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	time.Sleep(intPol)
}

func (s *GrpcScooterClient) Run(station models.Coordinate) error {
	//err, _ := s.SendAtStart(uid, int(s.Id))
	//if err != nil {
	//	fmt.Println(err)
	//}

	switch {
	case s.Latitude <= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude; s.Latitude,
		s.Longitude = s.Latitude+step,s.Longitude+step {
				fmt.Println(s)
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
		//err = db.SendAtEnd(tripId, s)
		//if err!=nil {
		//	fmt.Println(err)
		//}
	case s.Latitude >= station.Latitude && s.Longitude <= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude; s.Latitude,
			s.Longitude = s.Latitude-step,s.Longitude+step {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
		//err = db.SendAtEnd(tripId, s)
		//if err!=nil {
		//	fmt.Println(err)
		//}
	case s.Latitude >= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude; s.Latitude,
			s.Longitude = s.Latitude-step,s.Longitude-step  {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
		//err = db.SendAtEnd(tripId, s)
		//if err!=nil {
		//	fmt.Println(err)
		//}
	case s.Latitude <= station.Latitude && s.Longitude >= station.Longitude:
		for ; s.Latitude <= station.Latitude && s.Longitude <= station.Longitude; s.Latitude,
			s.Longitude = s.Latitude+step,s.Longitude-step {
			s.grpcScooterMessage()
		}
		fmt.Println("Trip finished. You are at the point")
		//err = db.SendAtEnd(tripId, s)
		//if err!=nil {
		//	fmt.Println(err)
		//}
	default:
		return fmt.Errorf("you are at this point now")
	}
	return nil
}

func InitClient(scooter models.Scooter, repo repositories.ScooterRepo) (*GrpcScooterClient, error) {
	scooterStatus, err := repo.GetScooterStatus(scooter.ID)
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

	return &GrpcScooterClient{
		Id: uint64(scooter.ID),
		Latitude: scooterStatus.Location.Latitude,
		Longitude: scooterStatus.Location.Longitude,
		stream: stream,
	}, nil
}

func (gss *GrpcScooterService) InitAndRun(scooterID int,
	coordinate models.Coordinate) error{
	scooter,err := gss.GetScooterById(scooterID)
	if err!= nil {
		fmt.Println(err)
		return err
	}
	// TODO по айди станции нахожу координаты
	client,err := InitClient(scooter, gss.ScooterRepo)
	if err!= nil {
		fmt.Println(err)
		return err
	}
	err = client.Run(coordinate)
	fmt.Println(err)
	return err
}