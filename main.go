package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

type AuthorsLocMap map[string]int

func (m *AuthorsLocMap) toalLoc() int {
	total := 0
	for _, lines := range *m {
		total += lines
	}
	return total
}

var bar *progressbar.ProgressBar

func main() {

	linesByAuthor, err := getChangedLinesByAuthor()
	if err != nil {
		fmt.Println(err)
		return
	}

	toalLoc := linesByAuthor.toalLoc()

	rows := make([]table.Row, 0)
	for author, lines := range linesByAuthor {
		rows = append(rows, table.Row{author, lines, fmt.Sprintf("%.2f%%", float32(lines*100)/float32(toalLoc))})
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SortBy([]table.SortBy{{Name: "LOC", Mode: table.DscNumeric}})
	t.AppendHeader(table.Row{"Author", "LOC", "Percentage"})
	t.AppendRows(rows)
	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", toalLoc})
	t.Render()

}

func getChangedLinesByAuthor() (AuthorsLocMap, error) {
	files, err := getAllGitFiles()
	if err != nil {
		fmt.Printf("get all files failed: %v\n", err)
		log.Fatal(err)
	}

	bar = progressbar.Default(int64(len(files)), "processing")

	authorsMap := make(AuthorsLocMap)
	for idx, file := range files {
		authors, err := getChangedLinesByAuthorForFile(file)
		if err != nil {
			return nil, err
		}
		for author, lines := range authors {
			authorsMap[author] += lines
		}
		bar.Add(idx + 1)
	}

	return authorsMap, nil
}

func getChangedLinesByAuthorForFile(file string) (AuthorsLocMap, error) {
	if len(file) == 0 {
		return nil, nil
	}
	// fmt.Printf("Processing file: %s\n", file)

	// check if file is under version control
	err := exec.Command("git", "ls-files", "--error-unmatch", file).Run()
	if err != nil {
		// file not under version control, skipping
		return nil, nil
	}

	cmd := exec.Command("git", "blame", "-w", "-M", "-e", file)
	out, err := cmd.Output()
	if err != nil {
		// fmt.Printf("git command failed: %v\n", err)
		return nil, err
	}

	linesByAuthor := make(AuthorsLocMap)
	for _, line := range strings.Split(string(out), "\n") {
		r, err := regexp.Compile(`<(.*?)>`)
		if err != nil {
			// fmt.Printf("regexp failed: %v\n", err)
			return nil, err
		}

		author := r.FindString(line)
		// fmt.Printf("line: %s\n", line)
		// fmt.Printf("author: %v\n", author)
		// fmt.Println()
		if len(author) == 0 {
			continue
		}

		linesByAuthor[author]++

	}

	// print info
	// toalLoc := linesByAuthor.toalLoc()
	// fmt.Printf("Total LOC: %d\n", toalLoc)
	// fmt.Println()

	return linesByAuthor, nil
}

func getAllGitFiles() ([]string, error) {
	output, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(output), "\n")
	return files, nil
}
