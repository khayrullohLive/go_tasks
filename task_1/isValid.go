package main

func isValid(s string) bool {
	var stack []rune

	for _, ch := range s {
		if ch == '{' || ch == '[' || ch == '(' {
			stack = append(stack, ch)
		} else {
			if len(stack) == 0 {
				return false
			}

			last := stack[len(stack)-1]

			if ch == ')' && last != '(' {
				return false
			}
			if ch == '}' && last != '{' {
				return false
			}
			if ch == ']' && last != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
