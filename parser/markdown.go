package parser

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const markdownTemplate = "# {{.Title}}\n\n" +
	"{{range .Structs}}\n" +
	"## {{.Name}}\n\n" +
	"{{if .Comment}}\n" +
	"{{.Comment}}\n" +
	"{{end}}\n\n" +
	"| Field | Type {{if $.ShowExported}}| Exported {{end}}| Comment |\n" +
	"|-------|------{{if $.ShowExported}}|----------{{end}}|----------|\n" +
	"{{range .Fields}}| {{.Name}} | {{if .IsStruct}}object **{{.RawType}}**{{else}}{{.Type}}{{end}} {{if $.ShowExported}}| {{if .Exported}}Yes{{else}}No{{end}} {{end}}| {{.Comment}} |\n" +
	"{{end}}\n\n" +
	"{{end}}\n"

// MarkdownData 用于模板渲染的数据结构
type MarkdownData struct {
	Title        string
	ShowExported bool
	Structs      []StructInfo
}

// GenerateMarkdown 生成Markdown文档
func (p *Parser) GenerateMarkdown(structs []StructInfo, outputPath string) error {
	// 解析模板
	tmpl, err := template.New("markdown").Parse(markdownTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// 处理注释中的换行
	for i := range structs {
		structs[i].Comment = strings.TrimSpace(structs[i].Comment)
		for j := range structs[i].Fields {
			structs[i].Fields[j].Comment = strings.ReplaceAll(
				strings.TrimSpace(structs[i].Fields[j].Comment),
				"\n",
				" ",
			)
		}
	}

	// 准备模板数据
	data := MarkdownData{
		Title:        p.docOptions.Title,
		ShowExported: p.docOptions.ShowExported,
		Structs:      structs,
	}

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
