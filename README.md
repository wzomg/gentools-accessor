### Description

> 基于GoAST为结构体自动生成Get和Set方法



### 工具安装

```shell
go install github.com/wzomg/gentools-accessor
```

### 用法示例

```go
package main

//go:generate gentools-accessor

type Student struct {
	Name      string `access:"r,w"`
	Age       int    `access:"r"`
	signature string `access:"w"`
	id        int
}
```

