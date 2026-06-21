FROM golang:latest AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags '-extldflags "-static"' -o /bank ./cmd/bank

FROM scratch
COPY --from=build /bank /usr/local/bin/bank
ENTRYPOINT ["bank"]
