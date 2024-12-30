# Load Balancer Prototype in Go

### Start the LB

Health Check interval is optional, default is 10s.

```
go run . -healthCheckSec 2
```

### Start mock BE servers

Default port is 3200

```
cd be

go run . -port 3200
go run . -port 3201
go run . -port 3202
```

### Test with one request

Curl the LB

```
curl http://localhost:3222
```

### Test with multiple requests in parralel

```
curl --parallel --parallel-immediate --parallel-max 3 --config reqs.txt
```
