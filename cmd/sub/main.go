package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/codegangsta/cli"
)

/*
	Eg:
	echo "foobarbaz" | sub '^foo(bar|biz)baz$' 'bar or biz: $1'
*/
func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "sub"
	app.Usage = "A command-line tool for substituting patterns from a stream."
	app.UsageText = "sub [pattern] [replacement]"
	app.Author = "Kevin Cantwell"
	app.Email = "kevin.cantwell@gmail.com"
	app.Action = func(c *cli.Context) error {
		pattern := c.Args().Get(0)
		replacement := []byte(c.Args().Get(1))

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			matches := regex.FindAllSubmatch(scanner.Bytes(), -1)
			for _, submatches := range matches {
				replaced, err := replace(replacement, submatches)
				if err != nil {
					return err
				}
				os.Stdout.Write(replaced)
				os.Stdout.Write([]byte("\n"))
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		exit(err.Error(), 1)
	}
}

func replace(replacement []byte, submatches [][]byte) ([]byte, error) {
	var replaced []byte
	for i := 0; i < len(replacement); i++ {
		char := replacement[i]
		switch char {
		case '\\':
			// If this is the last char, then just append it
			if i+1 >= len(replacement) {
				replaced = append(replaced, char)
				continue
			}
			nextChar := replacement[i+1]
			switch nextChar {
			// If this is an escape char, just append the escaped char and move the cursor
			case '$':
				replaced = append(replaced, nextChar)
				i++
			default:
				replaced = append(replaced, char)
			}
		case '$':
			if i+1 >= len(replacement) {
				return nil, errors.New("sub: no submatch index provided")
			}
			nextChar := replacement[i+1]
			switch nextChar {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				idx := int(nextChar - 48)
				if idx >= len(submatches) {
					return nil, errors.New("sub: no submatch for index: " + fmt.Sprint(idx))
				}
				replaced = append(replaced, submatches[idx]...) // '0' == byte(48)
				i++
			default:
				return nil, errors.New("sub: invalid index for submatch: " + string(nextChar))
			}
		default:
			replaced = append(replaced, char)
		}
	}
	return replaced, nil
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}
