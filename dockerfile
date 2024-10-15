FROM golang:1.22.5-alpine

WORKDIR /GoVersi

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /GoVersi/cmd

RUN go build -o /GoVersi/bin/GoVersi .

EXPOSE 8080
ENTRYPOINT [ "/GoVersi/bin/GoVersi" ]