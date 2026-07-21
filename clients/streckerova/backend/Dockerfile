FROM golang:1.25-alpine AS build
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /app .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app ./app
COPY frontend/ ./frontend/
EXPOSE 8080
ENTRYPOINT ["./app"]
