FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/karteikarten ./cmd/karteikarten

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /out/karteikarten /usr/local/bin/karteikarten
EXPOSE 8080
ENTRYPOINT ["karteikarten"]
