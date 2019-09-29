package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	size = 3 * 1 << 10
)

var (
	files []string
	sum   int64
)

func main() {
	start := time.Now()
	run()
	fmt.Println(time.Now().Sub(start))
}

func run() {
	flattenFiles()

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {

		go func(file string) {
			sumFile(file)
			wg.Done()

		}(file)

	}

	// wait for the files to be completely read
	wg.Wait()
	fmt.Println(sum)

}

// Creates an array from the files in the file dir
func flattenFiles() {
	for i := 0; i < 100; i++ {
		for j := 0; j < 10; j++ {
			files = append(files, fmt.Sprintf("./files/%06d-%06d/%06d.csv", (i*10)+1, (i*10)+10, (i*10)+j+1))
		}
	}
}

// reads a file, calls readNum and adds the sum to the global sum
func sumFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		panic((err))
	}

	n := readNum(f)
	f.Close()
	atomic.AddInt64(&sum, n)
}

// reads an io.Reader and retruns the sum of numbers in that reader
func readNum(f io.Reader) int64 {
	var b [size]byte
	var lastVal int
	var sum int64

	for {
		x, err := f.Read(b[:])

		if err == io.EOF {
			sum += int64(lastVal)
			break
		}

		if err != nil {
			panic(err)
		}

		var n int
		c := len(b[:x])
		for ii := 0; ii < c; ii++ {
			if b[ii] < '0' {
				i := concatenate(lastVal, b[n:ii])
				sum += int64(i)
				n = ii + 1
				lastVal = 0
			}
		}
		// save the last value after a `,`. not every read will end with a comma
		lastVal = (concatenate(lastVal, b[n:c]))
	}

	return sum
}

// converts a []byte to an int. []byte{'1','2'} => int(12)
func concatenate(prefix int, b []byte) int {
	n := prefix
	for _, c := range b {
		c -= '0'
		n = n*10 + int(c)
	}
	return n
}
