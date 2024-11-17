package util

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var (
	tpl = template.New("svg")

	cacheLoaded bool
	cacheMu     sync.RWMutex
	cache       = make(map[string]struct{})
)

func InitCache() {
	err := filepath.Walk(CacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".svg" {
			cache[path] = struct{}{}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("checking cache dir: %v\n", err)
		return
	}

	cacheMu.Lock()
	cacheLoaded = true
	cacheMu.Unlock()
}

func MakeProgressCircle(progress float64) (string, error) {
	progress = math.Max(0, math.Min(100, progress))
	filename := path.Join(CacheDir, fmt.Sprintf("%s.%.2f.svg", strings.Replace(Overrides.Color, "#", "", 1), progress))

	cacheMu.RLock()
	if cacheLoaded {
		_, exists := cache[filename]
		if exists {
			cacheMu.RUnlock()
			return filename, nil
		}
	}
	cacheMu.RUnlock()

	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}

	centerX := width / 2
	centerY := height / 2
	radius := float64(width)/2 - float64(strokeWidth) - float64(padding)
	baseWidth := int(math.Round(strokeWidth * 0.25))
	circumference := 2 * math.Pi * radius
	dashOffset := circumference * (1 - progress/100)

	data := svgParams{
		Width:         width,
		Height:        height,
		CenterX:       centerX,
		CenterY:       centerY,
		Radius:        radius,
		BaseWidth:     baseWidth,
		StrokeWidth:   strokeWidth,
		FgStrokeColor: Overrides.Color,
		BgStrokeColor: bgStrokeColor,
		Circumference: circumference,
		DashOffset:    dashOffset,
	}

	svgString, err := tpl.Parse(svgTemplate)
	if err != nil {
		return "", err
	}

	var svgBuffer bytes.Buffer
	err = svgString.Execute(&svgBuffer, data)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filename, svgBuffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("write SVG: %w", err)
	}

	return filename, nil
}
