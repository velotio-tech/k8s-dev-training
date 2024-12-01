FROM golang:alpine3.16

COPY . /res-crud

WORKDIR /res-crud

RUN go build -o main && \
    chmod 777 /res-crud

CMD ["./main"]