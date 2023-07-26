FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/go-video-text/main.go


FROM jrottenberg/ffmpeg:alpine AS FFmpeg

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=builder /app/main /app/main

COPY --from=FFmpeg / /
COPY ./configs /app/configs
COPY ./fonts /app/fonts

CMD ["./main"]
