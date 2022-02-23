package main

import (
	"sync"
	"time"
)

type Statistic struct {
	Start     time.Time
	Files     int64
	FileSkips int64
	Dirs      int64
	DirSkips  int64
	*sync.RWMutex
}

func (statistic *Statistic) AddFile() {
	statistic.Lock()
	defer statistic.Unlock()

	statistic.Files++
}

func (statistic *Statistic) AddFileSkip() {
	statistic.Lock()
	defer statistic.Unlock()

	statistic.FileSkips++
}

func (statistic *Statistic) AddDir() {
	statistic.Lock()
	defer statistic.Unlock()

	statistic.Dirs++
}

func (statistic *Statistic) AddDirSkip() {
	statistic.Lock()
	defer statistic.Unlock()

	statistic.DirSkips++
}
