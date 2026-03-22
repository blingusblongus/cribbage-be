FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o cribbage-be .

FROM alpine:3.21
COPY --from=build /app/cribbage-be /cribbage-be
EXPOSE 8080
CMD ["/cribbage-be"]
