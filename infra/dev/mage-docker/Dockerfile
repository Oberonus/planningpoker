FROM golang:1.17

RUN apt-get update && \
    apt-get install -y \
    --no-install-recommends \
    ca-certificates \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

# Download and install mage file into bin path
ARG MAGE_VERSION=1.11.0
RUN wget -qc https://github.com/magefile/mage/releases/download/v${MAGE_VERSION}/mage_${MAGE_VERSION}_Linux-64bit.tar.gz -O - | tar -xz -C /usr/bin mage

# Download and install golangci-lint into go bin path
ARG GOLANGCILINT_VERSION=1.43.0
RUN wget -qc https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh -O - | /bin/sh -s -- -b "$(go env GOPATH)/bin" v${GOLANGCILINT_VERSION}

ENV APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" && \
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
