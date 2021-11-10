# jini

### 简介

Golang轻量INI操作包

* 轻量
* 有序
* 增/删/改/查
* 支持顶级节点
* 键值支持多行
* 尽可能不破坏原格式

### 功能

```
// NewIni
// @Description:加载INI文件
// @param pFileName 文件名
// @param pCreate 当文件不存在时,是否创建文件
// @return rt
// @return err
func NewIni(pFileName string, pCreate bool) (rt *Ini, err error)

// GetFileName
// @Description: 获取当前文件名
// @receiver i
// @return string
func (i *Ini) GetFileName() string

// HasSection
// @Description: 查询节点是否存在
// @receiver i
// @param pSection 节点名
// @return bool
func (i *Ini) HasSection(pSection string) bool

// HasKey
// @Description: 查询键是否存在
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @return bool
func (i *Ini) HasKey(pSection string, pKey string) bool

// DelSection
// @Description: 删除节点
// @receiver i
// @param pSection 节点名
func (i *Ini) DelSection(pSection string)

// DelKey
// @Description: 删除键
// @receiver i
// @param pSection 节点名
// @param pKey 键名
func (i *Ini) DelKey(pSection string, pKey string)

// Set
// @Description: 写入配置项
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @param pValue 键值
func (i *Ini) Set(pSection string, pKey string, pValue string)

// Get
// @Description: 读取配置项
// @receiver i
// @param pSection 节点名
// @param pKey 键名
// @param pValue 默认值(当配置项不存在时,返回该值)
// @return string
// @return bool 配置项是否存在
func (i *Ini) Get(pSection string, pKey string, pValue string) (string, bool)

// SaveTo
// @Description: 保存到指定文件
// @receiver i
// @param pFileName 文件名
// @return err
func (i *Ini) SaveTo(pFileName string) (err error)

// Save
// @Description: 保存到当前文件
// @receiver i
// @return err
func (i *Ini) Save() (err error)
```

### 安装

```
go get github.com/gkertop/jini
```

### 示例

```
	ini, err := jini.NewIni("test.ini", true)
	if err != nil {
		panic(err.Error())
	}

	v1, ok := ini.Get("SectionA", "KeyA", "DefaultValueA")
	fmt.Printf("v1: value:%v has:%v\n", v1, ok)

	ini.Set("SectionB", "KeyB", "ValueB")

	v2, ok := ini.Get("SectionB", "KeyB", "DefaultValueB")
	fmt.Printf("v2: value:%v has:%v\n", v2, ok)

	err = ini.Save()
	if err != nil {
		panic(err.Error())
	}
```

```
v1: value:DefaultValueA has:false
v2: value:ValueB has:true
```