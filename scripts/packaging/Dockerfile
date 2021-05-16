FROM ubuntu:21.04

RUN apt-get update &&						\
    DEBIAN_FRONTEND=noninteractive apt-get install -y		\
    		       build-essential				\
    	    	       debhelper				\
		       devscripts				\
		       wget &&					\
    wget https://golang.org/dl/go1.16.4.linux-amd64.tar.gz &&	\
    tar -C /usr/local -xzf go1.16.4.linux-amd64.tar.gz &&  	\
    rm go1.16.4.linux-amd64.tar.gz &&		       		\
    rm -rf /var/cache/apt/archives

ADD package.sh /package.sh
RUN chmod +x package.sh

ENV PATH "${PATH}:/usr/local/go/bin"
ENV VERSION 0.25
ENV PACKAGE_NAME mountmond

VOLUME /src/
ENTRYPOINT ["/package.sh"]
