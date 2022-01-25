FROM golang:1.17-alpine AS builder

ENV CGO_ENABLED=0
ENV GOARCH=arm64

WORKDIR /build
COPY . .
RUN go build -o /build/kanarod

FROM gcr.io/distroless/base-debian11

ENV GOARCH=arm64

COPY --from=builder /build/kanarod /bin/

ENTRYPOINT ["/bin/kanarod"]
