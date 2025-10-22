FROM golang:1.25.1-alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base && apk add docker

RUN mkdir /app
WORKDIR /app
RUN mkdir -p /tmp/air && chmod 1777 /tmp/air

COPY . .
COPY .env .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go install github.com/air-verse/air@latest

EXPOSE 8080

CMD [ "air" ]
