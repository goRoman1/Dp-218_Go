FROM golang:1.17-alpine as builder
COPY go.mod go.sum /go/src/Dp218Go/
WORKDIR /go/src/Dp218Go
RUN go mod tidy
COPY . /go/src/Dp218Go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/scooterapp ./cmd/app

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/Dp218Go/migrations/. /home/Dp218Go/migrations
COPY --from=builder /go/src/Dp218Go/templates/. /home/Dp218Go/templates
COPY --from=builder /go/src/Dp218Go/build/scooterapp /usr/bin/scooterapp
EXPOSE 8080 8080
ENTRYPOINT [ "/usr/bin/scooterapp" ]