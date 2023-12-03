build_api_server:
	echo "building server..." & bash build_api.sh

build_kafka_producer:
	echo "building server..." & bash build_producer.sh

build_kafka_consumer:
	echo "building server..." & bash build_consumer.sh

run_server:
	echo "running server..." & ./build/app

run_producer:
	echo "running server..." & ./build/producer

run_consumer:
	echo "running server..." & ./build/consumer

start_server: build_api_server run_server
start_producer: build_kafka_producer run_producer
start_consumer: build_kafka_consumer run_consumer
