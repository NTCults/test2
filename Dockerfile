FROM golang:latest AS build-env

WORKDIR /src

COPY ./ /src

RUN go mod download

RUN go build -o /src/consumer2 /src/consumer

EXPOSE 5672

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait

RUN chmod +x /wait

CMD /wait && /src/consumer2