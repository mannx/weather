#syntax:docker/dockerfile:1

# build stage

#FROM golang:buster AS build
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

WORKDIR /
COPY --from=build /weather /weather
COPY static/ /static/

# copy the script used to import data on schedule
COPY newdata.sh /etc/periodic/15min/newdata

# copy the entry point startup script
COPY run.sh /run.sh

# make sure to include tzdata to use timezones correctly
# also added in curl for use in our cron script
RUN apk add tzdata curl

# need 
EXPOSE 8080

USER root:root

ENTRYPOINT ["/run.sh"]
