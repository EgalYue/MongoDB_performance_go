package main

import (
	"log"
	"os"
	"strconv"
)

func main(){
	if len(os.Args[:]) < 3 {
		log.Fatal("please enter count and average!")
	}

	count, err1 := strconv.Atoi( os.Args[1]) // number of operations e.g. 10 times insert...
	average, err2 := strconv.Atoi( os.Args[2]) // Loop index, used to compute the average time
	if err1 != nil || err2 != nil{
		log.Fatal("count or average has error!")
	}
	Compute_performance(count,average)

}