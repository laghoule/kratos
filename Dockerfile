FROM golang:1.17-alpine AS dep
WORKDIR /src/
COPY . .
RUN go get -d -v

FROM dep AS build
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o kratos .

FROM alpine:3.15
COPY --from=build /src/kratos /usr/bin/
ENTRYPOINT ["/usr/bin/kratos"]
CMD ["-h"]