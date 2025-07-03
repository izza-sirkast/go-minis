package main

import (
	"bufio"
	"encoding/csv"
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

func main() {
	userReader := bufio.NewReader(os.Stdin)
	programState := 10

	for programState != 0 {
		clearTerminal()

		fmt.Println("todo app")

		// get todo data
		file, err := os.Open("todo.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		todosData, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, record := range todosData {
			fmt.Println(record)
		}

		fmt.Print("\n\n\n" +
			"1) Add new todo\n" +
			"2) Delete a todo\n" +
			"3) Toggle completion status\n" +
			"4) Delete all completed todos\n" +
			"0) Exit program\n" +
			"What do you want to do [1/2/3/4/0]: ")

		userOptionPick, err := userReader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		userOptionPickInt, err := strconv.Atoi(strings.TrimSpace(userOptionPick))

		if err != nil {
			fmt.Println(err)
			return
		}

		programState = userOptionPickInt

		switch programState {
		case 1:
			// get user input of new todo description
			fmt.Print("Write the new todo: ")
			newTodoDesc, err := userReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			newTodoDesc = strings.TrimSpace(newTodoDesc)

			// == preparing new todo data structure
			lastTodoId, err := strconv.Atoi(todosData[len(todosData)-1][0]) // get last todo id
			newTodo := []string{
				strconv.Itoa(lastTodoId + 1),
				newTodoDesc,
				strconv.Itoa(0),
			}
			todosData = append(todosData, newTodo)

			// write new todos to todo.csv
			file, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(todosData); err != nil {
				fmt.Println(err)
				return
			}
		case 2:
			fmt.Println("remove a todo")

		case 3:
			fmt.Println("toggle todo completion")

		case 4:
			fmt.Println("remove all completed todos")

		case 0:
			fmt.Println("good bye...")

		}
	}

	// file, err := os.Create("todo.csv")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// data := [][]string{
	// 	{"id", "description", "status"},
	// 	{"1", "Learn go", "0"},
	// 	{"2", "Working on projects", "1"},
	// }

	// if err := writer.WriteAll(data); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
