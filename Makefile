docker-container = major7/mqtt-bq
dist_mqtt_bq = dist/mqtt-bq
version = 0.1

all: build-container

build-app:
	@go fmt *.go
	@go build -o $(dist_mqtt_bq) main.go devices.go

build-container:
	@go fmt *.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(dist_mqtt_bq) main.go devices.go
	@docker build -t $(docker-container):$(version) .

dev:
	docker run -ti --rm $(docker-container):$(version)

run:
	@go run main.go

clean:
	@go clean
	@rm -fv dist/*
	@-docker rmi $(docker-container) 2>/dev/null

deploy:
	docker push $(docker-container):$(version)
