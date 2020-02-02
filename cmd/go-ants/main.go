package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/the-wendell/go-ants/backend"
	"github.com/the-wendell/go-ants/render"
)

func main() {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	// cells := [][]render.Cell{
	// 	{{10, 0, 'A'}, {10, 0, 'A'}, {10, 0, 'A'}, {10, 0, 'A'}},
	// 	{{70, 20, 'A'}, {80, 30, 'A'}, {90, 40, 'A'}, {100, 50, 'A'}},
	// 	{{70, 20, 'A'}, {80, 30, 'A'}, {90, 40, 'A'}, {100, 50, 'A'}},
	// 	{{70, 20, 'A'}, {80, 30, 'A'}, {90, 40, 'A'}, {100, 50, 'A'}},
	// 	{{70, 20, 'A'}, {80, 30, 'A'}, {90, 40, 'A'}, {100, 50, 'A'}},
	// 	{{1, 10, 'A'}, {1, 10, 'A'}, {1, 10, 'A'}, {1, 10, 'A'}},
	// }

	wall := backend.GameObject{Solid: true, Sprite: backend.SpriteDirt, ForegroundColor: backend.ColorForegroundDirt, BackgroundColor: backend.ColorBackgroundDirt}
	tunn := backend.GameObject{Solid: false, Sprite: backend.SpriteTunnel, ForegroundColor: backend.ColorForegroundTunnel, BackgroundColor: backend.ColorBackgroundTunnel}

	gameState := [][]backend.GameObject{
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, wall},
		{wall, wall, wall, wall, wall, wall, wall, tunn, tunn, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, tunn, tunn, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, tunn, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, tunn, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, tunn, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, tunn, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
		{wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall, wall},
	}

	ants := []backend.Ant{
		{CurrentPosition: backend.Coords{2, 6}, PreviousPosition: backend.Coords{3, 6}},
		{CurrentPosition: backend.Coords{6, 6}, PreviousPosition: backend.Coords{5, 6}},
	}

	for row := 0; row < len(gameState); row++ {
		for col := 0; col < len(gameState[0]); col++ {
			gameState[row][col].Position = backend.Coords{col, row}
		}
	}

	timer := time.NewTimer(time.Second / 2)

	state := backend.GameState{Ants: ants, World: gameState}
	render.DrawScreen(s, state.RenderState())
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-timer.C:
			state.RunGameStep()
			render.DrawScreen(s, state.RenderState())
			timer.Reset(time.Second / 2)
		}
	}

	s.Fini()
}
