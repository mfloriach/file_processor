count = 1
cpu = 2
iterations = 25x
timeout = 50m
memProfile = profiles/mem.out

benchConc:
	go test -race -bench=BenchmarkConcurrency -count ${count} -cpu ${cpu} -benchtime ${iterations} -timeout ${timeout} -benchmem -memprofile=${memProfile} -run=^#

benchErrs:
	go test -race -bench=BenchmarkErrors -count ${count} -cpu ${cpu} -benchtime ${iterations} -timeout ${timeout} -benchmem -memprofile=${memProfile} -run=^#

bench:
	go test -race -bench=BenchmarkMain -count ${count} -cpu ${cpu} -benchtime ${iterations} -timeout ${timeout} -benchmem -memprofile=${memProfile} -run=^# 

run-sequencial:
	go run main.go -mode=sequencial

run-concurrent:
	go run main.go -workers=1

run-parallel:
	go run main.go -workers=10

run:
	go run main.go 
