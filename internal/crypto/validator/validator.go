package validator

func LuhnAlgorithm(order string) bool {
	total := 0
	parity := len(order) % 2
	for i := 0; i < len(order); i++ {
		digit := int(order[i] - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
	}
	return total%10 == 0
}
