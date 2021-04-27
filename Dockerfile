# syntax=docker/dockerfile:1.2
ARG version="local"
ARG app="gobra"

FROM golang:1.15.8 AS builder
LABEL maintainer="ericchou19831101@msn.com"

ARG author="Wen_Zhou"
ARG release=true

ENV GOOS=linux \
    GO111MODULE="on" \
    CGO_ENABLED=0

COPY . /src
WORKDIR /src/logic

RUN echo ${version} \
    && go version \
    && go mod download
RUN go build -ldflags="-w -s -X main.version=${version} -X main.author=${author}" -o ${app}



FROM scratch

ARG id=$id
ARG key=$key

ENV AWS_ACCESS_KEY_ID=$id
ENV AWS_SECRET_ACCESS_KEY=$key

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/logic/gobra /
COPY --from=builder /src/template/ template
COPY --from=builder /src/html/ html

EXPOSE 8080

ENTRYPOINT ["/${app}"]