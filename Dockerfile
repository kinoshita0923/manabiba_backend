FROM golang:latest

WORKDIR /var/www/manabiba/src

COPY ./ /var/www/manabiba

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
EXPOSE 8080

ENV BASH_ENV ~/.bashrc

ENV VOLTA_HOME /root/.volta

ENV PATH $VOLTA_HOME/bin:$PATH

RUN curl https://get.volta.sh | bash

RUN volta install node@18.12.1

RUN go get -u github.com/cosmtrek/air && go build -o /go/bin/air github.com/cosmtrek/air

RUN mkdir -p /etc/mysql/conf.d

CMD ["air", "-c", ".air.toml"]