package webapp

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
func GatherInfo() (map[string]string, []string) {
	// Open the data file
	inputFile, err := os.Open("output.jsonl")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return nil, nil
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

	return sectionsMap, roboticsHeaders
}

func AppResponse(userInput string, sectionsMap map[string]string, roboticsHeaders []string) string {
	//option to stop the chatbot
	if userInput == "exit" {
		return "Goodbye!"

		//option to list out all headers and see potential options for user input
	} else if userInput == "list" {
		// Print a list of all headers from the "Robotics" article
		response := "List of headers from the 'Robotics' article:\n"
		for _, header := range roboticsHeaders {
			response += header + "\n"
		}
		return response

		//If user input matches a header, the content mapped to that header will be produced
	} else {
		// Lookup and print the matching content
		content, exists := sectionsMap[userInput]
		if exists {
			return content

			//response if user input does not match any header
		} else {
			return "I'm sorry, I don't have any information on that topic. Please make sure your input is spelled correctly and matches the case formatting on the previous page. Thank you."
		}
	}
}

func removeBetween(str, start, end string) string {
	anyIncludingEndLine := fmt.Sprintf(`%s[\r\n\s\w]*%s`, start, end)
	return regexp.MustCompile(anyIncludingEndLine).ReplaceAllString(str, "")
}
