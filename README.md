# PostgresSQL to Elasticsearch Sync
![Blank diagram (2)](https://github.com/akshaykanawat/pgtoelastic/assets/39729121/f014d08e-e03d-4699-8344-76fdb3f80628)



This repository contains a system for real-time synchronization of data changes from a PostgresSQL database to Elasticsearch. The system leverages PostgresSQL triggers, Golang server, Kafka, and a Golang consumer for efficient and scalable synchronization.

## Overview

The system is designed to achieve real-time synchronization of changes made to a PostgresSQL database, propagating them to an Elasticsearch index. It uses PostgresSQL triggers to detect changes, a Golang server to listen for notifications, Kafka as a message broker, and a Golang consumer to sync data with Elasticsearch.

## Components

1. **PostgresSQL Database:**
   - Maintain the PostgresSQL database containing the data to be synchronized.

2. **PostgresSQL Trigger:**
   - Triggers within the PostgresSQL database notify the Golang server of changes.

3. **Golang Server:**
   - Listens for PostgresSQL notifications and acts as a Kafka producer.
   - Establishes a connection to the PostgresSQL database.

4. **Kafka:**
   - Serves as a message broker to decouple components and buffer messages.

5. **Go Consumer (Sync to Elasticsearch):**
   - Subscribes to the Kafka topic and processes messages.
   - Sync changes to the Elasticsearch index.

## Usage

### Prerequisites
- Go installed
- PostgresSQL server running
- Kafka broker accessible
- Elasticsearch cluster available

### Running the Docker Compose file to run kafka and elasticsearch locally

1. Make sure Docker is installed on your system.
2. Navigate to the root of the repository in the terminal.
3. Run the following command to start the Docker Compose services:

    ```bash
    docker-compose up -d
    ```

   This will start the PostgresSQL, Kafka, and other necessary services.

4. To stop the services, run:

    ```bash
    docker-compose down
    ```

### Creating a Kafka Topic

To create a Kafka topic, you can use the following command:

```bash
docker exec -it <your_kafka_container_id> kafka-topics --create --topic your_topic_name --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

### Build

Run the following commands to build and run the api server, kafka producer, and kafka consumer :

```bash
make start_server
make start_producer
make start_consumer
```
Run the following commands to generate dummy relational data in all tables in postgres  :

```bash
go run /_data_generation/dummy_data.go
```

### API Usage

Search Projects by User
```bash
curl --location 'localhost:8080/v1/projects/user/890'
```

Search Projects by Hashtag
```bash
curl --location 'localhost:8080/v1/projects/hashtags/Art'
```

Fuzzy Search Projects
```bash
curl --location 'localhost:8080/v1/projects/search' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "healthy-recipes",
    "description": "Description for"
}'

```