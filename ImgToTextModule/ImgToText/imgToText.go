package main

import (
	"bufio"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	convertPageImgToArrayOfBinaryLetters()
}

func convertPageImgToArrayOfBinaryLetters() {
	path := "ImgData/SamplePdfScreenshot1.png"

	// img, err := imgio.ImreadGray(path)
	// if err != nil {
	// 	fmt.Printf("Could not read image from path: %s", path)
	// 	panic(err)
	// }

	imgFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("File Opening err.Error(): %v\n", err.Error())
		panic(err)
	}

	reader := bufio.NewReader(imgFile)
	colorImage, err := png.Decode(reader)
	if err != nil {
		fmt.Printf("Image Reading err.Error(): %v\n", err.Error())
		panic(err)
	}

	// image, err := threshold.OtsuThreshold(img, threshold.ThreshBinary)
	// if err != nil {
	// 	fmt.Printf("Otsu Binarized Conversion err.Error(): %v\n", err.Error())
	// 	panic(err)
	// }

	fmt.Printf("colorImage.Bounds(): %v\n", colorImage.Bounds())
	fmt.Printf("colorImage.Bounds().Max: %v\n", colorImage.Bounds().Max)
	fmt.Printf("colorImage.Bounds().Min: %v\n", colorImage.Bounds().Min)

	// fmt.Printf("image.At(0, 0): %v\n", image.At(0, 0))
	// fmt.Println(image.At(0, 0).RGBA())
	// fmt.Printf("color.White: %v\n", color.White)
	// fmt.Printf("color.Black: %v\n", color.Black)
	// fmt.Println(color.White.RGBA())
	// fmt.Println(colors(image.At(0, 0)))

	pageLinesOccurances := []int{}
	imgBoundMaxX, imgBoundMaxY := colorImage.Bounds().Max.X, colorImage.Bounds().Max.Y

	fmt.Printf("imgBoundMaxX: %v\nimgBoundMaxY: %v\n", imgBoundMaxX, imgBoundMaxY)

	image := [][]int{}
	for y := 0; y < imgBoundMaxY; y++ {
		imgRow := []int{}
		for x := 0; x < imgBoundMaxX; x++ {
			imgRow = append(imgRow, convertToGreyscale(colorImage.At(x, y)))
		}
		image = append(image, imgRow)
	}

	fmt.Printf("image0Len: %v\nimageLen: %v\n", len(image[0]), len(image))

	for y := 0; y < imgBoundMaxY; y++ {
		pageLinesOccurances = append(pageLinesOccurances, 0)
		for x := 0; x < imgBoundMaxX; x++ {
			if image[y][x] < 177 {
				pageLinesOccurances[y] = 1
				continue
			}
		}
	}

	// fmt.Printf("pageLinesOccurances: %v\n", pageLinesOccurances)

	lineHeightToCount := map[int]int{}
	tempLineCount := 0
	lineHeightBounds := [][]int{}

	for i := 0; i < imgBoundMaxY; i++ {
		if pageLinesOccurances[i] == 0 {
			if i != 0 && pageLinesOccurances[i-1] == 1 {
				lineHeightToCount[tempLineCount] += 1
				lineHeightBounds = append(lineHeightBounds, []int{i - tempLineCount, i})
			}
			tempLineCount = 0
		} else {
			tempLineCount++
		}
	}

	if pageLinesOccurances[imgBoundMaxY-1] == 1 {
		lineHeightToCount[tempLineCount] += 1
		lineHeightBounds = append(lineHeightBounds, []int{imgBoundMaxY - tempLineCount, imgBoundMaxX})
	}

	fmt.Printf("lineHeightToCount: %v\n", lineHeightToCount)
	fmt.Printf("lineHeightBounds: %v\n", lineHeightBounds)
	fmt.Printf("len(lineHeightBounds): %v\n", len(lineHeightBounds))

	numLines := len(lineHeightBounds)

	pageCharacterOccurances := [][]int{}

	for l := 0; l < numLines; l++ {
		lineCharacterOccurances := []int{}
		for x := 0; x < imgBoundMaxX; x++ {
			lineCharacterOccurances = append(lineCharacterOccurances, 0)
			for y := lineHeightBounds[l][0]; y < lineHeightBounds[l][1]; y++ {
				if !(255-image[y][x] < 90) {
					lineCharacterOccurances[x] = 1
					continue
				}
			}
		}
		pageCharacterOccurances = append(pageCharacterOccurances, lineCharacterOccurances)
	}

	// fmt.Printf("pageCharacterOccurances[0]: %v\n", pageCharacterOccurances[0])

	pageContent := []string{}

	pageCharWidthToCount := []map[int]int{}
	pageCharWidthBounds := [][][]int{}

	for l := 0; l < numLines; l++ {

		charWidthToCount := map[int]int{}
		tempCharCount := 0
		tempSpaceCount := 0
		charWidthBounds := [][]int{}

		for i := 0; i < imgBoundMaxX; i++ {
			if pageCharacterOccurances[l][i] == 0 {
				if i != 0 && pageCharacterOccurances[l][i-1] == 1 {
					charWidthToCount[tempCharCount] += 1
					charWidthBounds = append(charWidthBounds, []int{i - tempCharCount, i})
					pageContent = append(pageContent, "character")
				}
				tempCharCount = 0
				tempSpaceCount += 1
			} else {
				if tempSpaceCount >= 5 {
					pageContent = append(pageContent, "wordbreak")
				}
				tempCharCount++
				tempSpaceCount = 0
			}
		}

		if pageCharacterOccurances[l][imgBoundMaxX-1] == 1 {
			charWidthToCount[tempCharCount] += 1
			charWidthBounds = append(charWidthBounds, []int{imgBoundMaxX - tempCharCount, imgBoundMaxX})
		}

		pageCharWidthToCount = append(pageCharWidthToCount, charWidthToCount)
		pageCharWidthBounds = append(pageCharWidthBounds, charWidthBounds)

		pageContent = append(pageContent, "newline")
	}

	fmt.Printf("charWidthToCount: %v\n", pageCharWidthToCount[0])
	fmt.Printf("charWidthBounds: %v\n", pageCharWidthBounds[0])
	fmt.Printf("len(charWidthBounds): %v\n", len(pageCharWidthBounds[0]))

	fmt.Printf("len(pageContent): %v\n", len(pageContent))

	pageContentCounter := 0
	singleLetters := [][][]int{}
	for i := 0; i < numLines; i++ {
		numChars := len(pageCharWidthBounds[i])
		for j := 0; j < numChars; j++ {
			letter := [][]int{}
			for y := lineHeightBounds[i][0]; y < lineHeightBounds[i][1]; y++ {
				row := []int{}
				for x := pageCharWidthBounds[i][j][0]; x < pageCharWidthBounds[i][j][1]; x++ {
					// _, _, _, a := image.At(x, y).RGBA()
					// row = append(row, a)
					// greyValue := convertToBinary(image.At(x, y).RGBA())
					greyValue := image[y][x]
					insertValue := 255 - greyValue
					// if insertValue < 155 {
					// 	insertValue = 0
					// }
					row = append(row, insertValue)

				}
				letter = append(letter, row)
			}
			for pageContent[pageContentCounter] != "character" {
				pageContentCounter++
			}
			singleLetters = append(singleLetters, letter)
			pageContent[pageContentCounter] = fmt.Sprint(len(singleLetters) - 1)
		}
	}

	fmt.Printf("pageContent[:50]: %v\n", pageContent[:50])

	splitUpOfLetterLengths := calcSplitUpOfLetterLengths(singleLetters)

	fmt.Println(splitUpOfLetterLengths)

	maxOccur := math.MinInt
	maxOccurLen := -1
	for letterLen, occurrences := range splitUpOfLetterLengths {
		// fmt.Printf("letterLen - occurrences: %v - %v\n", letterLen, occurrences)
		if occurrences > maxOccur {
			maxOccur = occurrences
			maxOccurLen = letterLen
		}
	}

	fmt.Printf("maxOccurLen: %v\n", maxOccurLen)

	fmt.Printf("So we know, if a \"letter\" has more than 2*maxOccurLen-1 (%v) length, it needs to be split by maxOccurLen\n", 2*maxOccurLen-1)

	newSplitUpSingleLetters := [][][]int{}
	for _, letter := range singleLetters {
		if len(letter[0]) < 2*maxOccurLen-2 {
			newSplitUpSingleLetters = append(newSplitUpSingleLetters, letter)
			continue
		}

		splitParts := splitApparentLetter(letter, maxOccurLen)
		// printCombinedLetterStrings(splitParts, 0, len(splitParts))

		newSplitUpSingleLetters = append(newSplitUpSingleLetters, splitParts...)

		// for _, splitLetter := range splitParts {
		// 	newSplitUpSingleLetters = append(newSplitUpSingleLetters, splitLetter)
		// }
	}

	printCombinedLetterStrings(singleLetters, 0, 50)

	printCombinedLetterStrings(newSplitUpSingleLetters, 0, 50)
	// printLetter(singleLetters, 0)
	// printLetter(singleLetters, 1)
	// printLetter(singleLetters, 2)
	// printLetter(singleLetters, 3)
	// printLetter(singleLetters, 4)
	// printLetter(singleLetters, 5)
	// printLetter(singleLetters, 6)
	// printLetter(singleLetters, 7)
	// printLetter(singleLetters, 8)
	// printLetter(singleLetters, 9)
	// printLetter(singleLetters, 10)
	// printLetter(singleLetters, 11)
	// printLetter(singleLetters, 12)
	// printLetter(singleLetters, 13)

	// printLetter(singleLetters, 19)
	// printLetter(singleLetters, 20)
	// printLetter(singleLetters, 21)
	// printLetter(singleLetters, 42)
	// printLetter(singleLetters, 88)
	// printLetter(singleLetters, 89)
	// printLetter(singleLetters, 90)
	// printLetter(singleLetters, 194)
	// printLetter(singleLetters, 195)
	// printLetter(singleLetters, 196)
	// printLetter(singleLetters, 197)
	// printLetter(singleLetters, 198)

	// println("=========UUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUU=======")
	// numChars := len(pageCharWidthBounds[4])
	// fmt.Printf("numChars: %v\n", numChars)
	// for y := lineHeightBounds[4][0]; y < lineHeightBounds[4][1]; y++ {
	// 	for x := pageCharWidthBounds[4][numChars/2][0]; x < pageCharWidthBounds[4][3*numChars/4][1]; x++ {
	// 		r, g, b, _ := image.At(x, y).RGBA()
	// 		fmt.Printf("%v, %v, %v| ", r, g, b)
	// 	}
	// 	println()
	// }
}

func splitApparentLetter(letter [][]int, splitLen int) [][][]int {
	splitParts := [][][]int{}

	ogHeight := len(letter)
	ogLen := len(letter[0])

	for i := 0; i < ogLen/splitLen; i++ {
		singleSplitLetter := [][]int{}
		for r := 0; r < ogHeight; r++ {
			splitLetterRow := []int{}
			for c := 0; c < splitLen && i*splitLen+c < ogLen; c++ {
				splitLetterRow = append(splitLetterRow, letter[r][i*splitLen+c])
			}
			singleSplitLetter = append(singleSplitLetter, splitLetterRow)
		}
		splitParts = append(splitParts, singleSplitLetter)
	}

	if ogLen%splitLen == 0 {
		return splitParts
	}

	singleSplitLetter := [][]int{}
	for r := 0; r < ogHeight; r++ {
		splitLetterRow := []int{}
		for c := 0; c < splitLen && c < ogLen%splitLen; c++ {
			splitLetterRow = append(splitLetterRow, letter[r][(ogLen/splitLen)*splitLen+c])
		}
		singleSplitLetter = append(singleSplitLetter, splitLetterRow)
	}
	splitParts = append(splitParts, singleSplitLetter)

	return splitParts
}

func calcSplitUpOfLetterLengths(singleLetters [][][]int) map[int]int {
	returnMap := map[int]int{}
	for _, letter := range singleLetters {
		letterLen := len(letter[0])
		numLetters, ok := returnMap[letterLen]
		if ok {
			returnMap[letterLen] = numLetters + 1
		} else {
			returnMap[letterLen] = 1
		}
	}
	return returnMap
}

func printCombinedLetterStrings(singleLetters [][][]int, start, end int) {
	numLetters := len(singleLetters)
	constructedStrings := []string{}
	for i := 0; i < numLetters && i >= start && i < end; i++ {
		for k, row := range singleLetters[i] {
			for _, c := range row {
				numString := fmt.Sprintf("%v", c)
				lengthNumber := len(numString)
				switch lengthNumber {
				case 1:
					numString = fmt.Sprintf("%v   ", numString)
					// fmt.Print(numString)
					if k >= len(constructedStrings) {
						constructedStrings = append(constructedStrings, numString)
					} else {
						constructedStrings[k] = constructedStrings[k] + numString
					}
				case 2:
					numString = fmt.Sprintf("%v  ", numString)
					// fmt.Print(numString)
					if k >= len(constructedStrings) {
						constructedStrings = append(constructedStrings, numString)
					} else {
						constructedStrings[k] = constructedStrings[k] + numString
					}
				case 3:
					numString = fmt.Sprintf("%v ", numString)
					// fmt.Print(numString)
					if k >= len(constructedStrings) {
						constructedStrings = append(constructedStrings, numString)
					} else {
						constructedStrings[k] = constructedStrings[k] + numString
					}
				default:
					fmt.Printf("%v|", numString)
				}
			}
			constructedStrings[k] = constructedStrings[k] + "     "
		}
	}

	for _, stringLetters := range constructedStrings {
		fmt.Printf("%v\n", stringLetters)
	}

}

func printLetter(singleLetters [][][]int, letterNum int) {
	fmt.Printf("singleLetters[%v]: \n", letterNum)
	for _, row := range singleLetters[letterNum] {
		for _, c := range row {
			numString := fmt.Sprintf("%v", c)
			lengthNumber := len(numString)
			switch lengthNumber {
			case 1:
				numString = fmt.Sprintf("%v   ", numString)
				fmt.Print(numString)
			case 2:
				numString = fmt.Sprintf("%v  ", numString)
				fmt.Print(numString)
			case 3:
				numString = fmt.Sprintf("%v ", numString)
				fmt.Print(numString)
			default:
				fmt.Printf("%v|", numString)
			}
		}
		fmt.Println()
	}
}

func convertToBinary(r, g, b, _ uint32) int {
	trld := uint32(0)
	condition := (r == trld && g == trld && b == trld)
	// condition := ((0.2126 * float32(r)) + (0.7152 * float32(g)) + (0.0722 * float32(b))) <= float32(200)
	if condition {
		// fmt.Printf("r: %v, g: %v\n", r, g)
		return 1
	}
	return 0
}

func convertToGreyscale(color color.Color) int {
	r, g, b, _ := color.RGBA()
	lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	// lum := (r + g + b) / 3
	return int(lum / 256)
}

func colors(alphaMultipliedColor color.Color) (uint32, uint32, uint32) {
	cR, cG, cB, cA := alphaMultipliedColor.RGBA()
	return cR * 255 / cA, cG * 255 / cA, cB * 255 / cA
}
