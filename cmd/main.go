package main

import (
	"flag"
	"fmt"
	"os"

	"compact_html_go/internal/compactor"
)

var version = "1.0.0"

func main() {
	outputFile := flag.String("o", "", "输出文件 (可选，默认在同目录生成)")
	showVersion := flag.Bool("v", false, "显示版本号")
	showHelp := flag.Bool("h", false, "显示帮助信息")

	flag.Parse()

	if *showVersion {
		fmt.Printf("compact_html_go v%s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("错误：必须指定输入文件")
		printHelp()
		os.Exit(1)
	}

	inputFile := args[0]

	err := compactor.CompactFile(inputFile, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "压实文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("处理成功!")
}

func printHelp() {
	fmt.Printf("\ncompact_html_go v%s - 将图片内嵌为 base64 的强大工具\n\n", version)
	fmt.Println("用法: compact_html_go [选项] <输入文件>")
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  <输入文件>     包含图片路径的输入文件，支持 .html 与 .md 文件 (必填)")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -o <文件路径>  指定输出文件 (可选，如果不填将输出到 input.compact.html/md)")
	fmt.Println("  -v             显示版本号")
	fmt.Println("  -h             显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  compact_html_go index.html")
	fmt.Println("  compact_html_go -o readme_offline.md readme.md")
}
