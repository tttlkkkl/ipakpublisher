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

// Track 发布频道
type Track string

const (
	TrackAlpha      Track = "alpha"
	TrackBeta       Track = "beta"
	TrackInternal   Track = "internal"
	TrackProduction Track = "production"
)

func (t Track) String() string {
	return string(t)
}

// Status release 的状态
type Status string

const (
	// StatusUnspecified 未定义的状态
	StatusUnspecified Status = "statusUnspecified"
	// StatusDraft 草稿
	StatusDraft Status = "draft"
	// StatusInProgress 根据用户评分，向一部分人发布
	StatusInProgress Status = "inProgress"
	// StatusHalted 暂停，这个版本将不会分发给用户
	StatusHalted Status = "halted"
	// StatusCompleted 版本已完成，用户可见
	StatusCompleted Status = "completed"
)

func (s Status) String() string {
	return string(s)
}
