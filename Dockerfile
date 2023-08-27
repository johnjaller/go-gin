FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w" -o myGolangApp

FROM alpine

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8004

ENTRYPOINT ["./myGolangApp"]
