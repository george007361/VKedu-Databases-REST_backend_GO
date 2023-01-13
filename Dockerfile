FROM golang:latest AS build

COPY . /app/

WORKDIR /app/

RUN go build cmd/main.go

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER george WITH SUPERUSER PASSWORD '12345678';" &&\
    createdb -O postgres george_forum_db &&\
    psql -f db/db_init.sql -d george_forum_db &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
COPY --from=build /app/main .

EXPOSE 5000

CMD service postgresql start && ./main