all: build-zookeeper build-kafka up

build-kafka:
	docker build -t kafka-jmx -f Dockerfile.kafka .

build-zookeeper:
	docker build -t zookeeper-jmx -f Dockerfile.zookeeper .

up:
	docker compose up -d

down:
	docker compose down
