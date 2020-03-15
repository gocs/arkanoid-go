package states

import (
	"fmt"

	"github.com/x-hgg-x/arkanoid-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// AboutState is the game over menu state
type AboutState struct {
	Score     int
	exitMenu  []ecs.Entity
	selection int
}

//
// Menu interface
//

func (st *AboutState) getSelection() int {
	return st.selection
}

func (st *AboutState) setSelection(selection int) {
	st.selection = selection
}

func (st *AboutState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Restart
		return states.Transition{TransType: states.TransSwitch, NewStates: []states.State{&GameplayState{}}}
	case 1:
		// Main Menu
		return states.Transition{TransType: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	case 2:
		// Exit
		return states.Transition{TransType: states.TransQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *AboutState) getMenuIDs() []string {
	return []string{"restart", "main_menu", "exit"}
}

func (st *AboutState) getCursorMenuIDs() []string {
	return []string{"cursor_restart", "cursor_main_menu", "cursor_exit"}
}

//
// State interface
//

// OnPause method
func (st *AboutState) OnPause(world w.World) {}

// OnResume method
func (st *AboutState) OnResume(world w.World) {}

// OnStart method
func (st *AboutState) OnStart(world w.World) {
	st.exitMenu = loader.LoadEntities("assets/metadata/entities/ui/about.toml", world)
}

// OnStop method
func (st *AboutState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.exitMenu...)
}

// Update method
func (st *AboutState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return updateMenu(st, world)
}
