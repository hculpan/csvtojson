package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var isWizardSpell bool

type WizardSpell struct {
	Level        string `json:"level"`
	Name         string `json:"name"`
	Reversible   string `json:"reversible"`
	School       string `json:"school"`
	Range        string `json:"range"`
	Components   string `json:"components"`
	Material     string `json:"material"`
	CastingTime  string `json:"casting_time"`
	Duration     string `json:"duration"`
	AreaOfEffect string `json:"area_of_effect"`
	SavingThrow  string `json:"saving_throw"`
}

type ClericSpell struct {
	Level        string `json:"level"`
	Name         string `json:"name"`
	Reversible   string `json:"reversible"`
	Material     string `json:"material"`
	Sphere       string `json:"sphere"`
	Range        string `json:"range"`
	Components   string `json:"components"`
	CastingTime  string `json:"casting_time"`
	Duration     string `json:"duration"`
	AreaOfEffect string `json:"area_of_effect"`
	SavingThrow  string `json:"saving_throw"`
}

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func convertToCleric(rec []string) (string, error) {
	if len(rec) != 11 {
		return strings.Join(rec, ","), errors.New("incorrect number of fields")
	}

	spell := ClericSpell{
		Level:        properTitle(rec[0]),
		Name:         properTitle(rec[1]),
		Reversible:   properTitle(rec[2]),
		Material:     properTitle(rec[3]),
		Sphere:       properTitle(rec[4]),
		Range:        properTitle(rec[5]),
		Components:   rec[6],
		CastingTime:  properTitle(rec[7]),
		Duration:     properTitle(rec[8]),
		AreaOfEffect: properTitle(rec[9]),
		SavingThrow:  properTitle(rec[10]),
	}

	result, err := json.MarshalIndent(spell, "  ", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func convertToWizard(rec []string) (string, error) {
	if len(rec) != 11 {
		return strings.Join(rec, ","), errors.New("incorrect number of fields")
	}

	spell := WizardSpell{
		Level:        properTitle(rec[0]),
		Name:         properTitle(rec[1]),
		Reversible:   properTitle(rec[2]),
		School:       properTitle(rec[3]),
		Range:        properTitle(rec[4]),
		Components:   rec[5],
		Material:     properTitle(rec[6]),
		CastingTime:  properTitle(rec[7]),
		Duration:     properTitle(rec[8]),
		AreaOfEffect: properTitle(rec[9]),
		SavingThrow:  properTitle(rec[10]),
	}

	result, err := json.MarshalIndent(spell, "  ", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func properTitle(input string) string {
	words := strings.Fields(strings.ToLower(input))
	smallwords := " a an on the to "
	for index, word := range words {
		word = strings.Replace(word, "&", "and", -1)
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}

	return strings.Join(words, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Requires file name")
		return
	}

	isWizardSpell = os.Args[1][0] == 'W'
	// open file
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	spells := []string{}
	row := 1
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line

		if isWizardSpell {
			s, err := convertToWizard(rec)
			if err != nil {
				fmt.Printf("ERROR line %d: %+s\n", row, err)
				return
			}

			spells = append(spells, s)
		} else {
			s, err := convertToCleric(rec)
			if err != nil {
				fmt.Printf("ERROR line %d: %+s\n", row, err)
				return
			}

			spells = append(spells, s)
		}

		row++
	}

	baseFileName := fileNameWithoutExtSliceNotation(os.Args[1])

	file, err := os.OpenFile(baseFileName+".json", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	datawriter.WriteString("[\n")
	for i, data := range spells {
		_, err = datawriter.WriteString(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		if i == len(spells)-1 {
			_, err = datawriter.WriteString("\n")
			if err != nil {
				fmt.Printf("ERROR line %d: %+s\n", row, err)
				return
			}
		} else {
			_, err = datawriter.WriteString(",\n")
			if err != nil {
				fmt.Printf("ERROR line %d: %+s\n", row, err)
				return
			}
		}
	}
	datawriter.WriteString("]\n")

	datawriter.Flush()
	file.Close()

	/*
		for _, s := range spells {
			fmt.Printf("%s\n", s)
		}
	*/
	fmt.Println("File created: " + baseFileName + ".json")
}
