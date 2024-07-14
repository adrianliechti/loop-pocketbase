# syntax=docker/dockerfile:1

FROM golang:1 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /pocketbase


FROM alpine:3 AS final

WORKDIR /

RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /pocketbase .

EXPOSE 8090

ENTRYPOINT ["/pocketbase", "serve", "--http=0.0.0.0:8090"]