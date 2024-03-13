FROM golang:1.21 AS build
WORKDIR /go/src/cloudmgm/
COPY . /go/src/cloudmgm/
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -installsuffix cgo -o /go/src/cloudmgm/build/cloudmgm /go/src/cloudmgm/cmd/app/main.go

FROM busybox AS runtime
WORKDIR /app
COPY --from=build /go/src/ewallet/build/cloudmgm /app/
EXPOSE 8080/tcp
ENTRYPOINT ["./ewallet"]
