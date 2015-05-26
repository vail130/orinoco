# Orinoco
A data stream management system with two services:

## Sieve
A pub-sub server to accept and distribute data streams and provide basic
statistics on the state of the stream.

## Tap
A client to consume data streams from Sieve.

## Litmus
A daemon to monitor data stream statistics through Sieve.

### Configuration
Orinoco can read from a config file. It supports everything that can be passed
in the command line (except config file paths), and is the only way to set up
custom event triggers.

```yaml
host: localhost
port: 9966
triggers:
  no_events:
    event: "*"
    metric: trailing_average_per_minute
    condition: ==0
    endpoint: http://example.com/events
  no_test_events:
    event: test_event
    metric: trailing_average_per_hour
    condition: <100
    endpoint: http://example.com/events

```
