# Home works for Golang course on Coursera

## Overview

This repo contains homeworks for course [Разработка веб-сервисов на Go - основы языка](https://www.coursera.org/learn/golang-webservices-1) on Coursera.

Each week includes programming assignment to consolidate knowledge achieved throughout the course.

All homeworks have minimal amount of tests to make sure program works correctly.

## Homeworks

### Week 1. `tree` command

**Problem**: implement simplified analog of Unix [`tree` command](https://linux.die.net/man/1/tree) that prints out structure of directory. Use `-f` flag to include files in output.

**Solution**: I implemented recursive and iterative approaches to solve this problem. Also output is colorized.

### Week 2. Crypto hash function

Week 2 is about using asynchronous functionality in Golang. Asynchrony in Golang is based on goroutines - lightweight threads that allow you to process operation in concurrent mode. Data between goroutines can be synchronized by channels. Channel is a data structure that works like a pipe and give you an ability to write data in one goroutine and read this data in other.

**Problem**: this problem consists of 2 parts that should be solved. 1) need to implement `ExecutePipeline` - is a function that accepts slice of jobs (`type job func(in, out chan interface{})`) and builds pipeline from them; 2) make calculations to work in parallel, because hash for single value calculates ~8s and we have some time limitations.

**Limitations**: all process shouldn't take more than 3s. `DataSignerCrc32` calculates 1s (yeah, sleep inside is intentional). `DataSignerMd5` can be called only once in the same time and takes 10ms, if it's being called in parallel - there will be overheat for 1s.

**Solution**: solution is put inside `hw2_signer/signer.go`, other files in this folder was provided by course, in these files implemented hash functions and tests. The main idea: all jobs communicate to each other by channels - we have one channel for sending data and one for receiving. Sender for job is receiver for the next one. I created wrapper for job that is controlled by wait group and closes sender when job is done.

### Week 3.

Main topics of 3rd week were dynamic data processing (handing JSON with `interface{}` type and reflection) and profiling of program using Golang tool `pprof` according to results of benchmark tests.

**Problem**: read file (`data/users.txt`) which contains stringified JSON objects, and find only users that use Android and MSIE browsers. Implement `FastSearch` method that will work faster than `SlowSearch` and show performance close to course solution.

**Solution**: base functionality is implemented in `SlowSearch` method. Before profiling I supposed what blocks of code should be performed. Main pieces of code that impacted performance and created memory overhead were:

- using `map[string]interface{}` to user entity. **Possible solution**: use determined structure to store user
- using regular expression in places where they can be avoided. **Possible solution**: replace regexp matching with `strings.Contains`
- reading all file content, creating users slice from it and storing it in memory. **Possible solution**: read file line by line and use single variable to store processed user structure (it reduces number of allocations and memory consumption)

In `FastSearch` I tried to minimize average time, consumed memory and number of allocations per operation. I created profiling files using `go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out` and explored them using `go tool pprof`. Using these tools I found out the most CPU and memory consuming places which were:

- [**CPU**] using `json.Unmarshal` to unpack `[]byte` into `map[string]interface{}`
- [**CPU + MEM**] regexp matching
- [**MEM**] reading full file and splitting it in lines
- [**MEM**] creating users slice before actual search and keeping it in memory
- [**MEM**] memory allocations for strings to generate result output

So I did some optimizations to improve all these benchmarks:

- read file line by line, using `bufio.Scanner`
- replaced browsers checking from regexp matching to `strings.Contains` method, because searchable string doesn't include any wildcards or complex matching conditions, only string literals
- provided `User` structure to reduce memory consumption
- created single `User` variable, and use it to store unpacked JSON object in it
- stopped allocating new strings and use `bytes.Buffer` to store result
- changed replacing of `@` from RegExp to direct writing to result buffer

To unmarshal JSON into `User` structure [easyjson package](https://github.com/mailru/easyjson) was used. This package was suggested in course. It's main idea is to generate marshal and unmarshal methods for structure from it's definition. Using of this package could be avoided, but this option also included JSON parsing, so I decided not to dive so deep.

After all my improvements I received next results:

| Test name                             | # operations | Operation Time | Operation Memory | # Allocations    |
| ------------------------------------- | ------------ | -------------- | ---------------- | ---------------- |
| BenchmarkSlow                         | 36           | 31923467 ns/op | 18661662 B/op    | 195778 allocs/op |
| BenchmarkFast (My)                    | 625          | 1842831 ns/op  | 542210 B/op      | 10160 allocs/op  |
| BenchmarkSolution-8 (Course Solution) | 500          | 2782432 ns/op  | 559910 B/op      | 10422 allocs/op  |

Last row is course solution benchmarks. Main goal is to have one of benchmarks lower than in solution (fast < solution) and one benchmark should be lower than solution \* 1.2 (fast < solution \* 1.2).
