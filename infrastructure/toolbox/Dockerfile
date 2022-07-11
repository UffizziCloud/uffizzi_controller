FROM docker:19.03.13

RUN apk update
RUN apk add --no-cache \
    python3-dev \
    py3-pip \
    libffi-dev \
    openssl-dev \
    gcc \
    libc-dev \
    make \
    jq \
    git \
    openssh \
    bash \
    vim \
    curl \
    groff \
    cargo

RUN pip3 install "urllib3<1.25"
RUN pip3 install pyrsistent==0.16.1
RUN pip3 install docker-compose

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
 && chmod +x ./kubectl \
 && mv ./kubectl /usr/local/bin/kubectl

ARG CLOUD_SDK_VERSION=266.0.0

ENV PATH /google-cloud-sdk/bin:$PATH
RUN apk --no-cache --update add curl python2 py-crcmod libc6-compat make less && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    ln -s /lib /lib64 && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image

RUN pip3 install --upgrade awscli==1.19.41
RUN pip3 install azure-cli==2.27.1

RUN curl -sL https://sentry.io/get-cli/ | bash

RUN mkdir /app
WORKDIR /app

