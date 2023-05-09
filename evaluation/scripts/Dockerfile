FROM ubuntu:23.04

# packages
RUN apt-get -y update \
    && apt-get -y install build-essential libssl-dev git zlib1g-dev 
RUN apt-get -y install unzip wget lua5.2 liblua5.2

# wrk2
RUN git clone https://github.com/giltene/wrk2.git wrk2 \
    && cd wrk2 \
    && make
RUN mv wrk2/wrk /usr/local/bin/wrk2

# amazon CLI
RUN wget "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.44.zip"
RUN	unzip awscli-exe-linux-x86_64-2.0.44.zip && ./aws/install

COPY . /scripts

# lua rocks
#RUN wget https://luarocks.org/releases/luarocks-3.9.2.tar.gz
#RUN tar zxpf luarocks-3.9.2.tar.gz \
#    && cd luarocks-3.9.2 \
#    && ./configure && make && make install \
RUN cd /scripts \
    && wget https://raw.githubusercontent.com/tiye/json-lua/main/JSON.lua

HEALTHCHECK CMD exit 0