FROM centos:7.3.1611
ENV WKHTML_URL https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download
ENV WKHTML_VERSION 0.12.5
ENV WKHTML_RPM wkhtmltox-${WKHTML_VERSION}-1.centos7.x86_64.rpm
ENV MYPROJECT https://github.com/happy2048/genpdf.git
RUN yum install epel-release -y
RUN yum install git \
	wget \
	curl \
	which \
 	xorg-x11-server-Xvfb \
	xorg-x11-fonts-Type1 \
	xorg-x11-fonts-100dpi -y
WORKDIR /tmp
RUN  wget -c ${WKHTML_URL}/${WKHTML_VERSION}/${WKHTML_RPM} && \
	yum localinstall -y  ${WKHTML_RPM}

WORKDIR /root

RUN git clone $MYPROJECT && \ 
	cd /root/genpdf && \
	cp scripts/wkhtmltopdf.sh /usr/bin && \
	cp -ar fonts /usr/share && \
	cp genpdf-server /usr/bin && \
	cp -ar html-template /usr/local && \
	chmod +x /usr/bin/genpdf-server && \
	chmod +x /usr/bin/wkhtmltopdf.sh
ENTRYPOINT ["genpdf-server"]
