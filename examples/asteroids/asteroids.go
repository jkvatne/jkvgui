package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
)

type Body struct {
	pos      f32.Pos
	vel      f32.Pos
	rotation float32
	radius   float32
	active   bool
	color    f32.Color
}

type Bullet = Body

type Asteroid struct {
	Body
	segments     int
	offsets      []float32
	points       []f32.Pos
	rotationStep float32 // = 1.0
}

type Player struct {
	Body
	engineOn       bool
	bullets        int
	bulletsLimit   int
	fuel           int
	fuelLimit      int
	cooldown       int
	cooldownFrames int
	points         []float32
}

type Message struct {
	frames int
	text   string
	size   int
	color  f32.Color
}
type Game struct {
	win       *sys.Window
	screen    f32.Pos
	player    Player
	bullets   []*Bullet
	asteroids []*Asteroid
	score     int
	highscore int
	ships     int
	level     int
	msg       Message
}

func (a *Asteroid) setup() {
	a.segments = 15 // int(a.radius/10) * 10
	a.offsets = make([]float32, a.segments)
	a.points = make([]f32.Pos, a.segments)
	for i := range a.segments {
		a.offsets[i] = a.radius + 15*(0.5-rand.Float32())
	}
}

func (p *Player) reset(screen f32.Pos) {
	p.fuelLimit = 9999
	p.bulletsLimit = 99
	p.bullets = p.bulletsLimit
	p.fuel = p.fuelLimit
	p.bullets, p.fuel = p.bulletsLimit, p.fuelLimit
	p.pos = f32.Pos{screen.X / 2, screen.Y / 2}
	p.vel = f32.Pos{}
	p.radius = 15
	p.rotation = -90
	p.cooldownFrames = 15
	p.active = true
}

// HandleInput will check for key presses and act accordingly
func (game *Game) HandleInput() {
	p := &game.player
	p.engineOn = false
	if game.win.LastKey == sys.KeyEscape {
		os.Exit(0)
	}
	if !p.active {
		return
	}
	if sys.KeyIsDown[sys.KeySpace] && p.cooldown <= 0 && p.active {
		game.addBullet()
		p.cooldown = p.cooldownFrames
	}
	if p.fuel >= 10 {
		if sys.KeyIsDown[sys.KeyUp] && p.active {
			angle := f32.Radians(p.rotation)
			p.vel.X += float32(math.Cos(float64(angle)) * 0.05)
			p.vel.Y += float32(math.Sin(float64(angle)) * 0.05)
			p.fuel -= 5
			p.engineOn = true
		}
	}
	if p.fuel >= 1 {
		if sys.KeyIsDown[sys.KeyLeft] {
			p.rotation -= 2
			p.fuel--
		}
		if sys.KeyIsDown[sys.KeyRight] {
			p.rotation += 2
			p.fuel--
		}
	}
}

// Update the game - moving things and drawing them
func (game *Game) Update() {
	game.msg.frames = max(0, game.msg.frames-1)
	p := &game.player
	p.cooldown = max(0, p.cooldown-1)
	p.pos.X += p.vel.X
	p.pos.Y += p.vel.Y
	p.pos.Wrap(game.screen)
	for _, b := range game.bullets {
		if !b.active {
			continue
		}
		b.pos.X += b.vel.X
		b.pos.Y += b.vel.Y
		if b.pos.X < 0 || b.pos.X > game.screen.X || b.pos.Y < 0 || b.pos.Y > game.screen.Y {
			b.active = false
		}
	}
	for _, a := range game.asteroids {
		if !a.active {
			continue
		}
		a.pos.X += a.vel.Y
		a.pos.Y += a.vel.Y
		a.pos.Wrap(game.screen)
		a.rotation += a.rotationStep
		if p.active && p.pos.Distance(a.pos) <= (p.radius+a.radius) {
			// Player/asteroids collision
			p.active = false
			a.active = false
			game.SplitAsteroid(a, p.vel.ScaleBy(0.5))
			game.player.reset(game.screen)
			game.score += 50
			game.ShowMessage("Your ship was destroyed.", f32.Red, 90)
			game.ships--
			if game.ships <= 0 {
				game.ships = 5
				game.score = 0
				game.asteroids = nil
				game.addAsteroids(10)
			}
		}
	}
	for _, b := range game.bullets {
		if !b.active {
			continue
		}
		for _, a := range game.asteroids {
			if !a.active {
				continue
			}
			if b.active && b.pos.Distance(a.pos) <= (b.radius+a.radius) {
				// Bullet/asteroid collision
				b.active = false
				a.active = false
				game.score += 100
				game.SplitAsteroid(a, b.vel.ScaleBy(0.2))
			}
		}
	}
	// Remve asteroids that are not active
	oldA := game.asteroids
	game.asteroids = nil
	for _, a := range oldA {
		if a.active {
			game.asteroids = append(game.asteroids, a)
		}
	}
	// Remove bullets that are not active
	oldB := game.bullets
	game.bullets = nil
	for _, b := range oldB {
		if b.active {
			game.bullets = append(game.bullets, b)
		}
	}
	if len(game.asteroids) == 0 {
		game.level++
		game.ships++
		game.player.reset(game.screen)
		game.addAsteroids(10)
		game.ShowMessage("YOU WIN", f32.Green, 90)
	}
	game.highscore = max(game.score, game.highscore)
}

func (game *Game) ShowMessage(text string, color f32.Color, frames int) {
	game.msg.text = text
	game.msg.color = color
	game.msg.frames = frames
}

func (game *Game) SplitAsteroid(a *Asteroid, vel f32.Pos) {
	if a.radius < 30 {
		return
	}
	shrinkFactor := 0.5 + 0.3*rand.Float32()
	a1 := Asteroid{Body: Body{active: true, radius: a.radius * shrinkFactor, vel: a.vel, pos: a.pos, rotation: a.rotation}, rotationStep: -a.rotationStep}
	a2 := Asteroid{Body: Body{active: true, radius: a.radius * (1 - shrinkFactor), vel: a.vel, pos: a.pos, rotation: a.rotation}, rotationStep: 2 * a.rotationStep}
	a1.Body.color = a.Body.color
	a2.Body.color = a.Body.color
	a1.vel = a1.vel.ScaleBy(shrinkFactor)
	a1.vel.X *= vel.X
	a1.vel.Y *= vel.Y
	a2.vel = a2.vel.ScaleBy(1 - shrinkFactor)
	a2.vel.X *= -vel.X
	a2.vel.Y *= -vel.Y
	a1.setup()
	a2.setup()
	game.asteroids = append(game.asteroids, &a1)
	game.asteroids = append(game.asteroids, &a2)
}

// Draw the game graphics (ship, bullets and asteroids)
func (game *Game) draw() {
	w, h := game.win.Window.GetSize()
	gd := game.win.Gd
	r := f32.Rect{0, 0, float32(w), float32(h)}
	gd.SolidRect(r, f32.Blue)
	game.DrawShip()
	for _, b := range game.bullets {
		gd.Circle(b.pos, 3, 0, f32.Yellow, f32.Yellow)
	}
	for _, a := range game.asteroids {
		game.DrawAsteroide(a)
	}
	f := font.Fonts[gpu.Normal16]
	label1 := fmt.Sprintf("Level: %d Ships: %d", game.level, game.ships)
	f.DrawText(game.win.Gd, 5, 16, f32.White, 0, 0, label1)
	label2 := fmt.Sprintf("Bullets: %d Fuel: %d", game.player.bullets, game.player.fuel)
	f.DrawText(game.win.Gd, game.win.WidthDp/2-33, 16, f32.White, 0, 0, label2)
	label3 := fmt.Sprintf("Score: %d  HighScore: %d", game.score, game.highscore)
	ww := f.Width(label3) + 5
	f.DrawText(game.win.Gd, game.win.WidthDp-ww, 16, f32.White, 0, 0, label3)
	f.DrawText(game.win.Gd, 20, game.win.HeightDp-20, f32.White, 0, 0, "Use arrows + space to control your ship. Use Escape to end the game.")
	// Display message for the given number of frames
	if game.msg.frames > 0 {
		font.Fonts[gpu.Bold28].DrawText(game.win.Gd, game.win.WidthDp*1/3, game.win.HeightDp/2, game.msg.color, 0, 0, game.msg.text)
	}
}

// DrawShip will draw a ship as an arrow
func (game *Game) DrawShip() {
	gd := &game.win.Gd
	p := &game.player
	if !p.active {
		return
	}
	angle := f32.Radians(p.rotation)
	p1 := p.pos.Offset(angle, p.radius*2)
	p2 := p.pos.Offset(angle+2.5, 0.6*p.radius*2)
	p3 := p.pos.Offset(angle, -0.3*p.radius*2)
	p4 := p.pos.Offset(angle-2.5, 0.6*p.radius*2)
	var pp = []f32.Pos{p1, p2, p3, p1, p3, p4}
	gd.Triangles(pp, f32.Yellow)
	if p.engineOn {
		engine := p.pos.Offset(angle+math.Pi, 0.7*p.radius)
		gd.Circle(engine, 8, 0, f32.Yellow, f32.Yellow)
	}
}

func (game *Game) DrawAsteroide(a *Asteroid) {
	if len(a.points) == 0 || !a.active {
		return
	}
	for i := range a.segments {
		angle := f32.Radians(a.rotation + float32(i)*360/float32(a.segments))
		o := f32.Pos{}
		p := o.Offset(angle, a.offsets[i])
		p.X += a.pos.X
		p.Y += a.pos.Y
		a.points[i] = p
	}
	game.win.Gd.Poly(a.points, a.color)
}

func (game *Game) addBullet() {
	if game.player.bullets <= 0 {
		return
	}
	game.player.bullets--
	angle := f32.Radians(game.player.rotation)
	game.bullets = append(game.bullets, &Bullet{
		pos:    game.player.pos,
		radius: 3,
		vel:    game.player.vel.Offset(angle, 10),
		active: true,
	})
}

func (game *Game) addAsteroids(count int) {
	npos := f32.Pos{}
	for range count {
		for {
			npos = f32.Random(game.screen)
			for _, a := range game.asteroids {
				if a.pos.Distance(npos) < (a.radius + 30) {
					continue
				}
			}
			if game.player.pos.Distance(npos) < 5*(game.player.radius+30) {
				continue
			}
			break
		}
		radius := 50 + 40*(0.5-rand.Float32())
		vv := f32.Pos{(rand.Float32() - rand.Float32()) + 0.1, (rand.Float32() - rand.Float32()) + 0.1}
		asteroid := Asteroid{
			Body: Body{
				pos:      npos,
				vel:      vv,
				radius:   radius,
				rotation: 360 * rand.Float32(),
				active:   true,
				color:    f32.Color{rand.Float32(), rand.Float32(), rand.Float32(), 0.9},
			},
			rotationStep: 2 * (0.5 - rand.Float32()),
		}
		asteroid.Body.color = asteroid.Body.color.Mute(0.6)
		asteroid.setup()
		game.asteroids = append(game.asteroids, &asteroid)
	}
}

func main() {
	sys.Init()
	defer sys.Shutdown()
	// Set up for 60 frames pr second
	sys.MinFrameDelay = time.Second / 60
	sys.MaxFrameDelay = time.Second / 60
	// Initialize game
	game := &Game{ships: 5, level: 1}
	game.win = sys.CreateWindow(0, 0, 0, 0, "Asteroids", 1, 1.0)
	game.screen.X = game.win.WidthDp
	game.screen.Y = game.win.HeightDp
	game.player.reset(game.screen)
	game.addAsteroids(10)
	game.win.MakeContextCurrent()
	for sys.Running() {
		game.win.StartFrame()
		game.draw()
		game.win.EndFrame()
		sys.PollEvents()
		game.HandleInput()
		game.Update()
	}
}
