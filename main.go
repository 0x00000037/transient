package transient

import (
	"fmt"
	"log/slog"
	"math"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowHeight                   = 1000
	WindowWidth                    = 1000
	ComponentPickerHeight          = 100
	ComponentPickerSeparatorHeight = 10
	CircuitHeight                  = 80
	CircuitWidth                   = 80
	CircuitPadding                 = 10
	CircuitFontSize                = 20
	ConnectionSize                 = 7
	ConnectionConnectAbleRadius    = 10
	WireSize                       = 5

	ComponentPickerStart = WindowWidth - ComponentPickerHeight - ComponentPickerSeparatorHeight
)

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func hypot(x, y int32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

type Direction int32

const (
	DirectionIn Direction = iota
	DirectionOut
)

type DisconnectedWire struct {
	I         *CircuitInstance
	Direction Direction
	N         int32
	X, Y      int32 // Lose end
	S         bool
}

type ConnectedWire struct {
	In, Out   *CircuitInstance
	InN, OutN int32
}

func (c *ConnectedWire) Query(g *Game) bool {
	return false
}

func (c *ConnectedWire) Set(g *Game, v bool) {

}

func (c *ConnectedWire) Delete(g *Game) {
}

func (c *ConnectedWire) Draw(g *Game) {
	x1, y1 := c.In.ConnectionPosition(DirectionOut, c.InN)
	x2, y2 := c.Out.ConnectionPosition(DirectionIn, c.OutN)
	rl.DrawLine(x1, y1, x2, y2, rl.Blue)
}

func (d *DisconnectedWire) Query(g *Game) bool {
	return d.S
}

func (d *DisconnectedWire) Set(g *Game, v bool) {
	d.S = v
}

func (d *DisconnectedWire) Delete(g *Game) {
	i := slices.Index(g.C[g.CIndex].DisconnectedWires, d)
	if i == -1 {
		slog.Warn("Wire isn't tracked so cant delete")
		return
	}
	g.C[g.CIndex].DisconnectedWires = slices.Delete(g.C[g.CIndex].DisconnectedWires, i, i+1)
}

func (d *DisconnectedWire) HoldingNotify(g *Game, x, y int32) {
	d.X = x
	d.Y = y
	//d.Draw(g)
}

/*
|
+-+
  |
// LATER
*/

func (d *DisconnectedWire) Draw(g *Game) {
	x, y := d.I.ConnectionPosition(d.Direction, d.N)
	rl.DrawLine(d.X, d.Y, x, y, rl.Blue)
	rl.DrawCircle(d.X, d.Y, ConnectionSize, rl.Blue)
}

func (d *DisconnectedWire) Release(g *Game, x, y int32) {
	d.X = x
	d.Y = y
	//g.C[g.CIndex].DisconnectedWires = append(g.C[g.CIndex].DisconnectedWires, d)
}

type WireI interface {
	Query(g *Game) bool
	Set(g *Game, v bool)
	Delete(g *Game)
}
type CircuitInstance struct {
	*Circuit
	X, Y int32

	InWires  []WireI
	OutWires []WireI
}

func (c *CircuitInstance) Delete(g *Game) {
	cc := g.C[g.CIndex]
	i := slices.Index(cc.CircuitInstance, c)
	if i == -1 {
		slog.Warn("CircuitInstance not found")
		return
	}
	cc.CircuitInstance = slices.Delete(cc.CircuitInstance, i, i+1)
	for _, w := range c.InWires {
		if w != nil {
			w.Delete(g)
		}
	}
	for _, w := range c.OutWires {
		if w != nil {
			w.Delete(g)
		}
	}
}

func (c *CircuitInstance) SetWire(d Direction, n int32, w WireI) {
	if d == DirectionIn {
		if n >= int32(len(c.InWires)) {
			panic(fmt.Sprintf("CircuitInstance.SetWire: in wire %v out of range\n", n))
		}
		c.InWires[n] = w
	} else {
		if n >= int32(len(c.OutWires)) {
			panic(fmt.Sprintf("CircuitInstance.SetWire: out wire %v out of range\n", n))
		}
		c.OutWires[n] = w
	}
}

func (c *CircuitInstance) GetWire(d Direction, n int32) WireI {
	if d == DirectionIn {
		if n >= int32(len(c.InWires)) {
			panic(fmt.Sprintf("CircuitInstance.SetWire: in wire %v out of range\n", n))
		}
		return c.InWires[n]
	}
	if n >= int32(len(c.OutWires)) {
		panic(fmt.Sprintf("CircuitInstance.SetWire: out wire %v out of range\n", n))
	}
	return c.OutWires[n]
}

func (c *CircuitInstance) AddOrPickDisconnectedWire(direction Direction, n, x, y int32) *DisconnectedWire {
	w := c.GetWire(direction, n)
	if w != nil {
		return w.(*DisconnectedWire)
	}
	w1 := &DisconnectedWire{
		I:         c,
		Direction: direction,
		N:         n,
		X:         x,
		Y:         y,
	}
	c.SetWire(direction, n, w1)
	c.DisconnectedWires = append(c.DisconnectedWires, w1)
	return w1
}

func (c *CircuitInstance) ConnectionPosition(direction Direction, n int32) (x, y int32) {
	var nn int32
	switch direction {
	case DirectionIn:
		nn = c.InWireCount
	case DirectionOut:
		nn = c.OutWireCount
	}
	y = c.Y + (n+1)*CircuitHeight/(nn+1)
	x = c.X
	if direction == DirectionOut {
		x += CircuitWidth
	}
	return
}

type Circuit struct {
	CircuitInstance           []*CircuitInstance
	InWireCount, OutWireCount int32
	DisconnectedWires         []*DisconnectedWire
	ConnectedWires            []*ConnectedWire
	Name                      string
}

type HoldingI interface {
	HoldingNotify(g *Game, x, y int32)
	Delete(g *Game)
}

type Game struct {
	C      []*Circuit
	CIndex int

	Holding HoldingI
}

func (c *CircuitInstance) PosConnectedTo(px, py int32) (ok bool, direction Direction, n int32) {
	for j := range c.InWireCount {
		x, y := c.ConnectionPosition(DirectionIn, j)
		d := hypot(abs(px-x), abs(py-y))
		if d < ConnectionConnectAbleRadius {
			return true, DirectionIn, j
		}
	}

	for j := range c.OutWireCount {
		x, y := c.ConnectionPosition(DirectionOut, j)
		d := hypot(abs(px-x), abs(py-y))
		if d < ConnectionConnectAbleRadius {
			return true, DirectionOut, j
		}
	}
	return false, DirectionIn, 0
}

func (c *CircuitInstance) HoldingNotify(g *Game, x, y int32) {
	c.X = x - CircuitWidth/2
	c.Y = y - CircuitHeight/2
}

func (c *Circuit) Draw(g *Game, x, y int32) {
	for i := range c.InWireCount {
		rl.DrawCircle(x, y+(i+1)*CircuitHeight/(c.InWireCount+1), ConnectionSize, rl.Blue)
	}

	for i := range c.OutWireCount {
		rl.DrawCircle(x+CircuitWidth, y+(i+1)*CircuitHeight/(c.OutWireCount+1), ConnectionSize, rl.Blue)
	}

	rl.DrawRectangle(x, y, CircuitWidth, CircuitHeight, rl.Black)
	rl.DrawText(c.Name, x, y, CircuitFontSize, rl.White)
}

func (g *Game) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	cc := g.C[g.CIndex]

	for _, c := range cc.CircuitInstance {
		c.Draw(g, c.X, c.Y)
	}

	for _, w := range cc.DisconnectedWires {
		w.Draw(g)
	}

	for _, w := range cc.ConnectedWires {
		w.Draw(g)
	}

	// Picker
	rl.DrawRectangle(0, ComponentPickerStart, WindowWidth, ComponentPickerSeparatorHeight, rl.Black)
	for i, c := range g.C {
		c.Draw(g, CircuitPadding+(CircuitWidth+CircuitPadding)*int32(i), WindowHeight-CircuitHeight-CircuitPadding)
	}

	p := rl.GetMousePosition()

	if g.Holding != nil {
		g.Holding.HoldingNotify(g, int32(p.X), int32(p.Y))
	}

	rl.EndDrawing()
}

func NewCircuitInstance(c *Circuit) *CircuitInstance {
	ci := &CircuitInstance{Circuit: c}
	for range c.InWireCount {
		ci.InWires = append(ci.InWires, nil)
	}
	for range c.OutWireCount {
		ci.OutWires = append(ci.OutWires, nil)
	}

	return ci
}

func (g *Game) Actions() {
	p := rl.GetMousePosition()
	px := int32(p.X)
	py := int32(p.Y)

	cc := g.C[g.CIndex]

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		if p.Y > ComponentPickerStart && p.Y < ComponentPickerStart+ComponentPickerHeight { // component picker
			x := int32(p.X - CircuitPadding)
			if x%(CircuitWidth+CircuitPadding) < CircuitWidth {
				xi := x / (CircuitWidth + CircuitPadding)
				if xi >= 0 && xi < int32(len(g.C)) {
					holding := NewCircuitInstance(g.C[xi])
					holding.X = px - CircuitWidth/2
					holding.Y = py - CircuitHeight/2
					g.Holding = holding
					cc.CircuitInstance = append(cc.CircuitInstance, holding)
					slog.Debug("Picked circuit", "circuit", holding.Name)
					return
				}
			}
		}

		for _, c := range cc.CircuitInstance {
			if !(py > c.Y && py < c.Y+CircuitWidth) {
				continue
			}

			if px > c.X && px < c.X+CircuitWidth {
				g.Holding = c
				break
			}

			ok, direction, n := c.PosConnectedTo(px, py)
			if ok {
				holding := c.AddOrPickDisconnectedWire(direction, n, px, py)
				g.Holding = holding

				slog.Debug("Connection grabbed", "circuit", c.Name, "direction", direction, "index", n)
			}
		}

		for _, w := range cc.DisconnectedWires {
			if hypot(abs(w.X-px), abs(w.Y-py)) < ConnectionConnectAbleRadius {
				g.Holding = w
				break
			}
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) { // Drop Component
		if g.Holding != nil {
			if p.Y >= ComponentPickerStart {
				g.Holding.Delete(g)
			}
			g.Holding = nil
		}
	}

}

func Main() (g *Game) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	rl.InitWindow(WindowWidth, WindowHeight, "transient")

	rl.SetTargetFPS(60)
	g = &Game{}
	g.C = append(g.C, &Circuit{Name: "Tran", InWireCount: 2, OutWireCount: 1})
	g.C = append(g.C, &Circuit{Name: "add"})
	g.C = append(g.C, &Circuit{Name: "add"})
	return g
}
