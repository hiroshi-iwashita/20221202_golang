### setup docker
$ docker network create golang_test_network
$ docker-compose build
$ docker-compose up

### create go mod file
$ go mod init ~

### connect to mysql
$ go get -u github.com/go-sql-driver/mysql@1.7.0

## migrate to mysql
$ go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@latest
### create migration files
$ migrate create -ext sql -dir build/db/migrations -seq create_users[table_name]
### create tables
$ migrate -path build/db/migrations -database "mysql://root:root_password@tcp(db:3306)/test_database?multiStatements=true&loc=Asia%2FTokyo" up 1
### drop tables
$ migrate -path build/db/migrations -database "mysql://root:root_password@tcp(db:3306)/test_database?multiStatements=true&loc=Asia%2FTokyo" down 1

- Uses [chi router](https://github.com/go-chi/chi)
$ go get -u github.com/go-chi/chi/v5