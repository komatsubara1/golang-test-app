#!/bin/bash

USER=root
PASSWORD=password
USER_HOST=user_mysql_1
MASTER_HOST=master_mysql

TABLES=`mysql -u${USER} -p${PASSWORD} -h${USER_HOST} user -Nse "SHOW TABLES;" | xargs echo`
for t in ${TABLES}
do
    mysqldump --no-data --skip-add-drop-table --skip-comments -u${USER} -p${PASSWORD} -h${USER_HOST} user $t | sed -e 's/CREATE TABLE/CREATE TABLE IF NOT EXISTS/g' | sed -e 's/ AUTO_INCREMENT=[0-9]*//g' > /db/user/$t.sql
done

TABLES=`mysql -u${USER} -p${PASSWORD} -h${MASTER_HOST} master -Nse "SHOW TABLES;"`
for t in ${TABLES}
do
        mysqldump --no-data --skip-add-drop-table --skip-comments -u${USER} -p${PASSWORD} -h${MASTER_HOST} master $t | sed -e 's/CREATE TABLE/CREATE TABLE IF NOT EXISTS/g' | sed -e 's/ AUTO_INCREMENT=[0-9]*//g' > /db/master/$t.sql
done
