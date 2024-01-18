# syntax=docker/dockerfile:1
FROM golang:1.20

WORKDIR /websmtp/pkg
COPY pkg ./

WORKDIR /websmtp/websmtp
COPY websmtp/go.mod websmtp/go.sum websmtp/*.go ./
RUN go mod download
RUN go build

ENV PORT 8080
ENV HOSTNAME "0.0.0.0"
EXPOSE 8080

CMD ["./websmtp"]
