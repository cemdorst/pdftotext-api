FROM golang:1.20 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...

FROM alpine:3.18
RUN apk add --no-cache poppler-utils
COPY --from=build /usr/local/bin/app /usr/bin/

EXPOSE 8080

CMD ["app"]
