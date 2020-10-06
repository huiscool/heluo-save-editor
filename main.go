package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Object map[string]interface{}

// func (s Object) Field(name string) Object {
// 	return Object(s[name].(map[string]interface{}))
// }

// func (s Object) SetField(name string, newValue interface{}) {
// 	s[name] = newValue
// }

var (
	Write    *bool
	SavePath *string
	OutPath  *string
	Items    *string
	Skills   *string
)

var (
	save Object
)

func main() {
	Write = flag.Bool("write", false, "")
	SavePath = flag.String("src", "save/Fast001.save", "save/Fast001.save")
	OutPath = flag.String("out", "Fast001.save", "Fast001.save")
	flag.Parse()
	if *Write {
		write()
	} else {
		read()
	}

}

func readIntoJson(filename string, objpath string) {
	jsonFile, err := os.OpenFile("player.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	objFile, err := os.OpenFile("obj.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	objFile.Truncate(0)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(srcFile)
	var line []byte
	cnt := 7
	for i := 0; i < cnt; i++ {
		line, err = reader.ReadBytes('\n')
		if err != io.EOF {
			continue
		}
	}

	err = json.Unmarshal(line, &save)
	if err != nil {
		panic(err)
	}
	encode := json.NewEncoder(objFile)
	encode.SetIndent("", "  ")
	err = encode.Encode(save)

	if err != nil {
		panic(err)
	}

	objFile.Close()
	srcFile.Close()

}

func write() {
	srcFile, err := os.Open(*SavePath)
	if err != nil {
		panic(err)
	}
	objFile, err := os.OpenFile("obj.json", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	dstFile, err := os.Create(*OutPath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(srcFile)
	var line []byte
	cnt := 7
	for i := 0; i < cnt; i++ {
		line, err = reader.ReadBytes('\n')
		if err != io.EOF {
			dstFile.Write(line)
			continue
		}
	}
	// 最后一行解析一下，然后改掉
	var save Object
	objDecoder := json.NewDecoder(objFile)
	err = objDecoder.Decode(&save)
	if err != nil {
		panic(err)
	}

	savebin, err := json.Marshal(save)
	if err != nil {
		panic(err)
	}
	_, err = dstFile.Write(savebin)
	if err != nil {
		panic(err)
	}
	srcFile.Close()
	objFile.Close()
	dstFile.Close()
}
