# Build stage
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GO111MODULE=on go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags musl --ldflags "-extldflags -static" -a -o golang-clean-architecture cmd/api/main.go

# deploy
FROM alpine

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/golang-clean-architecture .
COPY --from=builder /app/config . 

CMD ["./golang-clean-architecture"]
