package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
)

type Object map[string]interface{}

func (s Object) Field(name string) Object {
	return Object(s[name].(map[string]interface{}))
}

func (s Object) List(name string) []Object {
	l := s[name].([]interface{})
	out := make([]Object, 0, len(l))
	for _, v := range l {
		out = append(out, Object(v.(map[string]interface{})))
	}
	return out
}

func (s Object) SetField(name string, newValue interface{}) {
	s[name] = newValue
}
func (s Object) SetList(name string, list []Object) {
	s[name] = list
}

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
	SavePath = flag.String("path", "save/Fast001.save", "save/Fast001.save")
	OutPath = flag.String("out", "Fast001.save", "Fast001.save")
	Items = flag.String("items", "", "增加物品--items=it30010")
	Skills = flag.String("skills", "", "增加主角技能 --skills=it30010")
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
	jsonFile.Truncate(0)
	if err != nil {
		panic(err)
	}

	fields := strings.Split(objpath, "/")
	out := save
	for _, field := range fields {
		out = out.Field(field)
	}

	encode := json.NewEncoder(jsonFile)
	encode.SetIndent("", "  ")
	err = encode.Encode(out)
	if err != nil {
		panic(err)
	}

	jsonFile.Close()
}

func read() {
	srcFile, err := os.Open(*SavePath)
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

	readIntoJson("player.json", "Character/Player")
	readIntoJson("nicknames.json", "AvailableNicknames")
	readIntoJson("areafriendly.json", "AreaFriendly")
	srcFile.Close()

}

func writeJson(filename string, objpath string) {
	file, err := os.OpenFile("player.json", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	var obj Object
	playerDecoder := json.NewDecoder(file)
	err = playerDecoder.Decode(&obj)
	if err != nil {
		panic(err)
	}

	fields := strings.Split(objpath, "/")
	out := save
	for _, field := range fields[:len(fields)-1] {
		out = out.Field(field)
	}
	lastField := fields[len(fields)-1]
	out.SetField(lastField, obj)
	file.Close()
}

func write() {
	srcFile, err := os.Open(*SavePath)
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
	err = json.Unmarshal(line, &save)
	if err != nil {
		panic(err)
	}

	writeJson("player.json", "Character/Player")
	writeJson("nicknames.json", "AvailableNicknames")
	writeJson("areafriendly.json", "AreaFriendly")

	savebin, err := json.Marshal(save)
	if err != nil {
		panic(err)
	}
	_, err = dstFile.Write(savebin)
	if err != nil {
		panic(err)
	}
	srcFile.Close()

	dstFile.Close()
}

func AddItems(player Object, itemnames []string) {
	oriItems := player.List("Inventory")
	for _, name := range itemnames {
		oriItems = append(oriItems, NewDefaultItem(name, 1, 0))
	}
	player.SetList("Inventory", oriItems)
}

func NewDefaultItem(itemid string, count int, level int) Object {
	return (Object)(map[string]interface{}{
		"Count":           count,
		"Durability":      0,
		"EffectId":        []string{},
		"ForgeMaterials":  struct{}{},
		"Hurt":            struct{}{},
		"HurtDifference":  0,
		"Id":              nil,
		"IsNew":           false,
		"ItemId":          itemid,
		"Level":           level,
		"MaxDurability":   0,
		"QualityTitle":    "",
		"QuenchEffect":    []string{},
		"QuenchHoleCount": 0,
		"ReforgeType":     -1,
		"Stolen":          false,
		"Weight":          0,
	})
}

func AddSkill(player Object, names []string) {
	skill := player.Field("SkillTree").Field("LearnedSkill")
	for _, name := range names {
		skill.SetField(name, NewSkill())
	}
}

func NewSkill() Object {
	return (Object)(map[string]interface{}{
		"ClickedNodes": []string{},
		"CurrentExp":   100000,
		"IsBackup":     false,
		"IsExpert":     false,
		"MaxExp":       999999,
		"Name":         nil,
	})
}
