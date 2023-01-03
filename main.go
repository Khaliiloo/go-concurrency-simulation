package main

/*
	A small Go program that provides concurrent execution for defined processes and programs with simple instruction set
	Author: Khalil Abdulgawad
 	https://www.linkedin.com/in/kabdulgawad
*/

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type Process struct {
	ID   int
	Code []string
	Data map[string]interface{}
	PC   int
}

func NewProcess(id int, code []string) *Process {
	return &Process{
		ID:   id,
		Code: code,
		Data: map[string]interface{}{},
	}
}

type VirtualProcessor struct {
	processes []*Process
}

func NewVirtualProcessor(processes ...*Process) *VirtualProcessor {
	if processes == nil {
		return &VirtualProcessor{
			processes: make([]*Process, 0),
		}
	}
	return &VirtualProcessor{
		processes: processes,
	}
}

func (v *VirtualProcessor) Interpret(p *Process) {
	operand2 := ""
	n := 0
	for i := p.PC; i < len(p.Code); i++ {
		if n == 2 {
			return
		}
		p.PC++
		n++
		splittedInstruction := strings.Split(p.Code[i], " ")
		opCode, operand1 := splittedInstruction[0], splittedInstruction[1]
		if len(splittedInstruction) == 3 {
			operand2 = splittedInstruction[2]
		}

		switch opCode {
		case "print":
			if value, ok := p.Data[operand1]; ok {
				fmt.Printf("Process %v: %v\n", p.ID, value)
			} else {
				if operand1[0] == '\'' && operand1[len(operand1)-1] == '\'' {
					fmt.Printf("Process %v: %v\n", p.ID, operand1[1:len(operand1)-1])
				} else {
					fmt.Printf("Process %v: %v\n", p.ID, operand1)
				}
			}

		case "def":
			if _, ok := p.Data[operand1]; ok {
				fmt.Printf("%v varialbe is already defined", operand1)
			} else {
				if interfaceValue, ok := p.Data[operand2]; ok {
					p.Data[operand1] = interfaceValue
				} else {
					stringValue := operand2
					if stringValue[0] == '\'' && stringValue[len(stringValue)-1] == '\'' {
						p.Data[operand1] = stringValue[1 : len(stringValue)-1]
					} else {
						intValue, err := strconv.Atoi(stringValue)
						if err != nil {
							fmt.Printf("%v is not string nor integer", stringValue)
						} else {
							p.Data[operand1] = intValue
						}
					}

				}
			}
		case "add":
			if _, ok := p.Data[operand1]; !ok {
				fmt.Println(fmt.Errorf("%v varialbe is not defined", operand1))
			} else {
				intValue := 0
				value := p.Data[operand2]
				if value == nil {
					intValue, _ = strconv.Atoi(operand2)
				} else {
					intValue = value.(int)
				}
				p.Data[operand1] = p.Data[operand1].(int) + intValue
			}

		}
	}
}

func (v *VirtualProcessor) Run() {
	for len(v.processes) != 0 {
		currentProcess := v.processes[0]
		v.processes = v.processes[1:]
		v.Interpret(currentProcess)
		if currentProcess.PC < len(currentProcess.Code) {
			v.processes = append(v.processes, currentProcess)
		}
	}
	wg.Done()
}

func main() {
	code := []string{
		"print 'Hello-to-Go!'",
		"def num 50",
		"def num2 num",
		"add num2 10",
		"print num",
		"print num2",
		"add num num2",
		"print num",
		"print 100",
	}

	P3Code := []string{
		"print 'Process-3'",
		"def n 1",
		"add n 10",
		"print n",
	}

	P1 := NewProcess(1, code)
	P2 := NewProcess(2, code)
	VP := NewVirtualProcessor(P1, P2, NewProcess(3, P3Code))

	wg.Add(1)
	go VP.Run() // eqivalent to OS Thread
	wg.Wait()
}
