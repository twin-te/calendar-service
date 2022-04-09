FROM golang:1.18 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o calendar-service .

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=build /src/calendar-service /

EXPOSE 5000
CMD ["/calendar-service"]
