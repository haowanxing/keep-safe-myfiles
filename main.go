package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"keepsafemyfiles/enc"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	ignoreFiles = []string{"readme.md"}
	iv          = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	action      int
	key         []byte
	dir         string
)

func init() {
	flag.IntVar(&action, "a", 1, "1加密 2解密")
	flag.StringVar(&dir, "d", "./", "操作目录")
	flag.Parse()
}

func main() {
	ignoreFiles = append(ignoreFiles, os.Args[0])
	counter := 0
	fileLists, err := ListDir(dir)
	if err != nil {
		log.Fatalln("读取目录出错", err)
	}
	if len(fileLists) == 0 {
		log.Fatalln("没有文件可加密")
	}
	doStr := "加密"
	if action == 2 {
		doStr = "解密"
	}
	log.Printf("程序将自动为%d个文件进行'%s'操作！\n", len(fileLists), doStr)
	fmt.Print("Please input york code: ")
	_, err = fmt.Scanln(&key)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(err, fileLists, key, len(key))
	for _, filename := range fileLists {
		if action == 1 {
			log.Print("[Encoding]", filename)
			err := EncryptFile(filename, key)
			if err != nil {
				log.Fatalln(err)
			} else {
				counter++
				log.Println(" ok")
			}
		} else if action == 2 {
			log.Print("[Decoding]", filename)
			err := DecryptFile(filename, key)
			if err != nil {
				log.Fatalln(err)
			} else {
				counter++
				log.Println(" ok")
			}
		} else {
			log.Fatalln("Action not permitted", action)
			break
		}
	}
	log.Printf("Job done! we deal %d file in %d files", counter, len(fileLists))
}

func DecryptFile(filename string, key []byte) error {
	bt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	/*maxLen := base64.StdEncoding.DecodedLen(len(bt))
	log.Println("base64:",len(bt),maxLen)
	raw := make([]byte, maxLen)
	_, err = base64.StdEncoding.Decode(raw, bt)*/
	raw, err := base64.StdEncoding.DecodeString(string(bt))
	if err != nil {
		return err
	}
	orign, err := enc.AESDecrypt(raw, key, iv)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, orign, 0644)
	if err != nil {
		return err
	}
	return nil
}

func EncryptFile(filename string, key []byte) error {
	bt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	crypt, err := enc.AESEncrypt(bt, key, iv)
	if err != nil {
		return err
	}
	/*maxLen := base64.StdEncoding.EncodedLen(len(crypt))
	data := make([]byte, maxLen)
	base64.StdEncoding.Encode(data, crypt)*/
	data := []byte(base64.StdEncoding.EncodeToString(crypt))
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ListDir(folder string) (list []string, err error) {
	f, err := os.Stat(folder)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				if file.Name()[0:1] == "." {
					log.Println("skipDir:", file.Name())
					continue
				}
				subList, err := ListDir(folder + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				list = append(list, subList...)
			} else {
				if !CheckFile(file) {
					log.Println("skipFile:", file.Name())
					continue
				}
				strAbsPath, err := filepath.Abs(folder + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				list = append(list, strAbsPath)
			}
		}
	} else {
		absStr, err := filepath.Abs(folder)
		if err != nil {
			return nil, err
		}
		list = append(list, absStr)
	}

	return list, nil
}

//CheckFile 返回true的文件才会被操作
func CheckFile(file os.FileInfo) bool {
	if StrContain(file.Name(), ignoreFiles) {
		return false
	}
	ext := path.Ext(file.Name())
	if strings.EqualFold(ext, ".md") && !strings.EqualFold(file.Name(), "readme.md") {
		return true
	}
	return false
}

func StrContain(str string, slice []string) bool {
	for _, v := range slice {
		if str == v {
			return true
		}
	}
	return false
}
