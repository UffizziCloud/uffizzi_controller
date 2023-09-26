FROM docker:19.03.13

ENV CRYPTOGRAPHY_DONT_BUILD_RUST=1
ENV USE_GKE_GCLOUD_AUTH_PLUGIN=True

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

RUN pip3 install --upgrade pip
RUN pip3 install "urllib3<1.25"
RUN pip3 install pyrsistent==0.16.1
RUN pip3 install "cython<3.0.0" wheel && pip install pyyaml==5.4.1 --no-build-isolation
RUN pip3 install docker-compose
RUN echo "cryptography==3.3.2" > /tmp/requirements.txt \
    && pip3 install -U docker-compose -r /tmp/requirements.txt \
    && rm /tmp/requirements.txt

RUN wget -O /usr/local/bin/semver \
  https://raw.githubusercontent.com/fsaintjacques/semver-tool/master/src/semver \
  && chmod +x /usr/local/bin/semver  

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

ARG HELM_VERSION=v3.12.2

RUN wget https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && mv ./linux-amd64/helm /usr/local/bin/helm \
    && rm helm-${HELM_VERSION}-linux-amd64.tar.gz

ARG CLOUD_SDK_VERSION=412.0.0

ENV PATH /google-cloud-sdk/bin:$PATH
RUN apk --no-cache --update add curl python3 py-crcmod libc6-compat make less && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    ln -s /lib /lib64 && \
    gcloud components install gke-gcloud-auth-plugin && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image

RUN pip3 install --upgrade awscli==1.19.41

RUN curl -sL https://sentry.io/get-cli/ | bash

RUN mkdir /app
WORKDIR /app
