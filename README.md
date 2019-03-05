# rundeck-exporter

Rundeck exporter metrics to Prometheus format.

Now we are supporting the parse of metrics available on http://rundeck.local/metrics/metrics

## Metrics

The metrics will be exported prefixed with `rundeck_` name (Prometheus namespace).

## Changelog

`0.1.0`:

- Initial release with basic metrics translated from `/metrics/metrics`