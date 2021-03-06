// crash
package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	_TIME_LAYOUT    = "20060102150405"
	_BAIDU_PREFIX_1 = "at com.baidu."
	_BAIDU_PREFIX_2 = "at map."
)

var (
	regxDetail *regexp.Regexp = regexp.MustCompile(`\{([^\s\{\}=/w]{3,}?)=`)
	regxPages  *regexp.Regexp = regexp.MustCompile(`#\d+:([^#]+)`)
)

type ComsInfo string

func (self ComsInfo) Infos() []string {
	txt := strings.TrimSpace(string(self))
	if len(txt) == 0 {
		return nil
	}

	infos := strings.Split(txt, "-")
	rs := make([]string, 0, len(infos))
	for _, info := range infos {
		rs = append(rs, info)
	}

	return rs
}

type Pages string

func (self Pages) Pages() []string {
	txt := strings.TrimSpace(string(self))
	if len(txt) == 0 {
		return nil
	}

	couples := regxPages.FindAllStringSubmatchIndex(txt, -1)
	var index int
	var page string
	rs := make([]string, 0, len(couples))
	for _, couple := range couples {
		page = txt[couple[2]:couple[3]] // group(1)
		if index = strings.LastIndex(page, "|"); index != -1 && index != len(page)-1 {
			page = page[index+1:]
		}
		rs = append(rs, page)
	}

	return rs
}

type Detail string

func (self Detail) Feature() string {
	txt := strings.TrimSpace(string(self))
	if len(txt) == 0 {
		return ""
	}
	lines := strings.Split(txt, "<br>")
	buf := new(bytes.Buffer)
	fill := false
	for _, line := range lines {
		if strings.HasPrefix(line, _BAIDU_PREFIX_1) ||
			strings.HasPrefix(line, _BAIDU_PREFIX_2) {
			fill = true
			fmt.Fprintln(buf, line)
		}
	}

	if !fill {
		for _, line := range lines {
			if strings.HasPrefix(line, "at ") {
				index := strings.IndexByte(line, '(') // 去除文件名 + 行号
				if index == -1 {
					fmt.Fprintln(buf, line)
				} else {
					fmt.Fprintln(buf, line[:index])
				}
			}
		}
	}

	return buf.String()
}

type MemInfo string

// HeapMax:128,DvmTotal:43712,DvmFree:12302,Pss:50401,Private:34224,Shared:11860
func (self MemInfo) Info() map[string]int {
	items := strings.Split(string(self), ",")
	if len(items) != 6 {
		return nil
	}

	rs := make(map[string]int)
	index := -1
	for _, item := range items {
		index = strings.Index(item, ":")
		if index != -1 {
			rs[item[:index]], _ = strconv.Atoi(item[index+1:])
		}
	}

	return rs
}

type Crash struct {
	Time int64 // 精度为秒
	Tm   int64 // 精度为秒
	Sv   string
	Sw   int
	Sh   int
	Ov   string
	Ch   string
	Mb   string
	Cuid string
	Net  int

	Detail       Detail
	Meminfo      MemInfo
	ActiveThread int
	Locx         int
	Locy         int
	CpuAbi       string
	CpuAbi2      string
	ComsInfo     ComsInfo
	Pages        Pages

	Bgm int
	Bgt int
	Bgw int
	Fgm int
	Fgt int
	Fgw int
}

func NewCrash(text string) (rs *Crash, err error) {
	rs = &Crash{}
	text = (strings.TrimSpace(text))[1 : len(text)-1] // 去除两边的 [] 字符
	logs := strings.Split(text, "][")

	// 处理每一个 mainItem
	for _, log := range logs {
		if e := parseLog(rs, log); e != nil {
			return nil, e
		}
	}
	return
}

func CreateCrashTable(db *sql.DB) error {
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

func NewCrashInsertStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare(`INSERT INTO crash (
			time, tm, sv, sw, sh,
			ov, ch, mb, cuid, net,
			detail, mem_info, thread_num, locx, locy,
			cpu_abi, cpu_abi2, feature, coms_info, pages,
			bgm, bgt, bgw, fgm, fgt, fgw)
			VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?)`)
}

func (self *Crash) Insert(stmt *sql.Stmt) error {
	_, err := stmt.Exec(
		self.Time,
		self.Tm,
		self.Sv,
		self.Sw,
		self.Sh,
		self.Ov,
		self.Ch,
		self.Mb,
		self.Cuid,
		self.Net,

		string(self.Detail),
		string(self.Meminfo),
		self.ActiveThread,
		self.Locx,
		self.Locy,
		self.CpuAbi,
		self.CpuAbi2,
		self.Detail.Feature(),
		string(self.ComsInfo),
		string(self.Pages),

		self.Bgm,
		self.Bgt,
		self.Bgw,
		self.Fgm,
		self.Fgt,
		self.Fgw)
	return err
}

/*
[time=20150206221945]
[tm=1423225353.493]
[pd=map]
[sv=7.8.0]
[sw=720]
[sh=1280]
[os=android]
[ov=Android18]
[ch=1006822a]
[mb=SM-G7106]
[ver=2]
[cuid=2721739AC836A6492F8E8030B0A75E81|554514450874953]
[net=9]
[lt=1100]
[act=crashlog]
[ActParam=……]
*/
func parseLog(crash *Crash, item string) error {
	index := strings.Index(item, "=")
	if index == -1 {
		return errors.New("There is no '=' in main item")
	}

	key := item[:index]
	value := item[index+1:]
	var err error

	switch key {
	case "time":
		t, err := time.ParseInLocation(_TIME_LAYOUT, value, time.Local)
		if err != nil {
			log.Println("Error time:", err, value)
		} else {
			crash.Time = t.Unix()
		}
	case "tm":
		inner := strings.Index(value, ".")
		var second int64
		if inner == -1 {
			second, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Println("Error tm:", err, value)
			}
		} else {
			second, err = strconv.ParseInt(value[:inner], 10, 64)
			if err != nil {
				log.Println("Error tm:", err, value)
			}
		}
		if err == nil {
			crash.Tm = second
		}
	case "sv":
		crash.Sv = value
	case "sw":
		crash.Sw, _ = strconv.Atoi(value)
	case "sh":
		crash.Sh, _ = strconv.Atoi(value)
	case "ov":
		crash.Ov = value
	case "ch":
		crash.Ch = value
	case "mb":
		crash.Mb = value
	case "cuid":
		crash.Cuid = strings.TrimSpace(value)
	case "net":
		crash.Net, _ = strconv.Atoi(value)
	case "ActParam":
		used := strings.TrimSpace(value)
		parseActParam(crash, used)
	}

	return nil
}

/*
{nl_success=[app_BaiduMapBaselib, BDSpeechDecoder_V1, gnustl_shared, app_BaiduMapApplib, bds, app_BaiduNaviApplib, app_BaiduVIlib, bd_etts, cpu_features, etts_domain_data_builder]}
{mem_info=HeapMax:128,DvmTotal:43712,DvmFree:12302,Pss:50401,Private:34224,Shared:11860}
{bgm=0}
{detail=java.lang.UnsatisfiedLinkError: Cannot load library: soinfo_relocate(linker.cpp:992): cannot locate symbol "_ZN9_baidu_vi8CVStringC1EPKc" referenced by "libapp_Diagnose.so"...<br>at java.lang.Runtime.loadLibrary(Runtime.java:372)<br>at java.lang.System.loadLibrary(System.java:514)<br>at com.baidu.component.diagnose.a.<init>(Diagnose.java:15)<br>at com.baidu.component.diagnose.a.<init>(Diagnose.java:8)<br>at com.baidu.component.diagnose.a$a.<clinit>(Diagnose.java:34)<br>at com.baidu.component.diagnose.a.a(Diagnose.java:30)<br>at com.baidu.component.diagnose.DiagnoseEntity.<init>(DiagnoseEntity.java:21)<br>at java.lang.Class.newInstanceImpl(Native Method)<br>at java.lang.Class.newInstance(Class.java:1319)<br>at com.baidu.mapframework.component2.comcore.a.b.b.a(Hook.java:51)<br>at com.baidu.mapframework.component2.comcore.a.a.b(ComRuntime.java:102)<br>at com.baidu.mapframework.component2.comcore.a.b$2.a(ComRuntimeEngine.java:79)<br>at com.baidu.mapframework.component2.comcore.a.b$3.run(ComRuntimeEngine.java:156)<br>at com.baidu.mapframework.component2.a.e$a.run(MultipleTaskQueue.java:79)<br>at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1080)<br>at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:573)<br>at com.baidu.platform.comapi.util.h$1.run(NamedThreadFactory.java:34)}
{reason=java.lang.UnsatisfiedLinkError: Cannot load library: soinfo_relocate(linker.cpp:992): cannot locate symbol "_ZN9_baidu_vi8CVStringC1EPKc" referenced by "libapp_Diagnose.so"...}
{active_thread=52}
{fgm=0}
{net=1}
{locx=12832685}
{locy=4006638}
{maps=...}
{cpu_abi2=armeabi}
{bgt=0}
{nl_fail=[]}
{coms_info=map.android.baidu.advertctrl_1.0.4-map.android.baidu.aoi_1.4.0-map.android.baidu.bus_1.1.6-map.android.baidu.cater_1.6.5-map.android.baidu.citybus_1.0.4-map.android.baidu.diagnose_1.0.4-map.android.baidu.groupon_1.4.0-map.android.baidu.hotel_1.9.4-map.android.baidu.indoorguide_1.0.0-map.android.baidu.international_1.2.1-map.android.baidu.maplab_1.0.6-map.android.baidu.movie_1.6.4-map.android.baidu.pano_1.1.7-map.android.baidu.qrcode_1.1.1-map.android.baidu.rentcar_1.6.5-map.android.baidu.scenery_3.5.5-map.android.baidu.subway_1.0.6-map.android.baidu.takeout_1.9.2-map.android.baidu.taxi_3.4.3-map.android.baidu.trafficradio_1.1.1-map.android.baidu.violation_1.0.4-map.android.baidu.voice_1.4.1-map.android.baidu.websdk_1.4.0-map.android.baidu.weekend_1.0.9-}
{bgw=0}
{pages=#0:map.android.baidu.mainmap|com.baidu.baidumaps.MapsActivity@41f87930|com.baidu.baidumaps.base.MapFramePage}
{fgt=9}
{cpu_abi=armeabi-v7a}
{fgw=5}
*/
func parseActParam(crash *Crash, txt string) error {
	indexSlice := regxDetail.FindAllStringSubmatchIndex(txt, -1)
	size := len(indexSlice)
	var key, value string
	getValue := func(i, size int, slice [][]int, txt string) (string, string) {
		var k, v string

		k = txt[slice[i][2]:slice[i][3]]

		if i < size-1 {
			v = txt[slice[i][1]:slice[i+1][0]]
		} else {
			v = txt[slice[i][1]:]
		}

		index := strings.LastIndex(v, "}")
		if index != -1 {
			v = v[:index] // 去除 } 字符
		}

		return k, v
	}

	for i := range indexSlice {
		key, value = getValue(i, size, indexSlice, txt)
		switch key {
		case "detail":
			crash.Detail = Detail(value)
		case "mem_info":
			crash.Meminfo = MemInfo(value)
		case "active_thread":
			crash.ActiveThread, _ = strconv.Atoi(value)
		case "locx":
			crash.Locx, _ = strconv.Atoi(value)
		case "locy":
			crash.Locy, _ = strconv.Atoi(value)
		case "cpu_abi":
			crash.CpuAbi = value
		case "cpu_abi2":
			crash.CpuAbi2 = value
		case "bgm":
			crash.Bgm, _ = strconv.Atoi(value)
		case "bgt":
			crash.Bgt, _ = strconv.Atoi(value)
		case "bgw":
			crash.Bgw, _ = strconv.Atoi(value)
		case "fgm":
			crash.Fgm, _ = strconv.Atoi(value)
		case "fgt":
			crash.Fgt, _ = strconv.Atoi(value)
		case "fgw":
			crash.Fgw, _ = strconv.Atoi(value)
		case "coms_info":
			crash.ComsInfo = ComsInfo(value)
		case "pages":
			crash.Pages = Pages(value)
		}
	}
	return nil
}
