#!/bin/bash
set -ex
cd `dirname $0`

export MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
export MYSQL_PORT=${MYSQL_PORT:-3306}
export MYSQL_USER=${MYSQL_USER:-isuconp}
export MYSQL_DBNAME=${MYSQL_DBNAME:-isuconp}
export MYSQL_PWD=${MYSQL_PASS:-isuconp}
export LANG="C.UTF-8"

cat 0_Schema.sql | mysql --defaults-file=/etc/mysql/my.cnf -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER $MYSQL_DBNAME