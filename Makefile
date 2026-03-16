.PHONY: up down logs ps kafka-topics kafka-create-topics

up:
	docker compose up -d
	@echo "Wait for service readiness"
	@docker compose ps

down:
	docker compose down

down-clean:
	docker compose down -v
	@echo "All data deleted"

logs:
	docker compose logs -f

# Service's logs: make logs-kafka
logs-%:
	docker compose logs -f $*

#containers status
make ps:
	docker compose ps

# kafka's topic list
kafka-topics:
	docker exec clickstream-kafka /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--list

# for manual test
kafka-create-topic:
	docker exec clickstream-kafka /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--topic events \
		--partitions 3 \
		--replication-factor 1 \
		--if-not-exists

check:
	@echo "=== PostgreSQL ===" && \
	docker exec clickstream-postgres pg_isready -U clickstream && \
	echo "=== Redis ===" && \
	docker exec clickstream-redis redis-cli ping && \
	echo "=== Kafka ===" && \
	docker exec clickstream-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list && \
	echo "=== ClickHouse ===" && \
	curl -s http://localhost:8123/ping && \
	echo "\n All service are working"