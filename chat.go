package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// STRUCTURES
type Sections struct {
	Header  string `jsonl:"Header"`
	Content string `jsonl:"Content"`
}

type Article struct {
	Title    string     `jsonl:"Title"`
	Sections []Sections `jsonl:"Sections"`
}

// MAIN FUNCTION
func main() {
	// Open the data file
	inputFile, err := os.Open("output.jsonl")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// Create a scanner to read lines from the input file
	scanner := bufio.NewScanner(inputFile)

	// Create a map to store article sections by header
	sectionsMap := make(map[string]string)

	//Container to hold a list of all of article 1 headers
	roboticsHeaders := []string{}

	//boolean to check for article 1
	isRoboticsArticle := false

	// Scan each line
	for scanner.Scan() {
		line := scanner.Text()

		// Decode JSON and turn content into the article struct
		var article Article
		if err := json.Unmarshal([]byte(line), &article); err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		// Clean the text - remove everything within brackets '[]'
		for i := range article.Sections {
			article.Sections[i].Header = strings.ReplaceAll(article.Sections[i].Header, "[edit]", "")
			article.Sections[i].Content = removeBetween(article.Sections[i].Content, "[", "]")
		}

		//Turn content/header into the Sections struct
		if isRoboticsArticle {
			for _, section := range article.Sections {
				sectionsMap[section.Header] = section.Content
				roboticsHeaders = append(roboticsHeaders, section.Header)
			}
			break
		}

		//Check if this is the "Robotics" article (Article 1)
		if article.Title == "Robotics" {
			isRoboticsArticle = true
		}
	}

	//Handle errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// Now for the "chat" loop
	for {
		//initial question
		fmt.Print("What would you like to know more about? ")

		//variable for user input
		var userInput string

		//create scanner to scan user input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput = scanner.Text()

		//option to stop the chatbot
		if userInput == "exit" {
			break

			//option to list out all headers and see potential options for user input
		} else if userInput == "list" {
			// Print a list of all headers from the "Robotics" article
			fmt.Println("List of headers from the 'Robotics' article:")
			for _, header := range roboticsHeaders {
				fmt.Println(header)
			}

			//If user input matches a header, the content mapped to that header will be produced
		} else {
			// Lookup and print the matching content
			content, exists := sectionsMap[userInput]
			if exists {
				fmt.Printf("%s: %s\n", userInput, content)

				//response if user input does not match any header
			} else {
				fmt.Println("I'm sorry, I don't have any information on that topic.")
			}
		}
	}
}

// CLEAN TEXT - remove brackets '[]' and all text in between
func removeBetween(str, start, end string) string {
	anyIncludingEndLine := fmt.Sprintf(`%s[\r\n\s\w]*%s`, start, end)
	return regexp.MustCompile(anyIncludingEndLine).ReplaceAllString(str, "")
}
