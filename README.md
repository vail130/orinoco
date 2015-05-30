# Orinoco
A data stream management system with two services:

## Sieve
A pub-sub server to accept and distribute data streams and provide basic
statistics on the state of the stream.

## Pump
A client to feed data streams to Sieve.

## Tap
A client to consume data streams from Sieve.

## Litmus
A daemon to monitor data stream statistics through Sieve.

### Configuration
Orinoco can read from a config file. It supports everything that can be passed
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
