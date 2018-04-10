package base

import (
	"os"
	"log"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"io"
)

// 创建一个文件
func CreateFile(path string) (f *os.File, err error) {
	if CheckFileStat(path) {
		DeleteFile(path)
	}
	f, err = os.Create(path)
	if err != nil {
		log.Println("create file" + path + " faild!" )
	}
	return f,err
}

// 只创建文件不返回资源句柄
func CreateFileOnly(path string) (err error) {
	if CheckFileStat(path) {
		DeleteFile(path)
	}
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Println("create file" + path + " faild!" )
	}
	return err
}

// 删除一个文件
func DeleteFile(path string) error {
	if !CheckFileStat(path) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		log.Println("delete file " + path + " faild!")
	}
	return err
}

// 检测文件是否存在
func CheckFileStat(path string) bool {
	if _,err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// 检测某个 block 是否正常
func CheckBlockStat(path string,block *Block) bool {
	if CheckFileStat(path) {
		if int64(len(ReadFile(path))) == block.BlockSize {
			return true
		} else {
			return false
		}
	}
	return false
}

// 读取某个文件并返回 byte 数组
func ReadFile(path string) []byte {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatalln("read file error while open file!",err)
	}
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("read file error while read file!",err)
	}
	return fileBytes
}

// 向某个文件末尾添加内容
func AppendToFile(srcFile string, bytes []byte) error {
	f,err := os.OpenFile(srcFile,os.O_WRONLY,0644)
	defer f.Close()
	if err != nil {
		log.Println("open file " + srcFile + " failed!")
	} else {
		n,_ := f.Seek(0,os.SEEK_END)
		_,err = f.WriteAt(bytes,n)
	}
	return err
}

// 根据一个链接返回当前下载内容的文件名称
func GetUriName(url string) (prefixName, fullName string) {
	urlArr := []byte(url)
	qusIndex := strings.Index(url,"?")
	if qusIndex != -1 {
		urlArr = urlArr[0:qusIndex]
	}
	index := strings.LastIndex(url,"/")
	if index == -1 {
		fullName = string(urlArr[0:])
	} else if index == len(urlArr) - 1 {
		urlArr = urlArr[0:index]
		url = string(urlArr)
		index = strings.LastIndex(url,"/")
		if index == -1 {
			fullName = string(urlArr[0:])
		} else {
			fullName = string(urlArr[(index+1):])
		}
	}else {
		fullName = string(urlArr[(index+1):])
	}
	fullNameArr := []byte(fullName)
	pointIndex := strings.LastIndex(fullName,".")
	if pointIndex == -1 {
		prefixName = fullName
	} else {
		prefixName = string(fullNameArr[:pointIndex])
	}
	return prefixName,fullName
}

// 根据传入的字符串生成一个 32 位 MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 计算文件的 md5 值
func FileMd5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open file failed!", err)
		return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy file failed!", err)
		return ""
	}

	return string(md5hash.Sum(nil))
}