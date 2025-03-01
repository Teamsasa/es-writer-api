FROM golang:latest

WORKDIR /src

COPY . .
RUN go mod download
RUN apt-get update && apt-get install -y postgresql-client

WORKDIR app/cmd/server

RUN go build -o /main .

EXPOSE 8080

CMD ["/main"]
