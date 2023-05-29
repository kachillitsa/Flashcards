package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var stringLog strings.Builder

type Card struct {
	Term       string
	Definition string
	Mistakes   int
}

var deck = make(map[string]*Card)

func readLinesStr() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readLinesInt() int {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	res, _ := strconv.Atoi(strings.TrimSpace(line))
	return res
}

func checkTerm() string {
	term := readLinesStr()
	fmt.Fprintln(&stringLog, term)
	_, ok := deck[term]
	if ok {
		fmt.Printf("The term \"%s\" already exists. Try again:\n", term)
		fmt.Fprintf(&stringLog, "The term \"%s\" already exists. Try again:\n", term)
		term = checkTerm()
	}
	return term
}

func checkDef(term string) string {
	definition := readLinesStr()
	fmt.Fprintln(&stringLog, definition)
	for _, val := range deck {
		if val.Definition == definition {
			fmt.Printf("The definition \"%s\" already exists. Try again:\n", definition)
			fmt.Fprintf(&stringLog, "The definition \"%s\" already exists. Try again:\n", definition)
			definition = checkDef(term)
		}
	}
	return definition
}

func addNewCard() {
	fmt.Printf("The card term:\n")
	fmt.Fprintf(&stringLog, "The card term:\n")
	term := checkTerm()
	fmt.Printf("The definition:\n")
	fmt.Fprintf(&stringLog, "The definition:\n")
	definition := checkDef(term)
	newCard := Card{term, definition, 0}
	deck[term] = &newCard
	fmt.Printf("The pair (\"%s\":\"%s\") has been added\n", term, definition)
	fmt.Fprintf(&stringLog, "The pair (\"%s\":\"%s\") has been added\n", term, definition)
	return
}

func quizCheckAnswers(askTill0 int) int {
	var checked int
	for key, val := range deck {
		if askTill0 == 0 {
			return askTill0
		}
		askTill0--
		fmt.Printf("Print the definition of \"%s\"\n", key)
		fmt.Fprintf(&stringLog, "Print the definition of \"%s\"\n", key)
		answer := readLinesStr()
		fmt.Fprintln(&stringLog, answer)
		for key2, val2 := range deck {
			if val2.Definition == answer && val2.Term != key {
				fmt.Printf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\"\n", val.Definition, key2)
				fmt.Fprintf(&stringLog, "Wrong. The right answer is \"%s\", but your definition is correct for \"%s\"\n", val.Definition, key2)
				val2.Mistakes++
				checked = 1
			}
		}
		if checked == 0 {
			switch answer {
			case val.Definition:
				fmt.Printf("Correct!\n")
				fmt.Fprintf(&stringLog, "Correct!\n")
			default:
				val.Mistakes++
				fmt.Printf("Wrong. The right answer is \"%s\".\n", val.Definition)
				fmt.Fprintf(&stringLog, "Wrong. The right answer is \"%s\".\n", val.Definition)
			}
		}
		checked = 0
	}
	return askTill0
}

func quizStart() {
	var askTill0 int
	if len(deck) == 0 {
		fmt.Printf("There are no cards yet. Try to add one!\n")
		fmt.Fprintf(&stringLog, "There are no cards yet. Try to add one!\n")
		return
	}
	fmt.Printf("How many times to ask?\n")
	fmt.Fprintf(&stringLog, "How many times to ask?\n")
	askTill0 = readLinesInt()
	fmt.Fprintln(&stringLog, askTill0)
	for askTill0 != 0 {
		askTill0 = quizCheckAnswers(askTill0)
	}
	return
}

func removeCard() {
	fmt.Printf("Which card?\n")
	fmt.Fprintf(&stringLog, "Which card?\n")
	cardToRemove := readLinesStr()
	fmt.Fprintln(&stringLog, cardToRemove)
	_, ok := deck[cardToRemove]
	if ok {
		delete(deck, cardToRemove)
		fmt.Printf("The card has been removed.\n")
		fmt.Fprintf(&stringLog, "The card has been removed.\n")
		return
	}
	fmt.Printf("Can't remove \"%s\": there is no such card.\n", cardToRemove)
	fmt.Fprintf(&stringLog, "Can't remove \"%s\": there is no such card.\n", cardToRemove)
	return
}

func exportDeck(fileName string) {
	if fileName == "noFilename" {
		fmt.Printf("File name:\n")
		fmt.Fprintf(&stringLog, "File name:\n")
		fileName = readLinesStr()
		fmt.Fprintln(&stringLog, fileName)
	}

	deckJSON, _ := json.Marshal(deck)

	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer file.Close()
	fmt.Fprintf(file, string(deckJSON))
	fmt.Printf("%d cards have been saved.\n", len(deck))
	fmt.Fprintf(&stringLog, "%d cards have been saved.\n", len(deck))
	return
}

func importDeck(fileName string) {
	var countNewElements int
	var rawDeck = make(map[string]*Card)

	if fileName == "noFilename" {
		fmt.Printf("File name:\n")
		fmt.Fprintf(&stringLog, "File name:\n")
		fileName = readLinesStr()
		fmt.Fprintln(&stringLog, fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("File not found.\n")
		fmt.Fprintf(&stringLog, "File not found.\n")
		return
	}
	defer file.Close()

	data, _ := ioutil.ReadAll(file) // create a new Scanner for the file
	json.Unmarshal(data, &rawDeck)
	for key, val := range rawDeck {
		newCard := Card{val.Term, val.Definition, val.Mistakes}
		countNewElements++
		deck[key] = &newCard
	}
	fmt.Printf("%d cards have been loaded.\n", countNewElements)
	fmt.Fprintf(&stringLog, "%d cards have been loaded.\n", countNewElements)
	return
}

func saveLog() {
	fmt.Printf("File name:\n")
	fmt.Fprintf(&stringLog, "File name:\n")
	logFileName := readLinesStr()
	fmt.Fprintln(&stringLog, logFileName)

	logFile, _ := os.Create(logFileName)
	defer logFile.Close()

	fmt.Fprintln(logFile, stringLog.String())
	fmt.Printf("The log has been saved.\n")
	fmt.Fprintf(&stringLog, "The log has been saved.\n")
	return
}

func hardestCard() {
	var mostMistakes int
	var theCards []Card
	var theCardsStrArr []string

	for _, val := range deck {
		if val.Mistakes > mostMistakes {
			mostMistakes = val.Mistakes
			theCards = nil
		}
		if val.Mistakes == mostMistakes {
			theCards = append(theCards, *val)
		}
	}
	if mostMistakes == 0 {
		fmt.Printf("There are no cards with errors.\n")
		fmt.Fprintf(&stringLog, "There are no cards with errors.\n")
		return
	}
	if len(theCards) > 1 {
		for _, val := range theCards {
			theCardsStrArr = append(theCardsStrArr, "\""+val.Term+"\"")
		}
	} else {
		theCardsStrArr = append(theCardsStrArr, theCards[0].Term)
	}
	theCardsStr := strings.Join(theCardsStrArr, ", ")
	fmt.Printf("The hardest card is %s. You have %d errors answering it.\n", theCardsStr, mostMistakes)
	fmt.Fprintf(&stringLog, "The hardest card is %s. You have %d errors answering it.\n", theCardsStr, mostMistakes)
	return
}

func resetLog() {
	for _, val := range deck {
		val.Mistakes = 0
	}
	fmt.Printf("Card statistics have been reset.\n")
	fmt.Fprintf(&stringLog, "Card statistics have been reset.\n")
	return
}

func main() {
	var action string

	fileImport := flag.String("import_from", "noFilename", "Enter file name to import")
	fileExport := flag.String("export_to", "noFilename", "Enter file name to export")
	flag.Parse()
	if *fileImport != "noFilename" {
		importDeck(*fileImport)
	}

	for {
		fmt.Printf("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):\n")
		fmt.Fprintf(&stringLog, "Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):\n")
		action = readLinesStr()
		fmt.Fprintln(&stringLog, action)
		switch action {
		case "add":
			addNewCard()
		case "remove":
			removeCard()
		case "import":
			importDeck("noFilename")
		case "export":
			exportDeck("noFilename")
		case "ask":
			quizStart()
		case "exit":
			fmt.Printf("Bye bye!\n")
			fmt.Fprintf(&stringLog, "Bye bye!\n")
			exportDeck(*fileExport)
			return
		case "log":
			saveLog()
		case "hardest card":
			hardestCard()
		case "reset stats":
			resetLog()
		}
	}
}
