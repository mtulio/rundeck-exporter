# rundeck-exporter

Rundeck exporter metrics to Prometheus format.

This exporter uses the [GO Library](https://github.com/lusis/go-rundeck) for Rundeck to retrieve metrics from server and exposing it to Prometheus pulls.

Now we are supporting the parse of metrics available on http://rundeck.local/metrics/metrics.

## Metrics

The metrics will be exported prefixed with `rundeck_` name (Prometheus namespace).

## Dependencies

* Environment Variables (mandatory)

`RUNDECK_VERSION` : Version of the API to collect metrics (tested only in 18)

`RUNDECK_API_URL` : Rundeck URL - including the schema (http/s)

* Environment Variables (optional)

`RUNDECK_API_TOKEN=xzy1231212`

## Usage

1. Export the environment for your Rundeck:

> These vars is not supported yet as argument, and is required by the library.

```bash
export RUNDECK_VERSION=18
export RUNDECK_API_URL=http://rundeck.local/metrics/metrics
export RUNDECK_API_TOKEN=xzy1231212
```

Show Rundeck metrics:

```bash
./bin/rundeck-exporter -exporter -rundeck.user $RUN_USER -rundeck.pass $RUN_PASS -no-verify-ssl
```

> Sample output for `$ curl localhost:9802/metrics`:

```log
$ curl -s localhost:9802/metrics |grep rundeck
# HELP rundeck_api_requests_requestTimer Rundeck metrics Timer
# TYPE rundeck_api_requests_requestTimer gauge
rundeck_api_requests_requestTimer{type="Count"} 2438
rundeck_api_requests_requestTimer{type="M15Rate"} 0.01448732196209572
rundeck_api_requests_requestTimer{type="M1Rate"} 6.339432326718685e-06
rundeck_api_requests_requestTimer{type="M5Rate"} 0.0023874288268983237
rundeck_api_requests_requestTimer{type="Max"} 0.262770125
rundeck_api_requests_requestTimer{type="Mean"} 0.06049109607684825
rundeck_api_requests_requestTimer{type="MeanRate"} 0.022784876531781885
rundeck_api_requests_requestTimer{type="Min"} 0.046306241000000005
rundeck_api_requests_requestTimer{type="P50"} 0.05766555600000001
rundeck_api_requests_requestTimer{type="P75"} 0.063260466
rundeck_api_requests_requestTimer{type="P95"} 0.08199031544999998
rundeck_api_requests_requestTimer{type="P98"} 0.09199199426
rundeck_api_requests_requestTimer{type="P99"} 0.09786174672000002
rundeck_api_requests_requestTimer{type="P999"} 0.26257600741000003
rundeck_api_requests_requestTimer{type="Stddev"} 0.013228965044100493
# HELP rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_activeRequests Rundeck metrics Counter
# TYPE rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_activeRequests gauge
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_activeRequests 1
# HELP rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests Rundeck metrics Timer
# TYPE rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests gauge
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="Count"} 42725
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="M15Rate"} 1.0070459897192932
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="M1Rate"} 0.5436722643492509
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="M5Rate"} 0.8290200180890379
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="Max"} 2.121951113
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="Mean"} 0.05591648869455253
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="MeanRate"} 0.3959690854986382
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="Min"} 0.00025913
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P50"} 0.039484237000000005
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P75"} 0.045085134750000005
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P95"} 0.17360742895
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P98"} 0.23641499398
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P99"} 0.3665500127000004
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="P999"} 2.0943520583400037
rundeck_com_codahale_metrics_servlet_AbstractInstrumentedFilter_requests{type="Stddev"} 0.09086212691802696
[...]
```

## USAGE IN DOCKER

Show Rundeck metrics running in docker;.

```bash
docker run -p 9802:9802 -id mtulio/rundeck-exporter:latest \
    -rundeck.user=$RUN_USER \
    -rundeck.pass=$RUN_PASS
```

* Docker Compose definition

```YAML
# Rundeck exporter - https://github.com/mtulio/rundeck-exporter
  exporter:
    image: mtulio/rundeck-exporter:v0.1.1
    command:
        - -rundeck.email=myUser
        - -rundeck.password=myPass
    ports:
      - 9802:9802
    networks:
      - net
    deploy:
      resources:
        limits:
          cpus: "0.2"
          memory: 256M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9802/"]
      interval: 5s
      timeout: 2s
      retries: 3
```

## Changelog

See on each [release](https://github.com/mtulio/rundeck-exporter/releases).
