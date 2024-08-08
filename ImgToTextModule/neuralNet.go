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

	numIterations, alpha := 1, 0.3
	for itr := 0; itr < numIterations; itr++ {
		A1, A2, Z1, _ := feedForward(xTranspose, W1, W2, b1, b2)
		dW2, dB2, dW1, dB1 := calcErrors(yOneHot, A2, A1, Z1, xTranspose)
		W1, W2, b1, b2 = updateWeightsAndBias(W1, W2, b1, b2, dW2, dB2, dW1, dB1, alpha)
	}

}

func updateWeightsAndBias(W1, W2, b1, b2, dW2, dB2, dW1, dB1 [][]float64, alpha float64) ([][]float64, [][]float64, [][]float64, [][]float64) {
	newW1 := subtactMatrices(W1, multipyMatrixBy(dW1, alpha))
	newW2 := subtactMatrices(W2, multipyMatrixBy(dW2, alpha))
	newb1 := subtactMatrices(b1, multipyMatrixBy(dB1, alpha))
	newb2 := subtactMatrices(b2, multipyMatrixBy(dB2, alpha))

	return newW1, newW2, newb1, newb2
}

func multipyMatrixBy(A [][]float64, alpha float64) [][]float64 {
	mA, nA := len(A), len(A[0])

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nA; j++ {
			resultRow = append(resultRow, A[i][j]*alpha)
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
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

func feedForward(X, W1, W2, b1, b2 [][]float64) (A1, A2, Z1, Z2 [][]float64) {
	Z1 = addMatrices(crossProduct(W1, X), b1)
	A1 = activateNeurons(Z1, "reLU")
	Z2 = addMatrices(crossProduct(W2, A1), b2)
	A2 = activateNeurons(Z2, "softmax")

	fmt.Printf("feedForward - Z1Dim: (%v, %v)\n", len(Z1), len(Z1[0]))
	fmt.Printf("feedForward - A1Dim: (%v, %v)\n", len(A1), len(A1[0]))
	fmt.Printf("feedForward - Z2Dim: (%v, %v)\n", len(Z2), len(Z2[0]))
	fmt.Printf("feedForward - A2Dim: (%v, %v)\n", len(A2), len(A2[0]))

	return
}

func calcErrors(Y, A2, A1, Z1, X [][]float64) ([][]float64, [][]float64, [][]float64, [][]float64) {
	m := len(Y[0])

	dZ2 := subtactMatrices(A2, Y)

	fmt.Printf("calcErrors - dZ2Dim: (%v, %v)\n", len(dZ2), len(dZ2[0]))

	dW2 := crossProduct(dZ2, transposeMatrice(A1))
	dW2 = divideMatrixBy(dW2, float64(m))

	fmt.Printf("calcErrors - dW2Dim: (%v, %v)\n", len(dW2), len(dW2[0]))

	dB2 := addUpRows(dZ2)
	dB2 = divideMatrixBy(dB2, float64(m))
	dB2 = expandMatrix(dB2, m)

	fmt.Printf("calcErrors - dB2Dim: (%v, %v)\n", len(dB2), len(dB2[0]))

	gZ1 := undoActivationReLU(Z1)

	fmt.Printf("calcErrors - gZ1Dim: (%v, %v)\n", len(gZ1), len(gZ1[0]))

	dZ1 := dotProduct(crossProduct(transposeMatrice(dW2), dZ2), gZ1)

	fmt.Printf("calcErrors - dZ1Dim: (%v, %v)\n", len(dZ1), len(dZ1[0]))

	dW1 := crossProduct(dZ1, transposeMatrice(X))
	dW1 = divideMatrixBy(dW1, float64(m))

	fmt.Printf("calcErrors - dW1Dim: (%v, %v)\n", len(dW1), len(dW1[0]))

	dB1 := addUpRows(dZ1)
	dB1 = divideMatrixBy(dB1, float64(m))
	dB1 = expandMatrix(dB1, m)

	fmt.Printf("calcErrors - dB1Dim: (%v, %v)\n", len(dB1), len(dB1[0]))

	return dW2, dB2, dW1, dB1
}

func expandMatrix(A [][]float64, numCols int) [][]float64 {
	m := len(A)
	returnMatrix := [][]float64{}

	for i := 0; i < m; i++ {
		value := A[i][0]
		newRow := []float64{}
		for j := 0; j < numCols; j++ {
			newRow = append(newRow, value)
		}
		returnMatrix = append(returnMatrix, newRow)
	}

	return returnMatrix
}

func dotProduct(A, B [][]float64) [][]float64 {
	mA, nA, mB, nB := len(A), len(A[0]), len(B), len(B[0])

	if mA != mB || nA != nB {
		errDim := errors.New("Illegal Operation: Matrix Dot Product not possible due to dimentions")
		fmt.Printf("errDim.Error(): %v | mA: %v nA: %v mB: %v nB: %v\n", errDim.Error(), mA, nA, mB, nB)
		panic(errDim)
	}

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nB; j++ {
			resultRow = append(resultRow, A[i][j]*B[i][j])
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func undoActivationReLU(A [][]float64) [][]float64 {
	mA, nA := len(A), len(A[0])

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nA; j++ {
			derivative := 0.0
			if A[i][j] > 0 {
				derivative = 1.0
			}
			resultRow = append(resultRow, derivative)
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func addUpRows(A [][]float64) [][]float64 {
	mA, nA := len(A), len(A[0])

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		sum := 0.0
		for j := 0; j < nA; j++ {
			sum += A[i][j]
		}
		resultMatrix = append(resultMatrix, []float64{sum})
	}

	return resultMatrix
}

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
		fmt.Printf("errActivation.Error(): %v | activationName %v\n", errActivation.Error(), activationName)
		panic(errActivation)
	}

	return activatedNeurons
}

func divideMatrixBy(A [][]float64, denominator float64) [][]float64 {
	mA, nA := len(A), len(A[0])

	resultMatrix := [][]float64{}
	for i := 0; i < mA; i++ {
		resultRow := []float64{}
		for j := 0; j < nA; j++ {
			resultRow = append(resultRow, A[i][j]/denominator)
		}
		resultMatrix = append(resultMatrix, resultRow)
	}

	return resultMatrix
}

func crossProduct(A, B [][]float64) [][]float64 {
	mA, nA, mB, nB := len(A), len(A[0]), len(B), len(B[0])

	if nA != mB {
		errDim := errors.New("Illegal Operation: Dot Product not possible due to dimentions")
		fmt.Printf("errDim.Error(): %v | mA: %v nA: %v mB: %v nB: %v\n", errDim.Error(), mA, nA, mB, nB)
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
		fmt.Printf("errDim.Error(): %v | mA: %v nA: %v mB: %v nB: %v\n", errDim.Error(), mA, nA, mB, nB)
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
		fmt.Printf("errDim.Error(): %v | mA: %v nA: %v mB: %v nB: %v\n", errDim.Error(), mA, nA, mB, nB)
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
