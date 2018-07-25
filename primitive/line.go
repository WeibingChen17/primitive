package primitive

import (
    "fmt"
    "strings"
    "math"

    "github.com/fogleman/gg"
    "github.com/golang/freetype/raster"
)

type Line struct {
    Worker *Worker
    X1, Y1 float64
    X2, Y2 float64
    tX1, tY1 float64
    tX2, tY2 float64
    Width  float64
}

func NewRandomLine(worker *Worker) *Line {
    rnd := worker.Rnd
    x1 := rnd.Float64() * float64(worker.W)
    y1 := rnd.Float64() * float64(worker.H)
    x2 := rnd.Float64() * float64(worker.W)
    y2 := rnd.Float64() * float64(worker.H)
    if (x1 > x2) {
        x1, x2 = x2, x1
        y1, y2 = y2, y1
    }
    width := 1.0 / 2
    q := &Line{worker, x1, y1, x2, y2, x1, y1, x2, y2, width}
    q.Extend()
    q.Mutate()
    return q
}

func (q *Line) Extend() {
    W := float64(q.Worker.W)
    H := float64(q.Worker.H)
    if q.X1 == q.X2 {
        q.tY1 = 0
        q.tY2 = H
    } else {
        slope := W/H
        slope1 := (q.Y2 - q.Y1)/(q.X2 - q.X1)
        if slope < math.Abs(slope1)  {
            q.tX1 = q.X1 - q.Y1/slope1
            q.tX2 = q.X2 + (H - q.Y2)/slope1
            q.tY1 = 0
            q.tY2 = H
        } else if slope > math.Abs(slope1) {
            q.tY1 = q.Y1 - q.X1 * slope1
            q.tY2 = q.Y2 + (W - q.X2) * slope1
            q.tX1 = 0
            q.tX2 = W
        } else {
            q.tX1 = 0
            q.tY1 = 0
            q.tX2 = W
            q.tY2 = H
        }
    }
}

func (q *Line) Draw(dc *gg.Context, scale float64) {
    dc.DrawLine(q.tX1, q.tY1, q.tX2, q.tY2)
    dc.SetLineWidth(q.Width * scale)
    dc.Stroke()
}

func (q *Line) SVG(attrs string) string {
    // TODO: this is a little silly
    attrs = strings.Replace(attrs, "fill", "stroke", -1)
    return fmt.Sprintf(
        "<line %s fill=\"none\" x1=\"%f\" y1=\"%f\" x2=\"%f\" y2=\"%f\" stroke-width=\"%f\" />",
        attrs, q.tX1, q.tY1, q.tX2, q.tY2, q.Width)
    }

    func (q *Line) Copy() Shape {
        a := *q
        return &a
    }

    func (q *Line) Mutate() {
        const m = 16
        w := q.Worker.W
        h := q.Worker.H
        rnd := q.Worker.Rnd
        for {
            switch rnd.Intn(2) {
            case 0:
                q.X1 = clamp(q.X1+rnd.NormFloat64()*16, -m, float64(w-1+m))
                q.Y1 = clamp(q.Y1+rnd.NormFloat64()*16, -m, float64(h-1+m))
            case 1:
                q.X2 = clamp(q.X2+rnd.NormFloat64()*16, -m, float64(w-1+m))
                q.Y2 = clamp(q.Y2+rnd.NormFloat64()*16, -m, float64(h-1+m))
            case 2:
                q.Width = clamp(q.Width+rnd.NormFloat64(), 1, 4)
            }
            q.Extend()
            if q.Valid() {
                break
            }
        }
    }

    func (q *Line) Valid() bool {
        return (q.tX1 == 0 || q.tY1 == 0) 
    }

    func (q *Line) Rasterize() []Scanline {
        var path raster.Path
        p1 := fixp(q.tX1, q.tY1)
        p2 := fixp(q.tX2, q.tY2)
        path.Start(p1)
        path.Add1(p2)
        width := fix(q.Width)
        return strokePath(q.Worker, path, width, raster.RoundCapper, raster.RoundJoiner)
    }
