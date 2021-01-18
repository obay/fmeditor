package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	flag "github.com/spf13/pflag"
	"gopkg.in/djherbis/times.v1"
)

func check(e error) {
	if e != nil {
		err := sentry.Init(sentry.ClientOptions{})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		sentry.CaptureException(e)
		panic(e)
	}
}

// Frontmatter struct to hold all Hugo-related frontmatter fileds
type Frontmatter struct {
	title   string
	author  string
	date    time.Time
	updated time.Time
	draft   bool
}

func getFrontMatter(path string, author string, draft bool) string {
	t, _ := times.Stat(path)

	fm := Frontmatter{}
	fm.title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	fm.author = "\"" + author + "\""
	fm.date = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if t.HasBirthTime() {
		fm.date = t.BirthTime()
	}
	fm.updated = t.ModTime()
	fm.draft = draft

	frontMatterText := "---\n"
	frontMatterText = frontMatterText + "title: " + fm.title + "\n"
	frontMatterText = frontMatterText + "author: " + fm.author + "\n"
	frontMatterText = frontMatterText + "date: " + fm.date.Format(time.RFC3339) + "\n"
	frontMatterText = frontMatterText + "updated: " + fm.updated.Format(time.RFC3339) + "\n"
	frontMatterText = frontMatterText + "draft: " + strconv.FormatBool(fm.draft) + "\n"
	frontMatterText = frontMatterText + "---\n"

	return frontMatterText
}

func main() {
	var author *string = flag.String("author", "", "Set the post(s) author name.")
	var draft *bool = flag.Bool("draft", true, "Set the post(s) draft status to true or false.")
	var rootfolder *string = flag.String("rootfolder", "", "Root folder holding markdown files")
	flag.Parse()
	if *rootfolder == "" || *author == "" {
		flag.Usage()
		return
	}

	err := filepath.Walk(*rootfolder, func(path string, info os.FileInfo, err error) error {
		check(err)
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		contentb, err := ioutil.ReadFile(path)
		content := string(contentb)
		check(err)

		if strings.HasPrefix(content, "---") {
			fmt.Println("Skipping " + filepath.Base(path) + ". Already has front matter.")

			return nil
		}

		newContent := getFrontMatter(path, *author, *draft) + "\n" + content

		err = ioutil.WriteFile(path, []byte(newContent), 0644)
		check(err)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
