# Orinoco
A data stream management system with two services:

## Sieve
A pub-sub server to accept and distribute data streams and provide basic
statistics on the state of the stream.

## Pump
A client to feed data streams to Sieve.

### Data Sources
Pump will consume log files and POST the data from each line to the
associated URL.

### Configuration
Pump can read from a YAML config file.

```yaml
streams:
  /opt/event/test1.log: http://example.com/streams/test1
  /opt/event/test2.log: http://example.com/streams/test2
```

## Tap
A client to consume data streams from Sieve.

## Litmus
A daemon to monitor data stream statistics through Sieve.

### Configuration
Litmus can read from a YAML config file. It supports everything that can be passed
in the command line (except config file paths), and is the only way to set up
custom event triggers.

```yaml
url: http://localhost:9966
triggers:
  trailing_average_per_minute_all_events_at_zero:
    event: "*"
    metric: trailing_average_per_minute
    condition: "==0"
    endpoint: http://example.com/events
  trailing_average_per_hour_test_event_under_100:
    event: test_event
    metric: trailing_average_per_hour
    condition: "<100"
    endpoint: http://example.com/events
```
