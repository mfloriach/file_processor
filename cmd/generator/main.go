package main

import (
	"bufio"
	"os"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

const totalKb = 50000

func main() {
	// f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// if err != nil {
	// 	panic(err)
	// }

	f, err := os.Create("file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	bar := progressbar.NewOptions(totalKb,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Writing in file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for getSizeKb(f) < totalKb {

		writer.WriteString("1,3,4,5")
		writer.WriteString("\n")

		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}

		bar.Add(int(fi.Size()))
	}
}

func getSizeKb(f *os.File) float64 {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return float64(fi.Size()) / (1 << 10)
}
