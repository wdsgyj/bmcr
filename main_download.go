// main_download
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	_ERROR_URI = "http://mapmo.baidu.com/lbspm/apps/crashlog/protected/logdata/map/map_android_crash/%s"
)

func MainDownload(args []string) {
	cmd := flag.NewFlagSet("download", flag.ExitOnError)
	flagDate := cmd.String("date", "", "必须，日期，例如 2015-02-23")
	flagDb := cmd.String("db", "data.db", "可选，数据库文件位置")
	cmd.Parse(args)

	hasError := false

	if len(*flagDate) == 0 {
		fmt.Println("没有找到 -date 选项")
		fmt.Println("使用说明：")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", *flagDb)
	if err != nil {
		log.Fatalln(err)
	}

	if err = CreateCrashTable(db); err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf(_ERROR_URI, *flagDate)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("StatusCode:", resp.StatusCode)
		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Printf("%s\n", content)
		}

		os.Exit(1)
	}

	total, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	writed := 0
	var textGeted string
	var chars int

	scanner := bufio.NewScanner(resp.Body) // 逐行扫描

	fmt.Println("开始下载并录入数据……")

	tx, err := db.Begin() // 开始数据库事务
	if err != nil {
		log.Fatalln(err)
	}

	stmt, err := NewCrashInsertStmt(db)
	if err != nil {
		log.Fatalln(err)
	}

	stmt = tx.Stmt(stmt) // 转换为数据库事务的 stmt

	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				fmt.Println("\n出现运行时错误，数据库回滚……")
				if err = tx.Rollback(); err != nil {
					fmt.Println("回滚失败！", err)
				} else {
					fmt.Println("回滚成功！")
				}
			}
			fmt.Println(e)
		}
	}()

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
				hasError = true
			} else {
				err = crash.Insert(stmt)
				if err != nil {
					log.Println(err)
					hasError = true
				}
			}
		}
	}
	if total != 0 {
		fmt.Println()
	}

	// 检查，是否扫描过程中发生了错误
	if err = scanner.Err(); err != nil {
		log.Println(err)
		hasError = true
	}

	if hasError {
		fmt.Println("出现错误，数据库回滚……")
		if err = tx.Rollback(); err != nil {
			fmt.Println("回滚失败！", err)
		} else {
			fmt.Println("回滚成功！")
		}
		os.Exit(1)
	} else {
		fmt.Println("提交数据库事务……")
		if err = tx.Commit(); err != nil {
			fmt.Println("提交失败！")
			os.Exit(1)
		}
	}

	fmt.Println("处理完成！ --> ", *flagDb)
}
