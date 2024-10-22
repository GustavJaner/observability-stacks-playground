.PHONY: up down

up:
		docker-compose up -d

down:
		docker-compose down

get-otelc-internal-prom-metrics:
		curl localhost:8888/metrics

get-otelc-promexporter-metrics:
		curl localhost:8889/metrics
