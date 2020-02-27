package backend

import (
	"github.com/the-wendell/go-ants/render"
)

type GameState struct {
	World [][]GameObject
	Ants  []Ant
}

func (g *GameState) RenderState() [][]render.Cell {
	rows := len(g.World)
	cols := len(g.World[0])

	stateMap := make([][]render.Cell, rows)
	for i := 0; i < rows; i++ {
		stateMap[i] = make([]render.Cell, cols)
	}

	for row := 0; row < len(g.World); row++ {
		for col := 0; col < len(g.World[0]); col++ {
			stateMap[row][col] = g.World[row][col].toCell()
		}
	}

	for _, ant := range g.Ants {
		cell := stateMap[ant.CurrentPosition.Y][ant.CurrentPosition.X]
		cell.Forground = ColorAntForeground
		cell.Sprite = SpriteAnt
		stateMap[ant.CurrentPosition.Y][ant.CurrentPosition.X] = cell
	}

	return stateMap
}

func (g *GameState) RunGameStep() {
	for i, ant := range g.Ants {
		currentCell := g.getCell(ant.CurrentPosition)
		nearbyCells := g.getNeighbors(ant.CurrentPosition)
		currentSentTrail, currentSentTrailFound := currentCell.findSentTrail(ant.followingSentTrailID)
		a := &g.Ants[i]

		if currentSentTrailFound {
			a.followSentTrail(*currentSentTrail, *currentCell, nearbyCells)
		} else if currentCell.anySentTrails() {
			a.followSentTrail(*currentCell.strongestSentTrail(), *currentCell, nearbyCells)
		} else if nearbyCells.anySentTrails() {
			nextCell := nearbyCells.cellWithStrongestSentTrail()
			a.moveTo(nextCell.Position)
			a.followingSentTrailID = nextCell.strongestSentTrail().ID
		} else {
			a.wander(nearbyCells)
		}
	}
	// TODO: decayAntSentTrails
	// TODO: propigateHiveSentTrails
}

func (w *GameState) getCell(position Coords) *GameObject {
	return &w.World[position.Y][position.X]
}

func (w *GameState) getNeighbors(position Coords) surroundingCells {
	var left, up, right, down *GameObject

	if position.X > 0 {
		left = &w.World[position.Y][position.X-1]
	} else {
		left = &GameObject{EdgeOfWorld: true}
	}

	if position.Y < len(w.World) {
		up = &w.World[position.Y+1][position.X]
	} else {
		up = &GameObject{EdgeOfWorld: true}
	}

	if position.X < len(w.World[0]) {
		right = &w.World[position.Y][position.X+1]
	} else {
		right = &GameObject{EdgeOfWorld: true}
	}

	if position.Y > 0 {
		down = &w.World[position.Y-1][position.X]
	} else {
		down = &GameObject{EdgeOfWorld: true}
	}

	return surroundingCells{left, up, right, down}
}
