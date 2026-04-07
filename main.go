package main

import (
	"bufio"
	"fmt"
	"os"
	"yaligo/internal"
)

func repl() {
	env := internal.StandardEnv()
	parenSum := 0
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("yaligo => ")
		input := ""
		for scanner.Scan() {
			nextInput := scanner.Text() + "\n"
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
		if input == ".exit" {
			return
		}
		inputTokens := internal.LexTokens(internal.Tokenise(input))
		inputExp, _, err := internal.ReadFromTokens(inputTokens, 0)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		result, err := internal.Eval(inputExp, env)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		fmt.Println(result.String())
	}
}

func main() {
	repl()
}
