package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"git.dengqn.com/dqn/listversion/storage"
	"git.dengqn.com/dqn/listversion/util"
)

func main() {

	op := "h"
	if len(os.Args) > 1 {
		op = os.Args[1]
	}
	// for idx, arg := range os.Args {
	// 	fmt.Printf("%d: %s\n", idx, arg)
	// }

	switch op {
	case "h":
		h()
		break
	case "l":
		ls()
		break
	case "a":
		add()
		break
	case "d":
		del()
	case "r":
		rec()
	case "v":
		v()
		break
	}
}

func v() {
	path := ""

	if len(os.Args) > 2 {
		path = os.Args[2]
	}

	storage.All(path == "." || path == "./")
}

func h() {
	fmt.Println("________________________________________________________")
	fmt.Println("[h] print this")
	fmt.Println("[l] list verisons. [lv l xx.txt]")
	fmt.Println("[a] add file to cache. [lv a xx.txt xxx]")
	fmt.Println("[d] remove version from cache. [lv d xx.txt 1]")
	fmt.Println("[r] extract file version from cache. [lv r xx.txt 1]")
	fmt.Println("--------------------------------------------------------")
}

func rec() {
	// 获取绝对路径
	// 获取绝对路径
	fullPath, _ := filepath.Abs(os.Args[2])
	ver := os.Args[3]

	version, _ := strconv.ParseInt(ver, 10, 32)

	meta, err := storage.GetVersionList(util.ToHashHex(fullPath))
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}

	for _, v := range meta.Versions {
		func(vv storage.Version) {
			if v.Version == version {
				storage.Extract(fullPath, vv)
			}
		}(v)
	}
}

func ls() {
	// 获取绝对路径
	fullPath, _ := filepath.Abs(os.Args[2])

	meta, err := storage.GetVersionList(util.ToHashHex(fullPath))
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	fmt.Printf("[%s]\n______________________\n", meta.FileName)
	for _, v := range meta.Versions {
		fmt.Printf("[%d].%s\t%s\n", v.Version, time.UnixMilli(v.Created).Format("2006-01-02 15:04:05"), v.Desc)
	}
}

func add() {
	// xxx add xx "asdasd"
	fullPath, _ := filepath.Abs(os.Args[2])
	msg := os.Args[3]

	meta, err := storage.GetVersionList(util.ToHashHex(fullPath))
	if err != nil {
		// 新的
		meta = storage.NewMeta(fullPath)
	}

	max := int64(0)
	for _, v := range meta.Versions {
		if v.Version > int64(max) {
			max = v.Version
		}
	}

	meta.Versions = append(meta.Versions, storage.Version{
		Desc:    msg,
		Created: time.Now().UnixMilli(),
		Version: max + 1,
	})

	// fmt.Println("meta:", meta)

	storage.SaveMeta(util.ToHashHex(fullPath), meta)
	storage.CopyData(fullPath, meta.Versions[len(meta.Versions)-1])
}

func del() {
	// 获取绝对路径
	fullPath, _ := filepath.Abs(os.Args[2])
	ver := os.Args[3]

	version, _ := strconv.ParseInt(ver, 10, 32)

	meta, err := storage.GetVersionList(util.ToHashHex(fullPath))
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}

	tmp := make([]storage.Version, 0)
	for _, v := range meta.Versions {
		func(vv storage.Version) {
			if vv.Version != version {
				tmp = append(tmp, vv)
			} else {
				storage.DeleteData(fullPath, vv)
			}
		}(v)
	}
	meta.Versions = tmp
	storage.SaveMeta(util.ToHashHex(fullPath), meta)
}
