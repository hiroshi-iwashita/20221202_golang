FROM mysql:8.0 as dev

# MySQL設定ファイルをイメージ内にコピー
COPY ./my.cnf /etc/mysql/conf.d/my.cnf

# COPY ./init/Dump_db.sql /docker-entrypoint-initdb.d/

EXPOSE 3306

# docker runに実行される
CMD ["mysqld"]