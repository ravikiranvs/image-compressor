package main

import (
	"FileExplorer"
	"Imaging"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	size := 0
	if len(os.Args) > 1 {
		strSize := os.Args[1]
		intSize, err := strconv.Atoi(strSize)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		size = intSize
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	currentDir = currentDir + "\\"
	// currentDir := "C:\\PICS\\"
	err = os.Mkdir(currentDir+"Compressed", os.ModeDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	destinationDir := currentDir + "Compressed\\"

	files, err := FileExplorer.GetJPEGFiles(currentDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	filesCount := len(files)
	batchSize := 4
	//batchStart := time.Now()
	//totalTime := 0.0
	for index, file := range files {
		wg.Add(1)
		fileName := file.Name()
		go Imaging.ResizeImage(currentDir+fileName, destinationDir+fileName, size, &wg)
		oneIndex := index + 1
		if oneIndex%batchSize == 0 || oneIndex == filesCount {
			wg.Wait()
			//batchDuration := time.Since(batchStart).Seconds()
			//batchDurationStr := strconv.FormatFloat(time.Since(batchStart).Seconds(), 'f', 6, 64)
			//totalTime = totalTime + batchDuration
			//fmt.Println("Batch Duration: " + batchDurationStr)
			//batchStart = time.Now()
		}
	}

	//totalDurationStr := strconv.FormatFloat(totalTime, 'f', 6, 64)
	//fmt.Println("Total time: " + totalDurationStr)

	os.Exit(0)
}
