FROM golang:1.21.4-bookworm AS schemalex_build

# 最新のタグが随分古いので、masterのコミットハッシュを直接指定してバージョンの固定を行っている
RUN git clone https://github.com/schemalex/schemalex \
    && cd schemalex && git checkout -f bf2603ccbd8b4d7a012ec95827dcfd6689a741e7 \
    && make build \
    && cp artifacts/$(grep 'const Version' schemalex.go | awk '{print $4}' | sed -e s/\"//g)/schemalex_linux_amd64/* /usr/local/bin

RUN apt-get update && apt-get install --no-install-recommends -y lsb-release gnupg wget gettext \
    && wget https://dev.mysql.com/get/mysql-apt-config_0.8.22-1_all.deb \
    && DEBIAN_FRONTEND=noninteractive dpkg --force-confdef --force-confold -i mysql-apt-config_0.8.22-1_all.deb \
    && apt-get update --allow-insecure-repositories && apt-get install --no-install-recommends --allow-unauthenticated -y mysql-client

COPY ./docker/db_migrate_job/migrate.sh ./
COPY ./docker/db_migrate_job/export.sh ./
COPY ./db /db