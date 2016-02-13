# Orinoco !(Build Status)[https://travis-ci.org/vail130/orinoco.svg?branch=master]
Orinoco forks your data stream, maintains statistics on stream activity,
monitors the state of the data stream and triggers custom, configurable
events. Its main components include Sieve, which forks and forwards the
data stream, and Litmus which monitors the stream and triggers events.

## Sieve
A server hosting an HTTP API that keeps track of statistics on the state of
the data stream and forwards data to configured endpoints.

## Litmus
A client that monitors data stream statistics through Sieve. It gets
statistics from a designated Sieve server every second, and triggers events
based on configured conditions.

```yaml
port: 9966
streams:
  - /go/src/github.com/vail130/orinoco/artifacts/logs/
  - http://localhost:9966/streams/
  - ws://localhost:9966/streams/
triggers:
  -
    stream: test2
    metric: minute_to_date
    condition: ">10"
    endpoint: http://localhost:9966/streams/test2_stream_more_than_10_per_minute/
  -
    stream: test_stream
    metric: trailing_average_per_hour
    condition: "<100"
    endpoint: http://example.com/trailing_average_per_hour_test_stream_under_100/
```

# Testing
Install docker or boot2docker, then run tests from the make target:

`make test`
