FROM golang:1.17.7-buster

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

ENTRYPOINT ["./assignment1"]


