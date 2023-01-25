# syntax=docker/dockerfile:1
FROM golang:1.18-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

ADD https://beta-likeminds-media.s3.ap-south-1.amazonaws.com/environment/beta-environment ./.env

RUN go build -o ./kettle

CMD [ "./kettle" ]
