package services

import (
	"bufio"
	"os"
	"processor/types"
	"strings"
	"time"
)

type iterator struct {
	scanner     *bufio.Scanner
	file        *os.File
	paths       []string
	currentPath int
	position    int64
}

func NewIngestorFileReader(paths ...string) (types.Iterator, error) {
	file, err := os.OpenFile(paths[0], os.O_RDONLY, os.ModePerm)
	if err != nil {
		return &iterator{}, err
	}

	return &iterator{
			scanner:     bufio.NewScanner(file),
			file:        file,
			paths:       paths,
			currentPath: 0,
			position:    0,
		},
		nil
}

func (i *iterator) HasNext() bool {
	if !i.scanner.Scan() {
		i.file.Close()

		if i.currentPath != len(i.paths)-1 {
			i.currentPath += 1
			i.position = 0
			file, _ := os.OpenFile(i.paths[i.currentPath], os.O_RDONLY, os.ModePerm)
			i.file = file
			i.scanner = bufio.NewScanner(file)
			return true
		}

		return false
	}

	return true
}

func (i *iterator) Next() types.Movie {
	i.position += 1
	return i.ingestorFileReaderSerializer()
}

func (i *iterator) ingestorFileReaderSerializer() types.Movie {
	str := i.scanner.Text()

	fields := strings.Split(str, ",")

	return types.Movie{
		Title:       fields[0],
		Year:        fields[1],
		ReleaseDate: fields[2],
		Metadata: types.Metadata{
			StartProcessing: time.Now(),
			FileName:        i.paths[i.currentPath],
			Position:        i.position,
		},
	}
}
