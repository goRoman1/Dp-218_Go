package repositories

import (
	"Dp218Go/pkg/pb"
	"Dp218Go/pkg/postgres"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

const step = 0.001


func ClAdd() {
	conn, err := grpc.DialContext(context.Background(), ":8000", grpc.WithInsecure())
	// сделать переподключение, вместо падения
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	sclient := pb.NewScooterServiceClient(conn)
	stream, err := sclient.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	db, err := postgres.NewPostgres("postgres://scooteradmin:Megascooter!@localhost:5444/scooterdb")
	if err !=nil {
		fmt.Println(err)
	}
	cl := NewClient(4, 48.42332,35.03242, stream)
	cl.Run(1, 1, NewSc(db) )
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

var destination = []struct {
	latitude float64
	longitude float64
}{ {
	latitude: 48.4423,
	longitude: 35.0434,
	}, {
	latitude: 48.3233,
	longitude: 35.0434,
	}, {
	latitude: 48.0223,
	longitude: 35.0134,
	}, {
	latitude: 48.4223,
	longitude: 35.0234,
	},
}

var station = struct {
	latitude float64
	longitude float64
}{
	latitude: 48.4423,
	longitude: 35.0434,
}



func (s *Client) Run(interval, uid int, db *ScooterRepoDb) {
	err, tripId, locaId := db.SendAtStart(uid,s)
	if err != nil {
		fmt.Println(err)
	}
	//при старте отправляю координаты старта в базу
	switch {
	case s.Latitude <= station.latitude && s.Longitude <= station.longitude:
		for ; s.Latitude <=station.latitude && s.
			Longitude <= station.longitude; s.Latitude,s.Longitude = s.Latitude+step,s.Longitude+step {
			s.moving(1)
		}
		fmt.Println("Trip finished. You are at the point")
		db.SendAtEnd(tripId, locaId, s)
			if err!=nil {
				fmt.Println(err)
			}
	case s.Latitude >= station.latitude && s.Longitude <= station.longitude:
		for ; s.Latitude >= station.latitude && s.
			Longitude <= station.longitude; s.Latitude,s.Longitude = s.Latitude-step,s.Longitude+step {
			s.moving(1)
		}
		fmt.Println("Trip finished. You are at the point")
		db.SendAtEnd(tripId, locaId, s)
		if err!=nil {
			fmt.Println(err)
		}
	case s.Latitude >= station.latitude && s.Longitude >= station.longitude:
		for ; s.Latitude >= station.latitude && s.
			Longitude >= station.longitude; s.Latitude,s.Longitude = s.Latitude-step,s.Longitude-step {
			s.moving(1)
		}
		fmt.Println("Trip finished. You are at the point")
		db.SendAtEnd(tripId, locaId, s)
		if err!=nil {
			fmt.Println(err)
		}
	case s.Latitude <= station.latitude && s.Longitude >= station.longitude:
		for ; s.Latitude <=station.latitude && s.
			Longitude >= station.longitude; s.Latitude,s.Longitude = s.Latitude+step,s.Longitude-step {
			s.moving(1)
		}
		fmt.Println("Trip finished. You are at the point")
		db.SendAtEnd(tripId, locaId, s)
		if err!=nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("You are at this point now")
	}

	//при окончании поездки - отправляю координты завершения
}

func (s *Client) moving(interval int)  {
	intPol := time.Duration(interval) * time.Second

	fmt.Println("executing run in client")
	msg := &pb.ClientMessage{
		Id: s.Id,
		Latitude:  s.Latitude,
		Longitude:  s.Longitude,
	}
	err := s.stream.Send(msg)
	if err != nil {
		panic(err)
	}
	time.Sleep(intPol)
}
