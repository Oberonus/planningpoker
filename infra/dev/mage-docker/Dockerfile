FROM golang:1.17

ARG TARGET_ARCH
RUN echo Building mage CI image for "${TARGET_ARCH}" architecture

RUN apt-get update && \
    apt-get install -y \
    --no-install-recommends \
    ca-certificates \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Download and install mage file into bin path
ARG MAGE_VERSION=1.12.1

RUN case ${TARGET_ARCH} in \
         "amd64")  MAGE_ARCH=64bit  ;; \
         "x86_64") MAGE_ARCH=64bit  ;; \
         "arm64")  MAGE_ARCH=ARM64  ;; \
    esac && \
    wget -qc https://github.com/magefile/mage/releases/download/v${MAGE_VERSION}/mage_${MAGE_VERSION}_Linux-${MAGE_ARCH}.tar.gz -O - | tar -xz -C /usr/bin mage

# Download and install golangci-lint into go bin path
ARG GOLANGCILINT_VERSION=1.44.2
RUN wget -qc https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh -O - | /bin/sh -s -- -b "$(go env GOPATH)/bin" v${GOLANGCILINT_VERSION}

ENV APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1
RUN case ${TARGET_ARCH} in \
         "amd64")  DOCKER_ARCH=amd64  ;; \
         "x86_64") DOCKER_ARCH=amd64  ;; \
         "arm64")  DOCKER_ARCH=arm64  ;; \
    esac && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=${DOCKER_ARCH}] https://download.docker.com/linux/debian $(lsb_release -cs) stable" && \
    apt-get -y update && \
    apt-get install -y docker-ce \
    --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

RUN chmod 777 -R /go

RUN mkdir /home/mage && chmod 777 /home/mage
ENV HOME /home/mage
WORKDIR /home/mage

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]
