// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

var (
	colorList     = [...]string{"blue", "green", "yellow", "purple", "violet", "cyan"}
	ctx           js.Value
	startVelocity = Pos{x: 2, y: 6}
)

const (
	screenOffset      int     = 40
	FPS               int64   = 144
	startCircleRadius float32 = 14
	friction          float64 = 0.98
	gravity           float64 = 0.08
	fadeOutSpeed      float32 = 0.995
	amountOfCircles   int     = 60
)

type Circle struct {
	color    string
	pos      Pos
	velocity Pos
	radius   float32
}

type Pos struct {
	x float64
	y float64
}

func handleWindowResize(doc js.Value, canvasEl js.Value, bodySize Pos) {
	resizedW := doc.Get("body").Get("clientWidth").Float()
	resizedH := doc.Get("body").Get("clientHeight").Float()
	if resizedW != bodySize.x || resizedH != bodySize.y {
		bodySize.x, bodySize.y = resizedW, resizedH
		canvasEl.Set("width", bodySize.x)
		canvasEl.Set("height", bodySize.y)
	}
}

func getRandomColor() string {
	randomPos := rand.Intn(len(colorList) - 1)
	return colorList[randomPos]
}

func getRandomPosition(maxPos Pos) Pos {
	maxX := int(maxPos.x) - (screenOffset * 2)
	maxY := int(maxPos.y) - (screenOffset * 2)
	return Pos{
		x: float64(rand.Intn(maxX) + screenOffset),
		y: float64(rand.Intn(maxY) + screenOffset),
	}
}

func generateCircles(bodySize Pos) [amountOfCircles]Circle {
	circles := [amountOfCircles]Circle{}

	for i := range circles {
		circles[i].pos = getRandomPosition(bodySize)
		circles[i].color = getRandomColor()
		circles[i].velocity = startVelocity
		circles[i].radius = startCircleRadius
	}
	return circles
}

func drawArc(circle Circle) {
	ctx.Call("beginPath")
	ctx.Call("arc", circle.pos.x, circle.pos.y, circle.radius, 0, 2*math.Pi)
	ctx.Set("fillStyle", circle.color)
	ctx.Call("fill")
	ctx.Call("closePath")
}

func updateCirclePosition(c *Circle, bodySize Pos) {
	(*c).radius *= fadeOutSpeed
	(*c).velocity.x *= friction
	(*c).velocity.y *= friction
	(*c).velocity.y += gravity

	if c.pos.x < (bodySize.x/2)-float64(screenOffset) {
		(*c).pos.x -= (*c).velocity.x
	} else {
		(*c).pos.x += (*c).velocity.x
	}
	(*c).pos.y += (*c).velocity.y
}

func draw(doc js.Value, canvasEl js.Value, bodySize Pos, circles *[amountOfCircles]Circle) {
	handleWindowResize(doc, canvasEl, bodySize)
	for i := range circles {
		updateCirclePosition(&circles[i], bodySize)
		drawArc(circles[i])
	}

}

// https://codepen.io/chriscourses/pen/Vwamprd

func main() {
	doc := js.Global().Get("document")
	doc.Set("title", time.Now().Second())
	canvasEl := js.Global().Get("document").Call("getElementById", "mycanvas")

	bodySize := Pos{
		x: doc.Get("body").Get("clientWidth").Float(),
		y: doc.Get("body").Get("clientHeight").Float(),
	}
	canvasEl.Set("width", bodySize.x)
	canvasEl.Set("height", bodySize.y)

	ctx = canvasEl.Call("getContext", "2d")
	circles := generateCircles(bodySize)

	timer := time.Now().UnixMilli()
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := time.Now().UnixMilli()
		if now >= timer+(1000/FPS) {
			draw(doc, canvasEl, bodySize, &circles)
			timer = time.Now().UnixMilli()
		}

		js.Global().Call("requestAnimationFrame", renderFrame)

		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame, "asdasd")
	for {
		time.Sleep(1 * time.Hour)
	}
}
