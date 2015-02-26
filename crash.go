// crash
package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	_TIME_LAYOUT = "20060102150405"
)

var (
	regxDetail *regexp.Regexp = regexp.MustCompile(`\{([^\s\{\}=/w]{3,}?)=`)
	regxPages  *regexp.Regexp = regexp.MustCompile(`#\d+:([^#]+)`)
)

type Crash struct {
	Time time.Time
	Tm   time.Time
	Sv   string
	Sw   int
	Sh   int
	Ov   string
	Ch   string
	Mb   string
	Cuid string
	Net  int

	Detail       string
	Meminfo      string
	ActiveThread int
	locx         int
	locy         int
	CpuAbi       string
	CpuAbi2      string
	Feature      string
	ComsInfo     []string
	Pages        []string

	Bgm int
	Bgt int
	Bgw int
	Fgm int
	Fgt int
	Fgw int
}

func New(text string) (rs *Crash, err error) {
	rs = &Crash{}
	text = (strings.TrimSpace(text))[1 : len(text)-1] // 去除两边的 [] 字符
	mainItems := strings.Split(text, "][")

	// 处理每一个 mainItem
	for _, item := range mainItems {
		if e := parseMainItem(rs, item); e != nil {
			return nil, e
		}
	}
	return
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
func parseMainItem(crash *Crash, item string) error {
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
			crash.Time = t
		}
	case "tm":
		inner := strings.Index(value, ".")
		var second, nsec int64
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
			nsec, err = strconv.ParseInt(value[inner+1:]+"000000", 10, 64)
			if err != nil {
				log.Println("Error tm:", err, value)
			}
		}
		if err == nil {
			crash.Tm = time.Unix(second, nsec)
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
		if i != -1 {
			v = v[:index] // 去除 } 字符
		}

		return k, v
	}

	for i := range indexSlice {
		key, value = getValue(i, size, indexSlice, txt)
		switch key {
		case "detail":
			crash.Detail = value
			parseFeature(crash, value)
		case "mem_info":
			crash.Meminfo = value
		case "active_thread":
			crash.ActiveThread, _ = strconv.Atoi(value)
		case "locx":
			crash.locx, _ = strconv.Atoi(value)
		case "locy":
			crash.locy, _ = strconv.Atoi(value)
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
			parseComsInfo(crash, value)
		case "pages":
			parsePages(crash, value)
		}
	}
	return nil
}

const (
	_BAIDU_PREFIX_1 = "at com.baidu."
	_BAIDU_PREFIX_2 = "at map."
)

func parseFeature(crash *Crash, txt string) {
	txt = strings.TrimSpace(txt)
	if len(txt) == 0 {
		return
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

	crash.Feature = buf.String()
}

func parseComsInfo(crash *Crash, txt string) {
	txt = strings.TrimSpace(txt)
	if len(txt) == 0 {
		return
	}

	infos := strings.Split(txt, "-")
	for _, info := range infos {
		crash.ComsInfo = append(crash.ComsInfo, info)
	}
}

func parsePages(crash *Crash, txt string) {
	txt = strings.TrimSpace(txt)
	if len(txt) == 0 {
		return
	}

	couples := regxPages.FindAllStringSubmatchIndex(txt, -1)
	var index int
	var page string
	for _, couple := range couples {
		page = txt[couple[2]:couple[3]]
		if index = strings.LastIndex(page, "|"); index != -1 && index != len(page)-1 {
			page = page[index+1:]
		}
		crash.Pages = append(crash.Pages, page)
	}
}
