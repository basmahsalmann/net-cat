package handlers

import (
	"bufio"
	"os"
	"log"
)

func WelcomeMessage() []string {
    //Create empty array of strings to store linux logo
	fileText := make([]string, 0)//
	file, err := os.Open("../handlers/TextFiles/logo.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	//Reads from file
	scanner := bufio.NewScanner(file)

	//Return true if there is a line to read
	for scanner.Scan() {
		fileText = append(fileText, scanner.Text() + "\n") // appends line to fileText
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erorr reading file: %s", err)
	}

	return fileText
}



