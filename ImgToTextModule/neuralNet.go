package main

import(
	"os"
	"encoding/csv"
)

func main()  {
	file, errOs := os.Open("ImgToTextModule/MnistTrainingData/emnist-byclass-train.csv")
	if errOs != nil {
		errOs.Error()	
	}
}