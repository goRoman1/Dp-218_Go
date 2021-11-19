package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ITA-Dnipro/Dp-218_Go/db"
	"github.com/ITA-Dnipro/Dp-218_Go/handler"
)

const SERV_PORT = ":8080"

func main() {
	listener, err := net.Listen("tcp", SERV_PORT)
	if err != nil {
		log.Fatalf("Error occured: %s", err.Error())
	}
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Couldn't set up database: %v", err)
	}

	defer database.Conn.Close()

	httpHandler := handler.NewHandler(database)
	server := &http.Server{
		Handler: httpHandler,
	}

	go func() {
		server.Serve(listener)
	}()
	defer Stop(server)

	log.Printf("Started server on %s", SERV_PORT)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGINT)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")

}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Couldn't shutdown server correctly: %v \n", err)
		os.Exit(1)
	}
}
