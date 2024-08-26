package neuralnet

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func TestNeuralNet(testingFile string, numTestRecords int, W1, b1, W2, b2 [][]float64) {

	// W1, b1, W2, b2 := ImportSavedWeights(savedWeightsFileName)

	testFile, errOs := os.Open(testingFile)
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

	testInput := normalizeMatrix(transposeMatrice(testXs))

	testOutput := convertToOneHot(testYs)

	// bias1 := [][]float64{}
	// bias1 = append(bias1, transposeMatrice(b1)[testIndex])
	// bias1 = transposeMatrice(bias1)
	bias1 := expandMatrix(b1, numTestRecords)

	// bias2 := [][]float64{}
	// bias2 = append(bias2, transposeMatrice(b2)[testIndex])
	// bias2 = transposeMatrice(bias2)
	bias2 := expandMatrix(b2, numTestRecords)

	_, A2, _, _ := feedForward(testInput, W1, W2, bias1, bias2)
	// _, _, _, _, _ := calcErrors(testOutput, A2, A1, Z1, testInput)
	// totalErrors, accuracy, testPredictions, testActuals := calcFinalErrorsAndAccuracy(A2, testOutput)
	totalErrors, accuracy, _, _ := calcFinalErrorsAndAccuracy(A2, testOutput)

	fmt.Print("\n\n================================= TESTS =================================\n\n")

	fmt.Printf("TestOutput - errorRate: %v, accuracy: %v\n", totalErrors, accuracy)

	// fmt.Printf("TestActuals:	 		%v\n", testActuals)
	// fmt.Printf("TestPredictions: 	%v\n", testPredictions)
}

func ImportSavedWeights(savedWeightsFileName string) ([][]float64, [][]float64, [][]float64, [][]float64) {
	savedWeightsFile, errSavedWeightsFile := os.Open(savedWeightsFileName)
	if errSavedWeightsFile != nil {
		fmt.Printf("errSavedWeightsFile.Error(): %v\n", errSavedWeightsFile.Error())
		panic(errSavedWeightsFile)
	}

	defer savedWeightsFile.Close()

	fileScanner := bufio.NewScanner(savedWeightsFile)
	text := getText(fileScanner)

	numLayers, errNumLayers := strconv.Atoi(text)
	if errNumLayers != nil {
		fmt.Printf("errNumLayers.Error(): %v\n", errNumLayers.Error())
		panic(errNumLayers)
	}

	weightsDims := [][]int{}
	for i := 0; i < numLayers; i++ {
		weightDims := []int{}
		weightsDimStr := getText(fileScanner)
		weightsDimStrs := strings.Split(weightsDimStr, " ")
		for j := 0; j < 2; j++ {
			weightsDim, errWeightsDim := strconv.Atoi(weightsDimStrs[j])
			if errWeightsDim != nil {
				fmt.Printf("errWeightsDim.Error(): %v\n", errWeightsDim.Error())
				panic(errWeightsDim)
			}
			weightDims = append(weightDims, weightsDim)
		}
		weightsDims = append(weightsDims, weightDims)
	}

	biasesDims := []int{}
	for i := 0; i < numLayers; i++ {
		biasDim, errBiasDim := strconv.Atoi(getText(fileScanner))
		if errBiasDim != nil {
			fmt.Printf("errBiasDim.Error(): %v\n", errBiasDim.Error())
			panic(errBiasDim)
		}
		biasesDims = append(biasesDims, biasDim)
	}

	allWeights, allBiases := [][][]float64{}, [][][]float64{}

	for i := 0; i < numLayers; i++ {
		layerWeights := [][]float64{}
		for wd1 := 0; wd1 < weightsDims[i][0]; wd1++ {
			weightsRow := []float64{}
			weightsStr := strings.Split(getText(fileScanner), " ")
			for _, wStr := range weightsStr {
				weight := readFloat64(wStr)
				weightsRow = append(weightsRow, weight)
			}
			layerWeights = append(layerWeights, weightsRow)
		}

		layerBiases := [][]float64{}
		for bd1 := 0; bd1 < biasesDims[i]; bd1++ {
			biasesStr := strings.Split(getText(fileScanner), " ")
			for _, bStr := range biasesStr {
				bias := readFloat64(bStr)
				layerBiases = append(layerBiases, []float64{bias})
			}
		}
		allWeights = append(allWeights, layerWeights)
		allBiases = append(allBiases, layerBiases)
	}

	return allWeights[0], allBiases[0], allWeights[1], allBiases[1]
}

func getText(fileScanner *bufio.Scanner) string {
	fileScanner.Scan()
	text := fileScanner.Text()
	if len(text) > 0 {
		return text[:len(text)-1]
	}
	return text
}

func readFloat64(wStr string) float64 {
	weight, errWeight := strconv.ParseFloat(wStr, 64)
	if errWeight != nil {
		fmt.Printf("errWeight.Error(): %v\n", errWeight.Error())
		panic(errWeight)
	}
	return weight
}
