package main

import (
	neuralnet "WaveWeaver/ImgToTextModule/NeuralNet"
)

func main() {
	weightSavingFileName, trainingDataFileName, testingDataFileName := "NeuralNetWeights/testWeightsSaving_3.txt", "MnistTrainingData/emnist-digits-train.csv", "MnistTrainingData/emnist-digits-test.csv"
	neuralnet.TrainNeuralNetwork(trainingDataFileName, 1, 500, 0.1, -1, -1, -1, 10, 10, weightSavingFileName)
	W1, b1, W2, b2 := neuralnet.ImportSavedWeights(weightSavingFileName)
	neuralnet.TestNeuralNet(testingDataFileName, -1, W1, b1, W2, b2)
}
