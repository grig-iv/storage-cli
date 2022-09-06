package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("missing command")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]
	err := execute(cmd, args)
	if err != nil {
		fmt.Println(err)
	}
}

func execute(cmd string, args []string) error {
	var err error

	if cmd == "add" {
		err = add(args)
	} else {
		err = errors.New(fmt.Sprintf("Unknown command '%s'", cmd))
	}

	return err
}

func add(args []string) error {
	word := strings.TrimSpace(args[0])
	requestURL := fmt.Sprintf("https://www.ldoceonline.com/search/english/direct/?q=%s+", word)

	res, err := http.Get(requestURL)
	if err != nil {
		res, err = http.Get(requestURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer res.Body.Close()

	page, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile("<h1 class=\"pagetitle\">.*</h1>")
	match := string(reg.Find(page))
	match = strings.ReplaceAll(match, "<h1 class=\"pagetitle\">", "")
	match = strings.ReplaceAll(match, "</h1>", "")

	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    sqlStmt := `SELECT Word FROM Words WHERE word = ?`
    var qres string
    db.QueryRow(sqlStmt, word).Scan(&qres)

    fmt.Print(qres != "")
	return err
}

func t() (int, int) {
    return 1, 2
}

