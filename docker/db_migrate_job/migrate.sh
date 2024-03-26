#!/bin/sh

USER=root
PASSWORD=password
USER_HOST_1=user_mysql_1
USER_HOST_2=user_mysql_2
MASTER_HOST=master_mysql

echo "create database"
mysql -u${USER} -p${PASSWORD} -h${USER_HOST_1} -e "CREATE DATABASE IF NOT EXISTS user"
mysql -u${USER} -p${PASSWORD} -h${USER_HOST_2} -e "CREATE DATABASE IF NOT EXISTS user"
mysql -u${USER} -p${PASSWORD} -h${MASTER_HOST} -e "CREATE DATABASE IF NOT EXISTS master"

echo "combine sql file"
cat /db/user/* > user_combine.sql
cat /db/master/* > master_combine.sql

echo "deploy new table"
mysql -u${USER} -p${PASSWORD} -h${USER_HOST_1} user < user_combine.sql
mysql -u${USER} -p${PASSWORD} -h${USER_HOST_2} user < user_combine.sql
mysql -u${USER} -p${PASSWORD} -h${MASTER_HOST} master < master_combine.sql

echo "dump"
mysqldump --no-data -u${USER} -p${PASSWORD} -h${USER_HOST_1} user > user.sql
mysqldump --no-data -u${USER} -p${PASSWORD} -h${MASTER_HOST} master > master.sql

echo "migrate diff table"
schemalex -o user_diff.sql user.sql user_combine.sql
schemalex -o master_diff.sql master.sql master_combine.sql

echo "deploy"
mysql -u${USER} -p${PASSWORD} -h${USER_HOST_1} user < user_diff.sql
mysql -u${USER} -p${PASSWORD} -h${MASTER_HOST} master < master_diff.sql