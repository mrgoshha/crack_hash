FROM golang:alpine AS builder

WORKDIR usr/src/worker

#dependencies
COPY go.mod go.sum ./
RUN go mod download

#build
COPY . .
RUN go build -o /usr/local/bin/worker cmd/worker/main.go

FROM alpine AS runner
COPY --from=builder /usr/local/bin/worker /

#CMD ["./worker", "-baseurl=:5535"]
