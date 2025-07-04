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

var header []string = []string{
	"id", "description", "completed",
}

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

func showTodos(todos [][]string) {
	var longestId, longestDesc int = 0, 0
	for _, todo := range todos {
		if len(todo[0]) > longestId {
			longestId = len(todo[0])
		}
		if len(todo[1]) > longestDesc {
			longestDesc = len(todo[1])
		}
	}

	for i, todo := range todos {
		padId := strings.Repeat(" ", longestId-len(todo[0]))
		padDesc := strings.Repeat(" ", longestDesc-len(todo[1]))

		if i == 0 {
			fmt.Printf("| %s%s | %s%s | %s |\n", todo[0], padId, todo[1], padDesc, todo[2])

			sepId := strings.Repeat("-", longestId+2)
			sepDesc := strings.Repeat("-", longestDesc+2)
			sepComp := strings.Repeat("-", 11)
			fmt.Printf("|%s|%s|%s|\n", sepId, sepDesc, sepComp)
		} else {
			var mark string
			if todo[2] == "0" {
				mark = "❎"
			} else if todo[2] == "1" {
				mark = "✅"
			}

			fmt.Printf("| %s%s | %s%s |    %s     |\n", todo[0], padId, todo[1], padDesc, mark)
		}

	}
}

func printTitle() {
	fmt.Print(`
 ___       ___   __    __     ____     ________     ____     ______       ____    
(  (       )  ) (  \  /  )   (    )   (___  ___)   / __ \   (_  __ \     / __ \   
 \  \  _  /  /   \ (__) /    / /\ \       ) )     / /  \ \    ) ) \ \   / /  \ \  
  \  \/ \/  /     ) __ (    ( (__) )     ( (     ( ()  () )  ( (   ) ) ( ()  () ) 
   )   _   (     ( (  ) )    )    (       ) )    ( ()  () )   ) )  ) ) ( ()  () ) 
   \  ( )  /      ) )( (    /  /\  \     ( (      \ \__/ /   / /__/ /   \ \__/ /  
    \_/ \_/      /_/  \_\  /__(  )__\    /__\      \____/   (______/     \____/   
                                                                                  

`)
}

func addHeaderToTodos(todosData [][]string) [][]string {
	newTodos := append([][]string{header}, todosData...)
	return newTodos
}

func main() {
	userReader := bufio.NewReader(os.Stdin)
	programState := 10

	for programState != 0 {
		clearTerminal()

		printTitle()

		// get todo data
		file, err := os.Open("todo.csv")
		if os.IsNotExist(err) {
			_, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println("err")
				return
			}

			file, err = os.Open("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
		} else if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		todosDataAll, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if csv has header
		if len(todosDataAll) <= 0 || todosDataAll[0][0] != "id" {
			file, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			todosDataAll = addHeaderToTodos(todosDataAll)
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(todosDataAll); err != nil {
				fmt.Println(err)
				return
			}
		}

		showTodos(todosDataAll)

		var todosData [][]string
		if len(todosDataAll) > 1 {
			todosData = todosDataAll[1:]
		} else {
			todosData = [][]string{}
		}

		// get user input for next program action
		fmt.Print("\n\n\n" +
			"1) Add new todo\n" +
			"2) Delete a todo\n" +
			"3) Toggle completion status\n" +
			"4) Delete all completed todos\n" +
			"0) Exit program\n")

		var userOptionPickInt int
		for {
			fmt.Print("What do you want to do [1/2/3/4/0]: ")
			userOptionPick, err := userReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			userOptionPickInt, err = strconv.Atoi(strings.TrimSpace(userOptionPick))
			if err == nil {
				break
			}

			fmt.Println("Invalid input")
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

			// preparing new todo data structure
			lastTodoId := 0
			if len(todosData) > 0 { // if todos data not empty
				lastTodoId, err = strconv.Atoi(todosData[len(todosData)-1][0]) // get last todo id
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			newTodoId := strconv.Itoa(lastTodoId + 1)
			newTodo := []string{
				newTodoId,
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

			todosData = addHeaderToTodos(todosData)
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(todosData); err != nil {
				fmt.Println(err)
				return
			}
		case 2:
			var todoIdRemove int
			for {
				fmt.Print("Enter id of todo you want to remove: ")
				todoIdRemoveString, err := userReader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					return
				}
				todoIdRemoveString = strings.TrimSpace(todoIdRemoveString)

				todoIdRemove, err = strconv.Atoi(todoIdRemoveString)
				if err == nil {
					break
				}

				fmt.Println("Invalid input, please enter a number")
			}

			// update todos data (delete the todo the user is selected)
			var newTodosData [][]string
			newId := 1
			for _, todo := range todosData {
				todoId, err := strconv.Atoi(todo[0])
				if err != nil {
					fmt.Println(err)
					return
				}
				if todoId != todoIdRemove {
					todo[0] = strconv.Itoa(newId)
					newTodosData = append(newTodosData, todo)
					newId++
				}
			}

			// write new todos data to todo.csv
			file, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			newTodosData = addHeaderToTodos(newTodosData)
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(newTodosData); err != nil {
				fmt.Println(err)
				return
			}
		case 3: // toggle todo completion status
			var todoIdToggle int
			for {
				fmt.Print("Enter id of todo you want to toggle it's completion: ")
				todoIdToggleString, err := userReader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					return
				}
				todoIdToggleString = strings.TrimSpace(todoIdToggleString)

				todoIdToggle, err = strconv.Atoi(todoIdToggleString)
				if err == nil {
					break
				}

				fmt.Println("Invalid input, please enter a number")
			}

			// update todos data (toggle completion status of the todo the user is selected)
			for _, todo := range todosData {
				todoId, err := strconv.Atoi(todo[0])
				if err != nil {
					fmt.Println(err)
					return
				}
				if todoId == todoIdToggle {
					if todo[2] == strconv.Itoa(0) {
						todo[2] = strconv.Itoa(1)
					} else {
						todo[2] = strconv.Itoa(0)
					}
				}
			}

			// write new todos data to todo.csv
			file, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			todosData = addHeaderToTodos(todosData)
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(todosData); err != nil {
				fmt.Println(err)
				return
			}
		case 4: // remove all completed todos
			var newTodosData [][]string
			newTodoId := 1
			for _, todo := range todosData {
				if todo[2] == strconv.Itoa(0) {
					todo[0] = strconv.Itoa(newTodoId)
					newTodosData = append(newTodosData, todo)
					newTodoId++
				}
			}

			// write new todos data to todo.csv
			file, err := os.Create("todo.csv")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			newTodosData = addHeaderToTodos(newTodosData)
			writer := csv.NewWriter(file)
			defer writer.Flush()
			if err := writer.WriteAll(newTodosData); err != nil {
				fmt.Println(err)
				return
			}
		case 0:
			fmt.Println("good bye...")

		}
	}
}
