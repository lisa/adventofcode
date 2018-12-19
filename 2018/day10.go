package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
)

var (
	inputFile  = flag.String("input", "inputs/day10.txt", "input file")
	partB      = flag.Bool("partB", false, "do part b solution?")
	debug      = flag.Bool("debug", false, "debug?")
	debug2     = flag.Bool("debug2", false, "more debug")
	lineParser = regexp.MustCompile(`position=<\s?(-?\d+),\s{1,}(-?\d+)> velocity=<\s?(-?\d+),\s{1,}(-?\d+)>.*`)
)

func errorIf(msg string, e error) {
	if e != nil {
		fmt.Printf("%s\n")
		os.Exit(1)
	}
}

// Point - has current (X,Y) position and a X and Y component velocity
type Point struct {
	X, Y                 int
	XVelocity, YVelocity int
}

func (p *Point) Copy() *Point {
	return &Point{
		X: p.X, Y: p.Y,
		XVelocity: p.XVelocity, YVelocity: p.YVelocity,
	}
}

// EqualsTo equality test, based on (X,Y) comparison
func (p *Point) EqualTo(o *Point) bool {
	return ((p.X == o.X) && (p.Y == o.Y))
}

// The field containing all the +Point+s
type Field struct {
	Points []*Point
}

func (f *Field) Copy() *Field {
	ret := &Field{
		Points: make([]*Point, len(f.Points)),
	}
	for i := 0; i < len(f.Points); i++ {
		ret.Points[i] = f.Points[i].Copy()
	}
	return ret
}

// AddPoint to the Field
func (f *Field) AddPoint(x, y, xvel, yvel int) {
	f.Points = append(f.Points, &Point{
		X:         x,
		Y:         y,
		XVelocity: xvel,
		YVelocity: yvel,
	})
}

// Advance - move forward one tick
func (f *Field) Advance() {
	for _, point := range f.Points {
		point.X += point.XVelocity
		point.Y += point.YVelocity
	}
}

// CountOverlaps - how many Points overlap? Compare each point to each other. This could be used to find a relative minimum
func (f *Field) CountOverlaps() int {
	overlaps := 0
	for i, op := range f.Points {
		for j, ip := range f.Points {
			if i == j {
				continue
			} else {
				if op.EqualTo(ip) {
					overlaps++
				}
			}
		}
	}
	return overlaps
}

// Draw - Render the Points to an image based on the iteration
// returns (filename, error)
func (f *Field) Draw(i int) (string, error) {
	filename := fmt.Sprintf("day10a-iteration-%d.png", i)
	minX := f.Points[0].X
	minY := f.Points[0].Y
	maxX := f.Points[0].X
	maxY := f.Points[0].Y

	for _, point := range f.Points {
		if point.X < minX {
			minX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	if *debug {
		fmt.Printf("i=%d Image ranges: (%d,%d) -> (%d,%d)\n", i, minX-40, minY-40, 2*maxX, 2*maxY)
	}

	img := image.NewRGBA(image.Rect(minX-40, minY-40, 2*maxX, 2*maxY))

	for _, point := range f.Points {
		if *debug {
			fmt.Printf("  i=%d: Plotting (%d,%d)\n", i, point.X, point.Y)
		}
		img.Set(point.X, point.Y, color.RGBA{0, 0, 0, 255})
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		return "", err
	}
	png.Encode(file, img)
	return filename, nil

}

func NewField() *Field {
	return &Field{
		Points: make([]*Point, 0),
	}
}

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	errorIf("Can't open input file", err)

	defer input.Close()
	lineReader := bufio.NewScanner(input)

	field := NewField()

	for lineReader.Scan() {
		matches := lineParser.FindAllStringSubmatch(lineReader.Text(), -1)
		x, err := strconv.Atoi(matches[0][1])
		errorIf("Couldn't parse X\n", err)
		y, err := strconv.Atoi(matches[0][2])
		errorIf("Couldn't parse Y\n", err)
		xvel, err := strconv.Atoi(matches[0][3])
		errorIf("Couldn't parse X velocity\n", err)
		yvel, err := strconv.Atoi(matches[0][4])
		errorIf("Couldn't parse Y velocity\n", err)
		field.AddPoint(x, y, xvel, yvel)
	}
	input.Close()
	fields := make(map[int]*Field)
	hist := make(map[int]int)
	highscore := field.CountOverlaps()
	hist[0] = highscore
	if *debug {
		fmt.Printf("0 ")
	}
	for i := 1; i < 12000; i++ {

		field.Advance()
		hist[i] = field.CountOverlaps()

		if hist[i] > highscore {
			highscore = hist[i]
			if *debug {
				fmt.Printf("new highscore: %d -> %d\n", i, highscore)
			}
			fields[i] = field.Copy()
		}
		if *debug {
			fmt.Printf(".")
			if i%200 == 0 {
				fmt.Printf("\n%d ", i)
			}
		}
		if hist[i] == 0 {
			delete(hist, i)
		}
	}

	for i, bestField := range fields {
		fname, err := bestField.Draw(i)
		errorIf(fmt.Sprintf("Couldn't render image for i=%d\n", i), err)
		fmt.Printf("Rendered %s\n", fname)
	}
}

// highest overlaps = 10645
