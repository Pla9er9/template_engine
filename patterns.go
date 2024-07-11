package main

// correct if statment pattern should look like this
// [@if, isAdmin]
func IsIfStatmentPattern(claims []string) bool {
	if len(claims) != 2 {
		return false
	}

	if claims[0] != "@if" {
		return false
	}

	return true
}

// correct foreach statment pattern should look like this
// [@foreach, numbers, as, num]
func IsForeachStatmentPattern(claims []string) bool {
	if len(claims) != 4 {
		return false
	}

	if claims[0] != "@foreach" || claims[2] != "as" {
		return false
	}

	return true
}
