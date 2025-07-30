package util

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

// TrieNode 定义前缀树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// WordFilter 敏感词过滤器
type WordFilter struct {
	root *TrieNode
}

// 构造函数
func NewWordFilter() *WordFilter {
	var value = &WordFilter{root: &TrieNode{children: make(map[rune]*TrieNode)}}
	_ = value.LoadWordsFromFile("keywords.txt")
	return value

}

// 添加敏感词，统一转为小写
func (wf *WordFilter) AddWord(word string) {
	word = strings.ToLower(word)
	node := wf.root
	for _, r := range word {
		r = unicode.ToLower(r)
		if node.children[r] == nil {
			node.children[r] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[r]
	}
	node.isEnd = true
}

// 从文件加载敏感词
func (wf *WordFilter) LoadWordsFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wf.AddWord(word)
		}
	}
	return scanner.Err()
}

// 替换敏感词
func (wf *WordFilter) Replace(text string) string {
	runes := []rune(text)
	length := len(runes)
	result := []rune{}
	i := 0

	for i < length {
		node := wf.root
		j := i
		matchLen := 0

		for j < length {
			c := unicode.ToLower(runes[j])
			child, ok := node.children[c]
			if !ok {
				break
			}
			node = child
			if node.isEnd {
				matchLen = j - i + 1
			}
			j++
		}

		if matchLen > 0 {
			// 替换敏感词为 ***
			result = append(result, []rune("***")...)
			i += matchLen
		} else {
			result = append(result, runes[i])
			i++
		}
	}
	return string(result)
}
