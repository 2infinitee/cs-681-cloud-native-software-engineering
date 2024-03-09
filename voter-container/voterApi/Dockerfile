FROM golang:latest AS build-stage

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o voterApi

FROM alpine:latest AS run-stage

WORKDIR /
COPY --from=build-stage /app voterApi
EXPOSE 8080
ENV REDIS_URL=host.docker.internal:6379
CMD ["/voterApi"]
