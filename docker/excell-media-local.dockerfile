FROM golang:1.18-stretch

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends

ENV IMAGEMAGICK_VERSION=6.9.10-11

RUN cd && \
	wget https://github.com/ImageMagick/ImageMagick6/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib

# Setup document root
RUN mkdir -p /app/code
RUN mkdir -p /app/logs
RUN mkdir -p /app/k6
RUN mkdir -p /app/tmp
RUN mkdir -p /media

WORKDIR /app/code

COPY app/media/go.mod /app/code
COPY app/media/go.sum /app/code
COPY app/media/reflex.conf /app/code

RUN go mod download

RUN ["go", "install", "github.com/cespare/reflex@latest"]

COPY app/media/src src
COPY docker/env/app-local.env .env

EXPOSE 8080

ENTRYPOINT ["reflex", "-c", "./reflex.conf"]