FROM golang:alpine
WORKDIR /k8
COPY . .
RUN go build -o main
CMD [ "./main" ]