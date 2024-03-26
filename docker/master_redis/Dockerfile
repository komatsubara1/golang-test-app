FROM redis:6.2

RUN apt-get update && apt-get install -y redis-server

EXPOSE 6379

ENTRYPOINT ["/usr/bin/redis-server", "--protected-mode no"]
