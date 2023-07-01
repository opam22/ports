# Ports

There are 2 microservices on this repository.
1. Importer
2. Ports

### Importer
Importer service will read the json file (ports.json) and will decode it one by one and then sending it to Ports service through Ports client

### Ports
Ports service will process port data that received from Importer service and will store it in database (in-memory)

## How to run
Create your own config, copy env/config.sample to env/config

Run without docker:
Start the server with `make run-ports`
Then run the importer `make run-importer`

Run with docker:
`make docker-run`

Build docker image:
`make docker-build`

# Test
To run the tests, simple run `make test`

# Generate new protobuf
To generate new protobuf, run `make generate`

# Expected logs
If all is running well, we should be able to see log like this
```
ports-ports-1     | {
ports-ports-1     |   "level": "info",
ports-ports-1     |   "msg": "successfully storing port: \u0026{PortID:ZWUTA Name:Mutare City:Mutare Country:Zimbabwe Alias:[] Regions:[] Coordinates:[32.650351 -18.9757714] Province:Manicaland Timezone:Africa/Harare Unlocs:[ZWUTA] Code:}",
ports-ports-1     |   "time": "2023-07-01T08:20:08Z"
ports-ports-1     | }
```

