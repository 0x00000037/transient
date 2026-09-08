package main

import (
	"transient"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	g := transient.Main()
	defer rl.CloseWindow()
	Loop(g)
}

func Loop(g *transient.Game) {
	for !rl.WindowShouldClose() {
		g.Actions()
		g.Draw()
	}
}
