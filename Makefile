ENV=development
INPRODUCTION=false

DSN="user:password@tcp(db:3306)/test_db?allowNativePasswords=false&checkConnLiveness=false&collation=utf8mb4_unicode_ci&loc=Asia%2FTokyo&parseTime=true&maxAllowedPacket=0"
BINARY_NAME=api

API_PORT=8080

DB_DRIVER=mysql
DB_PORT=3306
MYSQL_DATABASE=test_db
MYSQL_HOST=db
MYSQL_USER=user
MYSQL_PASSWORD=password
MYSQL_ROOT_USER=root
MYSQL_ROOT_PASSWORD=root_password

# these three are for connecting to mysql, might be changed in production 
DB_CONNECT_RETRY=100
DB_TIMEOUT=5*time.Second
MAX_OPEN_DB_CONN=5 
MAX_IDLE_DB_CONN=5
MAX_DB_LIFE_TIME=5*time.Minute

PHPMYADMIN_PORT=8081

## build: Build binary
build:
	@echo "Building back end..."
	go build -o ${BINARY_NAME} ./cmd/api/
	@echo "Binary built!"

## run: builds and runs the application
run: build
	@echo "Starting back end..."
	@env DSN=${DSN} ENV=${ENV} go run ./cmd/api
	@echo "Back end started!"
## might be "@env DSN=${DSN} ./${BINARY_NAME} &" ?

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping back end..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped back end!"

## restart: stops and starts the running application
restart: stop start