package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Planet struct {
	Name          string
	Radius        float32
	SemiMajorAxis float32
	Eccentricity  float32
	Speed         float32
	Color         rl.Color
	Angle         float32
	Moons         []Moon
	TexturePath   string
	Texture       rl.Texture2D
	Mesh          rl.Mesh
	Material      rl.Material
}

type Moon struct {
	Radius      float32
	OrbitalDist float32
	Speed       float32
	Color       rl.Color
	Angle       float32
}

type OrbitCamera struct {
	Target   rl.Vector3
	Yaw      float32
	Pitch    float32
	Distance float32
	StartX   float32
	StartY   float32
	LastX    float32
	LastY    float32
	Dragging bool
}

type Star struct {
	Position rl.Vector2
	Radius   float32
	Alpha    uint8
	Parallax float32
}

type Comet struct {
	Active   bool
	Position rl.Vector3
	Velocity rl.Vector3
	Radius   float32
	Age      float32
	MaxAge   float32
	Texture  rl.Texture2D
	Mesh     rl.Mesh
	Material rl.Material
}

var (
	screenW int32 = 1280
	screenH int32 = 720
)
