package com

type SubDir string

const (
	SubDirAndroid SubDir = "android"
	SubDirIos     SubDir = "ios"
)

func (s SubDir) String() string {
	return string(s)
}
