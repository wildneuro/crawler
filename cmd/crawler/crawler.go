package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"crawler/pkg/worker"
)

func main() {

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Need to specify start pages")
		os.Exit(1)
	}

	worker := worker.NewWorker(10)
	worker.PrintResult()
	worker.Run()

	for _, link := range args {
		err := worker.Add(link)
		if err != nil {
			log.Fatalf("error link: `%+v`", err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
