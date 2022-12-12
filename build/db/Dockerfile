FROM mysql:8.0 as dev

# MySQL設定ファイルをイメージ内にコピー
COPY my.cnf /etc/mysql/conf.d/my.cnf

EXPOSE 3306

# docker runに実行される
CMD ["mysqld"]