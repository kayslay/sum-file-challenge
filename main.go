package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
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
	// sem limit the number of goroutines running. this prevent a "To Many Files open" error
	sem := make(chan struct{}, 500)

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		sem <- struct{}{}

		go func(file string) {
			sumFile(file)
			<-sem
			wg.Done()

		}(file)

	}

	// wait for the files to be completely read
	wg.Wait()
	fmt.Println(sum)

}

// 49952200434

func flattenFiles() {
	for i := 0; i < 100; i++ {
		for j := 0; j < 10; j++ {
			fileName := fmt.Sprintf("./files/%06d-%06d/%06d.csv", (i*10)+1, (i*10)+10, (i*10)+j+1)
			files = append(files, fileName)
		}
	}
}

func sumFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		panic((err))
	}
	n := readNum(f)
	atomic.AddInt64(&sum, n)
	f.Close()
}

func readNum(f io.Reader) int64 {
	var b [256]byte
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

		lastVal = (concatenate(lastVal, b[n:c]))
		// fmt.Println(off, n)
	}

	return sum
}

func concatenate(last int, b []byte) int {
	n := last
	for _, c := range b {
		c -= '0'
		n = n*10 + int(c)
	}
	return n
}
