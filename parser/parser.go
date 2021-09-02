package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CommandValue struct {
	ARG     bool
	handler func([]int)
}

const ARG = true
const NOARG = !ARG

type ParsingTable map[byte]CommandValue

func (table ParsingTable) initialize() {
	table['P'] = CommandValue{ARG, selectPen}
	table['D'] = CommandValue{NOARG, penDown}
	table['W'] = CommandValue{ARG, drawWest}
}

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
	if len(instr) < 1 {
		return // white space / blank lines shouldnt break anything
	}
	// first byte is the instr name, following bytes are the arguments (if any)
	instrName := instr[0]
	commandValue, ok := table[instrName]
	if !ok {
		log.Fatal(fmt.Sprintf("%b is not a valid instruction\n", instrName))
	}
	var args []int
	if commandValue.ARG {
		args = convertArgsToInt(instr)
	}
	commandValue.handler(args)
}

func convertArgsToInt(instr []byte) []int {
	if len(instr) < 3 { // less than 3 means there is no argument, even though there should be
		log.Fatal(fmt.Sprintf("Instruction %c requires argument(s)\n", instr[0]))
	}
	splitArgs := bytes.Split(instr[2:], []byte(" "))
	var convertedArgs []int
	for _, arg := range splitArgs {
		convertedArg, err := strconv.Atoi(string(arg))
		if err != nil {
			log.Fatal(err)
		}
		convertedArgs = append(convertedArgs, convertedArg)
	}
	return convertedArgs
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
	fmt.Printf("Selecting Pen #%v\n", penType[0])
}

func penDown(_ []int) {
	fmt.Printf("Putting the pen down\n")
}

func drawWest(distance []int) {
	fmt.Printf("Drawing west %v cm\n", distance[0])
}
