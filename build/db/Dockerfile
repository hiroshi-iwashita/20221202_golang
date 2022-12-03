FROM mysql:8.0 as dev

COPY Dump_db.sql /docker-entrypoint-initdb.d/

EXPOSE 3306

# MySQL設定ファイルをイメージ内にコピー
# COPY my.cnf /etc/mysql/conf.d/

# docker runに実行される
# CMD ["mysqld"]