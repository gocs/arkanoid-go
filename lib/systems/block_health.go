package gamesystem

import (
	gc "arkanoid/lib/components"
	"arkanoid/lib/resources"

	ecs "github.com/x-hgg-x/goecs"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// BlockHealthSystem manages block health
func BlockHealthSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	for _, blockCollisionEvent := range gameResources.Events.BlockCollisionEvents {
		block := gameComponents.Block.Get(blockCollisionEvent.Entity).(*gc.Block)
		sprite := world.Components.Engine.SpriteRender.Get(blockCollisionEvent.Entity).(*ec.SpriteRender)

		block.Health--
		if block.Health > 0 {
			sprite.SpriteNumber += 6
		} else {
			gameResources.CollisionWorld.DestroyBody(block.Body)
			world.Manager.DeleteEntity(blockCollisionEvent.Entity)
			gameResources.Events.ScoreEvents = append(gameResources.Events.ScoreEvents, resources.ScoreEvent{Score: 50})
		}
	}
	gameResources.Events.BlockCollisionEvents = nil

	if ecs.Join(gameComponents.Block).Empty() {
		gameResources.StateEvent = resources.StateEventLevelComplete
	}
}
