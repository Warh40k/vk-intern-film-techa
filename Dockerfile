FROM golang:1.22.1 AS build
WORKDIR /go/src/filmotecka/
COPY . /go/src/filmotecka/
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -installsuffix cgo -o /go/src/filmotecka/build/filmotecka /go/src/filmotecka/cmd/app/main.go

FROM busybox AS runtime
WORKDIR /app
COPY --from=build /go/src/filmotecka/build/filmotecka /app/
EXPOSE 8080/tcp
ENTRYPOINT ["./filmotecka"]
