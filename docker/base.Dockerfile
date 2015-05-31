FROM phusion/baseimage:0.9.16

RUN apt-get update && \
	apt-get install -y git golang && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/sbin/my_init"]
