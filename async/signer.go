package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const TH = 6

func ExecutePipeline(jobs ...job)  {
	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 1)
	for _, hashJob := range jobs {
		wg.Add(1)
		out := make(chan interface{}, 1)
		go func(hJob job, in chan interface{}, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			hJob(in, out)
		}(hashJob, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for data := range in {
		dat, err := data.(int)
		if err == false{
			panic(err)
		}
		wg.Add(1)
		data := strconv.Itoa(dat)
		go Single(data, wg, mu, out)
	}
	wg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for dat := range in {
		data, err := dat.(string)
		if err == false {
			panic(err)
		}
		wg.Add(1)
		go Multi(data, wg, out)
		}
	wg.Wait()
}

func CombineResults(in chan interface{}, out chan interface{})  {
	var mas []string
	for dat := range in {
		data, err := dat.(string)
		if err == false{
			panic(err)
		}
		mas = append(mas, data)
	}
	sort.Strings(mas)
	out <- strings.Join(mas, "_")
}

func Single (data string, wg *sync.WaitGroup, mu *sync.Mutex, out chan interface{}) {
	defer wg.Done()
	var res, resMd5 string
	wgS := &sync.WaitGroup{}
	wgS.Add(1)
	go func() {
		defer wgS.Done()
		mu.Lock()
		datMd5 := DataSignerMd5(data)
		mu.Unlock()
		resMd5 = DataSignerCrc32(datMd5)
	}()
	res = DataSignerCrc32(data)
	wgS.Wait()
	out <- res + "~" + resMd5
}

func Multi (data string, wg *sync.WaitGroup,  out chan interface{}) {
	defer wg.Done()
	var mas = make([]string, TH, TH)
	wgMu := &sync.WaitGroup{}
	for i := 0; i < TH; i++ {
		wgMu.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			mas[i] = DataSignerCrc32(strconv.Itoa(i) + data)
		}(i, wgMu)
	}
	wgMu.Wait()
	out <- strings.Join(mas, "")
}