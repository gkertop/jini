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

// NewIni
// @Description:加载INI文件
// @param pFileName 文件名
// @param pCreate 当文件不存在时,是否创建文件
// @return rt
// @return err
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
	if bytes.HasSuffix(dataLines[0], []byte("\r")) {
		rt.lineBreak = "\r\n"
	}
	for _, line := range dataLines {
		lineStr := string(line)
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
		} else if flagIndex := strings.Index(formatStr, "="); flagIndex > 0 && flagIndex < len(formatStr)-1 {
			flagIndex = strings.Index(lineStr, "=")
			kv := &keyValue{}
			kv.key = strings.TrimSpace(lineStr[:flagIndex])
			if strings.HasSuffix(lineStr, "\r") {
				kv.value = lineStr[flagIndex+1 : len(lineStr)-1]
			} else {
				kv.value = lineStr[flagIndex+1:]
			}
			rt.section[sectionIndex].item = append(rt.section[sectionIndex].item, kv)
			itemIndex++
		} else if itemIndex < 0 {
			kv := &keyValue{}
			kv.key = ";;"
			if strings.HasSuffix(lineStr, "\r") {
				kv.value = lineStr[:len(lineStr)-1]
			} else {
				kv.value = lineStr
			}
			rt.section[sectionIndex].item = append(rt.section[sectionIndex].item, kv)
			itemIndex++
		} else {
			if strings.HasSuffix(lineStr, "\r") {
				rt.section[sectionIndex].item[itemIndex].value += rt.lineBreak + lineStr[:len(lineStr)-1]
			} else {
				rt.section[sectionIndex].item[itemIndex].value += rt.lineBreak + lineStr
			}
		}
	}
	return
}

// GetFileName
// @Description: 获取当前文件名
// @receiver i
// @return string
func (i *Ini) GetFileName() string {
	return i.fileName
}

// HasSection
// @Description: 查询节点是否存在
// @receiver i
// @param pSection 节点名
// @return bool
func (i *Ini) HasSection(pSection string) bool {
	_, s := i.findSection(pSection)
	if s == nil {
		return false
	}
	return true
}

// HasKey
// @Description: 查询键是否存在
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @return bool
func (i *Ini) HasKey(pSection string, pKey string) bool {
	_, s := i.findSection(pSection)
	if s == nil {
		return false
	}
	_, item := i.findItem(s, pKey)
	if item == nil {
		return false
	}
	return true
}

// DelSection
// @Description: 删除节点
// @receiver i
// @param pSection 节点名
func (i *Ini) DelSection(pSection string) {
	i1, s := i.findSection(pSection)
	if s == nil {
		return
	}
	i.section = append(i.section[:i1], i.section[i1+1:]...)
	return
}

// DelKey
// @Description: 删除键
// @receiver i
// @param pSection 节点名
// @param pKey 键名
func (i *Ini) DelKey(pSection string, pKey string) {
	i1, s := i.findSection(pSection)
	if s == nil {
		return
	}
	i2, item := i.findItem(s, pKey)
	if item == nil {
		return
	}
	i.section[i1].item = append(i.section[i1].item[:i2], i.section[i1].item[i2+1:]...)
	return
}

// Set
// @Description: 写入配置项
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @param pValue 键值
func (i *Ini) Set(pSection string, pKey string, pValue string) {
	_, s := i.findSection(pSection)
	if s == nil {
		s = &section{key: pSection}
		i.section = append(i.section, s)
	}
	if s.item != nil {
		_, item := i.findItem(s, pKey)
		if item != nil {
			item.value = pValue
			return
		}
	}
	kv := &keyValue{key: pKey}
	kv.value = pValue
	s.item = append(s.item, kv)
}

// Get
// @Description: 读取配置项
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @param pValue 默认值(当配置项不存在时,返回该值)
// @return string
// @return bool 配置项是否存在
func (i *Ini) Get(pSection string, pKey string, pValue string) (string, bool) {
	_, s := i.findSection(pSection)
	if s != nil && s.item != nil {
		_, item := i.findItem(s, pKey)
		if item != nil {
			return item.value, true
		}
	}
	return pValue, false
}

// SaveTo
// @Description: 保存到指定文件
// @receiver i
// @param pFileName 文件名
// @return err
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

// Save
// @Description: 保存到当前文件
// @receiver i
// @return err
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

func (i *Ini) findSection(pSection string) (int, *section) {
	for i1, s1 := range i.section {
		if s1.key == pSection {
			return i1, s1
		}
	}
	return -1, nil
}

func (i *Ini) findItem(pSectionObj *section, pKey string) (int, *keyValue) {
	for i1, item := range pSectionObj.item {
		if item.key == pKey {
			return i1, item
		}
	}
	return -1, nil
}
