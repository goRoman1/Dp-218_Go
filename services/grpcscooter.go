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
	step = 0.0005
	dischargeStep = 1
	interval = 1
)

type GrpcScooterService struct {
	repositories.ScooterRepo
}

type GrpcScooterClient struct {
	ID uint64
	coordinate models.Coordinate
	batteryRemain float64
	stream protos.ScooterService_ReceiveClient
}

func NewGrpcScooterService(repo repositories.ScooterRepo) *GrpcScooterService {
	return &GrpcScooterService{
		repo,
	}
}

func NewGrpcScooterClient(id uint64, coordinate models.Coordinate, battery float64,
	stream protos.ScooterService_ReceiveClient) *GrpcScooterClient {
	return &GrpcScooterClient{
		ID: id,
		coordinate: coordinate,
		batteryRemain: battery,
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

	scooterStatus, err := gss.GetScooterStatus(scooterID)
	if err!= nil {
		fmt.Println(err)
		return err
	}

	if scooterStatus.BatteryRemain > 10 {
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
			scooterStatus.Location, scooter.BatteryRemain, stream)
		err = client.run(coordinate)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("trip finished")

		err = gss.SendCurrentStatus(int(client.ID), client.coordinate.Latitude, client.coordinate.Longitude,
			client.batteryRemain)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("status has been sent")

		if client.batteryRemain <= 0 {
			err = fmt.Errorf("scooter battery discharged. Trip is over")
			fmt.Println(err.Error())
			return err
		}
		return nil
	}

	err = fmt.Errorf("scooter battery is too low for trip. Choose another one")
	fmt.Println(err.Error())
	return err
}

func (s *GrpcScooterClient) grpcScooterMessage()  {
	intPol := time.Duration(interval) * time.Second

	fmt.Println("executing run in client")
	msg := &protos.ClientMessage{
		Id: s.ID,
		Latitude:  s.coordinate.Latitude,
		Longitude:  s.coordinate.Longitude,
	}
	err := s.stream.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(intPol)
}

func (s *GrpcScooterClient) run(station models.Coordinate) error {

	switch {
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.
			coordinate.Latitude,
		s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude+step,s.coordinate.Longitude+step,
		s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude-step,s.coordinate.Longitude+step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude-step,s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep  {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.
			Latitude,
			s.coordinate.Longitude, s.batteryRemain = s.coordinate.Latitude+step,s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
			fallthrough
	case s.coordinate.Latitude <= station.Latitude:
		for ; s.coordinate.Latitude <= station.Latitude && s.
			batteryRemain > 0; s.coordinate.Latitude, s.batteryRemain = s.coordinate.Latitude+step, s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Latitude >= station.Latitude:
		for ; s.coordinate.Latitude >= station.Latitude && s.
			batteryRemain > 0; s.coordinate.Latitude, s.batteryRemain = s.coordinate.Latitude-step, s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Longitude >= station.Longitude:
		for ; s.coordinate.Longitude >= station.Longitude && s.
			batteryRemain > 0; s.coordinate.Longitude, s.batteryRemain = s.coordinate.Longitude-step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
		fallthrough
	case s.coordinate.Longitude <= station.Longitude:
		for ; s.coordinate.Longitude <= station.Longitude && s.
			batteryRemain > 0; s.coordinate.Longitude, s.batteryRemain = s.coordinate.Longitude+step,
			s.batteryRemain-dischargeStep {
			s.grpcScooterMessage()
		}
	default:
		return fmt.Errorf("error happened")
	}
	return nil
}