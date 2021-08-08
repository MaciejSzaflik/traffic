package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
	"unsafe"

	"github.com/MaciejSzaflik/traffic/gfx"
	"github.com/MaciejSzaflik/traffic/gridMap"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 800
const windowHeight = 800

const checkTime = time.Millisecond
const frameDur = time.Second / 60

type updateFunc func(time.Duration)

type mainLoop struct {
	lastTime time.Time
	updaters []updateFunc
}

func (ml *mainLoop) addUpdater(update updateFunc) int {
	ml.updaters = append(ml.updaters, update)
	return len(ml.updaters) - 1
}

func (ml *mainLoop) update() {
	delta := time.Since(ml.lastTime)
	ml.lastTime = time.Now()
	for _, u := range ml.updaters {
		u(delta)
	}
}

func newMainLoop() *mainLoop {
	return &mainLoop{
		lastTime: time.Now(),
	}
}

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	rand.Seed(time.Now().Unix())

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "lol", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)

	err = programLoop(window)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTriangleVAO(vertices []float32, indices []uint32) uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return VAO
}

func loadShader() (*gfx.Program, error) {
	vertShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	return gfx.NewProgram(vertShader, fragShader)
}

func programLoop(window *glfw.Window) error {
	mainLoop := newMainLoop()

	shaderProgram, err := loadShader()
	if err != nil {
		return err
	}
	defer shaderProgram.Delete()

	gridSize := 100
	grid := gridMap.NewGridMap(1.5/float32(gridSize), gridSize)
	grid.VAO = CreateTriangleVAO(grid.Vertices, grid.Indices)
	grid.VAODirty = false

	fpsCounter := NewFpsCounter(100)

	agentDirector := gridMap.NewAgentDirectorRandom(grid, 1, 2)
	mainLoop.addUpdater(agentDirector.Update)

	mainLoop.addUpdater(fpsCounter.update)

	mainLoop.addUpdater(func(d time.Duration) {
		if grid.VAODirty {
			grid.VAO = CreateTriangleVAO(grid.Vertices, grid.Indices)
			grid.VAODirty = false
		}
	})

	mainLoop.addUpdater(func(d time.Duration) {
		glfw.PollEvents()

		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shaderProgram.Use()

		gl.BindVertexArray(grid.VAO)
		gl.DrawElements(gl.TRIANGLES, grid.VerticesCount, gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		window.SwapBuffers()
	})

	currentTime := time.Now()

	for range time.Tick(checkTime) {
		if window.ShouldClose() {
			break
		}

		if time.Since(currentTime) < frameDur {
			continue
		}

		currentTime = time.Now()

		mainLoop.update()
	}

	return nil
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
