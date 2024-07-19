# Batch processor
[![test](https://github.com/mfloriach/file_processor/actions/workflows/test.yml/badge.svg)](https://github.com/mfloriach/file_processor/actions/workflows/test.yml)
![go](https://img.shields.io/badge/go-1.22-blue)

Batch processor example to test concurrency patterns in golang.

Modules
- Metadata: get year, name, start, authors,....
- Compress: compress and move to S3 bucket
- Store: save into Mongo DB

## Visual tools

| Tool |  Description  |
|:-----|:--------:|
| [mongo](http://localhost:8082/)                   |Database management |
| [minio](http://127.0.0.1:9001)                    |S3 alternative manager  |
| [stats](http://localhost:18066/debug/statsview)   |Golang related statuses|
| [cadvisor](http://localhost:8080/docker)          |Hardware related statuses|

## How to use

### Start
```bash
$ docker compose up -d # start services
$ make run             # run in parallel mode
$ make run-sequencial  # run sequential
$ make run-concurrent  # run concurrent
$ make run-parallel    # run parallel
```

### Benchmark
```bash
$ make bench 
$ make benchErrs 
$ make benchConc 
```

## Architecture

### Sequencial
```mermaid
flowchart LR
    FA[read file] --> M[Metada]
    M[Metada] --> |job1|C[Compress]
    C[Compress] --> S[Store]
```

### Concurrent
```mermaid
flowchart LR
    FA[read file] -->|job3| Metadata
    subgraph Metadata
    direction RL
    M1[Worker]
    end
    Metadata --> |job2|Compress
    subgraph Compress
    direction RL
    C[Worker]
    end
    Compress --> |job1|Store
    subgraph Store
    direction RL
    S1[Worker]
    end
```

### Parallel
```mermaid
flowchart LR
    FA[read file] --> Metadata
    subgraph Metadata
    direction RL
    M[Worker job 6]
    M1[Worker job 5]
    end
    Metadata --> Compress
    subgraph Compress
    direction RL
    C[Worker job 3]
    C2[Worker job 4]
    end
    Compress --> Store
    subgraph Store
    direction RL
    S[Worker job 1]
    S1[Worker job 2]
    end
```

