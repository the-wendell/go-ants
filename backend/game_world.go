package backend

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/the-wendell/go-ants/render"
)

const SpriteAnt = 'a'
const SpriteTunnel = ' '
const SpriteDirt = ' '
const SpriteFood = '@'

const ColorAntForeground = 100
const ColorBackgroundDirt = 30
const ColorForegroundDirt = 40
const ColorBackgroundTunnel = 50
const ColorForegroundTunnel = 0

const sentFood = "FOOD"

var emptyUUID uuid.UUID

type Ant struct {
	CurrentPosition      Coords
	PreviousPosition     Coords
	followingSentTrailID uuid.UUID
	trailingSentTrailID  uuid.UUID
	HasFood              bool
}

func (a *Ant) moveTo(newPosition Coords) {
	a.PreviousPosition = a.CurrentPosition
	a.CurrentPosition = newPosition
}

func (a *Ant) wander(nearbyCells surroundingCells) {
	var unvisitedCells []GameObject

	// TODO: if currentCell hasFood? { pickUpFood; followSentTrail(foodStorageSentTrail); return }

	for _, cell := range nearbyCells.toSlice() {
		if cell.Position != a.PreviousPosition {
			unvisitedCells = append(unvisitedCells, *cell)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(unvisitedCells), func(i, j int) { unvisitedCells[i], unvisitedCells[j] = unvisitedCells[j], unvisitedCells[i] })

	for _, cell := range unvisitedCells {
		if cell.canWalkTo() {
			a.moveTo(cell.Position)
			return
		}
	}
	a.moveTo(a.PreviousPosition)
}

func (a *Ant) followSentTrail(trail sentTrail, currentCell GameObject, nearbyCells surroundingCells) {
	weakestSentTrail, _ := currentCell.findSentTrail(a.followingSentTrailID)
	cellWithWeakestSentTrail := &currentCell

	for _, cell := range nearbyCells.toSlice() {
		trail, found := cell.findSentTrail(a.followingSentTrailID)
		if found && trail.strength < weakestSentTrail.strength {
			cellWithWeakestSentTrail = cell
		}
	}

	if cellWithWeakestSentTrail != &currentCell {
		// TODO:
		// if a.HasFood { a.leaveSentTrail }
		a.moveTo(cellWithWeakestSentTrail.Position)
	}
	// TODO:
	// else if currentCell == foodStorageCell { dropFood; a.CurrentSentTrailID = emptyUUID }
	// else { pickUpFood; followSentTrail(foodStorageSentTrail) }
}

type Coords struct {
	X int
	Y int
}

type sentTrail struct {
	ID         uuid.UUID
	strength   int
	typeOfSent string
	hiveTrail  bool
}

type surroundingCells struct {
	left  *GameObject
	up    *GameObject
	right *GameObject
	down  *GameObject
}

func (c *surroundingCells) toSlice() []*GameObject {
	return []*GameObject{c.left, c.up, c.right, c.down}
}

func (c surroundingCells) anySentTrails() bool {
	return c.left.anySentTrails() || c.up.anySentTrails() || c.right.anySentTrails() || c.down.anySentTrails()
}

func (c *surroundingCells) cellWithStrongestSentTrail() *GameObject {
	cells := c.toSlice()
	strongest := cells[0]

	for i, cell := range cells {
		if cell.strongestSentTrail().strength > strongest.strongestSentTrail().strength {
			strongest = cells[i]
		}
	}

	return strongest
}

type GameObject struct {
	EdgeOfWorld     bool
	Solid           bool
	SentTrails      []sentTrail
	Integrity       int
	Sprite          rune
	ForegroundColor int
	BackgroundColor int
	Position        Coords
	FoodCount       int
}

func (g GameObject) toCell() render.Cell {
	return render.Cell{
		Forground:  g.ForegroundColor,
		Background: g.BackgroundColor,
		Sprite:     g.Sprite,
	}
}

func (g GameObject) canWalkTo() bool {
	return !g.EdgeOfWorld && !g.Solid
}

func (g *GameObject) findSentTrail(id uuid.UUID) (*sentTrail, bool) {
	for i, sentTrail := range g.SentTrails {
		if sentTrail.ID == id {
			return &g.SentTrails[i], true
		}
	}
	return nil, false
}

func (g *GameObject) anySentTrails() bool {
	return len(g.SentTrails) > 0
}

func (g *GameObject) strongestSentTrail() *sentTrail {
	strongestSentTrail := &g.SentTrails[0]

	for i, sentTrail := range g.SentTrails {
		if sentTrail.strength > strongestSentTrail.strength {
			strongestSentTrail = &g.SentTrails[i]
		}
	}

	return strongestSentTrail
}
