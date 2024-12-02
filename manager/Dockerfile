FROM golang:alpine AS builder

WORKDIR usr/src/manager

#dependencies
COPY go.mod go.sum ./
RUN go mod download

#build
COPY . .
RUN go build -o /usr/local/bin/manager cmd/manager/main.go

FROM alpine AS runner
COPY --from=builder /usr/local/bin/manager /
COPY configs/manager.toml configs/manager.toml

#CMD ["./manager", "http://localhost:5535"]
