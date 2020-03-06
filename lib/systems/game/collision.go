package gamesystem

import (
	"math"

	c "arkanoid/lib/components"
	"arkanoid/lib/ecs"
	w "arkanoid/lib/ecs/world"
	m "arkanoid/lib/math"
	"arkanoid/lib/resources"

	"github.com/ByteArena/box2d"
)

// CollisionSystem manages collisions
func CollisionSystem(world w.World) {
	paddles := ecs.Join(world.Components.Paddle, world.Components.Transform)
	if paddles.Empty() {
		return
	}
	firstPaddle := ecs.Entity(paddles.Next(-1))
	paddle := world.Components.Paddle.Get(firstPaddle).(*c.Paddle)
	paddleTranslation := world.Components.Transform.Get(firstPaddle).(*c.Transform).Translation

	// Set paddle body transform
	paddle.Body.SetTransform(box2d.MakeB2Vec2(paddleTranslation.X/resources.B2PixelRatio, paddleTranslation.Y/resources.B2PixelRatio), 0)

	// Set balls body transform
	ecs.Join(world.Components.Ball, world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)
		ballTranslation := world.Components.Transform.Get(entity).(*c.Transform).Translation
		ball.Body.SetTransform(box2d.MakeB2Vec2(ballTranslation.X/resources.B2PixelRatio, ballTranslation.Y/resources.B2PixelRatio), 0)
	}))

	// Find contacts
	collisionWorld := world.Resources.Game.CollisionWorld
	collisionWorld.M_contactManager.FindNewContacts()
	collisionWorld.M_contactManager.Collide()

	// Get list of contacts with normals and bodies
	contactsNormal := []box2d.B2Vec2{}
	contactsBodies := [][2]*box2d.B2Body{}
	for contactList := collisionWorld.GetContactList(); contactList != nil; contactList = contactList.GetNext() {
		wm := box2d.MakeB2WorldManifold()
		contactList.GetWorldManifold(&wm)
		// Test if normal is defined
		if (wm.Normal != box2d.B2Vec2{}) {
			contactsNormal = append(contactsNormal, wm.Normal)
			contactsBodies = append(contactsBodies, [2]*box2d.B2Body{contactList.GetFixtureA().GetBody(), contactList.GetFixtureB().GetBody()})
		}
	}

	// Loop on balls
	ecs.Join(world.Components.Ball, world.Components.StickyBall.Not(), world.Components.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		ball := world.Components.Ball.Get(entity).(*c.Ball)
		ballTranslation := &world.Components.Transform.Get(entity).(*c.Transform).Translation

		// Bounce at the top, left and right of the arena
		if ballTranslation.X <= ball.Radius {
			ball.Direction.X = math.Abs(ball.Direction.X)
		}
		if ballTranslation.X >= float64(world.Resources.ScreenDimensions.Width)-ball.Radius {
			ball.Direction.X = -math.Abs(ball.Direction.X)
		}
		if ballTranslation.Y >= float64(world.Resources.ScreenDimensions.Height)-ball.Radius {
			ball.Direction.Y = -math.Abs(ball.Direction.Y)
		}

		// Bounce at the paddle
		bounced := false
		for iContact := range contactsBodies {
			if contactsBodies[iContact] == [2]*box2d.B2Body{paddle.Body, ball.Body} || contactsBodies[iContact] == [2]*box2d.B2Body{ball.Body, paddle.Body} {
				bounced = true
				minValue := -math.Pi / 3
				maxValue := math.Pi / 3
				angle := math.Min(math.Max((paddleTranslation.X-ballTranslation.X)/paddle.Width*math.Pi, minValue), maxValue)
				ball.Direction = m.Vector2{X: math.Sin(-angle), Y: math.Cos(angle)}
			}
		}

		// Lose a life when ball reach the bottom of the arena
		if ballTranslation.Y <= ball.Radius && !bounced {
			entity.AddComponent(world.Components.StickyBall, &c.StickyBall{Period: 2})
			*ballTranslation = m.Vector2{X: paddleTranslation.X, Y: paddle.Height + ball.Radius}
		}

		// Bounce at the blocks
		blockNormals := []m.Vector2{}
		blockbodies := []*box2d.B2Body{}
		for iContact := range contactsNormal {
			// Normal is pointing towards block exterior
			var blockBody *box2d.B2Body
			if contactsBodies[iContact][0].GetUserData().(ecs.Entity).HasComponent(world.Components.Block) && contactsBodies[iContact][1] == ball.Body {
				blockBody = contactsBodies[iContact][0]
				blockNormals = append(blockNormals, m.Vector2{X: contactsNormal[iContact].X, Y: contactsNormal[iContact].Y})
			} else if contactsBodies[iContact][1].GetUserData().(ecs.Entity).HasComponent(world.Components.Block) && contactsBodies[iContact][0] == ball.Body {
				blockBody = contactsBodies[iContact][1]
				blockNormals = append(blockNormals, m.Vector2{X: -contactsNormal[iContact].X, Y: -contactsNormal[iContact].Y})
			}

			if blockBody != nil {
				blockbodies = append(blockbodies, blockBody)
			}
		}

		if len(blockNormals) == 0 {
			// No colliding blocks
			return
		} else if len(blockNormals) >= 3 {
			// 3 or more colliding blocks: reverse ball direction
			ball.Direction.X *= -1
			ball.Direction.Y *= -1
			return
		}

		// 1 or 2 colliding blocks: ball is reflected wrt the contact normal
		var incidenceAngle float64
		if len(blockNormals) == 1 {
			// 1 colliding block: use computed normal
			incidenceAngle = math.Atan2(-ball.Direction.Perp(blockNormals[0]), -ball.Direction.Dot(blockNormals[0]))
		} else if len(blockNormals) == 2 {
			// 2 colliding blocks: define normal as the perpendicular of the line between blocks center (towards ball)
			positions := []box2d.B2Vec2{blockbodies[0].GetPosition(), blockbodies[1].GetPosition()}
			positionDiff := m.Vector2{X: positions[1].X - positions[0].X, Y: positions[1].Y - positions[0].Y}
			positionDiffPerp := m.Vector2{X: -positionDiff.Y, Y: positionDiff.X}
			ballLocalWorldTranslation := m.Vector2{
				X: ballTranslation.X/resources.B2PixelRatio - positions[0].X,
				Y: ballTranslation.Y/resources.B2PixelRatio - positions[0].Y,
			}

			var normal m.Vector2
			if positionDiffPerp.Dot(ballLocalWorldTranslation) > 0 {
				normal = m.Vector2{X: positionDiffPerp.X, Y: positionDiffPerp.Y}
			} else {
				normal = m.Vector2{X: -positionDiffPerp.X, Y: -positionDiffPerp.Y}
			}

			// Normalize normal
			normalNorm := normal.Norm()
			normal.X /= normalNorm
			normal.Y /= normalNorm

			incidenceAngle = math.Atan2(-ball.Direction.Perp(normal), -ball.Direction.Dot(normal))
		}

		// Compute ball reflection
		sin, cos := math.Sincos(2 * incidenceAngle)
		ball.Direction = m.Vector2{
			X: -ball.Direction.X*cos + ball.Direction.Y*sin,
			Y: -ball.Direction.X*sin - ball.Direction.Y*cos,
		}

		// Renormalize ball direction
		ballNorm := ball.Direction.Norm()
		ball.Direction.X /= ballNorm
		ball.Direction.Y /= ballNorm
	}))
}
