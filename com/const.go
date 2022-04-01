package com

import "path/filepath"

type SubDir string

const (
	SubDirGoogle SubDir = "google"
	SubDirApple  SubDir = "apple"
)

func (s SubDir) String() string {
	return string(s)
}

var SubDirSlice []SubDir

func init() {
	SubDirSlice = []SubDir{
		SubDirGoogle,
		SubDirApple,
	}
}

type Workdir string

// 返回工作路径下面一个完整的子路径
func (w Workdir) SubPath(sub SubDir) string {
	return filepath.Join(w.String(), sub.String())
}
func (w Workdir) String() string {
	return string(w)
}
