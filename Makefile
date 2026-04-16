APP_NAME := weather
CMD_DIR := ./cmd/http
BIN_DIR := ./bin
PORT ?= 8080
CITY ?= Almaty
COUNTRY ?= Kazakhstan

.PHONY: run build fmt tidy curl-health curl-city curl-country curl-top

run:
	go run $(CMD_DIR)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

curl-city:
	curl -s "http://localhost:$(PORT)/weather/$(CITY)"

curl-country:
	curl -s "http://localhost:$(PORT)/weather/country/$(COUNTRY)"

curl-top:
	curl -s "http://localhost:$(PORT)/weather/country/$(COUNTRY)/top"
