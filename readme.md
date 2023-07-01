# Ports

There are 2 microservices on this repository.
1. Importer
2. Ports

### Importer
Importer service will read the json file (ports.json) and will decode it one by one and then sending it to Ports service through Ports client

### Ports
Ports service will process port data that received from Importer service and will store it in database (in-memory)


