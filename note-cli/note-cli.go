package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getNoteData() ([]byte, error) {
	// Read text from the input file
	noteData, errNoteData := os.ReadFile("note.txt")
	if errNoteData != nil {
		if os.IsNotExist(errNoteData) {
			emptyFile, err := os.Create("note.txt")
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("error creating new file")
			}
			emptyFile.Close()
			noteData = []byte{}
		} else {
			fmt.Println(errNoteData)
			return nil, fmt.Errorf("error reading input file")
		}
	}
	return noteData, nil
}

func printTitle() {
	fmt.Print(`
		
███████╗██╗███╗   ███╗██████╗ ██╗     ███████╗    ███╗   ██╗ ██████╗ ████████╗███████╗     ██████╗  ██████╗ 
██╔════╝██║████╗ ████║██╔══██╗██║     ██╔════╝    ████╗  ██║██╔═══██╗╚══██╔══╝██╔════╝    ██╔════╝ ██╔═══██╗
███████╗██║██╔████╔██║██████╔╝██║     █████╗      ██╔██╗ ██║██║   ██║   ██║   █████╗      ██║  ███╗██║   ██║
╚════██║██║██║╚██╔╝██║██╔═══╝ ██║     ██╔══╝      ██║╚██╗██║██║   ██║   ██║   ██╔══╝      ██║   ██║██║   ██║
███████║██║██║ ╚═╝ ██║██║     ███████╗███████╗    ██║ ╚████║╚██████╔╝   ██║   ███████╗    ╚██████╔╝╚██████╔╝
╚══════╝╚═╝╚═╝     ╚═╝╚═╝     ╚══════╝╚══════╝    ╚═╝  ╚═══╝ ╚═════╝    ╚═╝   ╚══════╝     ╚═════╝  ╚═════╝ 
                                                                                                            

`)
}

func printNote(noteLines []string) {
	printTitle()
	// format and print the note by appending index for every line
	for i, line := range noteLines {
		fmt.Printf("%d %s\n", i+1, line)
	}
}

// Get user input for program option
func getUserInputForProgramState(reader *bufio.Reader) (int, error) {
	fmt.Print("\n\n" +
		"What do you want to do?\n" +
		"1) Update a line\n" +
		"2) Delete a line\n" +
		"0) Exit\n")

	programState := 10

	invalidUserOptionPick := false

	for programState != 1 && programState != 2 && programState != 0 {
		if invalidUserOptionPick {
			fmt.Print("INVALID INPUT!\n")
		}

		fmt.Print("Enter your option [1/2/0]: ")

		userOptionPick, err := reader.ReadString('\n')
		if err != nil {
			return -1, err
		}
		userOptionPick = strings.TrimSpace(userOptionPick)

		userOptionPickInt, err := strconv.Atoi(userOptionPick)
		if err != nil {
			return -1, err
		}

		programState = userOptionPickInt
		invalidUserOptionPick = true
	}

	return programState, nil
}

func updateLine(reader *bufio.Reader, noteLines []string) {
	// Change / insert a line with user input
	fmt.Print("Which line to be updated: ")
	insertLineNum, errInsertLineNum := reader.ReadString('\n')
	if errInsertLineNum != nil {
		fmt.Println(errInsertLineNum)
		return
	}
	insertLineNum = strings.TrimSpace(insertLineNum)
	insertLineNumInt, errInsertLineNumInt := strconv.Atoi(insertLineNum)
	if errInsertLineNumInt != nil {
		fmt.Println(errInsertLineNumInt)
		return
	}

	fmt.Print("Type the new text: ")
	insertLine, errInsertLine := reader.ReadString('\n')
	if errInsertLine != nil {
		fmt.Println(errInsertLine)
		return
	}

	noteLines[insertLineNumInt-1] = insertLine
	var updatedNoteStrings string
	for _, line := range noteLines {
		updatedNoteStrings += fmt.Sprintf("%s\n", line)
	}

	// Write it to the file
	errWriteNoteData := os.WriteFile("note.txt", []byte(updatedNoteStrings), 0644)

	if errWriteNoteData != nil {
		fmt.Println(errWriteNoteData)
		return
	}
}

func deleteLine(reader *bufio.Reader, noteLines []string) {
	// Change / insert a line with user input
	fmt.Print("Which line to be deleted: ")
	deleteLineNum, errDeleteLineNum := reader.ReadString('\n')
	if errDeleteLineNum != nil {
		fmt.Println(errDeleteLineNum)
		return
	}
	deleteLineNum = strings.TrimSpace(deleteLineNum)
	deleteLineNumInt, errDeleteLineNumInt := strconv.Atoi(deleteLineNum)
	if errDeleteLineNumInt != nil {
		fmt.Println(errDeleteLineNumInt)
		return
	}
	noteLines = append(noteLines[:deleteLineNumInt-1], noteLines[deleteLineNumInt:]...)

	// Update the new text file
	var updatedNoteStrings string
	for _, line := range noteLines {
		updatedNoteStrings += fmt.Sprintf("%s\n", line)
	}

	// Write it to the file
	errWriteNoteData := os.WriteFile("note.txt", []byte(updatedNoteStrings), 0644)

	if errWriteNoteData != nil {
		fmt.Println(errWriteNoteData)
		return
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Application loop
	programState := 1

	for programState != 0 {
		clearTerminal()

		// get note data
		noteData, err := getNoteData()
		if err != nil {
			fmt.Println(err)
			return
		}
		noteLines := strings.Split(string(noteData), "\n")

		// show current note with title
		printNote(noteLines)

		// update program state from user input
		programState, err = getUserInputForProgramState(reader)
		if err != nil {
			fmt.Println(err)
			return
		}

		if programState == 1 {
			updateLine(reader, noteLines)
		} else if programState == 2 {
			deleteLine(reader, noteLines)
		} else if programState == 0 {
			fmt.Print("Program exited")
		}

		fmt.Println()
	}
}
