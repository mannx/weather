#syntax:docker/dockerfile:1

#
# build stage
#

FROM golang:alpine AS build

# need gcc tools for 1 or more libraries
RUN apk add build-base

ENV GOPATH /go/src

WORKDIR /go/src/github.com/mannx/weather

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./models ./models
COPY *.go ./

RUN go build -o /weather


#
# build react front end
#

FROM node:alpine AS react

WORKDIR /app
ENV NODE_ENV=production

COPY ["./frontend/package.json", "./frontend/package-lock.json", "./"]

RUN npm install --production

COPY ./frontend .

RUN npm run build

#
# Deploy stage
#

FROM alpine

# make sure required packages are installed
RUN apk update
RUN apk add tzdata

WORKDIR /

RUN mkdir data
RUN mkdir static

COPY --from=build /weather /weather
COPY --from=react /app/build /static

# create 2 symlinks for the static folder to work correctly
# TODO: fix this properly at some point
#WORKDIR /static
WORKDIR /
RUN ln -s static/js
RUN ln -s static/css

#WORKDIR /

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/weather"]
