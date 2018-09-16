package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"math/rand"
	"net/url"
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
	Size     Coordinate
	Image    []byte
}

func (t Tile) String() string {
	return fmt.Sprintf("%s (%d, %d) Position = %s, Size = %s", t.Tag, t.Row, t.Column, t.Position, t.Size)
}

type Page struct {
	PDF              *gofpdf.Fpdf
	Margin           float64
	Prefix           string
	PageDimensions   Coordinate
	Origin           Coordinate
	MaxtrixDiensions Coordinate
	Title            struct {
		Position   Coordinate
		Dimensions Coordinate
		Text       string
	}
	Rows  int
	Cols  int
	Tiles []*Tile
}

func (p Page) CreateTile(row int, col int, tag string) *Tile {
	t := new(Tile)
	fmt.Printf("Creating tag %s at (%d, %d)\n", tag, row, col)
	return t
}

func CreatePage(rows int, cols int, margin float64, prefix string, title string, filename string) (*Page, error) {
	const (
		minrows   = 2
		maxrows   = 10
		mincols   = 2
		maxcols   = 7
		minmargin = 10.0
		maxmargin = 25.0
	)
	var m = new(Page)
	if rows < minrows || rows > maxrows {
		return m, fmt.Errorf("Number of rows must be [%d..%d]", minrows, maxrows)
	}
	if cols < mincols || cols > maxcols {
		return m, fmt.Errorf("Number of columns must be [%d..%d]", mincols, maxcols)
	}
	if margin < minmargin || margin > maxmargin {
		return m, fmt.Errorf("Margin but be [%1.f..%.1f]mm", minmargin, maxmargin)
	}

	if _, err := url.Parse(prefix); err != nil {
		return m, err
	}

	m.Rows, m.Cols = rows, cols
	m.Margin = margin
	m.Prefix = prefix
	m.Title.Text = title
	m.Title.Dimensions.Width = 100.0
	m.Title.Dimensions.Height = 40.0
	m.Title.Position.Width = 0.0
	m.Title.Position.Height = 0.0
	m.PDF = gofpdf.New("P", "mm", "A4", "")
	m.PageDimensions.Width, m.PageDimensions.Height = m.PDF.GetPageSize()
	m.PDF.AddPage()
	m.PDF.SetFont("Arial", "B", 16)
	m.PDF.Cell(40, 10, "Sample QR Codes for Moien")
	m.PDF.SetFont("Arial", "", 12)
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Tiles = append(m.Tiles, m.CreateTile(row, col, GenerateTag()))
		}
	}
	if err := m.PDF.OutputFileAndClose(filename); err != nil {
		return m, err
	} else {

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
	const rows = 7 // Number of rows of QR codes to print
	const cols = 4 // Number of columns of QR codes to print
	const margin float64 = 10.0
	const prefix = "https://snapper.devops.ukfast.co.uk/"

	if p, err := CreatePage(rows, cols, margin, prefix, "Stuff", "foo.pdf"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(p)
	}
}
