package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(screenW, screenH, "Sistema Solar 3D")
	rl.SetTargetFPS(60)

	planets := []Planet{
		{
			Name: "Mercurio", Radius: 3.5, SemiMajorAxis: 60, Eccentricity: 0.2056, Speed: 4.15,
			Color:       rl.Color{R: 180, G: 170, B: 160, A: 255},
			TexturePath: "assets/mercury.png",
		},
		{
			Name: "Venus", Radius: 6, SemiMajorAxis: 100, Eccentricity: 0.0068, Speed: 1.62,
			Color:       rl.Color{R: 230, G: 200, B: 130, A: 255},
			TexturePath: "assets/venus_atmosphere.png",
		},
		{
			Name: "Terra", Radius: 6.5, SemiMajorAxis: 145, Eccentricity: 0.0167, Speed: 1.0,
			Color:       rl.Color{R: 80, G: 140, B: 220, A: 255},
			TexturePath: "assets/earth.png",
			Moons: []Moon{
				{Radius: 2, OrbitalDist: 18, Speed: 13.4, Color: rl.Color{R: 210, G: 210, B: 210, A: 255}},
			},
		},
		{
			Name: "Marte", Radius: 5, SemiMajorAxis: 200, Eccentricity: 0.0934, Speed: 0.53,
			Color:       rl.Color{R: 200, G: 100, B: 60, A: 255},
			TexturePath: "assets/mars.png",
			Moons: []Moon{
				{Radius: 1.5, OrbitalDist: 12, Speed: 20, Color: rl.Color{R: 180, G: 150, B: 130, A: 255}},
				{Radius: 1.2, OrbitalDist: 18, Speed: 8, Color: rl.Color{R: 170, G: 145, B: 120, A: 255}},
			},
		},
		{
			Name: "Jupiter", Radius: 20, SemiMajorAxis: 320, Eccentricity: 0.0489, Speed: 0.084,
			Color:       rl.Color{R: 220, G: 175, B: 120, A: 255},
			TexturePath: "assets/jupiter.png",
			Moons: []Moon{
				{Radius: 3, OrbitalDist: 35, Speed: 5.0, Color: rl.Color{R: 230, G: 220, B: 180, A: 255}},
				{Radius: 2.5, OrbitalDist: 47, Speed: 2.5, Color: rl.Color{R: 200, G: 180, B: 150, A: 255}},
				{Radius: 3.5, OrbitalDist: 60, Speed: 1.3, Color: rl.Color{R: 180, G: 160, B: 120, A: 255}},
				{Radius: 3, OrbitalDist: 74, Speed: 0.6, Color: rl.Color{R: 160, G: 140, B: 100, A: 255}},
			},
		},
		{
			Name: "Saturno", Radius: 17, SemiMajorAxis: 450, Eccentricity: 0.0565, Speed: 0.034,
			Color:       rl.Color{R: 210, G: 190, B: 140, A: 255},
			TexturePath: "assets/saturn.png",
			Moons: []Moon{
				{Radius: 4, OrbitalDist: 50, Speed: 0.7, Color: rl.Color{R: 230, G: 215, B: 180, A: 255}},
				{Radius: 2, OrbitalDist: 65, Speed: 0.3, Color: rl.Color{R: 200, G: 185, B: 155, A: 255}},
			},
		},
		{
			Name: "Urano", Radius: 12, SemiMajorAxis: 570, Eccentricity: 0.0457, Speed: 0.012,
			Color:       rl.Color{R: 130, G: 210, B: 220, A: 255},
			TexturePath: "assets/uranus.png",
		},
		{
			Name: "Netuno", Radius: 11, SemiMajorAxis: 670, Eccentricity: 0.0113, Speed: 0.006,
			Color:       rl.Color{R: 60, G: 90, B: 200, A: 255},
			TexturePath: "assets/neptune.png",
		},
	}

	// carrega as texturas
	for i := range planets {
		loadPlanetAssets(&planets[i])
		planets[i].Angle = float32(i) * 0.8
		for j := range planets[i].Moons {
			planets[i].Moons[j].Angle = float32(j) * 1.2
		}
	}
	earthIndex := -1
	for i, p := range planets {
		if p.Name == "Terra" {
			earthIndex = i
			break
		}
	}
	earthTexture := rl.Texture2D{}
	if earthIndex >= 0 {
		earthTexture = planets[earthIndex].Texture
	}

	// configura o sol
	sunMesh := rl.GenMeshSphere(1.0, 32, 32)
	sunMaterial := rl.LoadMaterialDefault()
	sunTexture := rl.LoadTexture("assets/sun.png")
	hasSunTexture := sunTexture.ID > 0
	if hasSunTexture {
		sunMaterial.GetMap(rl.MapDiffuse).Texture = sunTexture
	}
	sunAngle := float32(0)
	comet := Comet{}
	loadCometAssets(&comet, "assets/comet.png")

	impactTexture := rl.LoadTexture("assets/terra.png")

	ringTexture := rl.LoadTexture("assets/saturn_ring.png")

	// configurando estrelas
	stars := make([]Star, 520)
	for i := range stars {
		x := float32(rl.GetRandomValue(0, screenW))
		y := float32(rl.GetRandomValue(0, screenH))
		stars[i] = Star{
			Position: rl.Vector2{X: x, Y: y},
			Radius:   float32(rl.GetRandomValue(1, 3)) * 0.5,
			Alpha:    uint8(rl.GetRandomValue(90, 230)),
			Parallax: float32(rl.GetRandomValue(35, 100)) / 100,
		}
	}

	// configurando camera
	orbitCam := OrbitCamera{
		Target:   rl.Vector3{X: 0, Y: 0, Z: 0},
		Yaw:      2.35,
		Pitch:    0.55,
		Distance: 900,
	}
	speed := float32(1.0)
	paused := false
	showLabels := true
	showOrbits := true
	showGrid := true

	const (
		cameraModeFree = iota
		cameraModeFocus
		cameraModeFollowComet
	)
	cameraMode := cameraModeFree
	previousCameraMode := cameraModeFree
	focusedPlanet := -1

	cometButton := rl.Rectangle{
		X:      20,
		Y:      float32(screenH - 45), // Cast to float32
		Width:  310,
		Height: 30,
	}

	// programa rodando
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		clickReleased := false

		if comet.Active && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if rl.CheckCollisionPointRec(rl.GetMousePosition(), cometButton) {
				if cameraMode == cameraModeFollowComet {
					cameraMode = previousCameraMode
				} else {
					previousCameraMode = cameraMode
					cameraMode = cameraModeFollowComet
				}
			}
		}

		// zoom
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			orbitCam.Distance *= 1 - wheel*0.12
			orbitCam.Distance = clamp(orbitCam.Distance, 120, 1800)
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			mouse := rl.GetMousePosition()
			orbitCam.StartX = mouse.X
			orbitCam.StartY = mouse.Y
			orbitCam.LastX = mouse.X
			orbitCam.LastY = mouse.Y
			orbitCam.Dragging = true
		}
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			mouse := rl.GetMousePosition()
			dx := mouse.X - orbitCam.StartX
			dy := mouse.Y - orbitCam.StartY
			if cameraMode != cameraModeFollowComet {
				clickReleased = dx*dx+dy*dy < 36 && mouse.X > 350
			}
			orbitCam.Dragging = false
		}
		if orbitCam.Dragging {
			mouse := rl.GetMousePosition()
			orbitCam.Yaw -= (mouse.X - orbitCam.LastX) * 0.006
			orbitCam.Pitch += (mouse.Y - orbitCam.LastY) * 0.006
			orbitCam.Pitch = clamp(orbitCam.Pitch, -1.25, 1.25)
			orbitCam.LastX = mouse.X
			orbitCam.LastY = mouse.Y
		}

		panSpeed := orbitCam.Distance * dt * 0.65
		if rl.IsKeyDown(rl.KeyW) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.Z -= panSpeed
		}
		if rl.IsKeyDown(rl.KeyS) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.Z += panSpeed
		}
		if rl.IsKeyDown(rl.KeyA) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.X -= panSpeed
		}
		if rl.IsKeyDown(rl.KeyD) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.X += panSpeed
		}
		if rl.IsKeyDown(rl.KeyQ) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.Y += panSpeed
		}
		if rl.IsKeyDown(rl.KeyE) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			orbitCam.Target.Y -= panSpeed
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			paused = !paused
		}
		if rl.IsKeyPressed(rl.KeyL) {
			showLabels = !showLabels
		}
		if rl.IsKeyPressed(rl.KeyO) {
			showOrbits = !showOrbits
		}
		if rl.IsKeyPressed(rl.KeyG) {
			showGrid = !showGrid
		}
		if rl.IsKeyPressed(rl.KeyR) {
			cameraMode = cameraModeFree
			focusedPlanet = -1
			speed = 1.0
			comet = Comet{
				Texture:  comet.Texture,
				Mesh:     comet.Mesh,
				Material: comet.Material,
			}
			if earthIndex >= 0 && earthTexture.ID > 0 {
				planets[earthIndex].Texture = earthTexture
			}
			orbitCam = OrbitCamera{
				Target:   rl.Vector3{X: 0, Y: 0, Z: 0},
				Yaw:      2.35,
				Pitch:    0.55,
				Distance: 900,
			}
		}
		if rl.IsKeyDown(rl.KeyUp) {
			speed = float32(math.Min(float64(speed*1.05), 20))
		}
		if rl.IsKeyDown(rl.KeyDown) {
			speed = float32(math.Max(float64(speed*0.95), 0.05))
		}
		if rl.IsKeyPressed(rl.KeyZ) {
			speed = 1.0
		}

		if !paused {
			sunAngle += 0.05 * speed * dt
			for i := range planets {
				planets[i].Angle += planets[i].Speed * speed * dt
				for j := range planets[i].Moons {
					planets[i].Moons[j].Angle += planets[i].Moons[j].Speed * speed * dt
				}
			}
		}

		planetPositions := make([]rl.Vector3, len(planets))
		for i, p := range planets {
			planetPositions[i] = planetPosition(p.SemiMajorAxis, p.Eccentricity, p.Angle)
		}

		if earthIndex >= 0 {
			earthPosition := planetPositions[earthIndex]
			if rl.IsKeyPressed(rl.KeyC) {
				launchSpeed := speed
				if paused {
					launchSpeed = 0
				}
				comet = launchComet(planets[earthIndex], launchSpeed, comet)
			}
			if updateComet(&comet, earthPosition, planets[earthIndex].Radius, dt) && impactTexture.ID > 0 {
				planets[earthIndex].Texture = impactTexture
			}
		}

		switch cameraMode {
		case cameraModeFocus:
			if focusedPlanet >= 0 && focusedPlanet < len(planetPositions) {
				orbitCam.Target = planetPositions[focusedPlanet]
			}
		case cameraModeFollowComet:
			if comet.Active {
				direction := rl.Vector3Normalize(comet.Velocity)
				orbitCam.Target = comet.Position
				// Yaw baseado apenas no plano XZ para ficar exatamente atrás do cometa
				orbitCam.Yaw = float32(math.Atan2(float64(direction.X), float64(direction.Z))) + math.Pi
				// Pitch fixo em zero: câmera fica literalmente atrás, sem deslocamento vertical
				orbitCam.Pitch = 0.0
				orbitCam.Distance = 40.0
			} else {
				// Cometa parou (impacto ou saiu): volta para câmera livre apontando para o centro
				cameraMode = cameraModeFree
				orbitCam.Target = rl.Vector3{X: 0, Y: 0, Z: 0}
				orbitCam.Distance = 900
				orbitCam.Yaw = 2.35
				orbitCam.Pitch = 0.55
			}
		}

		// termina de configurar o sistema e começa de fato.

		camera := camera3D(orbitCam)
		if clickReleased {
			selected := pickPlanet(camera, planetPositions, planets)
			if selected >= 0 {
				focusedPlanet = selected
				cameraMode = cameraModeFocus
				orbitCam.Target = planetPositions[selected]
				camera = camera3D(orbitCam)
			}
		}

		// começa a desenhar as coisas na tela
		rl.BeginDrawing()
		rl.ClearBackground(rl.Color{R: 4, G: 5, B: 12, A: 255})
		drawStarBackground(stars, orbitCam)

		rl.BeginMode3D(camera)

		if showGrid {
			rl.DrawGrid(32, 50)
		}

		if showOrbits {
			for _, p := range planets {
				drawEllipticalOrbit(p.SemiMajorAxis, p.Eccentricity, 45)
			}
		}

		drawSun(sunMesh, sunMaterial, hasSunTexture, sunAngle)

		for i, p := range planets {
			position := planetPositions[i]

			if p.Name == "Saturno" {
				drawSaturnRings(position, p.Radius, ringTexture)
			}

			drawPlanet(p, position, p.Angle)

			for _, m := range p.Moons {
				moonPos := moonPosition(position, m.OrbitalDist, m.Angle)
				if showOrbits {
					drawMoonOrbit(position, m.OrbitalDist, 22)
				}
				rl.DrawSphere(moonPos, m.Radius, m.Color)
			}
		}
		if earthIndex >= 0 {
			drawComet(&comet)
		}

		rl.EndMode3D()

		if showLabels {
			drawLabel("Sol", rl.Vector3{X: 24, Y: 18, Z: 0}, camera, rl.Color{R: 255, G: 220, B: 80, A: 230})
			for i, p := range planets {
				labelPos := planetPositions[i]
				labelPos.Y += p.Radius + 5
				drawLabel(p.Name, labelPos, camera, rl.Color{R: 220, G: 220, B: 220, A: 220})
			}
		}

		rl.DrawRectangle(10, 10, 330, 192, rl.Color{R: 0, G: 0, B: 0, A: 150})
		rl.DrawRectangleLines(10, 10, 330, 192, rl.Color{R: 90, G: 90, B: 100, A: 190})

		statusStr := "Rodando"
		if paused {
			statusStr = "Pausado"
		}
		rl.DrawText("Sistema Solar 3D", 20, 18, 18, rl.Color{R: 255, G: 220, B: 80, A: 255})
		rl.DrawText(fmt.Sprintf("Status: %s", statusStr), 20, 46, 12, rl.Color{R: 220, G: 220, B: 220, A: 230})
		rl.DrawText(fmt.Sprintf("Velocidade: %.2fx", speed), 20, 62, 12, rl.Color{R: 220, G: 220, B: 220, A: 230})
		focusName := "Sol"
		if cameraMode == cameraModeFocus && focusedPlanet >= 0 && focusedPlanet < len(planets) {
			focusName = planets[focusedPlanet].Name
		} else if cameraMode == cameraModeFollowComet {
			focusName = "Cometa"
		}
		rl.DrawText(fmt.Sprintf("Foco: %s   Distancia: %.0f", focusName, orbitCam.Distance), 20, 78, 12, rl.Color{R: 220, G: 220, B: 220, A: 230})

		rl.DrawText("ESPAÇO  Pausar/Continuar", 20, 104, 11, rl.Color{R: 165, G: 205, B: 255, A: 230})
		rl.DrawText("Clique planeta Focar   Arraste Rotacionar", 20, 118, 11, rl.Color{R: 165, G: 205, B: 255, A: 230})
		rl.DrawText("WASD/QE Mover alvo da camera", 20, 132, 11, rl.Color{R: 165, G: 205, B: 255, A: 230})
		rl.DrawText("Setas   Velocidade   Z Reset velocidade", 20, 146, 11, rl.Color{R: 165, G: 205, B: 255, A: 230})
		rl.DrawText("L Labels   O Orbitas   G Grade   R Reset", 20, 160, 11, rl.Color{R: 165, G: 205, B: 255, A: 230})
		rl.DrawText("C Lancar cometa contra a Terra", 20, 174, 11, rl.Color{R: 255, G: 190, B: 120, A: 240})

		if comet.Active {
			buttonColor := rl.Color{R: 0, G: 0, B: 0, A: 150}
			if rl.CheckCollisionPointRec(rl.GetMousePosition(), cometButton) {
				buttonColor = rl.Color{R: 50, G: 50, B: 50, A: 150}
			}
			rl.DrawRectangleRec(cometButton, buttonColor)

			textColor := rl.Color{R: 255, G: 190, B: 120, A: 240}
			if cameraMode == cameraModeFollowComet {
				textColor = rl.Color{R: 120, G: 255, B: 120, A: 240}
			}

			text := "Mudar para camera do cometa"
			if cameraMode == cameraModeFollowComet {
				text = "Voltar para camera livre"
			}
			textWidth := float32(rl.MeasureText(text, 12)) // Cast result of MeasureText to float32
			rl.DrawText(
				text,
				int32(cometButton.X+cometButton.Width/2-textWidth/2), // Cast calculations to int32
				int32(cometButton.Y+cometButton.Height/2-6),          // Cast calculations to int32
				12,
				textColor,
			)
		}

		rl.DrawFPS(screenW-80, 10)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
