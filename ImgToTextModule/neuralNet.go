package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, errOs := os.Open("MnistTrainingData/emnist-byclass-train.csv")
	if errOs != nil {
		fmt.Printf("errOs.Error(): %v\n", errOs.Error())
		panic(errOs)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	allRecords, errRead := reader.ReadAll()
	if errRead != nil {
		fmt.Printf("errRead.Error(): %v\n", errRead.Error())
		panic(errRead)
	}

	fmt.Printf("allRecords[0]: %v\n", allRecords[0])

	totalTrainingM, totalTrainingN := len(allRecords), len(allRecords[0])-1

	fmt.Printf("totalTrainingM: %v\n", totalTrainingM)
	fmt.Printf("totalTrainingN: %v\n", totalTrainingN)

	allYs, allXs := extractToIntSlices(allRecords)

	fmt.Printf("len(allXs): %v\n", len(allXs))
	fmt.Printf("len(allXs[0]): %v\n", len(allXs[0]))

	xTranspose := transposeMatrice(allXs)

	fmt.Printf("allXsDim: (%v, %v)\n", len(allXs), len(allXs[0]))
	fmt.Printf("xTransposeDim: (%v, %v)\n", len(xTranspose), len(xTranspose[0]))
	fmt.Printf("allYsDim: (%v)\n", len(allYs))
}

func extractToIntSlices(allRecords [][]string) ([]int, [][]int) {
	allYs, allXs := []int{}, [][]int{}
	for _, row := range allRecords {
		y, errAtoiY := strconv.Atoi(row[0])
		if errAtoiY != nil {
			fmt.Printf("errAToIY.Error(): %v\n", errAtoiY.Error())
			panic(errAtoiY)
		}
		allYs = append(allYs, y)

		remainingRow := row[1:]
		singleRow := []int{}
		for _, entry := range remainingRow {
			singleEntry, errAtoiSingleEntry := strconv.Atoi(entry)
			if errAtoiSingleEntry != nil {
				fmt.Printf("errAtoiSingleEntry.Error(): %v\n", errAtoiSingleEntry.Error())
				panic(errAtoiSingleEntry)
			}
			singleRow = append(singleRow, singleEntry)
		}
		allXs = append(allXs, singleRow)
	}

	fmt.Printf("f: len(allXs): %v\n", len(allXs))
	fmt.Printf("f: len(allXs[0]): %v\n", len(allXs[0]))

	return allYs, allXs
}

func transposeMatrice(matrix [][]int) [][]int {
	returnMatrix := [][]int{}
	ogM, ogN := len(matrix), len(matrix[0])

	for i := 0; i < ogN; i++ {
		newRow := []int{}
		for j := 0; j < ogM; j++ {
			newRow = append(newRow, matrix[j][i])
		}
		returnMatrix = append(returnMatrix, newRow)
	}
	return returnMatrix
}
