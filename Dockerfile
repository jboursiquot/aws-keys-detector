FROM golang:1.13 as builder

WORKDIR /app
COPY . /app

RUN go get -d -v

# Static compilation
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

# Distroless container image
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]