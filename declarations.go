package main

import (
	"context"
	"embed"

	"github.com/go-vgo/robotgo"
)

type vdflib struct {
	Path                     string
	Label                    string
	ContentID                string
	TotalSize                string
	UpdateCleanBytesTally    string
	TimeLastUpdateCorruption string
	Apps                     []libapp
}

type libapp struct {
	AppID   string
	BuildID string
}

type operation func(ctx context.Context) error

//go:embed blank.mp3
var embeds embed.FS

var logLines []string
var currentTrack int

func (Track0) Play() {
	robotgo.KeyTap(robotgo.AudioPlay)
	currentTrack = 0
}
func (Track0) RW() {
	robotgo.KeyTap(robotgo.AudioPrev)
	currentTrack = 0
}
func (Track0) FF() {
	robotgo.KeyTap(robotgo.AudioNext)
	currentTrack = 0
}

type Track0 struct{}

func (Track1) Play() {
	robotgo.KeyTap(robotgo.AudioPlay)
	currentTrack = 1
}

func (Track1) FF() {
	robotgo.KeyTap(robotgo.AudioNext)
	currentTrack = 1
}

func (Track1) RW() {
	robotgo.KeyTap(robotgo.AudioPrev)
	currentTrack = 1
}

type Track1 struct{}

func (Track2) Play() {
	robotgo.KeyTap(robotgo.AudioPlay)
	currentTrack = 2
}

func (Track2) FF() {
	robotgo.KeyTap(robotgo.AudioNext)
	currentTrack = 2
}

func (Track2) RW() {
	robotgo.KeyTap(robotgo.AudioPrev)
	currentTrack = 2
}

type Track2 struct{}

var tick = make(chan bool)

var equalsZero []bool
var equalsOne []bool
var equalsTwo []bool
