// the REPL reads input, sends it to the interpreter for evaluation, prints the result/output of the
// interpreter and starts again. Read, Eval, Print, Loop.

package repl

import (
	"bufio"
	"fmt"
	"io"
	"mooshy/lexer"
	"mooshy/token"
)

const PROMT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan() // read input from the user
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for {
			tok := l.NextToken()
			if tok.Type == token.EOF {
				break // end of line
			}
			fmt.Printf("Type: %s, Literal: %s\n", tok.Type, tok.Literal)
		}
	}
}
