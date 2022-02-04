#syntax:docker/dockerfile:1

#
# build stage
#

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
COPY --from=react /app/react/build /static

# create 2 symlinks for the static folder to work correctly
# TODO: fix this properly at some point
WORKDIR /static
RUN ln -s static/js
RUN ln -s static/css

WORKDIR /

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/weather"]
