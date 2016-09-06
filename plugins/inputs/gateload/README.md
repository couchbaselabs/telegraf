# Telegraf Plugin: Couchbase Gateload

## Configuration:

```
# Read per-node from Couchbase SyncGateway
[[inputs.gateload]]
  ## specify gateload server via a url:
  ##  [protocol://address:port]
  ##  e.g.
  ##    http://localhost:9876/
  ##
  ## If no server is specified, then localhost is used as the host.
  server = ["http://localhost:9876"]
```

## Measurements:

### gateload_expvar

Fields:
 - user_active              (unit: count, example: 345.0)
 - user_awake               (unit: count, example: 134.0)
 - total_doc_pushed         (unit: count, example: 34.0)
 - total_doc_failed_to_push (unit: count, example: 5.0)
 - total_doc_pulled         (unit: bytes, example: 7464.0)
 - total_doc_failed_to_pull (unit: bytes, example: 11.0)
 
 ### gateload_expvar.ops
 
 Fields:
  - p25  (unit: count, example: 5.0)
  - p50  (unit: count, example: 24.0)
  - p75  (unit: count, example: 34.0)
  - p90  (unit: count, example: 74.0)
  - p95  (unit: count, example: 123.0)
  - p99  (unit: count, example: 2345.0)