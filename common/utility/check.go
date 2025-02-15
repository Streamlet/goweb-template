package utility

import "regexp"

func IsUrl(url string) bool {
	return regexp.MustCompile(`^https?://[\w-]+(\.[\w-]+)+`).MatchString(url)
}

func IsHost(url string) bool {
	return regexp.MustCompile(`^[\w-]+(\.[\w-]+)+`).MatchString(url)
}

func IsEmail(email string) bool {
	return regexp.MustCompile(`^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$`).MatchString(email)
}

func IsMobile(mobile string) bool {
	return regexp.MustCompile(`^(\+86)?1\d{10}$`).MatchString(mobile)
}

func IsNumeric(s string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(s)
}

func IsIdentifier(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString(s)
}
