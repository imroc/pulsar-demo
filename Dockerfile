FROM golang:1.22-alpine3.19 AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o pulsar-demo -ldflags '-w -extldflags "-static"' .

FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates

COPY --from=build /workspace/pulsar-demo /usr/local/bin/pulsar-demo

CMD ["keda-tencentcloud-clb-scaler"]

