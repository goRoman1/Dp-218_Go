package main

import (
	"Dp218Go/configs"
	"Dp218Go/protos"
	"Dp218Go/repositories/postgres"
	"Dp218Go/routing"
	"Dp218Go/routing/grpcserver"
	"Dp218Go/routing/httpserver"
	"Dp218Go/services"
	"fmt"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
)

var sessionKey = "secretkey"

func main() {

	var connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_DB)

	db, err := postgres.NewConnection(connectionString)
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer db.CloseDB()

	err = doMigrate(connectionString)
	if err != nil {
		log.Printf("app - Run - Migration issues: %v\n", err)
	}

	var userRoleRepoDB = postgres.NewUserRepoDB(db)
	var userService = services.NewUserService(userRoleRepoDB, userRoleRepoDB)

	var accRepoDb = postgres.NewAccountRepoDB(userRoleRepoDB, db)
	var accService = services.NewAccountService(accRepoDb, accRepoDb, accRepoDb)

	var scooterRepo = postgres.NewScooterRepoDB(db)
	var grpcScooterService = services.NewGrpcScooterService(scooterRepo)
	var scooterService = services.NewScooterService(scooterRepo)

	scL, err := scooterRepo.GetAllScooters()
	fmt.Println(scL, err)



	svr := grpcserver.NewServer()
	svr.Run()
	grpcServer := grpc.NewServer()

	protos.RegisterScooterServiceServer(grpcServer, svr)

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("grpc server started: 8000")
		log.Fatal(grpcServer.Serve(listener))
	}()

	http.HandleFunc("/scooter", svr.ScooterHandler)
	http.HandleFunc("/run", routing.StartScooterTrip)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("./templates/html/scooter-run.html")
		if err!= nil {
			fmt.Println(err)
		}
		err = tmpl.Execute(w, scL)
		if err!=nil {
			fmt.Println()
		}
	})

	sessStore := sessions.NewCookieStore([]byte(sessionKey))
	authService := services.NewAuthService(userRoleRepoDB, sessStore)

	handler := routing.NewRouter(authService)
	routing.AddUserHandler(handler, userService)
	routing.AddAccountHandler(handler, accService)
	routing.AddScooterHandler(handler,scooterService)
	routing.AddGrpcScooterHandler(handler, grpcScooterService)
	httpServer := httpserver.New(handler, httpserver.Port(configs.HTTP_PORT))


	fmt.Println("http server started: 9000")
	http.ListenAndServe(":9000", nil)



	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatalf("app - Run - httpServer.Notify: %v", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Fatalf("app - Run - httpServer.Shutdown: %v", err)
	}
}

func doMigrate(connStr string) error {
	migr, err := migrate.New("file://"+configs.MIGRATIONS_PATH, connStr+"?sslmode=disable")
	if err != nil {
		return err
	}

	if configs.MIGRATE_VERSION_FORCE > 0 {
		migr.Force(configs.MIGRATE_VERSION_FORCE)
	}

	if configs.MIGRATE_DOWN {
		migr.Down()
	}

	return migr.Up()
}
