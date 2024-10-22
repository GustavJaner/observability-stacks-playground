.PHONY: up down

up: # Builds, (re)creates, starts, and attaches to containers for all services
		docker-compose up -d # -d (Detached) flag runs the containers in the background

up-otel-collector:
		docker-compose up -d --no-deps otel-collector

up-prometheus:
		docker-compose up -d --no-deps prometheus

up-grafana:
		docker-compose up -d --no-deps grafana

down: # Stops and removes all the containers
		docker-compose down

down-otel-collector:
		docker-compose down otel-collector

down-prometheus:
		docker-compose down prometheus

down-grafana:
		docker-compose down grafana

get-otelc-internal-prom-metrics:
		curl localhost:8888/metrics
