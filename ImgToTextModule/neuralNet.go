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

	cmdLineNumItrs, cmdLineAlpha, cmdLineNumRecords, cmdLineNumTestRecords, cmdLineNumFeatures := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]

	numIterations, errAtoiNumIterations := strconv.Atoi(cmdLineNumItrs)
	if errAtoiNumIterations != nil {
		fmt.Printf("errAToINumIterations.Error(): %v\n", errAtoiNumIterations.Error())
		panic(errAtoiNumIterations)
	}

	alpha, errAtoiAlpha := strconv.ParseFloat(cmdLineAlpha, 64)
	if errAtoiAlpha != nil {
		fmt.Printf("errAToIAlpha.Error(): %v\n", errAtoiAlpha.Error())
		panic(errAtoiAlpha)
	}

	numRecords, errAtoiNumRecords := strconv.Atoi(cmdLineNumRecords)
	if errAtoiNumRecords != nil {
		fmt.Printf("errAToINumRecords.Error(): %v\n", errAtoiNumRecords.Error())
		panic(errAtoiNumRecords)
	}

	numTestRecords, errAtoiNumTestRecords := strconv.Atoi(cmdLineNumTestRecords)
	if errAtoiNumTestRecords != nil {
		fmt.Printf("errAToINumTestRecords.Error(): %v\n", errAtoiNumTestRecords.Error())
		panic(errAtoiNumTestRecords)
	}

	numFeatures, errAtoiNumFeatures := strconv.Atoi(cmdLineNumFeatures)
	if errAtoiNumFeatures != nil {
		fmt.Printf("errAToINumFeatures.Error(): %v\n", errAtoiNumFeatures.Error())
		panic(errAtoiNumFeatures)
	}

	fmt.Printf("numIterations: %v, alpha: %v, numRecords: %v, numTestRecords: %v, numFeatures: %v\n", numIterations, alpha, numRecords, numTestRecords, numFeatures)

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

	if numRecords != -1 {
		allRecords = allRecords[:numRecords]
	} else {
		numRecords = len(allRecords)
	}

	totalTrainingM, totalTrainingN := len(allRecords), len(allRecords[0])-1

	fmt.Printf("totalTrainingM: %v\n", totalTrainingM)
	fmt.Printf("totalTrainingN: %v\n", totalTrainingN)

	allYs, allXs := extractToIntSlices(allRecords)

	fmt.Printf("len(allXs): %v\n", len(allXs))
	fmt.Printf("len(allXs[0]): %v\n", len(allXs[0]))

	xTranspose := transposeMatrice(allXs)
	xTranspose = selectFeatures(xTranspose, numFeatures)
	xTranspose = normalizeMatrix(xTranspose)
	yOneHot := convertToOneHot(allYs)

	fmt.Printf("allXsDim: (%v, %v)\n", len(allXs), len(allXs[0]))
	fmt.Printf("xTransposeDim: (%v, %v)\n", len(xTranspose), len(xTranspose[0]))
	fmt.Printf("allYsDim: (%v)\n", len(allYs))
	fmt.Printf("yOneHotDim: (%v, %v)\n", len(yOneHot), len(yOneHot[0]))

	W1, W2, b1, b2 := initialize(len(xTranspose), len(xTranspose[0]), 10, 10)

	fmt.Printf("W1Dim: (%v, %v)\n", len(W1), len(W1[0]))
	fmt.Printf("W2Dim: (%v, %v)\n", len(W2), len(W2[0]))
	fmt.Printf("b1Dim: (%v, %v)\n", len(b1), len(b1[0]))
	fmt.Printf("b2Dim: (%v, %v)\n", len(b2), len(b2[0]))

	fmt.Println("Weights before iterations:")
	fmt.Printf("W1: %v\n", W1)
	fmt.Printf("W2: %v\n", W2)

	fmt.Println("Biases before iterations:")
	fmt.Printf("b1: %v\n", b1)
	fmt.Printf("b2: %v\n", b2)

	fmt.Println("Weights and Biases during iterations:")

	// numIterations, alpha := 100, 10.0
	for itr := 0; itr < numIterations; itr++ {

		fmt.Printf("\n\n================================= %v =================================\n\n", itr)

		A1, A2, Z1, _ := feedForward(xTranspose, W1, W2, b1, b2)
		_, dW2, dB2, dW1, dB1 := calcErrors(yOneHot, A2, A1, Z1, xTranspose)

		// fmt.Printf("dZ2: %v\n", dZ2)
		totalErrors, accuracy, predictions, actuals := calcFinalErrorsAndAccuracy(A2, yOneHot)
		fmt.Printf("ITR[%v] - errorRate: %v, accuracy: %v\n", itr, totalErrors, accuracy)

		W1, W2, b1, b2 = updateWeightsAndBias(W1, W2, b1, b2, dW2, dB2, dW1, dB1, alpha)

		// fmt.Printf("W1: %v\n", W1)
		// fmt.Printf("W2: %v\n", W2)
		// fmt.Printf("b1: %v\n", b1)
		// fmt.Printf("b2: %v\n", b2)

		fmt.Printf("Actuals: 			%v\n", actuals)
		fmt.Printf("Predictions: 	%v\n", predictions)

		// if itr <= 3 {
		// 	fmt.Printf("transpose(A2)[0]			: %v\n", transposeMatrice(A2)[0])
		// 	fmt.Printf("transpose(yOneHot)[0]	: %v\n", transposeMatrice(yOneHot)[0])
		// 	fmt.Printf("transpose(A2)[1]			: %v\n", transposeMatrice(A2)[1])
		// 	fmt.Printf("transpose(yOneHot)[1]	: %v\n", transposeMatrice(yOneHot)[1])
		// 	fmt.Printf("transpose(A2)[2]			: %v\n", transposeMatrice(A2)[2])
		// 	fmt.Printf("transpose(yOneHot)[2]	: %v\n", transposeMatrice(yOneHot)[2])
		// }

	}

	testFile, errOs := os.Open("MnistTrainingData/emnist-digits-test.csv")
	if errOs != nil {
		fmt.Printf("errOs.Error(): %v\n", errOs.Error())
		panic(errOs)
	}

	defer testFile.Close()

	testReader := csv.NewReader(testFile)
	testRecords, errRead := testReader.ReadAll()
	if errRead != nil {
		fmt.Printf("errRead.Error(): %v\n", errRead.Error())
		panic(errRead)
	}

	// fmt.Printf("testRecords[0]: %v\n", testRecords[0])

	// testIndex := 0

	if numTestRecords != -1 {
		testRecords = testRecords[:numTestRecords]
	} else {
		numTestRecords = len(testRecords)
	}

	// fmt.Printf("testRecords[testIndex]: %v\n", testRecords[testIndex])

	testYs, testXs := extractToIntSlices(testRecords)

	// testInput := [][]float64{}
	// testInput = append(testInput, testXs[testIndex])
	// testInput = transposeMatrice(testInput)

	// testOutput := convertToOneHot([]float64{testYs[testIndex]})

	testInput := normalizeMatrix(selectFeatures(transposeMatrice(testXs), numFeatures))

	testOutput := convertToOneHot(testYs)

	// bias1 := [][]float64{}
	// bias1 = append(bias1, transposeMatrice(b1)[testIndex])
	// bias1 = transposeMatrice(bias1)
	bias1 := expandMatrix(divideMatrixBy(addUpRows(b1), float64(numRecords)), numTestRecords)

	// bias2 := [][]float64{}
	// bias2 = append(bias2, transposeMatrice(b2)[testIndex])
	// bias2 = transposeMatrice(bias2)
	bias2 := expandMatrix(divideMatrixBy(addUpRows(b2), float64(numRecords)), numTestRecords)

	_, A2, _, _ := feedForward(testInput, W1, W2, bias1, bias2)
	// _, _, _, _, _ := calcErrors(testOutput, A2, A1, Z1, testInput)
	totalErrors, accuracy, testPredictions, testActuals := calcFinalErrorsAndAccuracy(A2, testOutput)

	fmt.Print("\n\n================================= TESTS =================================\n\n")

	fmt.Printf("TestOutput - errorRate: %v, accuracy: %v\n", totalErrors, accuracy)

	fmt.Printf("TestActuals:	 		%v\n", testActuals)
	fmt.Printf("TestPredictions: 	%v\n", testPredictions)
}

func selectFeatures(A [][]float64, numFeatures int) [][]float64 {
	returnMatrix := [][]float64{}

	if len(A) < numFeatures {
		errLen := errors.New("illegal Operation: Cant include numFeatures due to matrix being smaller")
		fmt.Printf("errLen.Error(): %v | len(A): %v numFeatures: %v\n", errLen.Error(), len(A), numFeatures)
		panic(errLen)
	}

	if numFeatures < 0 {
		numFeatures = len(A)
	}

	for i, row := range A {
		if i >= numFeatures {
			break
		}
		newRow := []float64{}
		newRow = append(newRow, row...)
		returnMatrix = append(returnMatrix, newRow)
	}

	// fmt.Printf("selectFeatures - returnMatrix: %v\n", returnMatrix)
	return returnMatrix
}

func normalizeMatrix(A [][]float64) [][]float64 {
	returnMatrix := [][]float64{}

	// maxValue := 255 // pixel values
	maxValue := float64(math.MinInt)

	for _, row := range A {
		maxValue = math.Max(maxValue, findMax(row))
	}

	if maxValue == 0 {
		return A
	}

	for _, row := range A {
		newRow := []float64{}
		for _, v := range row {
			newRow = append(newRow, v/float64(maxValue))
		}
		returnMatrix = append(returnMatrix, newRow)
	}

	return returnMatrix
}

func calcFinalErrorsAndAccuracy(A2, Y [][]float64) (float64, float64, []int, []int) {
	m, n := len(A2), len(A2[0])
	sumError, correct := 0.0, 0
	actuals, predictions := []int{}, []int{}

	for j := 0; j < n; j++ {
		sum, maxPossibility, predicted, actualAnswer := 0.0, float64(math.MinInt), -1, -1
		for i := 0; i < m; i++ {
			sum += math.Abs(A2[i][j] - Y[i][j])
			if math.Abs(A2[i][j]) > maxPossibility {
				maxPossibility = A2[i][j]
				predicted = i
			}
			if Y[i][j] == 1 {
				actualAnswer = i
				actuals = append(actuals, actualAnswer)
			}
		}
		sumError += sum / float64(m)
		predictions = append(predictions, predicted)
		if predicted == actualAnswer {
			correct += 1
		}
	}

	return sumError / float64(n), float64(correct) / float64(n), predictions, actuals
}

func updateWeightsAndBias(W1, W2, b1, b2, dW2, dB2, dW1, dB1 [][]float64, alpha float64) ([][]float64, [][]float64, [][]float64, [][]float64) {

	// fmt.Println("updateWeightsAndBias - ")
	// fmt.Printf("dW1: 		%v\n", dW1)

	newW1 := subtactMatrices(W1, multipyMatrixBy(dW1, alpha))
	newW2 := subtactMatrices(W2, multipyMatrixBy(dW2, alpha))
	newb1 := subtactMatrices(b1, multipyMatrixBy(dB1, alpha))
	newb2 := subtactMatrices(b2, multipyMatrixBy(dB2, alpha))

	// fmt.Printf("W1: 		%v\n", W1)
	// fmt.Printf("newW1: 	%v\n", newW1)
	// fmt.Printf("W2: 		%v\n", W2)
	// fmt.Printf("newW2: 	%v\n", newW2)

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
			weightRow = append(weightRow, -0.5+((rand.Float64()*10)+1)/10) //0.0
		}
		W1 = append(W1, weightRow)
	}

	for i := 0; i < l1; i++ {
		biasRow := []float64{}
		for i := 0; i < m; i++ {
			biasRow = append(biasRow, -0.5+((rand.Float64()*10)+1)/10) //0.0
		}
		b1 = append(b1, biasRow)
	}

	for j := 0; j < l2; j++ {
		weightRow := []float64{}
		for i := 0; i < l1; i++ {
			weightRow = append(weightRow, -0.5+((rand.Float64()*10)+1)/10) //0.0
		}
		W2 = append(W2, weightRow)
	}

	for i := 0; i < l2; i++ {
		biasRow := []float64{}
		for i := 0; i < m; i++ {
			biasRow = append(biasRow, -0.5+((rand.Float64()*10)+1)/10) //0.0
		}
		b2 = append(b2, biasRow)
	}

	return W1, W2, b1, b2
}

func feedForward(X, W1, W2, b1, b2 [][]float64) (A1, A2, Z1, Z2 [][]float64) {

	// fmt.Println("feedForward:")
	// fmt.Printf("X : %v\n", X)
	// fmt.Printf("W1: %v\n", W1)
	// fmt.Printf("W2: %v\n", W2)
	// fmt.Printf("b1: %v\n", b1)
	// fmt.Printf("b2: %v\n", b2)

	Z1 = addMatrices(crossProduct(W1, X), b1)
	// fmt.Printf("feedForward -\nZ1: %v\n", Z1)
	Z1 = normalizeMatrix(Z1)
	// fmt.Printf("feedForward -\nZ1 after normalizing: %v\n", Z1)
	A1 = activateNeurons(Z1, "reLU")
	// fmt.Printf("feedForward -\nA1: %v\n", A1)
	Z2 = addMatrices(crossProduct(W2, A1), b2)
	// fmt.Printf("feedForward -\nZ2: %v\n", Z2)
	Z2 = normalizeMatrix(Z2)
	// fmt.Printf("feedForward -\nZ2 after normalizing: %v\n", Z2)
	A2 = activateNeurons(Z2, "softmax")
	// fmt.Printf("feedForward -\nA2: %v\n", A2)

	fmt.Printf("feedForward - Z1Dim: (%v, %v)\n", len(Z1), len(Z1[0]))
	fmt.Printf("feedForward - A1Dim: (%v, %v)\n", len(A1), len(A1[0]))
	fmt.Printf("feedForward - Z2Dim: (%v, %v)\n", len(Z2), len(Z2[0]))
	fmt.Printf("feedForward - A2Dim: (%v, %v)\n", len(A2), len(A2[0]))

	return
}

func calcErrors(Y, A2, A1, Z1, X [][]float64) ([][]float64, [][]float64, [][]float64, [][]float64, [][]float64) {
	m := len(Y[0])

	dZ2 := subtactMatrices(A2, Y)

	// fmt.Printf("calcErrors - dZ2: %v\n", dZ2)

	// fmt.Printf("calcErrors - dZ2Dim: (%v, %v)\n", len(dZ2), len(dZ2[0]))

	dW2 := crossProduct(dZ2, transposeMatrice(A1))
	dW2 = divideMatrixBy(dW2, float64(m))

	// fmt.Printf("calcErrors - dW2Dim: (%v, %v)\n", len(dW2), len(dW2[0]))

	dB2 := addUpRows(dZ2)
	dB2 = divideMatrixBy(dB2, float64(m))
	dB2 = expandMatrix(dB2, m)

	// fmt.Printf("calcErrors - dB2Dim: (%v, %v)\n", len(dB2), len(dB2[0]))

	gZ1 := undoActivationReLU(Z1)

	// fmt.Printf("calcErrors -  gZ1: 		%v\n", gZ1)

	// fmt.Printf("calcErrors - gZ1Dim: (%v, %v)\n", len(gZ1), len(gZ1[0]))

	dZ1 := dotProduct(crossProduct(transposeMatrice(dW2), dZ2), gZ1)

	// fmt.Printf("calcErrors -  dZ1: 		%v\n", dZ1)

	// fmt.Printf("calcErrors - dZ1Dim: (%v, %v)\n", len(dZ1), len(dZ1[0]))

	dW1 := crossProduct(dZ1, transposeMatrice(X))
	// fmt.Printf("calcErrors -  dW1: 		%v\n", dW1)
	dW1 = divideMatrixBy(dW1, float64(m))
	// fmt.Printf("calcErrors -  dW1: 		%v\n", dW1)

	// fmt.Printf("calcErrors - dW1Dim: (%v, %v)\n", len(dW1), len(dW1[0]))

	dB1 := addUpRows(dZ1)
	dB1 = divideMatrixBy(dB1, float64(m))
	dB1 = expandMatrix(dB1, m)

	// fmt.Printf("calcErrors - dB1Dim: (%v, %v)\n", len(dB1), len(dB1[0]))

	return dZ2, dW2, dB2, dW1, dB1
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
		errDim := errors.New("illegal Operation: Matrix Dot Product not possible due to dimentions")
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
			// if j == 0 {
			// 	fmt.Printf("activateNeurons -\n")
			// }
			for i := 0; i < mZ; i++ {
				sum += math.Exp(Z[i][j])
				// if j == 0 {
				// 	fmt.Printf("Z[i][0]: %v, sum: %v\n", Z[i][j], sum)
				// }
			}
			denominators = append(denominators, sum)
			// if j == 0 {
			// 	fmt.Printf("activateNeurons -\ndenominators: %v\n", denominators)
			// }
		}

		for i := 0; i < mZ; i++ {
			singleNeuron := []float64{}
			for j := 0; j < nZ; j++ {
				singleNeuron = append(singleNeuron, math.Exp(Z[i][j])/denominators[j])
			}
			activatedNeurons = append(activatedNeurons, singleNeuron)
		}

	default:
		errActivation := errors.New("illegal Activation Function")
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
		errDim := errors.New("illegal Operation: Cross Product not possible due to dimentions")
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
		errDim := errors.New("illegal Operation: Matrix Addition not possible due to dimentions")
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

	// fmt.Printf("subtractMatrices -\nA: %v\nB: %v\n", A, B)

	if mA != mB || nA != nB {
		errDim := errors.New("illegal Operation: Matrix Subtraction not possible due to dimentions")
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
		// if i == 0 {
		// 	fmt.Printf("subtractMatrices -\nresultMatrix[0]: %v\n", resultMatrix)
		// }
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

	for j := 0; j < max(10, int(maxNum)+1); j++ {
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
		if math.Abs(num) > maxSoFar {
			maxSoFar = math.Abs(num)
		}
	}
	return maxSoFar
}
