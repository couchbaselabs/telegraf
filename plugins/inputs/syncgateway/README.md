# Telegraf Plugin: Couchbase Sync Gateway

## Configuration:

```
# Read per-node from Couchbase SyncGateway
[[inputs.syncgateway]]
  ## specify sync gateway server via a url:
  ##  [protocol://address:port]
  ##  e.g.
  ##    http://localhost:4985/
  ##
  ## If no server is specified, then localhost is used as the host.
  server = ["http://localhost:4985"]
```

## Measurements:

### syncgateway_expvar

Fields:
 - changesFeeds_total   (unit: count, example: 2345.0)
 - changesFeeds_active  (unit: count, example: 345.0)
 - requests_total       (unit: count, example: 234.0)
 - requests_active      (unit: count, example: 34.0)
 - revisionCache_hits   (unit: count, example: 2345.0)
 - revisionCache_misses (unit: count, example: 123.0)
 - memstatsAlloc        (unit: bytes, example: 202156957464.0)
 - memstatsSys          (unit: bytes, example: 212179309111.0)