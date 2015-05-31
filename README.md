# Orinoco
Orinoco allows monitoring of the state of a data stream, with
custom triggers based on stream activity. It includes 4 services to
customize your deployment.

## Sieve
A server hosting an HTTP API that keeps track of statistics on the state of
the data stream and forwards data to each subscriber over websockets.

## Pump
A client that feeds data streams to Sieve. Pump will consume log files and POST
the data from each line to the associated URL and stream name. It reads
stream mappings of log file paths to stream names from a YAML config file.

```yaml
url: http://localhost:9966
streams:
  /opt/stream/test1.log: test1
  /opt/stream/test2.log: test2
```

## Tap
A client that consumes data streams from Sieve. It connects to a sieve server
over a websocket and can print the data stream to stdout or log files, based
on a YAML config file.

```yaml
host: localhost
port: 9966
origin: http://localhost/
log_path: /go/src/github.com/vail130/orinoco/tap.log
```

## Litmus
A client that monitors data stream statistics through Sieve. It gets
statistics from a designated Sieve server every second, and evaluates each
trigger specified in a YAML config file.

```yaml
url: http://localhost:9966
triggers:
  trailing_average_per_minute_all_streams_at_zero:
    stream: "*"
    metric: trailing_average_per_minute
    condition: "==0"
    endpoint: http://example.com/trailing_average_per_minute_all_streams_at_zero
  trailing_average_per_hour_test_stream_under_100:
    stream: test_stream
    metric: trailing_average_per_hour
    condition: "<100"
    endpoint: http://example.com/trailing_average_per_hour_test_stream_under_100
```

# Testing
Install docker or boot2docker, then run tests from the make target:

`make test`
