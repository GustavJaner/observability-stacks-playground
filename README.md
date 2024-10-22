# otel-test
OpenTelemetry local test environment

## Instructions
1. Run `make up` in [`otelc-prom-grafana/docker`](./otelc-prom-grafana/docker)
2. Go to http://localhost:3000 and log in with username: "admin" password: "1337"
3. Run `make run` in [`telemetry-producers/otel-go-sdk/service-dice`](./telemetry-producers/otel-go-sdk/service-dice)
4. Go to http://localhost:8080/rolldice and refresh the page to record new metric data points (Increment metric counters)
5. Check how counters increase http://localhost:3000/d/fe1obkmx72hhcc/otel-metrics
