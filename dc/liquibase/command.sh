#!/bin/bash

export DB_URL="127.0.0.1:3306"
export DEFAULT_DB="go_backend"
export DB_USER="root"
export DB_PASSWORD="123456"
export MYSQL_DB_URL="jdbc:mysql://"$DB_URL"/"$DEFAULT_DB

COMMAND_SCRIPT_PATH=$(dirname $(readlink -f "${BASH_SOURCE[0]}"))
export liquibase_scripts=$COMMAND_SCRIPT_PATH"/scripts"
export liquibase_sql=$COMMAND_SCRIPT_PATH"/sqls"
export liquibase_master_xml=$COMMAND_SCRIPT_PATH"/master.xml"
export liquibase_driver="com.mysql.cj.jdbc.Driver"

if [ x"$1" = x ]; then
  echo "请输入操作参数：update"
  exit 1
fi

case $1 in
"update")
  echo "update..."
  if [ y"$2" = y ]; then
      echo "请输入本次更新的tag号"
      exit 1
  fi
  sh $liquibase_scripts"/update.sh" $2 $3
  ;;
"rollbackCount")
  echo "rollback..."
  if [ y"$2" = y ]; then
      echo "请输入回滚的次数"
      exit 1
  fi
  echo $liquibase_scripts"/rollback.sh"
  sh $liquibase_scripts"/rollbackCount.sh" $2
esac