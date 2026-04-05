package main

import (
	"bufio"
	"fmt"
	"os"
)

func repl() {
	env := globalEnv
	parenSum := 0
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("yaligo => ")
		input := ""
		for scanner.Scan() {
			nextInput := scanner.Text()
			for i := range nextInput {
				if nextInput[i] == '(' {
					parenSum += 1
				}
				if nextInput[i] == ')' {
					parenSum -= 1
				}
			}
			input += nextInput
			if parenSum == 0 {
				break
			}
		}
		inputTokens := LexTokens(tokenise(input))
		inputExp, _, err := readFromTokens(inputTokens, 0)
		if err != nil {
			fmt.Printf("%e\n", err)
		}
		result, err := Eval(inputExp, env)
		if err != nil {
			fmt.Printf("%e\n", err)
		}
		fmt.Println(result.String())
	}
}

func main() {
	repl()
}
