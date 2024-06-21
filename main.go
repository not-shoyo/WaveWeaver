package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DylanMeeus/GoAudio/wave"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Reference Documentation https://dylanmeeus.github.io/
// Reference Documentation https://github.com/u2takey/ffmpeg-go/tree/master

func main() {
	reader := bufio.NewReader(os.Stdin)

	initialOptionsPrompt :=
		`Please specify what you want to do:
		a) Remove a portion of your recording
		b) Replace a portion of your recording with another full recording
		c) Replace a portion of your recording with a portion of another recording
`
	initialOptionsAllowedResponses := []string{"a", "b", "c"}

	userInput, err := getUserInput(initialOptionsPrompt, reader)

	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	} else {
		// fmt.Printf("userInput: %v\n", userInput)
		switch {
		case slices.Contains(initialOptionsAllowedResponses, userInput):
			switch userInput {
			case "a":
				removeRecordingPortion(reader)
			case "b":
				replaceRecordingPortionWithFull(reader)
			case "c":
				replaceRecordingPortionWithPart(reader)
			}
		default:
			fmt.Println("Invalid Choice")
		}
	}

}

func removeRecordingPortion(reader *bufio.Reader) {
	ogFilePath, err := getUserInput("Provide your original recording file path: ", reader)

	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	} else if !isValidFilePath(ogFilePath) {
		fmt.Printf("Entered FilePath: %v is invalid\n", ogFilePath)
	} else {

		performableFilePath := switchFileToWavIfNecessary(ogFilePath)

		ogRecording, ogRecordingErr := wave.ReadWaveFile(performableFilePath)
		if ogRecordingErr != nil {
			fmt.Printf("Encountered Error Reading: %v\n", performableFilePath)
			panic(ogRecordingErr)
		}
		ogRecordingDuration := findDuration(ogRecording)
		fmt.Printf("Duration of selected recording (%v) is: %vs\n", performableFilePath, ogRecordingDuration)

		startTimestamp, errStart := getUserInput("Provide timestamp in HH:MM:SS format for the beginning of the portion you want to remove: ", reader)
		stopTimestamp, errStop := getUserInput("Provide timestamp in HH:MM:SS format for the end of the portion you want to remove: ", reader)

		if errStart != nil && errStop != nil {
			fmt.Println("Internal errors in both entered Timestamps:")
			fmt.Printf("errStart.Error(): %v\n", errStart.Error())
			fmt.Printf("errStop.Error(): %v\n", errStop.Error())
		} else {
			validStart, validStop := isValidTimestamp(startTimestamp), isValidTimestamp(stopTimestamp)
			if !validStart || !validStop {
				if !validStart {
					fmt.Printf("Entered Start Timestamp: %v is invalid\n", startTimestamp)
				}
				if !validStop {
					fmt.Printf("Entered Stop Timestamp: %v is invalid\n", stopTimestamp)
				}
			} else {
				removeStartTime, removeStopTime := convertValidTimestampToTime(startTimestamp), convertValidTimestampToTime(stopTimestamp)

				if removeStopTime < 0 || removeStartTime > ogRecordingDuration || removeStartTime > removeStopTime {
					fmt.Println("Invalid / Nonsensical start and stop timestamps with respect to original recording:")
					fmt.Printf("\tDuration of Original Recording: %v\n", ogRecordingDuration)
					fmt.Printf("\tTime at which the removal is supposed to Start: %v\n", removeStartTime)
					fmt.Printf("\tTime at which the removal is supposed to Stop: %v\n", removeStopTime)
				} else {
					if removeStartTime < 0 || removeStopTime > ogRecordingDuration {
						if removeStartTime < 0 {
							removeStartTime = 0
						}
						if removeStopTime > ogRecordingDuration {
							removeStopTime = ogRecordingDuration
						}
					}
					modifiedRecordingFrames, modifiedRecordingFmt := performRemoval(ogRecording, removeStartTime, removeStopTime)

					saveNewRecording(reader, modifiedRecordingFrames, modifiedRecordingFmt)
				}

			}
		}
	}

}

func replaceRecordingPortionWithFull(reader *bufio.Reader) {

	ogFilePath, ogFileErr := getUserInput("Provide your original recording file path: ", reader)

	if ogFileErr != nil {
		fmt.Printf("ogFileErr: %v\n", ogFileErr.Error())
	} else if !isValidFilePath(ogFilePath) {
		fmt.Printf("Entered FilePath: %v is invalid\n", ogFilePath)
	} else {

		performableFilePath := switchFileToWavIfNecessary(ogFilePath)

		ogRecording, ogRecordingErr := wave.ReadWaveFile(performableFilePath)
		if ogRecordingErr != nil {
			fmt.Printf("Encountered Error Reading: %v\n", performableFilePath)
			panic(ogRecordingErr)
		}
		ogRecordingDuration := findDuration(ogRecording)
		fmt.Printf("Duration of selected recording (%v) is: %vs\n", ogFilePath, ogRecordingDuration)

		insertionFilePath, insertionFileErr := getUserInput("Provide the file path to the recording you want to insert: ", reader)

		if insertionFileErr != nil {
			fmt.Printf("insertionFileErr: %v\n", insertionFileErr.Error())
		} else if !isValidFilePath(insertionFilePath) {
			fmt.Printf("Entered FilePath: %v is invalid\n", insertionFilePath)
		} else {

			performableInsertionFilePath := switchFileToWavIfNecessary(insertionFilePath)

			insertionRecording, insertionRecordingErr := wave.ReadWaveFile(performableInsertionFilePath)
			if insertionRecordingErr != nil {
				fmt.Printf("Encountered Error Reading: %v\n", performableInsertionFilePath)
				panic(insertionRecordingErr)
			}
			insertionRecordingDuration := findDuration(insertionRecording)
			fmt.Printf("Duration of selected recording (%v) is: %vs\n", insertionFilePath, insertionRecordingDuration)

			startTimestamp, errStart := getUserInput("Provide timestamp in HH:MM:SS format for the beginning of the portion you want to remove in the original: ", reader)
			stopTimestamp, errStop := getUserInput("Provide timestamp in HH:MM:SS format for the end of the portion you want to remove in the original: ", reader)

			if errStart != nil && errStop != nil {
				fmt.Println("Internal errors in both entered Timestamps:")
				fmt.Printf("errStart.Error(): %v\n", errStart.Error())
				fmt.Printf("errStop.Error(): %v\n", errStop.Error())
			} else {
				validStart, validStop := isValidTimestamp(startTimestamp), isValidTimestamp(stopTimestamp)
				if !validStart || !validStop {
					if !validStart {
						fmt.Printf("Entered Start Timestamp: %v is invalid\n", startTimestamp)
					}
					if !validStop {
						fmt.Printf("Entered Stop Timestamp: %v is invalid\n", stopTimestamp)
					}
				} else {
					removeStartTime, removeStopTime := convertValidTimestampToTime(startTimestamp), convertValidTimestampToTime(stopTimestamp)

					if removeStopTime < 0 || removeStartTime > ogRecordingDuration || removeStartTime > removeStopTime {
						fmt.Println("Invalid / Nonsensical start and stop timestamps with respect to original recording:")
						fmt.Printf("\tDuration of Original Recording: %v\n", ogRecordingDuration)
						fmt.Printf("\tTime at which the removal is supposed to Start: %v\n", removeStartTime)
						fmt.Printf("\tTime at which the removal is supposed to Stop: %v\n", removeStopTime)
					} else {
						if removeStartTime < 0 || removeStopTime > ogRecordingDuration {
							if removeStartTime < 0 {
								removeStartTime = 0
							}
							if removeStopTime > ogRecordingDuration {
								removeStopTime = ogRecordingDuration
							}
						}
						modifiedRecordingFrames, modifiedRecordingFmt := performInsertion(ogRecording, insertionRecording.Frames, removeStartTime, removeStopTime)

						saveNewRecording(reader, modifiedRecordingFrames, modifiedRecordingFmt)
					}
				}
			}
		}
	}
}

func replaceRecordingPortionWithPart(reader *bufio.Reader) {

	ogFilePath, ogFileErr := getUserInput("Provide your original recording file path: ", reader)

	if ogFileErr != nil {
		fmt.Printf("ogFileErr: %v\n", ogFileErr.Error())
	} else if !isValidFilePath(ogFilePath) {
		fmt.Printf("Entered FilePath: %v is invalid\n", ogFilePath)
	} else {
		performableFilePath := switchFileToWavIfNecessary(ogFilePath)

		ogRecording, ogRecordingErr := wave.ReadWaveFile(performableFilePath)
		if ogRecordingErr != nil {
			fmt.Printf("Encountered Error Reading: %v\n", performableFilePath)
			panic(ogRecordingErr)
		}
		ogRecordingDuration := findDuration(ogRecording)
		fmt.Printf("Duration of selected recording (%v) is: %vs\n", ogFilePath, ogRecordingDuration)

		insertionFilePath, insertionFileErr := getUserInput("Provide the file path to the recording containing the portion you want to insert: ", reader)

		if insertionFileErr != nil {
			fmt.Printf("insertionFileErr: %v\n", insertionFileErr.Error())
		} else if !isValidFilePath(insertionFilePath) {
			fmt.Printf("Entered FilePath: %v is invalid\n", insertionFilePath)
		} else {
			performableInsertionFilePath := switchFileToWavIfNecessary(insertionFilePath)

			insertionRecording, insertionRecordingErr := wave.ReadWaveFile(performableInsertionFilePath)
			if insertionRecordingErr != nil {
				fmt.Printf("Encountered Error Reading: %v\n", performableInsertionFilePath)
				panic(insertionRecordingErr)
			}
			insertionRecordingDuration := findDuration(insertionRecording)
			fmt.Printf("Duration of selected recording (%v) is: %vs\n", insertionFilePath, insertionRecordingDuration)

			ogStartTimestamp, errOgStart := getUserInput("Provide timestamp in HH:MM:SS format for the beginning of the portion you want to remove in the original recording: ", reader)
			ogStopTimestamp, errOgStop := getUserInput("Provide timestamp in HH:MM:SS format for the end of the portion you want to remove in the original recording: ", reader)

			extractStartTimestamp, errExtractStart := getUserInput("Provide timestamp in HH:MM:SS format for the beginning of the portion you want to insert from the insertion recording: ", reader)
			extractStopTimestamp, errExtractStop := getUserInput("Provide timestamp in HH:MM:SS format for the end of the portion you want to insert from the insertion recording: ", reader)

			flawlessOg, ogRemoveStartTime, ogRemoveStopTime := timespaceComparisonsCheck(errOgStart, errOgStop, ogStartTimestamp, ogStopTimestamp, ogRecordingDuration)
			flawlessExtraction, extractStartTime, extractStopTime := timespaceComparisonsCheck(errExtractStart, errExtractStop, extractStartTimestamp, extractStopTimestamp, insertionRecordingDuration)

			if flawlessOg && flawlessExtraction {
				extractedRecordingFrames := performExtraction(insertionRecording, extractStartTime, extractStopTime)
				modifiedRecordingFrames, modifiedRecordingFmt := performInsertion(ogRecording, extractedRecordingFrames, ogRemoveStartTime, ogRemoveStopTime)
				saveNewRecording(reader, modifiedRecordingFrames, modifiedRecordingFmt)
			}
		}
	}
}

// call Insert Frame function to remove portions of original file
func performRemoval(ogRecording wave.Wave, removeStartTime, removeStopTime float64) ([]wave.Frame, wave.WaveFmt) {
	emptyWave := []wave.Frame{}
	sampleRate := ogRecording.SampleRate
	numChannels := ogRecording.NumChannels
	newFmt := ogRecording.WaveFmt

	return insertFrame(ogRecording.Frames, emptyWave, time.Duration(removeStartTime), time.Duration(removeStopTime), sampleRate, numChannels), newFmt
}

// insert the new recording into the original reocrding, replacing the section between startTime and stopTime with the new recording
func performInsertion(ogRecording wave.Wave, newRecordingFrames []wave.Frame, insertionStartTime, insertionStopTime float64) ([]wave.Frame, wave.WaveFmt) {
	sampleRate := ogRecording.SampleRate
	numChannels := ogRecording.NumChannels
	newFmt := ogRecording.WaveFmt

	return insertFrame(ogRecording.Frames, newRecordingFrames, time.Duration(insertionStartTime), time.Duration(insertionStopTime), sampleRate, numChannels), newFmt
}

// Extract the required duration from the recording
func performExtraction(recording wave.Wave, extractStartTime, extractStopTime float64) []wave.Frame {
	sampleRate := recording.SampleRate
	numChannels := recording.NumChannels

	return extractFrame(recording.Frames, time.Duration(extractStartTime), time.Duration(extractStopTime), sampleRate, numChannels)

}

// saves the new recording after validating the path
func saveNewRecording(reader *bufio.Reader, modifiedRecordingFrames []wave.Frame, modifiedRecordingFmt wave.WaveFmt) {
	resultFilePath, err := getUserInput("Provide file path to save new recording: ", reader)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
	} else {
		splitFilePath := strings.Split(resultFilePath, string(os.PathSeparator))
		if !isValidFilePath(strings.Join(splitFilePath[:len(splitFilePath)-1], string(os.PathSeparator))) {
			fmt.Printf("Entered FilePath: %v is invalid\n", resultFilePath)
		} else {
			resultFileErr := wave.WriteFrames(modifiedRecordingFrames, modifiedRecordingFmt, resultFilePath)
			if resultFileErr != nil {
				panic(resultFileErr)
			}
		}
	}
}

// compare timestamps and throw errors, checking for internal errors as well as format validity and comparison with actual record duration.
func timespaceComparisonsCheck(errStart, errStop error, startTimestamp, stopTimestamp string, recordingDuration float64) (bool, float64, float64) {
	flawless := true
	if errStart != nil && errStop != nil {
		fmt.Println("Internal errors in both entered Timestamps:")
		fmt.Printf("errStart.Error(): %v\n", errStart.Error())
		fmt.Printf("errStop.Error(): %v\n", errStop.Error())
		flawless = false
	} else {
		validStart, validStop := isValidTimestamp(startTimestamp), isValidTimestamp(stopTimestamp)
		if !validStart || !validStop {
			if !validStart {
				fmt.Printf("Entered Start Timestamp: %v is invalid\n", startTimestamp)
				flawless = false
			}
			if !validStop {
				fmt.Printf("Entered Stop Timestamp: %v is invalid\n", stopTimestamp)
				flawless = false
			}
		} else {
			startTime, stopTime := convertValidTimestampToTime(startTimestamp), convertValidTimestampToTime(stopTimestamp)

			if stopTime < 0 || startTime > recordingDuration || startTime > stopTime {
				fmt.Println("Invalid / Nonsensical start and stop timestamps with respect to original recording:")
				fmt.Printf("\tDuration of Original Recording: %v\n", recordingDuration)
				fmt.Printf("\tTime at which the removal is supposed to Start: %v\n", startTime)
				fmt.Printf("\tTime at which the removal is supposed to Stop: %v\n", stopTime)
				flawless = false
			} else {
				if startTime < 0 || stopTime > recordingDuration {
					if startTime < 0 {
						startTime = 0
					}
					if stopTime > recordingDuration {
						stopTime = recordingDuration
					}
				}
			}
			return flawless, startTime, stopTime
		}
	}
	return flawless, 0, 0
}

// Checks if the file specified is a wav. If not, we change it to a wav.
func switchFileToWavIfNecessary(ogFilePath string) string {
	performableFilePath := ogFilePath
	splitOgFilePath := strings.Split(ogFilePath, string(os.PathSeparator))
	splitFileName := strings.Split(splitOgFilePath[len(splitOgFilePath)-1], ".")
	if splitFileName[len(splitFileName)-1] != "wav" {
		fullWorkingDirectory, err := os.Getwd()
		if err != nil {
			fmt.Printf("file name split err.Error(): %v\n", err.Error())
		} else {
			remainingPath := strings.Join(splitOgFilePath[:len(splitOgFilePath)-1], string(os.PathSeparator))
			tempOutputileNameWoutExt := strings.Join(splitFileName[:len(splitFileName)-1], string(os.PathSeparator))
			basePathWoutExt := fullWorkingDirectory + string(os.PathSeparator) + remainingPath + string(os.PathSeparator) + tempOutputileNameWoutExt
			performableFilePath = basePathWoutExt + ".wav"

			err := ffmpeg.Input(basePathWoutExt + "." + splitFileName[len(splitFileName)-1]).
				Output(performableFilePath).
				OverWriteOutput().ErrorToStdOut().Run()

			if err != nil {

				panic(err.Error())
			} else {
				fmt.Println("Temp WAV File Created Successfully")
				fmt.Printf("performableFilePath: %v\n", performableFilePath)
			}
		}
	}
	return performableFilePath
}

// ensures the file exists and can be reached
func isValidFilePath(filePath string) bool {

	_, errStat := os.Stat(filePath)
	if errors.Is(errStat, os.ErrNotExist) {
		fmt.Printf("Specified FilePath (%v) does not exist: \n%v\n", filePath, errStat.Error())
	} else if errors.Is(errStat, os.ErrPermission) {
		fmt.Printf("Specified FilePath (%v) does not have appropriate permissions: \n%v\n", filePath, errStat.Error())
	} else if errStat != nil {
		fmt.Printf("FilePath (%v) gives errStat.Error(): \n%v\n", filePath, errStat.Error())
	} else {
		return true
	}

	return false
}

// checks for correct formatting
func isValidTimestamp(timestamp string) bool {

	timestampParts, timestampPartErr := retrieveValidTimestampParts(timestamp, 3)
	if timestampPartErr {
		return false
	}

	for i, part := range timestampParts {
		if strings.ContainsAny(part, "abcdefghijklmnopqrstuvwxyz!@#$%^&*()-=_+[]\\{}|;':\",./<>?`~") {
			fmt.Printf("part #%v contains invalid non digit character: %v\n", i+1, part)
			return false
		}

		intValue, intErr := strconv.ParseInt(part, 10, 16)
		if intErr != nil {
			fmt.Printf("part #%v gives parsing intErr.Error(): %v\n", i, intErr.Error())
			return false
		}

		switch i {
		case 0:
			if intValue < 0 {
				fmt.Printf("Hour value %v should be a positive parsed 16 bit integer: %v\n", part, intValue)
				return false
			}
		case 1:
			if intValue < 0 || intValue >= 60 {
				fmt.Printf("Minutes value %v should be a positive parsed 16 bit integer or less than 60: %v\n", part, intValue)
				return false
			}
		case 2:
			if intValue < 0 || intValue >= 60 {
				fmt.Printf("Seconds value %v should be a positive parsed 16 bit integer or less than 60: %v\n", part, intValue)
				return false
			}
		default:
			fmt.Printf("Invalid case part #%v: %v\n", i, part)
			return false
		}
	}

	return true
}

// confirmed whether timestamp has expected parts and returns them
func retrieveValidTimestampParts(timestamp string, expectedParts int) ([]string, bool) {
	timestampParts := strings.Split(timestamp, ":")
	if len(timestampParts) != expectedParts {
		fmt.Printf("found invalid number of timestamp parts (Expecting 3): %v\n", len(timestampParts))
		return nil, true
	}
	return timestampParts, false
}

// finds duration of the recording in seconds
func findDuration(recording wave.Wave) float64 {
	numSamples, sampleRate, numChannels := len(recording.Frames), recording.SampleRate, recording.NumChannels
	return roundFloat(float64(numSamples)/(float64(sampleRate)*float64(numChannels)), 3)
}

// converts Timestamp to time in Seconds
func convertValidTimestampToTime(timestamp string) float64 {

	timestampParts, _ := retrieveValidTimestampParts(timestamp, 3)
	totalTime := float64(0)

	for i, part := range timestampParts {
		intValue, _ := strconv.ParseInt(part, 10, 16)

		switch i {
		case 0:
			totalTime += float64(intValue) * 60 * 60
		case 1:
			totalTime += float64(intValue) * 60
		case 2:
			totalTime += float64(intValue)
		}

	}

	return totalTime
}

// round floating number decimal places
func roundFloat(val float64, decimalPrecision uint) float64 {
	ratio := math.Pow(10, float64(decimalPrecision))
	return math.Round(val*ratio) / ratio
}

// reads the string input by the user until when the \n character has been passed and returns the value
func getUserInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	name, err := reader.ReadString('\n')

	return strings.TrimSpace(name), err
}

func _() {
	fmt.Printf("\"Hello World\": %v\n", "Hello World")

	inputWaveFile, outputWaveFile1, outputWaveFile2 := "AudioData/TestData1.wav", "AudioData/TestOutputAudio1.wav", "AudioData/TestOutputAudio2.wav"

	sampleInputWave, sampleInputWaveErr := wave.ReadWaveFile(inputWaveFile)

	if sampleInputWaveErr != nil {
		panic(sampleInputWaveErr)
	}

	sampleInsertWaveFmt := sampleInputWave.WaveFmt

	sampleInsertFrames := createSampleFramesToInsert(sampleInsertWaveFmt)

	sampleOutputWave1Err := wave.WriteFrames(sampleInsertFrames, sampleInsertWaveFmt, outputWaveFile1)

	if sampleOutputWave1Err != nil {
		panic(sampleOutputWave1Err)
	}

	insertStart, insertEnd := 5, 15

	sampleModifiedWaveFmt := sampleInputWave.WaveFmt
	sampleModifiedFrames := insertFrame(sampleInputWave.Frames, sampleInsertFrames, time.Duration(insertStart), time.Duration(insertEnd), sampleModifiedWaveFmt.SampleRate, sampleModifiedWaveFmt.NumChannels)

	sampleOutputWave2Err := wave.WriteFrames(sampleModifiedFrames, sampleModifiedWaveFmt, outputWaveFile2)

	if sampleOutputWave2Err != nil {
		panic(sampleOutputWave2Err)
	}
}

/**************************BEGIN IMPORTANT**********************************/

func insertFrame(og, addOn []wave.Frame, insertStart, insertEnd time.Duration, sampleRate, numChannels int) []wave.Frame {
	resultFrames := []wave.Frame{}

	ogLen := len(og)

	for i := 0; i <= sampleRate*numChannels*int(insertStart) && i < ogLen; i += 1 {
		resultFrames = append(resultFrames, og[i])
	}

	resultFrames = append(resultFrames, addOn...)

	for i := sampleRate * numChannels * int(insertEnd); i < ogLen; i += 1 {
		resultFrames = append(resultFrames, og[i])
	}

	return resultFrames
}

func extractFrame(og []wave.Frame, extractStart, extractEnd time.Duration, sampleRate, numChannels int) []wave.Frame {
	resultFrames := []wave.Frame{}

	ogLen := len(og)

	for i := sampleRate * numChannels * int(extractStart); i <= sampleRate*numChannels*int(extractEnd) && i < ogLen; i += 1 {
		resultFrames = append(resultFrames, og[i])
	}

	return resultFrames
}

/**************************END IMPORTANT**********************************/

func createSampleFramesToInsert(waveFmt wave.WaveFmt) []wave.Frame {
	resultFrames := []wave.Frame{}

	start := float64(1.0)
	end := float64(1.0e-4)

	duration := 10
	frequency := 440
	numSamples := duration * waveFmt.SampleRate * waveFmt.NumChannels

	totalAngle := math.Pi * 2
	angle := totalAngle / float64(numSamples)
	decayfac := math.Pow(end/start, 1.0/float64(numSamples))

	for i := 0; i < numSamples; i++ {
		sample := math.Sin(angle * float64(frequency) * float64(i))
		sample *= start
		start *= decayfac
		resultFrames = append(resultFrames, wave.Frame(sample))
	}

	return resultFrames
}
