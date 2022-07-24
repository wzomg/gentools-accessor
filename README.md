### Description

> 基于GoAST为结构体自动生成Get和Set方法



### 工具安装

```shell
go install github.com/wzomg/gentools-accessor@v0.0.2
```

### 用法示例

法1：将`//go:generate gentools-accessor`写入待解析的文件中，用goland或vscode提供的标志执行即可。

提供了tag标识：`access`，对应tag-value只识别`r`（Getter）、`w`（Setter）。不写这个tag，解析语法树时都带上，默认都解析。

```go
package main

//go:generate gentools-accessor

type Student struct {
	Name      string `access:"r,w"`
	Age       int    `access:"w"`
	signature string `access:"r"`
	id        int
}
```
注意：若不加tag:`access`，首字母小写的字段，默认只提供Getter方法，并且对应的方法名为如`Getid`，非`GetId`。这样处理是为了避免struct里有两个字段名几乎一模一样，仅仅因一个字符大小写的区别而产生冲突！

<img src="./img/goland_exec.png" width="50%" alt="goland执行" /><img src="./img/vscode_exec.png" width="50%" alt="vscode执行" />

法2：命令行执行：`gentools-accessor -file=文件名`（支持相对路径和绝对路径）
<img src="./img/absolute_path.png" width="50%" alt="相对路径" /><img src="./img/relative_path.png" width="50%" alt="绝对路径" />


