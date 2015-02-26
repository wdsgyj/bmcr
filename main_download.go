// main_download
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_ERROR_URI = "http://mapmo.baidu.com/lbspm/apps/crashlog/protected/logdata/map/map_android_crash/%s"
	_ASCII_    = 8
)

func MainDownload(args []string) {
	cmd := flag.NewFlagSet("download", flag.ExitOnError)
	flagDate := cmd.String("date", "", "日期，例如 2015-02-23")
	flagDb := cmd.String("db", "data.db", "数据库")
	cmd.Parse(args)

	if len(*flagDate) == 0 {
		log.Fatalln("no date found!")
	}

	db, err := sql.Open("sqlite3", *flagDb)
	if err != nil {
		log.Fatalln(err)
	}

	if err = createTableIfNeeded(db); err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf(_ERROR_URI, *flagDate)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	total, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	writed := 0
	var textGeted string
	var chars int

	crashs := make([]*Crash, 0, 512)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		textGeted = scanner.Text()
		writed += len(textGeted)
		if total != 0 {
			if chars != 0 {
				fmt.Printf("%c", '\u000d')
				chars = 0
			}
			chars, _ = fmt.Printf("进度：%.0f%%", float32(writed)*100/float32(total))
		}

		text := strings.TrimSpace(textGeted)
		if len(text) > 0 {
			crash, err := NewCrash(text)
			if err != nil {
				log.Println(err)
			} else {
				crashs = append(crashs, crash)
			}
		}
	}
	if total != 0 {
		fmt.Println()
	}

	if err = scanner.Err(); err != nil {
		log.Println(err)
	}
}

func createTableIfNeeded(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS crash (
		id INTEGER PRIMARY KEY,
		time INTEGER,
		tm INTEGER,
		sv TEXT,
		sw INTEGER,
		sh INTEGER,
		ov TEXT,
		ch TEXT,
		mb TEXT,
		cuid TEXT,
		net INTEGER,

		detail TEXT,
		mem_info TEXT,
		thread_num INTEGER,
		locx INTEGER,
		locy INTEGER,
		cpu_abi TEXT,
		cpu_abi2 TEXT,
		feature TEXT,
		coms_info TEXT,
		pages TEXT,

		bgm INTEGER,
		bgt INTEGER,
		bgw INTEGER,
		fgm INTEGER,
		fgt INTEGER,
		fgw INTEGER,

		UNIQUE (time, tm, cuid) ON CONFLICT IGNORE);`)
	return err
}
