FROM golang:1.26 AS build

WORKDIR /app
#COPY go.mod go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o ./server github.com/robindittmar/dttmr-api/cmd/api-server

FROM debian:latest

WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
