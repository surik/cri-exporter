FROM golang:1.22-alpine AS builder
WORKDIR /app

ARG TAG

COPY go.mod go.sum ./
COPY . ./

RUN apk add --no-cache make
RUN go mod download
RUN make

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin/cri-exporter .

ENTRYPOINT [ "./cri-exporter" ]