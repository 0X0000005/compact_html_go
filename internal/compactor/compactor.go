package compactor

import (
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	// Matches HTML img tags: <img ... src="URL" ...>
	htmlImgRegex = regexp.MustCompile(`(?i)<img[^>]+src=["'](.*?)["']`)
	// Matches Markdown images: ![alt](URL)
	mdImgRegex = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
)

// CompactFile reads an input file (HTML or OS), finds image references, and embeds them as base64.
func CompactFile(inputFile, outputFile string) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputFile))
	isMarkdown := ext == ".md"

	contentStr := string(content)
	
	// Replace in HTML tags
	contentStr = replaceImagesFunc(contentStr, htmlImgRegex, inputFile)

	// Additionally replace Markdown image syntax if it's a markdown file (or just generically attempt)
	if isMarkdown {
		contentStr = replaceImagesFunc(contentStr, mdImgRegex, inputFile)
	}

	// Write output
	if outputFile == "" {
		if ext != "" {
			outputFile = strings.TrimSuffix(inputFile, ext) + ".compact" + ext
		} else {
			outputFile = inputFile + ".compact"
		}
	}

	err = os.WriteFile(outputFile, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	return nil
}

func replaceImagesFunc(content string, re *regexp.Regexp, inputFile string) string {
	inputDir := filepath.Dir(inputFile)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		rawUrl := submatches[1]

		base64Data, err := FetchAndEncodeImage(rawUrl, inputDir)
		if err != nil {
			fmt.Printf("[Warning] Failed to fetch image %s: %v\n", rawUrl, err)
			return match // Keep original if fetching fails
		}

		if base64Data == "" {
			// e.g. already base64
			return match
		}

		// replace the url within the match carefully to preserve the rest of the tag/markdown syntax
		return strings.Replace(match, rawUrl, base64Data, 1)
	})
}

// FetchAndEncodeImage fetches image from local or remote, and returns the data URI scheme representation
func FetchAndEncodeImage(rawUrl, baseDir string) (string, error) {
	if strings.HasPrefix(rawUrl, "data:image/") {
		return "", nil // Already embedded
	}

	// 1. HTML entity decoding (e.g. &amp; -> &)
	decodedUrl := html.UnescapeString(rawUrl)

	// 2. We don't want to over-unescape if it's a web URL, because web URLs need to be encoded.
	// But local paths might have space represented as %20 or something else, we try Unescape.
	// If it's a web URL, parsing it with url.Parse handles it. Let's separate remote and local.

	var imgBytes []byte
	var mimeType string

	isRemote := strings.HasPrefix(decodedUrl, "http://") || strings.HasPrefix(decodedUrl, "https://")

	if isRemote {
		parsedUrl, err := url.Parse(decodedUrl)
		if err != nil {
			return "", err
		}
		imgBytes, mimeType, err = fetchRemoteImage(parsedUrl.String())
		if err != nil {
			return "", err
		}
	} else {
		// Local file
		// Decode URL encoding usually used in markdown for local spaces like "my%20file.png"
		unescapedPath, err := url.QueryUnescape(decodedUrl)
		if err == nil {
			decodedUrl = unescapedPath
		}

		imgPath := decodedUrl
		if !filepath.IsAbs(imgPath) {
			// 先尝试在输入文件所在目录查找
			primaryPath := filepath.Join(baseDir, imgPath)
			if _, err := os.Stat(primaryPath); err == nil {
				imgPath = primaryPath
			} else {
				// 如果没找到，尝试在当前工作目录查找
				if _, err := os.Stat(imgPath); err == nil {
					// imgPath 本身就是相对路径，相对于 CWD
					// 这里什么都不做，直接用 imgPath 去读取
				} else {
					// 如果都找不到，默认使用 primaryPath 以便报错信息更清晰
					imgPath = primaryPath
				}
			}
		}
		
		imgBytes, err = os.ReadFile(imgPath)
		if err != nil {
			return "", err
		}

		ext := strings.ToLower(filepath.Ext(imgPath))
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			// guess by content
			mimeType = http.DetectContentType(imgBytes)
		}
	}

	if imgBytes == nil || len(imgBytes) == 0 {
		return "", fmt.Errorf("empty image data")
	}

	// Extract standard mime like image/png, etc without charset
	mimeType = strings.Split(mimeType, ";")[0]
	if !strings.HasPrefix(mimeType, "image/") {
		// Just force it for safety, though DetectContentType might return application/octet-stream
		mimeType = "image/jpeg" 
	}

	b64 := base64.StdEncoding.EncodeToString(imgBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), nil
}

func fetchRemoteImage(urlStr string) ([]byte, string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("bad status: %s", resp.Status)
	}

	mimeType := resp.Header.Get("Content-Type")

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	return data, mimeType, nil
}
