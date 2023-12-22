# syntax=docker/dockerfile:1
FROM golang:1.20

WORKDIR /oauth2/pkg
COPY pkg ./

WORKDIR /oauth2/oauth2
COPY oauth2/go.mod oauth2/go.sum oauth2/*.go ./
RUN go mod download
RUN go build

ENV PORT 8080
EXPOSE 8080

CMD ["./oauth2"]
