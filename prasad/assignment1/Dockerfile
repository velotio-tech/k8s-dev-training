FROM golang:1.16-alpine

ENV GO111MODULE auto

ENV APP_HOME go/src/assignment1
WORKDIR $APP_HOME

COPY . .

RUN go mod verify
RUN go mod download
RUN go mod tidy
RUN go build
CMD ["./assignment1"]