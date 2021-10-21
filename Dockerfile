FROM golang:1.17-alpine3.14 AS dep
WORKDIR /src/
COPY . .
RUN go get -d -v

FROM dep AS build
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o k8dep .

FROM alpine:3.14
COPY --from=build /src/k8dep /usr/bin/
ENTRYPOINT ["/usr/bin/k8dep"]
CMD ["-h"]