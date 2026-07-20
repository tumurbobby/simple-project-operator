FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o controller .

FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY --from=builder /app/controller /controller

ENTRYPOINT ["/controller"]