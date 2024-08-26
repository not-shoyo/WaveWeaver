package main

import (
	neuralnet "WaveWeaver/ImgToTextModule/NeuralNet"
)

func main() {
	neuralnet.TrainNeuralNetwork("MnistTrainingData/emnist-digits-train.csv", 1, 500, 0.1, -1, -1, -1, 10, 10, "NeuralNetWeights/testWeightsSaving_3.txt")
	neuralnet.TestNeuralNet("MnistTrainingData/emnist-digits-test.csv", -1, "NeuralNetWeights/testWeightsSaving_3.txt")
}
