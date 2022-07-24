### Description

> 基于GoAST为结构体自动生成Get和Set方法



### 工具安装

```shell
go install github.com/wzomg/gentools-accessor@v0.0.2
```

### 用法示例

法1：将`//go:generate gentools-accessor`写在待解析的文件，用goland或vscode提供的标志执行即可

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

<img src="./img/goland_exec.png" width="50%" alt="goland执行" /><img src="./img/vscode_exec.png" width="50%" alt="vscode执行" />

法2：命令行执行：`gentools-accessor -file=文件名`（支持相对路径和绝对路径）
<img src="./img/absolute_path.png" width="50%" alt="相对路径" /><img src="./img/relative_path.png" width="50%" alt="绝对路径" />


