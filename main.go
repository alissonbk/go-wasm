// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

var (
	ctx          js.Value
	circleRadius       = 20
	offset             = float64(circleRadius * 2)
	colorList          = [...]string{"asd", "qwe", "asdasd", "qweqwe", "asdasd"}
	FPS          int64 = 10
)

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

func getRandomColor() {

}

func getRandomPosition(maxPos Pos) Pos {
	maxX := (maxPos.x - offset) - 20
	maxY := (maxPos.y - offset) - 20
	return Pos{
		x: rand.Float64()*(maxX+offset) + offset,
		y: rand.Float64()*(maxY+offset) + offset,
	}
}

func draw(doc js.Value, canvasEl js.Value, bodySize Pos) {
	handleWindowResize(doc, canvasEl, bodySize)
	randomPos := getRandomPosition(bodySize)

	ctx.Call("beginPath")
	ctx.Call("arc", randomPos.x, randomPos.y, circleRadius, 0, 2*math.Pi)
	ctx.Call("fill")
	ctx.Call("restore")
}

//https://codepen.io/chriscourses/pen/Vwamprd

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

	doc.Get("body").Get("style").Set("backgroundColor", "red")

	ctx = canvasEl.Call("getContext", "2d")

	timer := time.Now().UnixMilli()
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := time.Now().UnixMilli()
		if now >= timer+(1000/FPS) {
			draw(doc, canvasEl, bodySize)
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
