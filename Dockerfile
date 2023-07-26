FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/go-video-text/main.go


FROM scratch

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY ./configs /app/configs

CMD ["./main"]
