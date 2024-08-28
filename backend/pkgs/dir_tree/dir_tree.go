// Package dir_tree provides a way to generate a directory tree.
//
// Example usage:
//
//	tree, err := directory_tree.NewTree("/home/me")
//
// I did my best to keep it OS-independent but truth be told I only tested it
// on OS X and Debian Linux so YMMV. You've been warned.
package dir_tree

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
}

// Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.DirEntry) *FileInfo {
	fi, _ := v.Info()
	if v.IsDir() {
		return &FileInfo{v.Name(), 0, v.Type(), fi.ModTime(), v.IsDir()}
	}
	return &FileInfo{v.Name(), fi.Size(), v.Type(), fi.ModTime(), v.IsDir()}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}
type Node1 struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	IsLeaf   bool
}

func NewDirList(root string) ([]*Node1, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node1, 0)
	ds, err := os.ReadDir(absRoot)
	if err != nil {
		return nil, err
	}
	for _, file := range ds {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		node := &Node1{
			FullPath: "",
			Info:     nil,
			IsLeaf:   false,
		}
		node.FullPath = filepath.Join(absRoot, file.Name())
		node.Info = fileInfoFromInterface(file)
		if file.IsDir() {
			node.IsLeaf = false
		} else {
			node.IsLeaf = true
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// Create directory hierarchy.
func NewTree(root string) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		parents[path] = &Node{
			FullPath: path,
			Info:     fileInfoFromInterface(info),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.WalkDir(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
