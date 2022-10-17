package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

// Global setting
var (
	nameFile string
)

// This can be get flag setting
func init() {
	nameFile = "data.tsv"
}

// Films Structure file
// tconst|titleType|primaryTitle|originalTitle|isAdult|startYear|endYear|runtimeMinutes|genres
type Films struct {
	PrimaryTitle   string `json:"primaryTitle"`
	OoriginalTitle string `json:"originalTitle"`
	Tconst         string `json:"tconst"`
	TitleType      string `json:"titleType"`
	IsAdult        string `json:"IsAdult"`
	StartYear      string `json:"startYear"`
	EndYear        string `json:"endYear"`
	RuntimeMinutes string `json:"runtimeMinutes"`
	Genres         string `json:"genres"`
}

// Main process
func main() {
	ch := make(chan []Films, 10000)

	// Read from chanel and do
	go ReadResult(ch)

	csvFile, _ := os.Open(nameFile)
	defer csvFile.Close()

	// Building composite filter
	// Example filtering
	filters := []string{"Phantom of the Opera", "Heinrich bringt", "Khoshbakhti", "Off for the Day", "Almuñécar"}

	// Filtering by filters
	processCSV(ch, csvFile, filters)
}

// BuilFilter Building filter
func BuilderFilter(txtfilter []string) []*regexp.Regexp {
	regs := []*regexp.Regexp{}
	for _, tf := range txtfilter {
		regs = append(regs, regexp.MustCompile("([^/])("+tf+"((/.*)|()))"))
	}
	return regs
}

// ReadResult : Reading results from chanel
func ReadResult(ch chan []Films) {

	// Can do any channel
	// select {
	// case msg1 := <- c1:
	// fmt.Println(msg1)

	for r := range ch {
		fmt.Println("> ", r[0])
		<-ch
		// Do testing with getting information...
	}
}

// **************************************************
// Read from file with output to chanel
// **************************************************
func processCSV(ch chan []Films, rc io.Reader, txtfilter []string) {
	
	// Filtered creation
	regs := BuilderFilter(txtfilter)
	
	r := csv.NewReader(rc)

	// read header
	if _, err := r.Read(); err != nil {
		fmt.Println("Error : Header: ", err)
	}

	// Setting is reader
	r.Comma = '\t'
	r.FieldsPerRecord = -1
	start := time.Now()

	// Cycle about file
	for {
		line, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error: Read file :", err)
		}
		go func() {
			films := make([]Films, 0)

			// for _, re := range []*regexp.Regexp{re1, re2} {
			for _, re := range regs {
				s := line[2]

				if re.MatchString(s) {
					// Add search film by filter
					films = append(films, Films{PrimaryTitle: line[0],
						OoriginalTitle: line[1],
						Tconst:         line[2],
						TitleType:      line[3],
						IsAdult:        line[4],
						StartYear:      line[5],
						EndYear:        line[6],
						RuntimeMinutes: line[7],
						Genres:         line[8],
					})
					ch <- films
				}
			}
		}()
	}

	// Show work time
	finish := time.Now()

	fmt.Println("Duration :", finish.Sub(start))
}
