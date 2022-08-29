FROM golang:1.18-alpine AS builder

ENV CGO_ENABLED=0
ENV GOARCH=arm64

WORKDIR /build
COPY cmd/ .
COPY go.mod .
RUN go build -o /build/kanarod

FROM gcr.io/distroless/base-debian11

ENV GOARCH=arm64

COPY --from=builder /build/kanarod /bin/

ENTRYPOINT ["/bin/kanarod"]
