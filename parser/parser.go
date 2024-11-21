package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// StructField 表示结构体字段信息
type StructField struct {
	Name       string
	Type       string
	RawType    string // 原始类型名（用于链接）
	IsStruct   bool   // 是否为结构体类型
	Exported   bool
	Comment    string
	Tags       map[string]string
}

// StructInfo 表示结构体信息
type StructInfo struct {
	Name     string
	Comment  string
	Fields   []StructField
	Exported bool
}

// DocOptions 文档生成选项
type DocOptions struct {
	Title        string
	ShowExported bool
}

// Parser 结构体解析器
type Parser struct {
	fset         *token.FileSet
	fieldNameTag string // 用于指定字段名来源的标签
	docOptions  DocOptions
}

// NewParser 创建新的解析器实例
func NewParser() *Parser {
	return &Parser{
		fset: token.NewFileSet(),
		docOptions: DocOptions{
			Title:        "Go Structs Documentation",
			ShowExported: true,
		},
	}
}

// SetFieldNameTag 设置用于字段名的标签
func (p *Parser) SetFieldNameTag(tag string) {
	p.fieldNameTag = tag
}

// SetDocOptions 设置文档选项
func (p *Parser) SetDocOptions(opts DocOptions) {
	p.docOptions = opts
}

// parseStructTag 解析结构体标签
func (p *Parser) parseStructTag(tag string) map[string]string {
	tags := make(map[string]string)
	if tag == "" {
		return tags
	}

	structTag := reflect.StructTag(strings.Trim(tag, "`"))
	for _, key := range []string{"json", "xml", "yaml", "db"} {
		if value, ok := structTag.Lookup(key); ok {
			// 处理类似 json:"name,omitempty" 的情况
			parts := strings.Split(value, ",")
			tags[key] = parts[0]
		}
	}
	return tags
}

// getFieldName 根据设置获取字段名
func (p *Parser) getFieldName(field *ast.Field, origName string) string {
	if p.fieldNameTag == "" {
		return origName
	}

	if field.Tag != nil {
		tags := p.parseStructTag(field.Tag.Value)
		if tagValue, ok := tags[p.fieldNameTag]; ok && tagValue != "" {
			return tagValue
		}
	}
	return origName
}

// getTypeInfo 获取类型信息
func (p *Parser) getTypeInfo(expr ast.Expr) (typeName string, rawType string, isStruct bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, t.Name, t.Obj != nil && t.Obj.Kind == ast.Typ
	case *ast.StarExpr:
		baseType, rawT, isS := p.getTypeInfo(t.X)
		if isS {
			return "*" + baseType, rawT, true
		}
		return "*" + baseType, rawT, false
	case *ast.ArrayType:
		baseType, rawT, isS := p.getTypeInfo(t.Elt)
		return "[]" + baseType, rawT, isS
	case *ast.MapType:
		keyType, _, _ := p.getTypeInfo(t.Key)
		valueType, rawT, isS := p.getTypeInfo(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType), rawT, isS
	case *ast.SelectorExpr:
		pkg, _, _ := p.getTypeInfo(t.X)
		return fmt.Sprintf("%s.%s", pkg, t.Sel.Name), t.Sel.Name, false
	default:
		return fmt.Sprintf("%T", expr), "", false
	}
}

// ParsePackage 解析指定路径下的Go包
func (p *Parser) ParsePackage(path string) ([]StructInfo, error) {
	var structs []StructInfo

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// 解析Go源文件
		file, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %v", path, err)
		}

		// 遍历AST
		ast.Inspect(file, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			// 只处理导出的结构体
			if !typeSpec.Name.IsExported() {
				return true
			}

			structInfo := StructInfo{
				Name:     typeSpec.Name.Name,
				Exported: true,
				Fields:   make([]StructField, 0),
			}

			// 获取结构体的注释
			if typeSpec.Doc != nil {
				structInfo.Comment = typeSpec.Doc.Text()
			}

			// 解析字段
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					typeName, rawType, isStruct := p.getTypeInfo(field.Type)
					fieldInfo := StructField{
						Name:       p.getFieldName(field, name.Name),
						Type:       typeName,
						RawType:    rawType,
						IsStruct:   isStruct,
						Exported:   name.IsExported(),
						Tags:       make(map[string]string),
					}

					if field.Tag != nil {
						fieldInfo.Tags = p.parseStructTag(field.Tag.Value)
					}

					if field.Comment != nil {
						fieldInfo.Comment = strings.TrimSpace(field.Comment.Text())
					}

					structInfo.Fields = append(structInfo.Fields, fieldInfo)
				}
			}

			structs = append(structs, structInfo)
			return true
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return structs, nil
}
