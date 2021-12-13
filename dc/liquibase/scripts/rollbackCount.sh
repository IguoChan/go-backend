#!/bin/bash

echo "rollback start ......"
if [ x"$1" = x ]; then
    echo "请输入回滚的次数"
    exit 1
fi

echo "rollback $1 ......"
liquibase \
  --driver=$liquibase_driver \
  --url=$MYSQL_DB_URL \
  --username=$DB_USER \
  --password=$DB_PASSWORD \
  --classpath=/ \
  --changeLogFile=$liquibase_master_xml \
  --logLevel=info \
  rollbackCount $1

echo "rollback $1 end"