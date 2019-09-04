FROM amazonlinux:latest

RUN amazon-linux-extras install epel
RUN yum -y update && yum -y install clamav

CMD ["/usr/bin/freshclam"]
