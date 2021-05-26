# FROM golang:1.15.11-alpine3.12 AS build
# WORKDIR /doolah
# COPY . .
# RUN go mod download
# RUN GOOS=linux go build -ldflags="-s -w" -o ./bin ./app/main.go

# FROM alpine:3.12
# WORKDIR /usr/bin
# RUN mkdir /app
# COPY --from=build /doolah/bin /app
# EXPOSE 80
# RUN cd /app
# RUN mkdir storage
# ENTRYPOINT /app


# FROM golang:1.15.11
# WORKDIR /go/src/app
# COPY . .
# RUN go mod download
# CMD ["go", "run", "app/main.go"]

# EXPOSE 80

FROM golang:alpine AS build

WORKDIR /go/src/app

COPY . .
RUN go mod download
RUN GOOS=linux go build -ldflags="-s -w" -o ./server ./app/main.go

FROM alpine:3.10
WORKDIR /usr/bin
COPY --from=build /go/src/app/server .
RUN chmod +x ./server
EXPOSE 80
ENTRYPOINT ["./server"]