FROM golang:1.16.3 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o calendar-service .

FROM alpine

COPY --from=build /src/calendar-service /

EXPOSE 5000
CMD ["/calendar-service"]
