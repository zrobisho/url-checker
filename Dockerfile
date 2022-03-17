FROM golang:1.17 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN ls cmd
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-checker cmd/url-checker/main.go

FROM alpine
COPY --from=builder /app/url-checker .

ENTRYPOINT ["/url-checker"]