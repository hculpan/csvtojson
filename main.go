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
		Level:        rec[0],
		Name:         rec[1],
		Reversible:   rec[2],
		Material:     rec[3],
		Sphere:       rec[4],
		Range:        rec[5],
		Components:   rec[6],
		CastingTime:  rec[7],
		Duration:     rec[8],
		AreaOfEffect: rec[9],
		SavingThrow:  rec[10],
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
		Level:        rec[0],
		Name:         rec[1],
		Reversible:   rec[2],
		School:       rec[3],
		Range:        rec[4],
		Components:   rec[5],
		Material:     rec[6],
		CastingTime:  rec[7],
		Duration:     rec[8],
		AreaOfEffect: rec[9],
		SavingThrow:  rec[10],
	}

	result, err := json.MarshalIndent(spell, "  ", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
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
		_, _ = datawriter.WriteString(data)
		if i == len(spells)-1 {
			_, _ = datawriter.WriteString("\n")
		} else {
			_, _ = datawriter.WriteString(",\n")
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
