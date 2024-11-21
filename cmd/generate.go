package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"godocs/parser"
)

var (
	sourcePath   string
	outputPath   string
	fieldNameTag string
	docTitle     string
	showExported bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate documentation for Go structs",
	Long: `Generate Markdown documentation for Go structs in the specified package.
The documentation includes struct definitions, field information, and comments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sourcePath == "" {
			return fmt.Errorf("please provide a source path using --path flag")
		}

		// 解析绝对路径
		absPath, err := filepath.Abs(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}

		// 创建解析器实例
		p := parser.NewParser()
		if fieldNameTag != "" {
			p.SetFieldNameTag(fieldNameTag)
		}

		// 设置文档选项
		p.SetDocOptions(parser.DocOptions{
			Title:        docTitle,
			ShowExported: showExported,
		})

		// 解析包
		structs, err := p.ParsePackage(absPath)
		if err != nil {
			return fmt.Errorf("failed to parse package: %v", err)
		}

		// 生成文档
		if err := p.GenerateMarkdown(structs, outputPath); err != nil {
			return fmt.Errorf("failed to generate documentation: %v", err)
		}

		log.Printf("Documentation generated successfully at: %s\n", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// 添加命令行标志
	generateCmd.Flags().StringVarP(&sourcePath, "path", "p", "", "Path to the Go package to generate documentation for")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "docs.md", "Output markdown file path")
	generateCmd.Flags().StringVarP(&fieldNameTag, "field-tag", "t", "", "Tag to use for field names (e.g., json)")
	generateCmd.Flags().StringVarP(&docTitle, "title", "", "Go Structs Documentation", "Documentation title")
	generateCmd.Flags().BoolVarP(&showExported, "show-exported", "e", true, "Show exported field information")

	// 标记必需的标志
	generateCmd.MarkFlagRequired("path")
}
