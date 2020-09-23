package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const MultiHashThreadsCount = 6

func toString(x interface{}) string {
	return fmt.Sprintf("%v", x)
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := make(chan interface{})

	for _, currentJob := range jobs {
		wg.Add(1)
		go func(in, out chan interface{}, myJob job) {
			defer wg.Done()
			myJob(in, out)
			close(out)
		}(in, out, currentJob)

		in = out
		out = make(chan interface{})
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for rawData := range in {
		wg.Add(1)
		go func(wg *sync.WaitGroup, mu *sync.Mutex, out chan interface{}, rawData interface{}) {
			defer wg.Done()
			data := toString(rawData)
			firstPart := make(chan string)
			secondPart := make(chan string)

			go func(out chan string, data string) {
				out <- DataSignerCrc32(data)
			}(firstPart, data)

			go func(out chan string, data string) {
				mu.Lock()
				md5 := DataSignerMd5(data)
				mu.Unlock()

				res := DataSignerCrc32(md5)

				out <- res
			}(secondPart, data)

			out <- (<-firstPart + "~" + <-secondPart)
		}(wg, mu, out, rawData)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for rawData := range in {
		wg.Add(1)
		go func(wg *sync.WaitGroup, rawData interface{}) {
			defer wg.Done()
			innerWg := &sync.WaitGroup{}
			data := toString(rawData)
			resultCh := make(chan struct{})
			totalRes := make([]string, MultiHashThreadsCount)

			for i := 0; i < MultiHashThreadsCount; i++ {
				innerWg.Add(1)
				go func(out chan struct{}, th int, data string) {
					defer innerWg.Done()

					totalRes[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(resultCh, i, data)
			}

			innerWg.Wait()

			out <- strings.Join(totalRes, "")
		}(wg, rawData)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0, MultiHashThreadsCount)

	for hash := range in {
		hashes = append(hashes, toString(hash))
	}

	sort.Strings(hashes)

	out <- strings.Join(hashes, "_")
}
