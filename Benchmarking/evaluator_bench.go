package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/SebastiaanWouters/verigo/lexer"
	"github.com/SebastiaanWouters/verigo/object"
	"github.com/SebastiaanWouters/verigo/parser"
	"github.com/SebastiaanWouters/verigo/repl"
)

const (
	AMOUNT       = 20000
	ADDITION     = 0
	SUBSTRACTION = 1
	MULTIPLY     = 2
	DIVIDE       = 3
	ST           = 4
	GT           = 5
	EQ           = 6
	NEQ          = 7
	PRIME        = 8
	SIN          = 9
	TAN          = 10
	RAND         = 11
	POW          = 12
	SQRT         = 13
	LEN          = 14
	FIB          = 15
	CONCAT       = 16
)

var WEIGHTS = map[int]float32{
	ADDITION:     1.0025,
	SUBSTRACTION: 1.0,
	MULTIPLY:     1.0019,
	DIVIDE:       1.0049,
	ST:           1,
	GT:           1,
	EQ:           1,
	NEQ:          1,
	PRIME:        1.265,
	SIN:          1.2385,
	TAN:          1.2385,
	RAND:         2.2342,
	POW:          1.4593,
	SQRT:         1.2185,
	LEN:          1.2646,
	FIB:          1.7498,
	CONCAT:       1.3115,
}

type testProgram struct {
	Program   string
	Operation string
}

func main() {
	ops := flag.Int("ops", 1000000, "# of operations per run")
	runs := flag.Int("runs", 10, "# of runs")
	flag.Parse()

	f, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := f.WriteString("operation,amount,time\n"); err != nil {
		fmt.Println(err)
	}

	f.Close()

	benchmark_scripts(100)
	benchmark(*runs, *ops)
	benchmark_normal(*runs, *ops)
	benchmark_middle(*runs, *ops)
	benchmark_simple(*runs, *ops)

}

func opChanMonitor(c chan int) {
	for {
		<-c
	}
}
func rChanMonitor(c chan object.Result) {
	for {
		<-c
	}
}

func writeTofile(amount int, operation string, duration int64, fileName string) {
	f, err := os.OpenFile("results"+fileName+".csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.WriteString(operation + "," + fmt.Sprintf("%d", amount) + "," + fmt.Sprintf("%d", duration) + "\n"); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(operation + "," + fmt.Sprintf("%d", amount) + "," + fmt.Sprintf("%d", duration) + "\n")
	f.Close()
}

func benchmark_normal(runs int, ops int) {
	for r := 0; r < runs; r++ {
		var elapsed int64 = 0
		start := time.Now()
		for o := 0; o < ops; o++ {
			_ = o + o
		}
		elapsed += time.Since(start).Microseconds()
		writeTofile(ops, "addition", elapsed, "_normal")

		elapsed = 0
		start = time.Now()
		for o := 0; o < ops; o++ {
			_ = o - o
		}
		elapsed += time.Since(start).Microseconds()
		writeTofile(ops, "substraction", elapsed, "_normal")

		elapsed = 0
		start = time.Now()
		for o := 0; o < ops; o++ {
			_ = o / 30
		}
		elapsed += time.Since(start).Microseconds()
		writeTofile(ops, "division", elapsed, "_normal")

		elapsed = 0
		start = time.Now()
		for o := 0; o < ops; o++ {
			_ = o * 30
		}
		elapsed += time.Since(start).Microseconds()
		writeTofile(ops, "multiplication", elapsed, "_normal")

	}

}

func benchmark_scripts(runs int) {
	opChan := make(chan int)
	rChan := make(chan object.Result)
	go opChanMonitor(opChan)
	go rChanMonitor(rChan)

	usefulString := ReadVGFile("../useful.vg")
	l1 := lexer.New(usefulString)
	p1 := parser.New(l1)
	program1 := p1.ParseProgram()
	env1 := object.NewEnvironment()

	attackerString := ReadVGFile("../attacker.vg")
	l2 := lexer.New(attackerString)
	p2 := parser.New(l2)
	program2 := p2.ParseProgram()
	env2 := object.NewEnvironment()

	for i := 0; i <= runs; i++ {
		start1 := time.Now()
		repl.EvalParsed(program1, env1, rChan, opChan)
		elapsed1 := time.Since(start1).Microseconds()
		writeTofile(100000, "useful", elapsed1, "_scripts")
	}

	for i := 0; i <= runs; i++ {
		start2 := time.Now()
		repl.EvalParsed(program2, env2, rChan, opChan)
		elapsed2 := time.Since(start2).Microseconds()
		writeTofile(100000, "attacker", elapsed2, "_scripts")
	}

}

func benchmark(runs int, ops int) {
	opChan := make(chan int)
	rChan := make(chan object.Result)
	go opChanMonitor(opChan)
	go rChanMonitor(rChan)

	addProg := testProgram{Program: generateInfix(ops, "+", object.INTEGER_OBJ), Operation: "addition"}
	subProg := testProgram{Program: generateInfix(ops, "-", object.INTEGER_OBJ), Operation: "substraction"}
	mulProg := testProgram{Program: generateInfix(ops, "*", object.INTEGER_OBJ), Operation: "multiplication"}
	divProg := testProgram{Program: generateInfix(ops, "/", object.INTEGER_OBJ), Operation: "division"}
	stProg := testProgram{Program: generateInfix(ops, "<", object.INTEGER_OBJ), Operation: "smaller than"}
	gtProg := testProgram{Program: generateInfix(ops, ">", object.INTEGER_OBJ), Operation: "greater than"}
	eqProg := testProgram{Program: generateInfix(ops, "==", object.INTEGER_OBJ), Operation: "equals"}
	neqProg := testProgram{Program: generateInfix(ops, "!=", object.INTEGER_OBJ), Operation: "not equals"}
	primeProg := testProgram{Program: generatePrefix(ops, 1, "isPrime", []object.ObjectType{object.INTEGER_OBJ}), Operation: "isPrime"}
	sinProg := testProgram{Program: generatePrefix(ops, 1, "sin", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sin"}
	tanProg := testProgram{Program: generatePrefix(ops, 1, "tan", []object.ObjectType{object.INTEGER_OBJ}), Operation: "tan"}
	randProg := testProgram{Program: generatePrefix(ops, 1, "rand", []object.ObjectType{object.INTEGER_OBJ}), Operation: "rand"}
	fibProg := testProgram{Program: generatePrefix(ops, 1, "fib", []object.ObjectType{object.INTEGER_OBJ}), Operation: "fib"}
	sqrtProg := testProgram{Program: generatePrefix(ops, 1, "sqrt", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sqrt"}
	powProg := testProgram{Program: generatePrefix(ops, 2, "pow", []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ}), Operation: "pow"}
	lenProg := testProgram{Program: generatePrefix(ops, 1, "len", []object.ObjectType{object.STRING_OBJ}), Operation: "len"}
	concatProg := testProgram{Program: generateInfix(ops, "+", object.STRING_OBJ), Operation: "concat"}

	progList := []testProgram{addProg, subProg, mulProg, divProg, stProg, gtProg, eqProg, neqProg, primeProg, sinProg, tanProg, randProg, fibProg, powProg, lenProg, sqrtProg, concatProg}

	for j := 0; j < 2; j++ {
		for _, prog := range progList {
			l := lexer.New(prog.Program)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()

			for i := 0; i < runs; i++ {
				start := time.Now()
				repl.EvalParsed(program, env, rChan, opChan)
				elapsed := time.Since(start).Microseconds()
				writeTofile(ops, prog.Operation, elapsed, "")
			}
			fmt.Println(prog.Operation + " evaluated")

		}
	}

}

func benchmark_middle(runs int, ops int) {

	addProg := testProgram{Program: generateInfix(ops, "+", object.INTEGER_OBJ), Operation: "addition"}
	subProg := testProgram{Program: generateInfix(ops, "-", object.INTEGER_OBJ), Operation: "substraction"}
	mulProg := testProgram{Program: generateInfix(ops, "*", object.INTEGER_OBJ), Operation: "multiplication"}
	divProg := testProgram{Program: generateInfix(ops, "/", object.INTEGER_OBJ), Operation: "division"}
	stProg := testProgram{Program: generateInfix(ops, "<", object.INTEGER_OBJ), Operation: "smaller than"}
	gtProg := testProgram{Program: generateInfix(ops, ">", object.INTEGER_OBJ), Operation: "greater than"}
	eqProg := testProgram{Program: generateInfix(ops, "==", object.INTEGER_OBJ), Operation: "equals"}
	neqProg := testProgram{Program: generateInfix(ops, "!=", object.INTEGER_OBJ), Operation: "not equals"}
	primeProg := testProgram{Program: generatePrefix(ops, 1, "isPrime", []object.ObjectType{object.INTEGER_OBJ}), Operation: "isPrime"}
	sinProg := testProgram{Program: generatePrefix(ops, 1, "sin", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sin"}
	tanProg := testProgram{Program: generatePrefix(ops, 1, "tan", []object.ObjectType{object.INTEGER_OBJ}), Operation: "tan"}
	randProg := testProgram{Program: generatePrefix(ops, 1, "rand", []object.ObjectType{object.INTEGER_OBJ}), Operation: "rand"}
	fibProg := testProgram{Program: generatePrefix(ops, 1, "fib", []object.ObjectType{object.INTEGER_OBJ}), Operation: "fib"}
	sqrtProg := testProgram{Program: generatePrefix(ops, 1, "sqrt", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sqrt"}
	powProg := testProgram{Program: generatePrefix(ops, 2, "pow", []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ}), Operation: "pow"}
	lenProg := testProgram{Program: generatePrefix(ops, 1, "len", []object.ObjectType{object.STRING_OBJ}), Operation: "len"}
	concatProg := testProgram{Program: generateInfix(ops, "+", object.STRING_OBJ), Operation: "concat"}

	progList := []testProgram{addProg, subProg, mulProg, divProg, stProg, gtProg, eqProg, neqProg, primeProg, sinProg, tanProg, randProg, fibProg, powProg, lenProg, sqrtProg, concatProg}

	for j := 0; j < 2; j++ {
		for _, prog := range progList {
			l := lexer.New(prog.Program)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			var opCount *int = new(int) // Create a new integer pointer

			// Assign a value to the pointer
			*opCount = 42

			for i := 0; i < runs; i++ {
				start := time.Now()
				repl.EvalParsed_Middle(program, env, opCount)
				elapsed := time.Since(start).Microseconds()
				writeTofile(ops, prog.Operation, elapsed, "_middle")
			}
			fmt.Println(prog.Operation + " evaluated")

		}
	}

}

func benchmark_simple(runs int, ops int) {

	addProg := testProgram{Program: generateInfix(ops, "+", object.INTEGER_OBJ), Operation: "addition"}
	subProg := testProgram{Program: generateInfix(ops, "-", object.INTEGER_OBJ), Operation: "substraction"}
	mulProg := testProgram{Program: generateInfix(ops, "*", object.INTEGER_OBJ), Operation: "multiplication"}
	divProg := testProgram{Program: generateInfix(ops, "/", object.INTEGER_OBJ), Operation: "division"}
	stProg := testProgram{Program: generateInfix(ops, "<", object.INTEGER_OBJ), Operation: "smaller than"}
	gtProg := testProgram{Program: generateInfix(ops, ">", object.INTEGER_OBJ), Operation: "greater than"}
	eqProg := testProgram{Program: generateInfix(ops, "==", object.INTEGER_OBJ), Operation: "equals"}
	neqProg := testProgram{Program: generateInfix(ops, "!=", object.INTEGER_OBJ), Operation: "not equals"}
	primeProg := testProgram{Program: generatePrefix(ops, 1, "isPrime", []object.ObjectType{object.INTEGER_OBJ}), Operation: "isPrime"}
	sinProg := testProgram{Program: generatePrefix(ops, 1, "sin", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sin"}
	tanProg := testProgram{Program: generatePrefix(ops, 1, "tan", []object.ObjectType{object.INTEGER_OBJ}), Operation: "tan"}
	randProg := testProgram{Program: generatePrefix(ops, 1, "rand", []object.ObjectType{object.INTEGER_OBJ}), Operation: "rand"}
	fibProg := testProgram{Program: generatePrefix(ops, 1, "fib", []object.ObjectType{object.INTEGER_OBJ}), Operation: "fib"}
	sqrtProg := testProgram{Program: generatePrefix(ops, 1, "sqrt", []object.ObjectType{object.INTEGER_OBJ}), Operation: "sqrt"}
	powProg := testProgram{Program: generatePrefix(ops, 2, "pow", []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ}), Operation: "pow"}
	lenProg := testProgram{Program: generatePrefix(ops, 1, "len", []object.ObjectType{object.STRING_OBJ}), Operation: "len"}
	concatProg := testProgram{Program: generateInfix(ops, "+", object.STRING_OBJ), Operation: "concat"}

	progList := []testProgram{addProg, subProg, mulProg, divProg, stProg, gtProg, eqProg, neqProg, primeProg, sinProg, tanProg, randProg, fibProg, powProg, lenProg, sqrtProg, concatProg}

	for j := 0; j < 2; j++ {
		for _, prog := range progList {
			l := lexer.New(prog.Program)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()

			for i := 0; i < runs; i++ {
				start := time.Now()
				repl.EvalParsed_Simple(program, env)
				elapsed := time.Since(start).Microseconds()
				writeTofile(ops, prog.Operation, elapsed, "_simple")
			}
			fmt.Println(prog.Operation + " evaluated")

		}
	}

}

func generateRandomNumber(from int, to int) string {
	nb := strconv.Itoa(rand.Intn(to) + from)
	return nb
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func addLine(input string, line string) string {
	if input == "" {
		return line + "\n"
	}
	return input + line + "\n"
}

func generateInfix(amount int, operator string, argtype object.ObjectType) string {
	nb := amount / 100
	program := ""
	if argtype == object.INTEGER_OBJ {
		for i := 0; i < nb; i++ {
			num1 := generateRandomNumber(100, 2000)
			num2 := generateRandomNumber(100, 2000)
			newLine := num1 + operator + num2
			program = addLine(program, newLine)
		}
	} else if argtype == object.STRING_OBJ {
		for i := 0; i < nb; i++ {
			str1 := generateRandomString(20)
			str2 := generateRandomString(20)
			newLine := "\"" + str1 + "\"" + operator + "\"" + str2 + "\""
			program = addLine(program, newLine)
		}
	}
	counter := 0
	initProgram := program
	for counter < 100 {
		program = program + initProgram
		counter += 1
	}

	return program
}

func generatePrefix(amount int, args int, operator string, argtypes []object.ObjectType) string {
	nb := amount / 100
	program := ""
	for i := 0; i < nb; i++ {
		newLine := generatePrefixOp(operator, args, argtypes)
		program = addLine(program, newLine)
	}
	counter := 0
	initProgram := program
	for counter < 100 {
		program = program + initProgram
		counter += 1
	}

	return program
}

func generatePrefixOp(operator string, args int, argtypes []object.ObjectType) string {
	op := operator + "("
	for i := 0; i < args; i++ {
		if argtypes[i] == object.INTEGER_OBJ {
			if i > 0 {
				//not first argument
				if operator == "pow" {
					op = op + "," + generateRandomNumber(1, 7)
				} else {
					op = op + "," + generateRandomNumber(100, 1000)
				}
			} else {
				//first argument
				if operator == "pow" {
					op = op + generateRandomNumber(1, 20)
				} else {
					op = op + generateRandomNumber(100, 1000)
				}
			}
		} else if argtypes[i] == object.STRING_OBJ {
			if i > 0 {
				op = op + "," + "\"" + generateRandomString(rand.Intn(20)) + "\""
			} else {
				op = op + "\"" + generateRandomString(rand.Intn(20)) + "\""
			}
		}
	}
	op += ")"
	return op
}

func ReadVGFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	return string(data)
}
