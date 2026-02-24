FROM golang AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o docker-stats-exporter .

FROM scratch

WORKDIR /app

COPY --from=builder /app/docker-stats-exporter /app/docker-stats-exporter

EXPOSE 9100

ENTRYPOINT ["/app/docker-stats-exporter"]
