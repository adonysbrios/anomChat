package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type ChatInfo struct {
	Host string
	Key  string
}

func ShowChatInfo(chat ChatInfo, showPassword bool) {
	fmt.Print("Host: ", chat.Host, "\r\n")

	if !showPassword {
		fmt.Print("Key: ", HidePassword(chat.Key), "\r\n")
		return
	}
	fmt.Print("Key: ", chat.Key, "\r\n")
}

func HidePassword(pwd string) string {
	return strings.Repeat("x", len(pwd))
}

func GenerateRandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.,"
	length := 40
	// buffer para la contraseña
	password := make([]byte, length)

	// llenar con caracteres aleatorios
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

func GenerateRandomUsername() string {
	Sustantive := []string{
		"Jhon", "Micro",
		"Cat", "Dog", "Fish", "Bird", "Lion", "Tiger", "Bear", "Wolf", "Fox", "Rabbit",
		"Sky", "Cloud", "Star", "Moon", "Sun", "Planet", "Comet", "Asteroid", "Galaxy",
		"Apple", "Banana", "Car", "House", "Tree", "Mountain", "River", "Ocean", "Book", "Computer",
		"Phone", "Chair", "Table", "Door", "Window", "Bridge", "Road", "City", "Village", "Forest",
		"Desert", "Island", "Castle", "Sword", "Shield", "Crown", "Ring", "Potion", "Spell", "Dragon",
		"Eagle", "Shark", "Elephant", "Giraffe", "Panda", "Koala", "Penguin", "Dolphin", "Whale", "Turtle",
		"Rose", "Lily", "Oak", "Pine", "Maple", "Gold", "Silver", "Diamond", "Ruby", "Emerald",
		"Thunder", "Storm", "Wind", "Rain", "Snow", "Fire", "Ice", "Light", "Shadow", "Dream",
	}
	Colors := []string{
		"Red", "Blue", "Green", "Yellow", "Purple", "Orange", "Pink", "Black", "White", "Gray",
		"Brown", "Cyan", "Magenta", "Lime", "Teal", "Indigo", "Violet", "Maroon", "Navy", "Olive",
		"Silver", "Gold", "Crimson", "Coral", "Turquoise", "Lavender", "Beige", "Khaki", "Salmon", "Plum",
	}
	num1, err := rand.Int(rand.Reader, big.NewInt(int64(len(Sustantive))))
	if err != nil {
		return Sustantive[0] + Colors[0]
	}
	num2, err := rand.Int(rand.Reader, big.NewInt(int64(len(Colors))))
	if err != nil {
		return Sustantive[0] + Colors[0]
	}
	return Sustantive[num1.Int64()] + Colors[num2.Int64()]
}
