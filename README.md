# jini

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
	ini.Set("Section", "Key", "Value")
	v1, ok := ini.Get("SectionA", "KeyA", "DefaultValue")
	fmt.Printf("A value:%v,has:%v\n", v1, ok)
	v2, ok := ini.Get("SectionB", "KeyB", "DefaultValue")
	fmt.Printf("B value:%v,has:%v\n", v2, ok)
	//ini.SaveTo("test.ini")
	ini.Save()
```