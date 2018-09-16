package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"math/rand"
	"net/url"
	"strings"
)

type Coordinate struct {
	Width, Height float64
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", c.Width, c.Height)
}

type Tile struct {
	Tag      string // e.g. W1X2Y3Z4
	Row      int
	Column   int
	Position Coordinate
	Image    []byte
}

func (t Tile) String() string {
	return fmt.Sprintf("%s (%d, %d) Position = %s", t.Tag, t.Row, t.Column, t.Position)
}

type Page struct {
	PDF               *gofpdf.Fpdf
	Margin            Coordinate
	Prefix            string
	PageDimensions    Coordinate
	Origin            Coordinate
	MaxtrixDimensions Coordinate
	TileDimensions    Coordinate
	Rows              int
	Cols              int
	Tiles             []*Tile
}

func (p Page) String() string {
	bob := strings.Builder{}
	bob.WriteString(fmt.Sprintf("(%d x %d) %s", p.Rows, p.Cols, p.TileDimensions))
	for _, r := range p.Tiles {
		bob.WriteString(fmt.Sprintf("\n%s", *r))
	}
	return bob.String()
}

func (p Page) CreateTile(row int, col int, tag string) *Tile {
	t := new(Tile)
	fmt.Printf("Creating tag %s at (%d, %d)\n", tag, row, col)
	t.Row = row
	t.Column = col
	t.Tag = tag
	t.Position.Height = p.Origin.Height + float64(t.Row)*p.TileDimensions.Height
	t.Position.Width = p.Origin.Width + float64(t.Column)*p.TileDimensions.Width
	return t
}

func CreatePage(rows int, cols int, margin Coordinate, prefix string, filename string) (*Page, error) {
	const (
		minrows   = 2
		maxrows   = 10
		mincols   = 2
		maxcols   = 7
		minmargin = 10.0
		maxmargin = 50.0
	)
	var m = new(Page)
	if rows < minrows || rows > maxrows {
		return m, fmt.Errorf("Number of rows must be [%d..%d]", minrows, maxrows)
	}
	if cols < mincols || cols > maxcols {
		return m, fmt.Errorf("Number of columns must be [%d..%d]", mincols, maxcols)
	}
	if margin.Width < minmargin || margin.Height < minmargin || margin.Width > maxmargin || margin.Height > maxmargin {
		return m, fmt.Errorf("Margins must be [%1.f..%.1f]mm", minmargin, maxmargin)
	}

	if _, err := url.Parse(prefix); err != nil {
		return m, err
	}

	m.Rows, m.Cols = rows, cols
	m.Margin.Width = margin.Width
	m.Margin.Height = margin.Height
	m.Prefix = prefix
	m.PDF = gofpdf.New("P", "mm", "A4", "")
	m.PageDimensions.Width, m.PageDimensions.Height = m.PDF.GetPageSize()
	m.MaxtrixDimensions.Height = m.PageDimensions.Height - 2.0*m.Margin.Height
	m.MaxtrixDimensions.Width = m.PageDimensions.Width - 2.0*m.Margin.Width
	m.Origin.Height = m.Margin.Height
	m.Origin.Width = m.Margin.Width
	m.TileDimensions.Height = m.MaxtrixDimensions.Height / float64(m.Rows)
	m.TileDimensions.Width = m.MaxtrixDimensions.Width / float64(m.Cols)
	m.PDF.AddPage()
	m.PDF.SetFont("Arial", "", 12)
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			tile := m.CreateTile(row, col, GenerateTag())
			m.Tiles = append(m.Tiles, tile)
			m.PDF.MoveTo(tile.Position.Width, tile.Position.Height)
			m.PDF.Cell(m.TileDimensions.Width, m.TileDimensions.Height, tile.Tag)
		}
	}
	if err := m.PDF.OutputFileAndClose(filename); err != nil {
		return m, err
	}
	return m, nil
}

func Intersects(a, b []rune) bool {
	// Returns true if any rune in b is contained in a (i.e. intersection of the two sets)
	if len(a) == 0 || len(b) == 0 {
		fmt.Printf("Empty set")
	}
	for _, r := range a {
		for _, s := range b {
			if r == s {
				return true
			}
		}
	}
	return false
}

func GenerateTag() string {
	/*
	** Generate a plausible service tag (i.e. `length` random alphanumeric characters).
	 */
	const length = 8
	letters := []rune("ABCDEFGHIJKLMNPQRSTUVWXYZ")
	digits := []rune("0123456789")

	candidates := []rune{}
	candidates = append(candidates, letters...)
	candidates = append(candidates, digits...)
	tag := make([]rune, length)

	for {
		for i := 0; i < length; i++ {
			tag[i] = candidates[rand.Intn(len(candidates))]
		}
		if Intersects(tag, letters) && Intersects(tag, digits) {
			break
		}
	}
	return string(tag)
}

func main() {
	var rows = 7 // Number of rows of QR codes to print
	var cols = 4 // Number of columns of QR codes to print
	var margin = Coordinate{20.0, 20.0}
	var prefix = "https://snapper.devops.ukfast.co.uk/"

	if p, err := CreatePage(rows, cols, margin, prefix, "foo.pdf"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(p)
	}
}
