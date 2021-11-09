package jini

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Ini struct {
	fileName  string
	lineBreak string
	section   []*section
}

type section struct {
	key  string
	item []*keyValue
}

type keyValue struct {
	key   string
	value string
}

func NewIni(pFileName string, pCreate bool) (rt *Ini, err error) {
	dataBytes, err := readFile(pFileName, pCreate)
	if err != nil {
		return
	}
	rt = &Ini{fileName: pFileName, lineBreak: getLineBreak()}
	rt.section = append(rt.section, &section{})
	if len(dataBytes) <= 0 {
		return
	}
	dataLines := bytes.Split(dataBytes, []byte("\n"))
	sectionIndex := 0
	itemIndex := -1
	crNum := 0
	for _, line := range dataLines {
		lineStr := string(line)
		if strings.HasSuffix(lineStr, "\r") {
			crNum++
		}
		formatStr := strings.TrimSpace(lineStr)
		if strings.HasPrefix(formatStr, ";") {
			kv := &keyValue{}
			kv.key = ";"
			kv.value = formatStr[1:]
			rt.section[sectionIndex].item = append(rt.section[sectionIndex].item, kv)
			itemIndex++
		} else if strings.HasPrefix(formatStr, "[") && strings.HasSuffix(formatStr, "]") {
			s := &section{}
			s.key = formatStr[1 : len(formatStr)-1]
			rt.section = append(rt.section, s)
			sectionIndex++
			itemIndex = -1
		} else if flagIndex := strings.Index(lineStr, "="); flagIndex > 0 && flagIndex < len(lineStr)-1 {
			kv := &keyValue{}
			kv.key = strings.TrimSpace(lineStr[0:flagIndex])
			kv.value = lineStr[flagIndex+1:]
			rt.section[sectionIndex].item = append(rt.section[sectionIndex].item, kv)
			itemIndex++
		} else if itemIndex < 0 {
			kv := &keyValue{}
			kv.key = ";;"
			kv.value = lineStr
			rt.section[sectionIndex].item = append(rt.section[sectionIndex].item, kv)
			itemIndex++
		} else {
			rt.section[sectionIndex].item[itemIndex].value += "\n" + lineStr
		}
	}
	if crNum > 0 && crNum >= (len(dataLines)-1)/2 {
		rt.lineBreak = "\r\n"
	} else {
		rt.lineBreak = "\n"
	}
	return
}

func readFile(pFileName string, pCreate bool) (rt []byte, err error) {
	flag := os.O_RDONLY
	if pCreate {
		flag = flag | os.O_CREATE
	}
	fileObj, err := os.OpenFile(pFileName, flag, 0666)
	if err != nil {
		return
	}
	defer func(fileObj *os.File) {
		err := fileObj.Close()
		if err != nil {
			panic(err.Error())
		}
	}(fileObj)
	rt, err = ioutil.ReadAll(fileObj)
	return
}

func (i *Ini) GetFileName() string {
	return i.fileName
}

func (i *Ini) Set(pSection string, pKey string, pValue string) {
	s := (*section)(nil)
	for _, s1 := range i.section {
		if s1.key == pSection {
			s = s1
			break
		}
	}
	if s == nil {
		s = &section{key: pSection}
		i.section = append(i.section, s)
	}
	if s.item != nil {
		for _, item := range s.item {
			if item.key == pKey {
				item.value = pValue
				return
			}
		}
	}
	kv := &keyValue{key: pKey}
	kv.value = pValue
	s.item = append(s.item, kv)
}

func (i *Ini) Get(pSection string, pKey string, pValue string) (string, bool) {
	for _, s1 := range i.section {
		if s1.key == pSection {
			if s1.item == nil {
				continue
			}
			for _, item := range s1.item {
				if item.key == pKey {
					return item.value, true
				}
			}
		}
	}
	return pValue, false
}

func (i *Ini) SaveTo(pFileName string) (err error) {
	var dataObj strings.Builder
	for _, s1 := range i.section {
		if s1.key != "" {
			dataObj.WriteString("[" + s1.key + "]" + i.lineBreak)
		}
		if s1.item == nil {
			continue
		}
		for _, item := range s1.item {
			if item.key == "" {
				continue
			}
			if item.key == ";" {
				dataObj.WriteString(item.key + item.value + i.lineBreak)
			} else if item.key == ";;" {
				dataObj.WriteString(item.value + i.lineBreak)
			} else {
				dataObj.WriteString(item.key + "=" + item.value + i.lineBreak)
			}
		}
	}
	err = i.writeFile(pFileName, dataObj.String())
	return
}

func (i *Ini) Save() (err error) {
	err = i.SaveTo(i.fileName)
	return
}

func (i *Ini) writeFile(pFileName string, pData string) (err error) {
	fileObj, err := os.OpenFile(pFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer func(fileObj *os.File) {
		err := fileObj.Close()
		if err != nil {
			panic(err.Error())
		}
	}(fileObj)
	_, err = fileObj.WriteString(pData)
	return
}
