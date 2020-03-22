FROM 4brp6gm1.mirror.aliyuncs.com/library/golang:1.14

COPY . /GetAnything-Server
WORKDIR /GetAnything-Server/
ENV GOPROXY=https://goproxy.io
RUN go build ./
CMD chmod +x ./GetAnything-Server && ./GetAnything-Server