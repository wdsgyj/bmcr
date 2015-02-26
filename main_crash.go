// main_crash
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

func MainCrash(args []string) {
	cmd := flag.NewFlagSet("crash", flag.ExitOnError)
	flagFile := cmd.String("crash", "", "Crash File")

	err := cmd.Parse(args)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(*flagFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	crashs := make([]*Crash, 0, 512)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) > 0 {
			crash, err := NewCrash(text)
			if err != nil {
				log.Fatalln(err)
			}
			crashs = append(crashs, crash)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	for _, c := range crashs {
		for _, page := range c.Pages {
			log.Println(page)
		}
		log.Println()
	}
}
