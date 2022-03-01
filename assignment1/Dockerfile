FROM golang:latest

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod verify
RUN go mod download
RUN go build

CMD ["./assignment1"]