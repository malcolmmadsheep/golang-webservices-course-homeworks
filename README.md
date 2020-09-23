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
