.PHONY: build test test-coverage test-integration lint clean build-rpi-zero install deploy build-rpi-zero-deploy

DEPLOY_HOST ?= dsmr.local
DEPLOY_USER ?= root
INSTALL_PATH ?= /usr/local/bin
SYSTEMD_PATH ?= /etc/systemd/system
SSH_TARGET   := $(DEPLOY_USER)@$(DEPLOY_HOST)

build:
	go build ./...

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-integration:
	go test -tags integration ./...

lint:
	golangci-lint run ./...

clean:
	rm -f coverage.out dsmr

build-rpi-zero:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dsmr .

install:
	scp dsmr.service $(SSH_TARGET):$(SYSTEMD_PATH)/dsmr.service
	ssh $(SSH_TARGET) 'chmod 0644 $(SYSTEMD_PATH)/dsmr.service'
	ssh $(SSH_TARGET) 'systemctl daemon-reload'
	ssh $(SSH_TARGET) 'systemctl enable dsmr'
	ssh $(SSH_TARGET) 'systemctl restart dsmr'

deploy:
	ssh $(SSH_TARGET) 'systemctl stop dsmr'
	scp dsmr $(SSH_TARGET):$(INSTALL_PATH)/dsmr
	ssh $(SSH_TARGET) 'chmod +x $(INSTALL_PATH)/dsmr'
	ssh $(SSH_TARGET) 'systemctl start dsmr'

build-rpi-zero-deploy: build-rpi-zero deploy
