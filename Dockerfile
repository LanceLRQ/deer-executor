# 这是基于Ubuntu 18.04评测机运行环境，自带了GCC包
# GCC版本是5.4，并且安装了Java 1.8 (32位), Python3.6,
# Python2.7, Go 1.10, NodeJs 8.10, ruby 2.5, php 7.2
FROM buildpack-deps:bionic-scm

RUN sed -i 's/http:\/\/archive\.ubuntu\.com\/ubuntu\//http:\/\/mirrors\.ustc\.edu\.cn\/ubuntu\//g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y \
    python3 python3-pip language-pack-zh-hans language-pack-zh-hans-base \
	vim supervisor

ENV LANG zh_CN.UTF-8

RUN export DEBIAN_FRONTEND=noninteractive && apt-get install -y tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

RUN dpkg --add-architecture i386 && apt-get update && apt-get install -y lib32ncurses5 lib32z1 libc6:i386 libc6-i386 lib32z1 lib32stdc++6 openjdk-8-jdk:i386

RUN apt-get update && apt-get install -y golang-go
RUN apt-get update && apt-get install -y nodejs
RUN apt-get update && apt-get install -y ruby
RUN apt-get update && apt-get install -y php7.2 \
	php7.2-json php7.2-mbstring php7.2-common php7.2-readline php7.2-xml

#RUN wget -q https://packages.microsoft.com/config/ubuntu/18.04/packages-microsoft-prod.deb && dpkg -i packages-microsoft-prod.deb
#RUN apt-get install -y apt-transport-https && apt-get update && apt-get install -y dotnet-sdk-2.1

RUN rm /usr/bin/java && ln -s /usr/lib/jvm/java-8-openjdk-i386/bin/java /usr/bin/java

ADD lib/testlib /testlib

CMD ["supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
