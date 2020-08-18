package osu

import (
	"bytes"
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"
)

type Format struct {
	FormatVersion int
	General
	Editor
	Metadata
	Difficulty
	Events
	TimingPoints
	Colours
	HitObjects
}

type General struct { // delimiter:(space)
	AudioFilename            string
	AudioLeadIn              int
	AudioHash                string // deprecated
	PreviewTime              int
	Countdown                int // nofloat
	SampleSet                string
	StackLeniency            float64
	Mode                     int // nofloat
	LetterboxInBreaks        bool
	StoryFireInFront         bool // deprecated
	UseSkinSprites           bool
	AlwaysShowPlayfield      bool // deprecated
	OverlayPosition          string
	SkinPreference           string
	EpilepsyWarning          bool
	CountdownOffset          int
	SpecialStyle             bool
	WidescreenStoryboard     bool
	SamplesMatchPlaybackRate bool
}
type Editor struct { // delimiter:(space)
	Bookmarks       []int // delimiter,
	DistanceSpacing float64
	BeatDivisor     float64
	GridSize        int
	TimelineZoom    float64
}
type Metadata struct { // delimiter:
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string // delimiter(space)
	BeatmapID     int
	BeatmapSetID  int
}
type Difficulty struct { // delimiter:
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}
type Events []Event
type TimingPoints []TimingPoint
type Colours struct { // manual
	Combos              [8]color.RGBA
	SliderTrackOverride color.RGBA
	SliderBorder        color.RGBA
}
type HitObjects []HitObject

func Parse(path string) (*Format, error) {
	o := Format{
		General: General{
			PreviewTime:      -1,
			Countdown:        1,
			SampleSet:        "Normal",
			StackLeniency:    0.7,
			StoryFireInFront: true,
			OverlayPosition:  "NoChange",
		},
		Events:       make([]Event, 0),
		TimingPoints: make([]TimingPoint, 0),
		HitObjects:   make([]HitObject, 0),
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return &o, err
	}
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))

	var section string
	for _, l := range bytes.Split(dat, []byte("\n")) {
		l = bytes.TrimSpace(l)
		line := string(l)
		if isPass(line) {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			continue
		}
		switch section {
		case "General":
			kv := strings.Split(line, `: `)
			switch kv[0] {
			case "AudioFilename":
				o.General.AudioFilename = kv[1]
			case "AudioLeadIn":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.General.AudioLeadIn = int(f)
			case "AudioHash":
				o.General.AudioHash = kv[1]
			case "PreviewTime":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.General.PreviewTime = int(f)
			case "Countdown":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.Countdown = i
			case "SampleSet":
				o.General.SampleSet = kv[1]
			case "StackLeniency":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.General.StackLeniency = f
			case "Mode":
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.Mode = i
			case "LetterboxInBreaks":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.LetterboxInBreaks = b
			case "StoryFireInFront":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.StoryFireInFront = b
			case "UseSkinSprites":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.UseSkinSprites = b
			case "AlwaysShowPlayfield":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.AlwaysShowPlayfield = b
			case "OverlayPosition":
				o.General.OverlayPosition = kv[1]
			case "SkinPreference":
				o.General.SkinPreference = kv[1]
			case "EpilepsyWarning":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.EpilepsyWarning = b
			case "CountdownOffset":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.General.CountdownOffset = int(f)
			case "SpecialStyle":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.SpecialStyle = b
			case "WidescreenStoryboard":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.WidescreenStoryboard = b
			case "SamplesMatchPlaybackRate":
				b, err := strconv.ParseBool(kv[1])
				if err != nil {
					return &o, err
				}
				o.General.SamplesMatchPlaybackRate = b
			}
		case "Editor":
			kv := strings.Split(line, `: `)
			switch kv[0] {
			case "Bookmarks":
				slice := make([]int, 0)
				for _, s := range strings.Split(kv[1], ",") {
					i, err := strconv.Atoi(s)
					if err != nil {
						return &o, err
					}
					slice = append(slice, i)
				}
				o.Editor.Bookmarks = slice
			case "DistanceSpacing":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.DistanceSpacing = f
			case "BeatDivisor":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.BeatDivisor = f
			case "GridSize":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.GridSize = int(f)
			case "TimelineZoom":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Editor.TimelineZoom = f
			}
		case "Metadata":
			kv := strings.Split(line, `:`)
			switch kv[0] {
			case "Title":
				o.Metadata.Title = kv[1]
			case "TitleUnicode":
				o.Metadata.TitleUnicode = kv[1]
			case "Artist":
				o.Metadata.Artist = kv[1]
			case "ArtistUnicode":
				o.Metadata.ArtistUnicode = kv[1]
			case "Creator":
				o.Metadata.Creator = kv[1]
			case "Version":
				o.Metadata.Version = kv[1]
			case "Source":
				o.Metadata.Source = kv[1]
			case "Tags":
				slice := make([]string, 0)
				for _, s := range strings.Split(kv[1], " ") {
					slice = append(slice, s)
				}
				o.Metadata.Tags = slice
			case "BeatmapID":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					// return &o, err
					f = -1
				}
				o.Metadata.BeatmapID = int(f)
			case "BeatmapSetID":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					// return &o, err
					f = -1
				}
				o.Metadata.BeatmapSetID = int(f)
			}
		case "Difficulty":
			kv := strings.Split(line, `:`)
			switch kv[0] {
			case "HPDrainRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.HPDrainRate = f
			case "CircleSize":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.CircleSize = f
			case "OverallDifficulty":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.OverallDifficulty = f
			case "ApproachRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.ApproachRate = f
			case "SliderMultiplier":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.SliderMultiplier = f
			case "SliderTickRate":
				f, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return &o, err
				}
				o.Difficulty.SliderTickRate = f
			}
		case "Events":
			e, err := newEvent(line)
			if err != nil {
				// return &o, err
				continue
			}
			o.Events = append(o.Events, e)
		case "TimingPoints":
			tp, err := newTimingPoint(line)
			if err != nil {
				return &o, err
			}
			o.TimingPoints = append(o.TimingPoints, tp)
		case "Colours":
			kv := strings.Split(line, ` : `)
			rgb := newRGB(kv[1])
			switch kv[0] {
			case "Combo1":
				o.Colours.Combos[0] = rgb
			case "Combo2":
				o.Colours.Combos[1] = rgb
			case "Combo3":
				o.Colours.Combos[2] = rgb
			case "Combo4":
				o.Colours.Combos[3] = rgb
			case "Combo5":
				o.Colours.Combos[4] = rgb
			case "Combo6":
				o.Colours.Combos[5] = rgb
			case "Combo7":
				o.Colours.Combos[6] = rgb
			case "Combo8":
				o.Colours.Combos[7] = rgb
			case "SliderTrackOverride":
				o.Colours.SliderTrackOverride = rgb
			case "SliderBorder":
				o.Colours.SliderBorder = rgb
			}
		case "HitObjects":
			ho, err := newHitObject(line)
			if err != nil {
				return &o, err
			}
			o.HitObjects = append(o.HitObjects, ho)
		}
	}
	return &o, nil
}
func isPass(line string) bool {
	return len(line) == 0 || len(line) >= 2 && line[:2] == "//"
}
func isSection(line string) bool {
	if len(line) == 0 {
		return false
	}
	return string(line[0]) == "[" && string(line[len(line)-1]) == "]"
}

func newRGB(s string) color.RGBA {
	var rgb color.RGBA
	for i, c := range strings.Split(s, `,`) {
		f, err := strconv.ParseFloat(c, 64)
		if err != nil {
			f = 0
		}
		switch i {
		case 0:
			rgb.R = uint8(f)
		case 1:
			rgb.G = uint8(f)
		case 2:
			rgb.B = uint8(f)
		}
	}
	rgb.A = 255
	return rgb
}