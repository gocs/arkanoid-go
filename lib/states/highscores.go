package states

import (
	"fmt"
	"strings"

	"github.com/x-hgg-x/arkanoid-go/lib/components"
	"github.com/x-hgg-x/arkanoid-go/lib/loader"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// HighscoresState is the game over menu state
type HighscoresState struct {
	goToMenu  []ecs.Entity
	selection int
}

//
// Menu interface
//

func (st *HighscoresState) getSelection() int {
	return st.selection
}

func (st *HighscoresState) setSelection(selection int) {
	st.selection = selection
}

func (st *HighscoresState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Main Menu
		return states.Transition{TransType: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}
	// will pass here if the menuIds/toml are inconsistent
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *HighscoresState) getMenuIDs() []string {
	return []string{"main_menu"}
}

func (st *HighscoresState) getCursorMenuIDs() []string {
	return []string{"cursor_main_menu"}
}

//
// State interface
//

// OnPause method
func (st *HighscoresState) OnPause(world w.World) {}

// OnResume method
func (st *HighscoresState) OnResume(world w.World) {}

// OnStart method
func (st *HighscoresState) OnStart(world w.World) {
	st.goToMenu = loader.LoadEntities("assets/metadata/entities/ui/highscores.toml", world)
	
	scores := getHighscore()

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "highscores" {
			
			var sb strings.Builder
			for i, score := range scores {
				sb.WriteString(fmt.Sprintf("%d: %s\n", i, score))
			}
			text.Text = fmt.Sprintf("Scores\n%s", sb.String())
		}
	}))
}

func getHighscore() []string {
	p, err := components.NewPersist("arkanoid.db", "Score")
	if err != nil {
		fmt.Println("error from creating", err)
		return nil
	}
	defer p.Close()

	d, err := p.ViewList([]byte("scores"))
	if err != nil {
		fmt.Println("error from getting list", err)
		return nil
	}
	var scores []string
	for _, v := range d {
		scores = append(scores, string(v))
	}
	return scores
}

// OnStop method
func (st *HighscoresState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.goToMenu...)
}

// Update method
func (st *HighscoresState) Update(world w.World, screen *ebiten.Image) states.Transition {
	return updateMenu(st, world)
}
