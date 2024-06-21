# WaveWeaver

WaveWeaver is a project where users can interact with the CLI to edit their audio files. The repository is maintained in Go.

## User Actions

Users can:
- Remove portions of an audio file
- Replace portions of an audio file with a new file
- Replace portions of an audio file with a snippet of another file.

In order to perform the actions, users will be prompted to:
- Provide file paths to existing audio files users wish to use.
- Provide file paths to the output files users wish to procure.
- Provide start and stop timestamps for any snipping where necessary in `HH:MM:SS` format.

## Motivation
The project was originally meant to make life easier for people who perform reading and recording as services for others who are disadvantaged. This way, the reader can continue recording despite any disturbances or mistakes. They can also now afford to make mistakes when recording. These problems can be fixed now in the review phase of the service.

#### A next stage would be to either
1) Convert the program into a more accessible form, instead of a CLI program you have to manually execute. Currently considering using text/chat based models to get repeated prompts and give easier workflow for the user. 
2) Integrating a text to speech element to the project, where users can provide few samples of their recordings, have a model learn from it and be able to then read from images provided to the same. Users may then use the existing tools to simply correct the mistakes of the model in review phase.

Open to suggestions for both at akheelsaajid@gmail.com
