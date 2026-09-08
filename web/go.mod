module transient_web

go 1.26.1

require github.com/gen2brain/raylib-go/raylib v0.60.1

replace (
	github.com/BrownNPC/Raylib-Go-Wasm/wasm-runtime => ./Raylib-Go-Wasm/wasm-runtime
	github.com/gen2brain/raylib-go/raygui => ./Raylib-Go-Wasm/raygui
	github.com/gen2brain/raylib-go/raylib => ./Raylib-Go-Wasm/raylib
)

require (
	github.com/BrownNPC/Raylib-Go-Wasm/wasm-runtime v0.0.0-00010101000000-000000000000 // indirect
	github.com/BrownNPC/wasm-ffi-go v1.3.0 // indirect
)
