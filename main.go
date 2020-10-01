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

func read() {
	srcFile, err := os.Open(*SavePath)
	if err != nil {
		panic(err)
	}
	playerFile, err := os.OpenFile("player.json", os.O_CREATE|os.O_RDWR, 0755)
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
	var save Object
	err = json.Unmarshal(line, &save)
	if err != nil {
		panic(err)
	}
	player := save.Field("Character").Field("Player")

	encode := json.NewEncoder(playerFile)
	encode.SetIndent("", "  ")
	err = encode.Encode(player)
	if err != nil {
		panic(err)
	}

	playerFile.Close()
	srcFile.Close()

}

func write() {
	srcFile, err := os.Open(*SavePath)
	if err != nil {
		panic(err)
	}
	playerFile, err := os.OpenFile("player.json", os.O_RDONLY, 0755)
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
	err = json.Unmarshal(line, &save)
	if err != nil {
		panic(err)
	}
	var player Object
	playerDecoder := json.NewDecoder(playerFile)
	err = playerDecoder.Decode(&player)
	if err != nil {
		panic(err)
	}
	// 修改 player
	itemnames := strings.Split(*Items, ",")
	skillnames := strings.Split(*Skills, ",")
	AddItems(player, itemnames)
	AddSkill(player, skillnames)

	save.Field("Character").SetField("Player", player)
	savebin, err := json.Marshal(save)
	if err != nil {
		panic(err)
	}
	_, err = dstFile.Write(savebin)
	if err != nil {
		panic(err)
	}
	srcFile.Close()
	playerFile.Close()
	dstFile.Close()
}

func AddItems(player Object, itemnames []string) {
	oriItems := player.List("Inventory")
	for _, name := range itemnames {
		oriItems = append(oriItems, NewDefaultItem(name))
	}
	player.SetList("Inventory", oriItems)
}

func NewDefaultItem(itemid string) Object {
	return (Object)(map[string]interface{}{
		"Count":           1,
		"Durability":      0,
		"EffectId":        []string{},
		"ForgeMaterials":  struct{}{},
		"Hurt":            struct{}{},
		"HurtDifference":  0,
		"Id":              nil,
		"IsNew":           false,
		"ItemId":          itemid,
		"Level":           0,
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
