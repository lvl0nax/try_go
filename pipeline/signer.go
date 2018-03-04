package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func myMd5(data string, out chan string, quotaCh chan struct{}) {
	quotaCh <- struct{}{}

	out <- DataSignerMd5(data)
	close(out)

	<-quotaCh
}

func myCrc32(data string, out chan string) {
	out <- DataSignerCrc32(data)
	close(out)
}

func calculateSingleHash(data string, out chan interface{}, quotaCh chan struct{}, myWg *sync.WaitGroup) {
	defer myWg.Done()

	md5Ch := make(chan string, 1)
	crc32Ch1 := make(chan string, 1)
	crc32Ch2 := make(chan string, 1)

	go myCrc32(data, crc32Ch1)
	go myMd5(data, md5Ch, quotaCh)
	go myCrc32(<-md5Ch, crc32Ch2)

	out <- (<-crc32Ch1 + "~" + <-crc32Ch2)
}

func calculateMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var chans = []chan string{
		make(chan string, 1),
		make(chan string, 1),
		make(chan string, 1),
		make(chan string, 1),
		make(chan string, 1),
		make(chan string, 1),
	}

	for i := 0; i <= 5; i++ {
		go myCrc32(strconv.Itoa(i)+data, chans[i])
	}

	out <- (<-chans[0] + <-chans[1] + <-chans[2] + <-chans[3] + <-chans[4] + <-chans[5])
}

func SingleHash(in, out chan interface{}) {
	quotaCh := make(chan struct{}, 1)
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go calculateSingleHash(fmt.Sprintf("%v", i), out, quotaCh, wg)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go calculateMultiHash(fmt.Sprintf("%v", i), out, wg)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var arr []string

	for i := range in {
		arr = append(arr, fmt.Sprintf("%v", i))
	}

	sort.Strings(arr)
	str := strings.Join(arr, "_")

	out <- str
}

func startWorker(worker job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer close(out)
	defer wg.Done()

	worker(in, out)
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)

	for _, worker := range jobs {
		wg.Add(1)

		go startWorker(worker, in, out, wg)

		in = out
		out = make(chan interface{}, 100)
	}

	wg.Wait()
}
