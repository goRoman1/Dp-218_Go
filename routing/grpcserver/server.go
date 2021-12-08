package grpcserver

import (
	"Dp218Go/protos"
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

type Client struct {
	w    io.Writer
	done chan struct{}
}

type Server struct {
	client map[int]*Client
	taken  map[int]bool
	codes  map[int]int
	in     chan *protos.ClientMessage
	*protos.UnimplementedScooterServiceServer
}


func GISHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./presentation/views/templates/html/index.html")
}

func (ss *Server) ScooterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new client connected")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := &Client{
		w:    w,
		done: make(chan struct{}),
	}
	ss.AddClient(client)

	<-client.done
	fmt.Println("connection closed")
}

func NewServer() *Server {
	return &Server{
		client: make(map[int]*Client),
		taken:  make(map[int]bool),
		codes:  make(map[int]int),
		in:     make(chan *protos.ClientMessage),
	}
}

func (ss *Server) AddClient(c *Client) {
	ss.client[1] = c
}

func (ss *Server) Register(msg *protos.ClientRequest, stream protos.ScooterService_RegisterServer) error {
	return nil
}

func (ss *Server) Receive(stream protos.ScooterService_ReceiveServer) error {
	var err error

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			err = status.Errorf(codes.Internal, "unexpected error %v", err)
			break
		}

		ss.in <- msg

	}

	return err
}

func (ss *Server) Run() {
	go func() {
		for {
			select {
			case msg := <-ss.in:
				var buf bytes.Buffer
				json.NewEncoder(&buf).Encode(msg)

				for _, client := range ss.client {

					go func(c *Client) {
						if _, err := fmt.Fprintf(c.w, "data: %v\n\n", buf.String()); err != nil {
							fmt.Println(err)
							delete(ss.client, 1)
							close(c.done)
							return
						}

						if f, ok := c.w.(http.Flusher); ok {
							f.Flush()
						}
						fmt.Printf("data: %v\n", buf.String())
					}(client)
				}
			}
		}
	}()
}
