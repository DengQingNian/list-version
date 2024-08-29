package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"git.dengqn.com/dqn/listversion/util"
)

/**
VersionMeta

data root --> file path bash	/
								/ file version1
								/ file version2
								/ meta.json --> VersionMeta[]
		  --> other config.json
*/

type FileMeta struct {
	FileName     string    `json:"fileName"`
	AbsolutePath string    `json:"absolutePath"`
	NameHash     string    `json:"nameHash"`
	Versions     []Version `json:"versions"`
}

type Version struct {
	Version int64  `json:"version"`
	Created int64  `json:"created"`
	Desc    string `json:"desc"`
}

func All(current bool) (all []FileMeta) {
	root, _ := os.UserConfigDir()

	appRoot, err := os.Open(path.Join(root, "list-version"))
	if err != nil {
		fmt.Println("%s", "打开文件失败0"+err.Error())
		if errors.Is(err, os.ErrNotExist) {
			os.Mkdir(path.Join(root, "list-version"), os.ModePerm)
		}
	}

	defer func() {
		if appRoot != nil {
			appRoot.Close()
		}
	}()

	sub, _ := appRoot.Readdir(0)

	// filter this wd
	pwd, _ := os.Getwd()
	// log.Println("pwd:", pwd)
	for _, fi := range sub {
		var meta FileMeta

		fileMeta, err := os.Open(path.Join(root, "list-version", fi.Name(), "meta.json"))
		// not saved yet
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		// read meta.json

		buf, _ := io.ReadAll(fileMeta)
		json.Unmarshal(buf, &meta)
		fileMeta.Close()

		if current {
			if !strings.HasPrefix(meta.AbsolutePath, pwd) {
				continue
			}
		}

		if len(meta.Versions) == 0 {
			// rm
			err = os.RemoveAll(path.Join(root, "list-version", fi.Name()))
			if err != nil {
				log.Panicln(err.Error())
			}
			continue
		}

		fmt.Printf("[%d] %s\n", len(meta.Versions), meta.AbsolutePath)
	}

	return all
}

func Extract(fullPath string, version Version) {
	root, _ := os.UserConfigDir()
	pathHex := util.ToHashHex(fullPath)
	os.Remove(fullPath)
	origin, _ := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	target, _ := os.OpenFile(path.Join(root, "list-version", pathHex, strconv.FormatInt(version.Version, 10)), os.O_RDONLY, os.ModePerm)
	io.Copy(origin, target)
	origin.Close()
	target.Close()
}

func SaveMeta(pathHex string, meta FileMeta) {
	root, _ := os.UserConfigDir()

	appRoot, err := os.Open(path.Join(root, "list-version"))
	if err != nil {
		fmt.Println("%s", "打开文件失败0"+err.Error())
		if errors.Is(err, os.ErrNotExist) {
			os.Mkdir(path.Join(root, "list-version"), os.ModePerm)
		}
	}

	defer func() {
		if appRoot != nil {
			appRoot.Close()
		}
	}()

	_, err = os.Open(path.Join(root, "list-version", pathHex))
	// not saved yet
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.Mkdir(path.Join(root, "list-version", pathHex), os.ModePerm)
		}
	}

	os.Remove(path.Join(root, "list-version", pathHex, "meta.json.bak"))
	os.Rename(path.Join(root, "list-version", pathHex, "meta.json"), path.Join(root, "list-version", pathHex, "meta.json.bak"))
	// meta.json
	metaJsonFile, err := os.OpenFile(path.Join(root, "list-version", pathHex, "meta.json"), os.O_WRONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// 直接写
			metaJsonFile, err = os.Create(path.Join(root, "list-version", pathHex, "meta.json"))
		}
	}
	buf, err := json.Marshal(meta)
	if err != nil {
		fmt.Println("%s", "json"+err.Error())
	}
	metaJsonFile.Write(buf)
	metaJsonFile.Sync()
	metaJsonFile.Close()
}

func DeleteData(originFile string, version Version) {
	root, _ := os.UserConfigDir()
	pathHex := util.ToHashHex(originFile)
	os.Remove(path.Join(root, "list-version", pathHex, strconv.FormatInt(version.Version, 10)))
}

func CopyData(originFile string, version Version) {
	root, _ := os.UserConfigDir()
	pathHex := util.ToHashHex(originFile)

	origin, _ := os.OpenFile(originFile, os.O_RDONLY, os.ModePerm)
	target, _ := os.Create(path.Join(root, "list-version", pathHex, strconv.FormatInt(version.Version, 10)))
	io.Copy(target, origin)
	origin.Close()
	target.Close()
}

func GetVersionList(pathHex string) (meta FileMeta, err error) {

	root, _ := os.UserConfigDir()

	appRoot, err := os.Open(path.Join(root, "list-version"))
	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(path.Join(root, "list-version"), os.ModePerm)
	}
	defer func() {
		if appRoot != nil {
			appRoot.Close()
		}
	}()

	fileMeta, err := os.Open(path.Join(root, "list-version", pathHex, "meta.json"))
	// not saved yet
	if errors.Is(err, os.ErrNotExist) {
		return meta, errors.New("no versions")
	}
	defer fileMeta.Close()
	// read meta.json

	buf, _ := io.ReadAll(fileMeta)
	json.Unmarshal(buf, &meta)

	return meta, nil
}

func NewMeta(filePath string) (meta FileMeta) {

	meta.AbsolutePath = filePath

	st, _ := os.Stat(filePath)

	meta.FileName = st.Name()
	meta.NameHash = util.ToHashHex(filePath)
	meta.Versions = make([]Version, 0)

	return meta
}
