package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 256
	screenHeight = 256

	dotRadius = 4
	dotSpeed  = 2
)

type Game struct {
	x           float32 // 白い点のX座標
	direction   float32 // 移動方向 (+1: 右, -1: 左)
	prevPressed bool    // 前フレームで押されていたか
}

// pressed はタップ長押し・マウス左ボタン・スペースキーのいずれかが
// 押されているかを返す。
func pressed() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(ebiten.AppendTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

func (g *Game) Update() error {
	cur := pressed()

	// 押している間だけ現在の方向へ動かす
	if cur {
		g.x += g.direction * dotSpeed
		// 画面内に収める
		if g.x < dotRadius {
			g.x = dotRadius
		}
		if g.x > screenWidth-dotRadius {
			g.x = screenWidth - dotRadius
		}
	}

	// 離した瞬間に移動方向を反転する (右→左→右…と交互)
	if g.prevPressed && !cur {
		g.direction = -g.direction
	}
	g.prevPressed = cur

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 白い点を画面中央の高さに表示する
	vector.FillCircle(screen, g.x, screenHeight/2, dotRadius, color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle(fmt.Sprintf("%s (v%s)", title, version))
	game := &Game{
		x:         screenWidth / 2, // 画面中央からスタート
		direction: 1,               // 最初は右へ
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
