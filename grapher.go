package main

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
)

type Grapher struct {
	left   image.Image
	right  image.Image
	center image.Image

	bgNormal image.Image
	bgMedium image.Image
	bgRare   image.Image
	bgSuper  image.Image

	font *truetype.Font
}

func NewGrapher(pathPrefix string) *Grapher {
	return &Grapher{
		left:   loadImage(pathPrefix + "/g24-left.png"),
		right:  loadImage(pathPrefix + "/g24-right.png"),
		center: loadImage(pathPrefix + "/g24-center.png"),

		bgNormal: loadImage(pathPrefix + "/normal.png"),
		bgMedium: loadImage(pathPrefix + "/medium.png"),
		bgRare:   loadImage(pathPrefix + "/rare.png"),
		bgSuper:  loadImage(pathPrefix + "/super.png"),

		font: LoadFont(pathPrefix + "/JetBrainsMono-ExtraBold.ttf"),
	}
}

func (g *Grapher) Draw(domain string) *gg.Context {
	dc := gg.NewContext(1024, 1024)
	dc.SetColor(color.Black)

	domainLen := len(domain)

	if domainLen <= 4 {
		dc.DrawImage(g.bgSuper, 0, 0)
	} else if domainLen <= 7 {
		dc.DrawImage(g.bgRare, 0, 0)
	} else if domainLen <= 10 {
		dc.DrawImage(g.bgMedium, 0, 0)
	} else {
		dc.DrawImage(g.bgNormal, 0, 0)
	}

	domainName := domain + ".ton"

	left := 60
	shadowWidth := 120 - left

	canvasWith := dc.Width()

	fontSize := float64(46)
	dc.SetFontFace(SizedFont(g.font, fontSize*2))
	w, _ := dc.MeasureString(domainName)

	for int(w) > canvasWith-left*4 && fontSize > 20 {
		fontSize -= 2
		dc.SetFontFace(SizedFont(g.font, fontSize*2))
		w, _ = dc.MeasureString(domainName)
	}

	for int(w) > canvasWith-left*4 && len(domainName) > 20 {
		middle := len(domainName) / 2
		left := domainName[:middle-5]
		right := domainName[middle+5:]
		domainName = left + "..." + right
		w, _ = dc.MeasureMultilineString(domainName, 1)
	}

	bg := g.getBox(int(int(w)+left*4+shadowWidth*2), 470)
	dc.DrawImageAnchored(bg.Image(), 512, 512, .5, .5)
	dc.DrawStringAnchored(domainName, 512, 512-fontSize/2.25, 0.5, 0.5)

	return dc
}

func (g *Grapher) getBox(w int, h int) *gg.Context {
	left := g.left
	right := g.right
	center := g.center

	imageHeight := left.Bounds().Size().Y
	sideWidth := left.Bounds().Size().X
	zoom := float64(imageHeight) / float64(h)
	rZoom := 1 / zoom

	minWidth := sideWidth + sideWidth

	width := int(float64(w) * zoom)
	if width < minWidth {
		width = minWidth
	}

	out := gg.NewContext(width, imageHeight)
	out.DrawImage(left, 0, 0)
	out.DrawImage(right, width-sideWidth, 0)
	if (width - minWidth) > 0 {
		out.SetFillStyle(gg.NewSurfacePattern(center, gg.RepeatX))
		out.Push()
		out.DrawRectangle(float64(sideWidth), 0, float64(width-sideWidth-sideWidth), float64(imageHeight))
		out.Fill()
		out.Pop()
	}

	res := gg.NewContext(w, h)
	res.Scale(rZoom, rZoom)
	res.DrawImage(out.Image(), 0, 0)
	return res
}
