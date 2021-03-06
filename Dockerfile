FROM golang:1.15.8 AS builder
LABEL maintainer="ericchou19831101@msn.com"

ARG version="local"
ARG author="Wen Zhou"
ARG app="gobra"
ARG release=true

ENV GOOS=linux \
    GO111MODULE="on" \
    CGO_ENABLED=0

COPY . /src
WORKDIR /src/logic

RUN echo ${version} \
    && go version \
    && go mod download  \
    && go build -ldflags "-w -s -X main.version=${version} -X main.author=${author} -X main.release=${release} -o ${app}"

FROM scratch 
COPY --from=builder /src/logic/${app} /
COPY --from=builder /src/template/ template
EXPOSE 8080
ENTRYPOINT ["/${app}"]