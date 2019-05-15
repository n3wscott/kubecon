package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

func Calculate(exp string) (string, bool) {
	x := strings.Split(strings.TrimSpace(exp), " ")
	var result []string

	done := true

	for i, t := range x {
		if t == "=" {
			if i > 0 {
				return x[i-1], done
			}
			return "0", done
		}

		if i < 2 {
			continue
		}

		a, err := strconv.Atoi(x[i-2])
		if err != nil {
			return err.Error(), done
		}
		b, err := strconv.Atoi(x[i-1])
		if err != nil {
			return err.Error(), done
		}
		c := 0

		defer func() {
			// Catch divide by zero, etc.
			if r := recover(); r != nil {
				fmt.Printf("Panic! %s", r)
			}
		}()

		switch t {
		case "+":
			c = a + b
		case "−", "-":
			c = a - b
		case "×", "*":
			c = a * b
		case "÷", "/":
			c = a / b
		default:
			continue
		}
		done = false // keep calculating.

		result = append(x[:i-2], strconv.Itoa(c))
		result = append(result, x[i+1:]...)
		break
	}
	return strings.Join(result, " "), done
}
