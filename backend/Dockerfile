# BUILD STAGE
FROM golang:1.26 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git \
    curl \
    bash \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .


# DEV STAGE
FROM golang:1.26 AS dev

WORKDIR /workspace

RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/air-verse/air@latest

CMD ["sleep", "infinity"]