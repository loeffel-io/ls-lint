FROM golang:1.18 as builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

FROM gcr.io/distroless/base

COPY --from=builder /usr/src/app/ls-lint-linux /ls-lint
VOLUME /data
WORKDIR /data
ENTRYPOINT ["/ls-lint"]