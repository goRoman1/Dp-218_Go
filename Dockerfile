FROM golang:1.17-alpine as builder
COPY go.mod go.sum /go/src/github.com/ITA-Dnipro/Dp-218_Go/
WORKDIR /go/src/github.com/ITA-Dnipro/Dp-218_Go
RUN go mod download
COPY . /go/src/github.com/ITA-Dnipro/Dp-218_Go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/scooterapp github.com/ITA-Dnipro/Dp-218_Go

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/ITA-Dnipro/Dp-218_Go/build/scooterapp /usr/bin/scooterapp
EXPOSE 8088 8080
ENTRYPOINT [ "/usr/bin/scooterapp" ]