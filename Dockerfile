FROM golang

WORKDIR /app

COPY . /app

RUN go build -o api

EXPOSE 8081

ENTRYPOINT [ "./api" ]