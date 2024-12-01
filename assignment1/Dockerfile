FROM golang:latest

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build

CMD ["./assignment1"]