package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	// file, errOs := os.Open("MnistTrainingData/emnist-byclass-train.csv")
	file, errOs := os.Open("MnistTrainingData/emnist-digits-train.csv")
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
	yOneHot := convertToOneHot(allYs)

	fmt.Printf("allXsDim: (%v, %v)\n", len(allXs), len(allXs[0]))
	fmt.Printf("xTransposeDim: (%v, %v)\n", len(xTranspose), len(xTranspose[0]))
	fmt.Printf("allYsDim: (%v)\n", len(allYs))
	fmt.Printf("yOneHotDim: (%v, %v)\n", len(yOneHot), len(yOneHot[0]))

	W1, W2, b1, b2 := initialize(len(allXs[0]), len(allXs), 10, 10)

	fmt.Printf("W1Dim: (%v, %v)\n", len(W1), len(W1[0]))
	fmt.Printf("W2Dim: (%v, %v)\n", len(W2), len(W2[0]))
	fmt.Printf("b1Dim: (%v, %v)\n", len(b1), len(b1[0]))
	fmt.Printf("b2Dim: (%v, %v)\n", len(b2), len(b2[0]))

	numIterations := 1
	for itr := 0; itr < numIterations; itr++ {
		// _, A2 := feedForward(xTranspose, W1, W2, b1, b2)
		// calcErrors(yOneHot, A2)
	}

}

func initialize(n, m, l1, l2 int) ([][]float64, [][]float64, [][]float64, [][]float64) {
	W1, W2 := [][]float64{}, [][]float64{}
	b1, b2 := [][]float64{}, [][]float64{}

	for j := 0; j < l1; j++ {
		weightRow := []float64{}
		for i := 0; i < n; i++ {
			weightRow = append(weightRow, ((rand.Float64()*10)+1)/10) //0.0
		}
		W1 = append(W1, weightRow)
	}

	for i := 0; i < l1; i++ {
		biasRow := []float64{}
		for i := 0; i < m; i++ {
			biasRow = append(biasRow, ((rand.Float64()*10)+1)/10) //0.0
		}
		b1 = append(b1, biasRow)
	}

	for j := 0; j < l2; j++ {
		weightRow := []float64{}
		for i := 0; i < l1; i++ {
			weightRow = append(weightRow, ((rand.Float64()*10)+1)/10) //0.0
		}
		W2 = append(W2, weightRow)
	}

	for i := 0; i < l2; i++ {
		biasRow := []float64{}
		for i := 0; i < m; i++ {
			biasRow = append(biasRow, ((rand.Float64()*10)+1)/10) //0.0
		}
		b2 = append(b2, biasRow)
	}

	return W1, W2, b1, b2
}

func feedForward(X, W1, W2, b1, b2 [][]float64) (A1, A2 [][]float64) {
	Z1 := addMatrices(dotProduct(W1, X), b1)
	A1 = activateNeurons(Z1, "reLU")
	Z2 := addMatrices(dotProduct(W2, A1), b2)
	A2 = activateNeurons(Z2, "softmax")
	return
}

// func calcErrors(Y, A2 [][]float64) [][]float64 {
// 	cost := subtactMatrices(Y, A2)

// 	return cost
// }

func activateNeurons(Z [][]float64, activationName string) [][]float64 {
	mZ, nZ := len(Z), len(Z[0])
	activatedNeurons := [][]float64{}
	switch activationName {
	case "reLU":
		for i := 0; i < mZ; i++ {
			singleNeuron := []float64{}
			for j := 0; j < nZ; j++ {
				singleNeuron = append(singleNeuron, math.Max(Z[i][j], 0))
			}
			activatedNeurons = append(activatedNeurons, singleNeuron)
		}

	case "softmax":
		denominators := []float64{}
		for j := 0; j < nZ; j++ {
			sum := 0.0
			for i := 0; i < mZ; i++ {
				sum += math.Exp(Z[i][j])
			}
			denominators = append(denominators, sum)
		}

		for i := 0; i < mZ; i++ {
			singleNeuron := []float64{}
			for j := 0; j < nZ; j++ {
				singleNeuron = append(singleNeuron, math.Exp(Z[i][j])/denominators[j])
			}
			activatedNeurons = append(activatedNeurons, singleNeuron)
		}

	default:
		errActivation := errors.New("Illegal Activation Function")
		fmt.Printf("errActivation.Error(): %v\n | activationName %v", errActivation.Error(), activationName)
		panic(errActivation)
	}

	return activatedNeurons
}

func dotProduct(A, B [][]float64) [][]float64 {
	mA, nA, mB, nB := len(A), len(A[0]), len(B), len(B[0])

	if nA != mB {
		errDim := errors.New("Illegal Operation: Dot Product not possible due to dimentions")
		fmt.Printf("errDim.Error(): %v\n | mA: %v nA: %v mB: %v nB: %v", errDim.Error(), mA, nA, mB, nB)
		panic(errDim)
	}

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nB; j++ {
			sum := 0.0
			for k := 0; k < nA; k++ {
				sum += A[i][k] * B[k][j]
			}
			resultRow = append(resultRow, sum)
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func addMatrices(A, B [][]float64) [][]float64 {
	mA, nA, mB, nB := len(A), len(A[0]), len(B), len(B[0])

	if mA != mB || nA != nB {
		errDim := errors.New("Illegal Operation: Matrix Addition not possible due to dimentions")
		fmt.Printf("errDim.Error(): %v\n | mA: %v nA: %v mB: %v nB: %v", errDim.Error(), mA, nA, mB, nB)
		panic(errDim)
	}

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nB; j++ {
			resultRow = append(resultRow, A[i][j]+B[i][j])
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func subtactMatrices(A, B [][]float64) [][]float64 {
	mA, nA, mB, nB := len(A), len(A[0]), len(B), len(B[0])

	if mA != mB || nA != nB {
		errDim := errors.New("Illegal Operation: Matrix Subtraction not possible due to dimentions")
		fmt.Printf("errDim.Error(): %v\n | mA: %v nA: %v mB: %v nB: %v", errDim.Error(), mA, nA, mB, nB)
		panic(errDim)
	}

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nB; j++ {
			resultRow = append(resultRow, A[i][j]-B[i][j])
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func extractToIntSlices(allRecords [][]string) ([]float64, [][]float64) {
	allYs, allXs := []float64{}, [][]float64{}
	for _, row := range allRecords {
		y, errAtoiY := strconv.Atoi(row[0])
		if errAtoiY != nil {
			fmt.Printf("errAToIY.Error(): %v\n", errAtoiY.Error())
			panic(errAtoiY)
		}
		allYs = append(allYs, float64(y))

		remainingRow := row[1:]
		singleRow := []float64{}
		for _, entry := range remainingRow {
			singleEntry, errAtoiSingleEntry := strconv.Atoi(entry)
			if errAtoiSingleEntry != nil {
				fmt.Printf("errAtoiSingleEntry.Error(): %v\n", errAtoiSingleEntry.Error())
				panic(errAtoiSingleEntry)
			}
			singleRow = append(singleRow, float64(singleEntry))
		}
		allXs = append(allXs, singleRow)
	}

	return allYs, allXs
}

func transposeMatrice(matrix [][]float64) [][]float64 {
	returnMatrix := [][]float64{}
	ogM, ogN := len(matrix), len(matrix[0])

	for i := 0; i < ogN; i++ {
		newRow := []float64{}
		for j := 0; j < ogM; j++ {
			newRow = append(newRow, matrix[j][i])
		}
		returnMatrix = append(returnMatrix, newRow)
	}
	return returnMatrix
}

func convertToOneHot(list []float64) [][]float64 {
	returnMatrix := [][]float64{}
	maxNum, numEntries := findMax(list), len(list)

	for j := 0; j < int(maxNum)+1; j++ {
		oneHotRow := []float64{}
		for i := 0; i < numEntries; i++ {
			if list[i] == float64(j) {
				oneHotRow = append(oneHotRow, 1)
			} else {
				oneHotRow = append(oneHotRow, 0)
			}
		}
		returnMatrix = append(returnMatrix, oneHotRow)
	}

	return returnMatrix
}

func findMax(list []float64) float64 {
	maxSoFar := float64(math.MinInt)
	for _, num := range list {
		if num > maxSoFar {
			maxSoFar = num
		}
	}
	return maxSoFar
}
