#syntax:docker/dockerfile:1

# build stage

FROM golang:alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /weather


#
# deploy stage

#FROM gcr.io/distroless/base-debian10 
FROM alpine

RUN apk update

# make sure curl is installed to retrieve weather updates
RUN apk add curl tzdata

WORKDIR /
COPY --from=build /weather /weather
COPY static/ /static/

RUN mkdir data

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/weather"]
