FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /privaquest ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /privaquest /app/privaquest
ENV PORT=8080
ENV ADMIN_API_KEY=change-me-admin-key
ENV DATABASE_URL=file:/data/privaquest.db?_pragma=foreign_keys(1)
VOLUME /data
EXPOSE 8080
CMD ["/app/privaquest"]
