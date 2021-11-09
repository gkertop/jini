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
	ini.Set("SectionA", "KeyA", "ValueA")
	v1, ok := ini.Get("SectionA", "KeyA", "DefaultValueA")
	fmt.Printf("A value:%v,has:%v\n", v1, ok)
	v2, ok := ini.Get("SectionB", "KeyB", "DefaultValueB")
	fmt.Printf("B value:%v,has:%v\n", v2, ok)
	//ini.SaveTo("test.ini")
	ini.Save()
```

```
A value:ValueA,has:true
B value:DefaultValueB,has:false
```