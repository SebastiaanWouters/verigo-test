package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
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
	ops := flag.Int("ops", 100000, "# of operations per run")
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

	benchmarkNormal(*runs, *ops)
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

func benchmarkNormal(runs int, ops int) {
	for r := 0; r < runs; r++ {
		var elapsed int64 = 0
		start := time.Now()
		for o := 0; o < ops; o++ {
			_ = o + o
		}
		elapsed += time.Since(start).Microseconds()
		writeTofile(ops, "addition", elapsed, "_normal")
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
				writeTofile(ops, prog.Operation, elapsed, strconv.Itoa(j))
			}
			fmt.Println(prog.Operation + " evaluated")

		}
	}

}

// Benchmarks
func BenchmarkFile(b *testing.B) {
	f, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := f.WriteString("operation,amount,time\n"); err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	counter := 0
	opChan := make(chan int)
	rChan := make(chan object.Result)
	go opChanMonitor(opChan)
	go rChanMonitor(rChan)

	addStr := generateInfix(AMOUNT, "+", object.INTEGER_OBJ)
	subStr := generateInfix(AMOUNT, "-", object.INTEGER_OBJ)
	mulStr := generateInfix(AMOUNT, "*", object.INTEGER_OBJ)
	divStr := generateInfix(AMOUNT, "/", object.INTEGER_OBJ)
	primeStr := generatePrefix(AMOUNT, 1, "isPrime", []object.ObjectType{object.INTEGER_OBJ})
	sinStr := generatePrefix(AMOUNT, 1, "sin", []object.ObjectType{object.INTEGER_OBJ})
	tanStr := generatePrefix(AMOUNT, 1, "tan", []object.ObjectType{object.INTEGER_OBJ})
	randStr := generatePrefix(AMOUNT, 1, "rand", []object.ObjectType{object.INTEGER_OBJ})
	fibStr := generatePrefix(AMOUNT, 1, "fib", []object.ObjectType{object.INTEGER_OBJ})
	powStr := generatePrefix(AMOUNT, 2, "pow", []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ})
	lenStr := generatePrefix(AMOUNT, 1, "len", []object.ObjectType{object.STRING_OBJ})
	sqrtStr := generatePrefix(AMOUNT, 1, "sqrt", []object.ObjectType{object.INTEGER_OBJ})
	concatStr := generateInfix(AMOUNT, "+", object.STRING_OBJ)

	b.Run("addition", func(b *testing.B) {
		counter += 1
		l := lexer.New(addStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()

		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}

		if counter%3 == 0 {
			f, err := os.OpenFile("results.csv",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			if _, err := f.WriteString("addition," + fmt.Sprintf("%f", float32(b.N*AMOUNT)) + "," + fmt.Sprintf("%f", float32(b.Elapsed().Microseconds())) + "\n"); err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[ADDITION], "work units or", float32(b.N*AMOUNT)*WEIGHTS[ADDITION]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("substraction", func(b *testing.B) {
		l := lexer.New(subStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[SUBSTRACTION], "work units or", float32(b.N*AMOUNT)*WEIGHTS[SUBSTRACTION]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("multiply", func(b *testing.B) {
		l := lexer.New(mulStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[MULTIPLY], "work units or", float32(b.N*AMOUNT)*WEIGHTS[MULTIPLY]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("divide", func(b *testing.B) {
		l := lexer.New(divStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[DIVIDE], "work units or", float32(b.N*AMOUNT)*WEIGHTS[DIVIDE]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("isPrime", func(b *testing.B) {
		l := lexer.New(primeStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[PRIME], "work units or", float32(b.N*AMOUNT)*WEIGHTS[PRIME]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("sin", func(b *testing.B) {
		l := lexer.New(sinStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[SIN], "work units or", float32(b.N*AMOUNT)*WEIGHTS[SIN]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("tan", func(b *testing.B) {
		l := lexer.New(tanStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[TAN], "work units or", float32(b.N*AMOUNT)*WEIGHTS[TAN]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("rand", func(b *testing.B) {
		l := lexer.New(randStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[RAND], "work units or", float32(b.N*AMOUNT)*WEIGHTS[RAND]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("fib", func(b *testing.B) {
		l := lexer.New(fibStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[FIB], "work units or", float32(b.N*AMOUNT)*WEIGHTS[FIB]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("pow", func(b *testing.B) {
		l := lexer.New(powStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[POW], "work units or", float32(b.N*AMOUNT)*WEIGHTS[POW]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("len", func(b *testing.B) {
		l := lexer.New(lenStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[LEN], "work units or", float32(b.N*AMOUNT)*WEIGHTS[LEN]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("sqrt", func(b *testing.B) {
		l := lexer.New(sqrtStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[SQRT], "work units or", float32(b.N*AMOUNT)*WEIGHTS[SQRT]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})
	b.Run("concat", func(b *testing.B) {
		l := lexer.New(concatStr)
		p := parser.New(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()

		for i := 0; i < b.N; i++ {
			repl.EvalParsed(program, env, rChan, opChan)
		}
		fmt.Println()
		fmt.Println("SIMPLE: Received", float32(b.N*AMOUNT), "work units, or", float32(b.N*AMOUNT)/float32(b.Elapsed().Microseconds()), "work units per microsecond")
		fmt.Println("ADVANCED: Received", float32(b.N)*WEIGHTS[CONCAT], "work units or", float32(b.N*AMOUNT)*WEIGHTS[CONCAT]/float32(b.Elapsed().Microseconds()), "work units per microsecond")
	})

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
