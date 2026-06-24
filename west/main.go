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

	bulletWidth  = 3 // 弾(長方形)の幅
	bulletHeight = 8 // 弾(長方形)の高さ
	bulletSpeed  = 4 // 弾の上方向への速度
	fireInterval = 6 // 連射間隔 (フレーム数。小さいほど高速連射)
)

var (
	gray = color.RGBA{0x0a, 0x0a, 0x0a, 0xff}
)

// bullet は自機から発射される弾を表す。x, y は長方形の中心座標。
type bullet struct {
	x, y float32
}

type Game struct {
	x           float32  // 白い点(自機)のX座標
	direction   float32  // 移動方向 (+1: 右, -1: 左)
	prevPressed bool     // 前フレームで押されていたか
	bullets     []bullet // 飛んでいる弾
	fireCount   int      // 連射クールダウン用カウンタ
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

	// 自機から長方形の弾を自動連射する。
	// fireInterval フレームに 1 発、自機の位置から発射する。
	g.fireCount--
	if g.fireCount <= 0 {
		g.bullets = append(g.bullets, bullet{x: g.x, y: screenHeight / 2})
		g.fireCount = fireInterval
	}

	// 弾を上方向へ移動し、画面外に出たものを取り除く。
	alive := g.bullets[:0]
	for _, b := range g.bullets {
		b.y -= bulletSpeed
		if b.y+bulletHeight/2 < 0 {
			continue // 画面上端より外に出たので消す
		}
		alive = append(alive, b)
	}
	g.bullets = alive

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を明示的に黒で塗りつぶす。
	screen.Fill(gray)

	// 自機(白い点)を境界に、進行方向側の背景をグレーで塗る。
	// (右へ進行中なら右側、左へ進行中なら左側がグレー)
	if g.direction > 0 {
		vector.DrawFilledRect(screen, g.x, 0, screenWidth-g.x, screenHeight, color.Black, false)
	} else {
		vector.DrawFilledRect(screen, 0, 0, g.x, screenHeight, color.Black, false)
	}

	// 飛んでいる弾(長方形)を描画する。x, y を中心とするので左上に補正する。
	for _, b := range g.bullets {
		vector.DrawFilledRect(screen, b.x-bulletWidth/2, b.y-bulletHeight/2, bulletWidth, bulletHeight, color.White, true)
	}

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
