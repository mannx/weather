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
# build react front end
#

FROM node:alpine AS react

WORKDIR /app/react

COPY ./frontend ./

RUN npm install
RUN npm run build


#
# Deploy stage
#
FROM alpine

# make sure required packages are installed
RUN apk update
RUN apk add curl tzdata

WORKDIR /

RUN mkdir data
RUN mkdir static

COPY --from=build /weather /weather
COPY --from=react /app/react/build /static

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/weather"]
