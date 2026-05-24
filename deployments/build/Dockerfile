# Builder
FROM golang:1.24-alpine as builder
ARG appName
ARG appPath
ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache gcc musl-dev

RUN mkdir /${appName}
COPY . /${appName}
WORKDIR /${appName}
COPY go.mod .
COPY go.sum .

RUN --mount=type=ssh GO111MODULE=on go mod download
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -a -installsuffix cgo -o ${appName} ${appPath}

# Cert
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# Distribution
FROM alpine:latest
ARG appName
WORKDIR /${appName}

COPY --from=builder /${appName}/${appName} ./app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /${appName}/web ./web
ENV TZ=Asia/Jakarta
RUN apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone
ENTRYPOINT ["./app"]