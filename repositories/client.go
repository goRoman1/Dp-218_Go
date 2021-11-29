package repositories

import (
	"Dp218Go/pkg/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

const step = 0.0001

func ClAdd() {
	conn, err := grpc.DialContext(context.Background(), ":8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	sclient := pb.NewScooterServiceClient(conn)
	stream, err := sclient.Receive(context.Background())
	if err != nil {
		panic(err)
	}

	cl := NewClient(1, 48.423,35.032, stream)
	cl.Run(1)
}

type Client struct {
	Id uint64
	Latitude  float64
	Longitude  float64
	// In  chan ServerMessage
	stream pb.ScooterService_ReceiveClient
}

func NewClient(id uint64, latitude, longitude float64, stream pb.ScooterService_ReceiveClient) *Client {
	return &Client{
		Id: id,
		Latitude:  latitude,
		Longitude:  longitude,
		// In:  in,
		stream: stream,
	}
}

func (s *Client) Run(interval int) {

	intPol := time.Duration(interval) * time.Second

	fmt.Println("executing run in client")
	// x, y := randomStep()

	for {

		//TODO change direction make it random

		s.Latitude, s.Longitude = s.Latitude+step, s.Longitude+step
		// send location to server
		msg := &pb.ClientMessage{
			Id: s.Id,
			Latitude:  s.Latitude,
			Longitude:  s.Longitude,
		}
		err := s.stream.Send(msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("after send client")
		time.Sleep(intPol)
	}
}
