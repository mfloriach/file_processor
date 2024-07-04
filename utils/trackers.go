package utils

import (
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/go-echarts/statsview"
)

func AddTracers() (func(), error) {
	go statsview.New().Start()

	traceFile, err := os.Create("profiles/trace.out")
	if err != nil {
		return func() {}, err
	}

	if err := trace.Start(traceFile); err != nil {
		return func() {}, err
	}

	f, err := os.Create("profiles/cpu.out")
	if err != nil {
		return func() {}, err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return func() {}, err
	}

	return func() {
		pprof.StopCPUProfile()
		traceFile.Close()
		trace.Stop()
		f.Close()
	}, nil
}
