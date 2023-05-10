FROM ubuntu:23.04

# packages
RUN apt-get -y update \
    && apt-get -y install build-essential libssl-dev git zlib1g-dev 
RUN apt-get -y install unzip wget lua5.2 liblua5.2 vim curl

# wrk2
RUN git clone https://github.com/giltene/wrk2.git wrk2 \
    && cd wrk2 \
    && make \
    && mv wrk /usr/local/bin/wrk2

# go 1.18.10
RUN wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz 
RUN tar -C /usr/local -xpf go1.18.10.linux-amd64.tar.gz \
     && rm go1.18.10.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"
# goimports
RUN go install golang.org/x/tools/cmd/goimports@latest

# amazon CLI
RUN wget "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.44.zip"
RUN	unzip awscli-exe-linux-x86_64-2.0.44.zip && ./aws/install && rm awscli-exe-linux-x86_64-2.0.44.zip

# serverless framework
RUN curl -o- -L https://slss.io/install | bash 

RUN export PATH="$HOME/.serverless/bin:$PATH" >> $HOME/.bashrc

# copy entire project
COPY . /nubes

RUN cd /nubes/evaluation/scripts \
    && wget https://raw.githubusercontent.com/tiye/json-lua/main/JSON.lua

HEALTHCHECK CMD exit 0