#! /usr/bin/env bash

mkdir bench

go test -bench . -benchmem -cpuprofile cpu.out -memprofile mem.out -outputdir ./bench