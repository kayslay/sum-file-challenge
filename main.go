package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	files []string
)

func main() {
	flattenFiles()
	// sem limit the number of goroutines running. this prevent a "To Many Files open" error
	sem := make(chan struct{}, 50)

	var wg sync.WaitGroup
	wg.Add(len(files))

	var sum int64

	for _, file := range files {
		sem <- struct{}{}

		go func(file string) {
			f, err := os.Open(file)
			if err != nil {
				panic((err))
			}

			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic((err))
			}
			f.Close()

			var numSum int64
			lines := strings.Split(string(b), "\n")

			for _, line := range lines {

				nums := strings.Split(line, ",")

				for _, n := range nums {
					i, err := strconv.Atoi(n)
					if err != nil {
						panic((err))
					}
					numSum += int64(i)
				}
			}

			atomic.AddInt64(&sum, numSum)
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
