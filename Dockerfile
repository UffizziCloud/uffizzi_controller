FROM golang:1.17

RUN apt-get update -y

RUN apt-get install -y \
    vim

# Install kubectl for handy debug purposes only (https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl-binary-with-curl-on-linux)
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.44.0

ARG SENTRY_RELEASE
ENV SENTRY_RELEASE=${SENTRY_RELEASE:-}

WORKDIR /app
COPY . .
RUN go install ./cmd/controller/...

EXPOSE 8080

CMD ["bash", "-c", "/go/bin/controller"]
