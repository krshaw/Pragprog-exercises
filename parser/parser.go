package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type CommandValue struct {
	ARG     bool
	handler func([]int)
}

type ParsingTable map[byte]CommandValue

func (table ParsingTable) initialize() {
	table['P'] = CommandValue{ARG, selectPen}
}

const ARG = true
const NOARG = !ARG

func executeInstrs(scanner *bufio.Scanner, table ParsingTable) {
	for scanner.Scan() {
		instr := scanner.Bytes()
		execute(instr, table)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

//TODO: handle error where the command is not in the table!
func execute(instr []byte, table ParsingTable) {
	// first byte is the instr name, following bytes are the arguments (if any)
	fmt.Println(string(instr))
	instrName := instr[0]
	commandValue := table[instrName]
	args := make([]int, 5) // 5 should be sufficient, realistically wont be more than
	if commandValue.ARG {
		splitArgs := bytes.Split(instr[2:], []byte(" "))
		rawArgs := bytes.Join(splitArgs, []byte(""))
		fmt.Println(string(rawArgs))
	}
	commandValue.handler(args)
}

func main() {
	// next step: create a test file with this command, read in the file,
	// seperate by line, match against the first byte of the line with the table
	table := make(ParsingTable)
	table.initialize()
	file, err := os.Open("./testfile")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	executeInstrs(scanner, table)
}

func selectPen(penType []int) {
	if len(penType) != 1 {
		fmt.Printf("selectPen requires 1 argument, not %v\n", len(penType))
		return
	}
	fmt.Printf("Selecting Pen #%v\n", penType[0])
}
