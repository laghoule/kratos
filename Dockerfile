FROM golang:1.18-alpine AS dep
WORKDIR /src/
COPY . .
RUN go get -d -v

FROM dep AS build
ARG VERSION "devel"
ARG BUILD_DATE ""
ARG GIT_COMMIT ""
ARG GIT_REF ""
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/laghoule/kratos/cmd.version=$VERSION' -X 'github.com/laghoule/kratos/cmd.buildDate=$BUILD_DATE' -X 'github.com/laghoule/kratos/cmd.gitCommit=$GIT_COMMIT' -X 'github.com/laghoule/kratos/cmd.gitRef=$GIT_REF'" -o kratos .

FROM alpine:3.15
LABEL org.opencontainers.image.source https://github.com/laghoule/kratos
COPY --from=build /src/kratos /usr/bin/
ENTRYPOINT ["/usr/bin/kratos"]
CMD ["-h"]