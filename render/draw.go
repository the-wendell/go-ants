package render

import "github.com/gdamore/tcell"

type Cell struct {
	Forground  int
	Background int
	Sprite     rune
}

func DrawScreen(s tcell.Screen, cells [][]Cell) {
	for row := 0; row < len(cells); row++ {
		for col := 0; col < len(cells[0]); col++ {
			cell := cells[row][col]

			st := tcell.StyleDefault
			st = st.Foreground(tcell.Color(cell.Forground))
			st = st.Background(tcell.Color(cell.Background))

			sprite, _, style, _ := s.GetContent(row, col)
			if sprite == cell.Sprite && style == st {
				continue
			}

			s.SetContent(col, row, cell.Sprite, []rune{}, st)
		}
	}
	s.Show()
}
