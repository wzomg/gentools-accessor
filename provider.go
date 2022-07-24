package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
)

const AccessRead = "r"
const AccessWrite = "w"
const AccessTagName = "access"

var (
	fileName = flag.String("file", "", "a parsed filename; must be set")
)

func main() {
	log.SetFlags(0) // 设置日志的抬头信息
	log.SetPrefix("gentools-accessor: ")
	flag.Parse() //把用户传递的命令行参数解析为对应变量的值

	var inputName string
	if len(*fileName) > 0 {
		inputName = *fileName
	} else {
		//尝试获取当前被执行的文件
		inputName = os.Getenv("GOFILE")
	}

	if len(inputName) == 0 {
		log.Fatalf("请输入一个正确的待解析的文件路径！")
		return
	}
	log.Println("当前文件名为：", inputName)

	g := Generator{
		buf: bytes.NewBufferString(""),
	}
	g.generate(inputName)
	var src = (g.buf).Bytes()
	outputName := strings.TrimSuffix(inputName, ".go") + ".accessor.go"
	err := ioutil.WriteFile(outputName, src, 0644) // ignore_security_alert
	if err != nil {
		log.Fatalf("writing output: %s\n", err)
		return
	}
	log.Printf("\n\tautomatic code generation finished!\n\toutput_name: %s\n", outputName)
}

type StructFieldInfo struct {
	Name   string   // 字段名
	Type   string   // 类型名
	Access []string // tag对应的value，即为r，w
}

type StructFieldInfoArr = []StructFieldInfo

type Generator struct {
	buf *bytes.Buffer // Accumulated output.
}

func (g *Generator) myPrintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) generate(fileName string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil || f == nil {
		log.Fatalln("解析文件失败")
		return
	}
	structInfo, err := parseAllStructInSingleFile(f, fset, AccessTagName)
	if err != nil {
		log.Fatalln("解析文件中结构体失败: ", err)
		return
	}
	g.myPrintf("// Code generated by \"gentools-accessor\"; DO NOT EDIT.\n")
	g.myPrintf("\n")
	g.myPrintf("package %s\n", f.Name)
	g.myPrintf("\n")
	for stName, info := range structInfo {
		for _, field := range info {
			for _, access := range field.Access {
				switch access {
				case AccessWrite:
					g.myPrintf("%s\n", genSetter(stName, field.Name, field.Type))
				case AccessRead:
					g.myPrintf("%s\n", genGetter(stName, field.Name, field.Type))
				}
			}
		}
	}
}

func parseAllStructInSingleFile(file *ast.File, fileSet *token.FileSet, tagName string) (structMap map[string]StructFieldInfoArr, err error) {
	structMap = make(map[string]StructFieldInfoArr)

	collectStructs := func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		// 获取结构体名称
		structName := ts.Name.Name

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		fileInfos := make([]StructFieldInfo, 0)
		for _, field := range s.Fields.List {
			name := field.Names[0].Name // 字段名称
			info := StructFieldInfo{Name: name}
			var typeNameBuf bytes.Buffer
			err := printer.Fprint(&typeNameBuf, fileSet, field.Type)
			if err != nil {
				log.Println("获取字段类型失败:", err)
				return true
			}
			info.Type = typeNameBuf.String()
			if field.Tag != nil { // 有tag
				tag := field.Tag.Value
				tag = strings.Trim(tag, "`")
				tags, err := structtag.Parse(tag)
				if err != nil {
					return true
				}
				access, err := tags.Get(tagName)
				if err == nil {
					access.Options = append(access.Options, access.Name)
					for i, v := range access.Options {
						if v == AccessRead || v == AccessWrite {
							continue
						}
						// 剔除除 r,w 之外的 tag-value
						access.Options = append(access.Options[:i], access.Options[i+1:]...)
					}
				}
				info.Access = access.Options
			} else {
				firstChar := name[0:1]
				if strings.ToUpper(firstChar) == firstChar { //大写，封装对外可读可写
					info.Access = []string{AccessRead, AccessWrite}
				} else { // 小写，只封装可读方法，字段名的首字母也不改成大写
					info.Access = []string{AccessRead}
				}
			}
			fileInfos = append(fileInfos, info)
		}
		structMap[structName] = fileInfos
		return false
	}

	ast.Inspect(file, collectStructs)

	return structMap, nil
}

func genSetter(structName, fieldName, typeName string) string {
	tpl := `func ({{.Receiver}} *{{.Struct}}) Set{{.Field}}(param {{.Type}}) {
	{{.Receiver}}.{{.Field}} = param
}`
	t := template.New("setter")
	t = template.Must(t.Parse(tpl))
	res := bytes.NewBufferString("")
	_ = t.Execute(res, map[string]string{
		"Receiver": strings.ToLower(structName[0:1]),
		"Struct":   structName,
		"Field":    fieldName,
		"Type":     typeName,
	})
	return res.String()
}

func genGetter(structName, fieldName, typeName string) string {
	tpl := `func ({{.Receiver}} *{{.Struct}}) Get{{.Field}}() (v0 {{.Type}}) {
	if {{.Receiver}} == nil {
		return
	}
	return {{.Receiver}}.{{.Field}}
}`
	t := template.New("getter")
	t = template.Must(t.Parse(tpl))
	res := bytes.NewBufferString("")
	_ = t.Execute(res, map[string]string{
		"Receiver": strings.ToLower(structName[0:1]),
		"Struct":   structName,
		"Field":    fieldName,
		"Type":     typeName,
	})
	return res.String()
}
