package jwt

// заглушка для создания jwt токена
func GenerateToken(username string) string {
	return username
}

// заглушка для проверки jwt токена
func CheckToken(token string) (string, bool) {
	return "my_login75", true
}
