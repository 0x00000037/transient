package main

import (
	"transient"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	g := transient.Main()

	var update = func() {
		g.Actions()
		g.Draw()
	}

	rl.SetMainLoop(update)
	for !rl.WindowShouldClose() {
		update()
	}
}
