#syntax:docker/dockerfile:1

# build stage

FROM golang:buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /weather


#
# deploy stage

FROM gcr.io/distroless/base-debian10 

WORKDIR /
COPY --from=build /weather /weather
COPY static/ /static/

EXPOSE 8080

USER root:root

ENTRYPOINT ["/weather"]
