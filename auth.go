package main

import (
	"unicode"
	"strings"
	"math/rand"
)

func CheckPasswordStrength(password string) int {
	result := 0

	hasDigit := false
	hasLower := false
	hasUpper := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if len(password) < 6 {
		result += 2
	}
	if !hasDigit {
		result += 4
	}
	if !hasLower {
		result += 8
	}
	if !hasUpper {
		result += 16
	}
	if !hasSpecial {
		result += 32
	}

	return result
}

func DescribePasswordStrength(strength int) string {
	issues := []string{}

	if strength&2 == 2 {
		issues = append(issues, "Password length less than 6")
	}
	if strength&4 == 4 {
		issues = append(issues, "Does not contain digits")
	}
	if strength&8 == 8 {
		issues = append(issues, "Does not contain lowercase letters")
	}
	if strength&16 == 16 {
		issues = append(issues, "Does not contain uppercase letters")
	}
	if strength&32 == 32 {
		issues = append(issues, "Does not contain special characters")
	}

	if len(issues) == 0 {
		return "Password is strong"
	}

	return strings.Join(issues, " and ")
}

func GeneratePassword() string {
	
	// Символьные группы
	digits := "0123456789"
	lowers := "abcdefghijklmnopqrstuvwxyz"
	uppers := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specials := "!@#$%^&*()-_=+[]{}|;:,.<>?/"
	
	// Объединяем все символы
	allChars := digits + lowers + uppers + specials
	
	// Гарантируем наличие хотя бы одного символа из каждой группы
	var password strings.Builder
	// rand.Seed(time.Now().UnixNano())
	length := rand.Intn(6) + 6
	password.WriteByte(digits[rand.Intn(len(digits))])
	password.WriteByte(lowers[rand.Intn(len(lowers))])
	password.WriteByte(uppers[rand.Intn(len(uppers))])
	password.WriteByte(specials[rand.Intn(len(specials))])

	// Добавляем оставшиеся символы случайным образом
	for i := 4; i < length; i++ {
		password.WriteByte(allChars[rand.Intn(len(allChars))])
	}

	// Перемешиваем пароль для случайности
	runes := []rune(password.String())
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}