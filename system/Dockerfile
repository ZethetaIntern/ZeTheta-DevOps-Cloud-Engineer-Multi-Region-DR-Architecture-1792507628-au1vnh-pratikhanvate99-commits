FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY app ./app
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/paysecure ./app

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/paysecure /paysecure
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/paysecure"]

