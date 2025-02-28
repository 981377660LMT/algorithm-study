// 588. 设计内存文件系统
// https://leetcode.cn/problems/design-in-memory-file-system/solutions/2388114/she-ji-nei-cun-wen-jian-xi-tong-by-leetc-0wia/
//
// 1 <= path.length, filePath.length <= 100
// 1 <= content.length <= 50
// ls, mkdir, addContentToFile, readContentFromFile 最多被调用 300 次
// !path 和 filePath 都是绝对路径，除非是根目录 ‘/’ 自身，其他路径都是以 ‘/’ 开头且 不 以 ‘/’ 结束。
// 你可以假定所有操作的参数都是有效的，即用户不会获取不存在文件的内容，或者获取不存在文件夹和文件的列表。
// 你可以假定所有文件夹名字和文件名字都只包含小写字母，且同一文件夹下不会有相同名字的文件夹或文件。
// 你可以假定 addContentToFile 中的文件的父目录都存在。
//
// 方法 1：使用单独的目录和文件列表(文件和目录的结构不同)
// !方法 2：使用统一的目录和文件列表(更好)

package main

import (
	"sort"
	"strings"
)

func main() {

}

// File 表示一个文件系统中的节点，可以是目录或者文件。
type File struct {
	isFile   bool
	content  string
	children map[string]*File
}

func NewFile() *File {
	return &File{children: map[string]*File{}}
}

type FileSystem struct {
	root *File
}

func Constructor() FileSystem {
	return FileSystem{root: NewFile()}
}

// Ls 列出 path 路径下的内容：
// - 如果 path 指向一个文件，则返回只包含该文件名称的列表。
// - 如果 path 指向目录，则返回该目录下所有子目录和文件的名称（字典顺序排列）。
func (fs *FileSystem) Ls(path string) []string {
	node, _ := fs.traverse(path)
	if node.isFile {
		return []string{path}
	}
	res := []string{}
	for name := range node.children {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

func (fs *FileSystem) Mkdir(path string) {
	fs.traverse(path)
}

// AddContentToFile 将内容添加到 filePath 文件中：
// - 如果文件不存在，则先创建文件节点，再写入内容；
// - 如果文件存在，则将内容追加到原有内容后面。
func (fs *FileSystem) AddContentToFile(filePath string, content string) {
	node, _ := fs.traverse(filePath)
	node.isFile = true
	node.content += content
}

func (fs *FileSystem) ReadContentFromFile(filePath string) string {
	node, _ := fs.traverse(filePath)
	if node == nil || !node.isFile {
		return ""
	}
	return node.content
}

// traverse 根据 path 分割逐级查找或创建节点，返回最后一个节点以及其名称（用于 ls 返回文件名称）.
// 如果 path 为 "/" 则返回 root 节点.
// 如果子节点存在则直接进入，否则创建一个新节点（默认创建目录节点）.
func (fs *FileSystem) traverse(path string) (*File, string) {
	if path == "/" {
		return fs.root, ""
	}
	parts := strings.Split(path, "/")
	node := fs.root
	name := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		name = part
		if child, ok := node.children[part]; ok {
			node = child
		} else {
			child := NewFile()
			node.children[part] = child
			node = child
		}
	}
	return node, name
}
