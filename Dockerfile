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

# copy over required directories
COPY ./models ./models
COPY ./api ./api
COPY *.go ./
COPY go-build.sh ./

RUN ./go-build.sh

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

COPY --from=build /weather /weather
COPY --from=react /app/build /static

WORKDIR /

# copy in the city.json list that should be found in ./data/city.json, along with startup script
COPY ./city.list.min.json ./city.json
COPY run.sh ./

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/run.sh"]
