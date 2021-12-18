FROM golang AS builder

WORKDIR /epitar

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o epitar

FROM debian
WORKDIR /epitar
COPY --from=builder /epitar/epitar /epitar/epitar
COPY --from=builder /epitar/modules /epitar/modules
ENTRYPOINT ["/epitar/epitar"]
