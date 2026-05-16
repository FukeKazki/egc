package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/NotoSansJP-Bold.ttf
var fontData []byte

const (
	imgW = 500
	imgH = 500
)

func main() {
	colorName := flag.String("color", "pink", "文字の色: pink, yellow, black, red, green, blue")
	flag.StringVar(colorName, "c", "pink", "文字の色 (短縮形)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c COLOR] TEXT\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	raw := strings.ReplaceAll(flag.Arg(0), `\n`, "\n")
	lines := strings.Split(raw, "\n")

	col, err := pickColor(*colorName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ttf, err := opentype.Parse(fontData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse font:", err)
		os.Exit(1)
	}

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	size := fitSize(ttf, lines, imgW, imgH)
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create face:", err)
		os.Exit(1)
	}
	defer face.Close()

	n := len(lines)
	_, maxLineH := lineBounds(face, lines)
	// Tight line spacing keeps glyphs large; ~1.15× the tallest glyph leaves
	// just enough breathing room without letting Noto Sans JP's generous
	// line gap shrink the type.
	lineSpacing := maxLineH + maxLineH/8

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}
	for i, line := range lines {
		bounds, _ := font.BoundString(face, line)
		minX, minY := bounds.Min.X.Floor(), bounds.Min.Y.Floor()
		maxX, maxY := bounds.Max.X.Ceil(), bounds.Max.Y.Ceil()
		w := maxX - minX
		h := maxY - minY
		originX := (imgW-w)/2 - minX
		originY := (imgH-h)/2 - minY
		yOffset := (2*i - (n - 1)) * lineSpacing / 2
		drawer.Dot = fixed.P(originX, originY+yOffset)
		drawer.DrawString(line)
	}

	name := strings.ReplaceAll(raw, "\n", "") + ".png"
	out, err := os.Create(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create file:", err)
		os.Exit(1)
	}
	defer out.Close()
	if err := png.Encode(out, img); err != nil {
		fmt.Fprintln(os.Stderr, "failed to encode png:", err)
		os.Exit(1)
	}
}

func fitSize(ttf *sfnt.Font, lines []string, maxW, maxH int) float64 {
	size := 500.0
	const step = 5.0
	for i := 0; i < 200 && size > step; i++ {
		face, err := opentype.NewFace(ttf, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingNone})
		if err != nil {
			return size
		}
		maxLineW, maxLineH := lineBounds(face, lines)
		face.Close()
		lineSpacing := maxLineH + maxLineH/8
		totalH := lineSpacing*(len(lines)-1) + maxLineH
		if maxLineW <= maxW && totalH <= maxH {
			return size
		}
		size -= step
	}
	return size
}

func lineBounds(face font.Face, lines []string) (maxW, maxH int) {
	for _, line := range lines {
		bounds, _ := font.BoundString(face, line)
		w := (bounds.Max.X - bounds.Min.X).Ceil()
		h := (bounds.Max.Y - bounds.Min.Y).Ceil()
		if w > maxW {
			maxW = w
		}
		if h > maxH {
			maxH = h
		}
	}
	return
}

func pickColor(name string) (color.Color, error) {
	switch name {
	case "pink":
		return color.RGBA{255, 0, 255, 255}, nil
	case "yellow":
		return color.RGBA{255, 255, 0, 255}, nil
	case "black":
		return color.RGBA{0, 0, 0, 255}, nil
	case "red":
		return color.RGBA{255, 0, 0, 255}, nil
	case "green":
		return color.RGBA{0, 255, 0, 255}, nil
	case "blue":
		return color.RGBA{0, 0, 255, 255}, nil
	default:
		return nil, fmt.Errorf("指定できる色 : pink, yellow, black, red, green, blue")
	}
}
