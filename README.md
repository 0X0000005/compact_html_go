# compact_html_go

使用 Go 语言重写并强化的图片 Base64 内嵌辅助工具。
此程序可将 HTML 或 Markdown 文件中的图片链接（包括本地文件与远程网络连接）统一抓取并转化为 Base64 内嵌格式。特别优化了路径解析，能够从容支持含有特殊符号、URL 编码空格 `%20`，以及因网页转义产生的 `&amp;` 的各类图片路径。

## 特色与优势
- **更强的鲁棒性**：完美兼容 URL 编码与 HTML ASCII 转义编码，解决 C++ 原版读取路径极易奔溃的问题。
- **智能相对路径解析**：双重 fallback 机制，首先基于输入文件所在目录查找图片，若未找到则自动退回至当前工作目录查找，完美兼顾标准网页结构和跨目录调用的场景。
- **跨平台与极简部署**：依托 Go 跨平台编译，借助 `upx` 进行单文件压缩，开箱即用，无需任何多余的动态链接库。

## 如何使用

可以直接在命令行后直接加上被处理文件：
```bash
./compact_html sample.html
```

或处理 markdown 文档并将输出保存为特定文件名：
```bash
./compact_html -o Offline_README.md README.md
```

## 构建方法

Windows系统请执行 `build.bat`。  
Linux/macOS 系统请执行 `build.sh`。  

生成的 `compact_html.exe` 即为立等可用的二进制程序。

## 致谢

- 特别感谢原项目 [EdenHell/compact_html](https://github.com/EdenHell/compact_html) 提供的基础思路与实现启发。
