package states

import (
	"fmt"
	"strconv"

	"github.com/x-hgg-x/arkanoid-go/lib/binconv"
	"github.com/x-hgg-x/arkanoid-go/lib/components"
	"github.com/x-hgg-x/arkanoid-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct {
	Score             int
	levelCompleteMenu []ecs.Entity
	selection         int
}

//
// Menu interface
//

func (st *LevelCompleteState) getSelection() int {
	return st.selection
}

func (st *LevelCompleteState) setSelection(selection int) {
	st.selection = selection
}

func (st *LevelCompleteState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Main Menu
		return states.Transition{TransType: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *LevelCompleteState) getMenuIDs() []string {
	return []string{"main_menu"}
}

func (st *LevelCompleteState) getCursorMenuIDs() []string {
	return []string{"cursor_main_menu"}
}

//
// State interface
//

// OnPause method
func (st *LevelCompleteState) OnPause(world w.World) {}

// OnResume method
func (st *LevelCompleteState) OnResume(world w.World) {}

// OnStart method
func (st *LevelCompleteState) OnStart(world w.World) {
	st.levelCompleteMenu = loader.LoadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world)

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "score" {
			text.Text = fmt.Sprintf("SCORE: %d", st.Score)
		}
	}))

	addToDB(st.Score)
}

func addToDB(score int) {
	p, err := components.NewPersist("arkanoid.db", "Score")
	if err != nil {
		fmt.Println("error from creating", err)
		return
	}
	defer p.Close()

	key := []byte("scores")

	d, err := p.ViewList(key)
	if err != nil {
		fmt.Println("error from getting list", err)
		return
	}
	arrKey := append(key, binconv.Itob(len(d))...)
	if p.Update(arrKey, []byte(strconv.Itoa(score))) != nil {
		fmt.Println("error from updating list", err)
		return
	}
	d, err = p.ViewList(key)
	if err != nil {
		fmt.Println("error from getting list", err)
		return
	}
	for i, v := range d {
		fmt.Printf("data[%d]: %s\n", i, string(v))
	}
}

// OnStop method
func (st *LevelCompleteState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.levelCompleteMenu...)
}

// Update method
func (st *LevelCompleteState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return updateMenu(st, world)
}
