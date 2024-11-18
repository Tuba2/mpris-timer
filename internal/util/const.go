package util

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	AppId   = "io.github.efogdev.mpris-timer"
	AppName = "Play Timer"

	width         = 256
	height        = 256
	padding       = 16
	strokeWidth   = 32
	bgStrokeColor = "#535353"
)

const svgTemplate = `
<svg width="{{.Width}}" height="{{.Height}}">
  <style>{{if .HasShadow}}#progress{ filter: drop-shadow(-5px 8px 6px rgb(16 16 16 / 0.35)); }{{end}}</style>
  <circle cx="{{.CenterX}}" cy="{{.CenterY}}" r="{{.Radius}}" fill="none" stroke="{{.BgStrokeColor}}" stroke-width="{{.BaseWidth}}"{{if .HasRoundedCorners}} stroke-linecap="round"{{end}} />
  <circle cx="{{.CenterX}}" cy="{{.CenterY}}" r="{{.Radius}}" fill="none" stroke="{{.FgStrokeColor}}" stroke-width="{{.StrokeWidth}}" stroke-dasharray="{{.Circumference}}" stroke-dashoffset="{{.DashOffset}}" transform="rotate(-90 {{.CenterX}} {{.CenterY}})" id="progress"{{if .HasRoundedCorners}} stroke-linecap="round"{{end}} />
</svg>`

var (
	CacheDir string
	DataDir  string
)

type svgParams struct {
	Width             int
	Height            int
	CenterX           int
	CenterY           int
	Radius            float64
	FgStrokeColor     string
	BgStrokeColor     string
	BaseWidth         int
	StrokeWidth       int
	Circumference     float64
	DashOffset        float64
	HasShadow         bool
	HasRoundedCorners bool
}

func init() {
	DataDir = glib.GetUserDataDir()
	if !strings.Contains(DataDir, AppId) {
		DataDir = path.Join(DataDir, AppId)
	}

	CacheDir, _ = os.UserHomeDir()
	CacheDir = path.Join(CacheDir, ".var", "app", AppId, "cache")

	_ = os.MkdirAll(CacheDir, 0755)
	_ = os.MkdirAll(DataDir, 0755)

	// because backward compatibility
	go func() {
		_ = filepath.Walk(CacheDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(info.Name()) == ".svg" {
				_ = os.Remove(path)
			}

			return nil
		})
	}()
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
