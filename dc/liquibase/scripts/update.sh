#!/bin/bash
echo "update start ......"
if [ x"$1" = x ]; then
    echo "请输入本次更新的Tag号"
    exit 1
fi
echo $liquibase_master_xml
if [ y"$2" != y ]; then
    echo "update -- label=$2 tag $1 ......"
    liquibase \
      --driver=$liquibase_driver \
      --url=$MYSQL_DB_URL \
      --username=$DB_USER \
      --password=$DB_PASSWORD \
      --classpath=/ \
      --changeLogFile=$liquibase_master_xml \
      --labels=$2 \
      --logLevel=info \
      update
else
    echo "update ......"
    liquibase \
      --driver=$liquibase_driver \
      --url=$MYSQL_DB_URL \
      --username=$DB_USER \
      --password=$DB_PASSWORD \
      --classpath=/ \
      --changeLogFile=$liquibase_master_xml \
      --logLevel=info \
      update
fi
echo "update end ......"

# 构建版本
echo "buildTag start ......"
echo "tag $1 ......"
liquibase \
  --driver=$liquibase_driver \
  --url=$MYSQL_DB_URL \
  --username=$DB_USER \
  --password=$DB_PASSWORD \
  --classpath=/ \
  --changeLogFile=$liquibase_master_xml \
  --logLevel=info \
  tag $1

echo "构建版本成功. 版本号：$1"
echo "buildTag end ......"