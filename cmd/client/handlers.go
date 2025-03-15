package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(state routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(state)
		log.Println("Pause handler triggered; acknowledging")
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(move gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)

		if outcome == gamelogic.MoveOutComeSafe || outcome == gamelogic.MoveOutcomeMakeWar {
			log.Println("Move handler acknowledges; outcome:", outcome)
			return pubsub.Ack
		} else if outcome == gamelogic.MoveOutcomeSamePlayer {
			log.Println("Move handler rejects and discards (same player); outcome:", outcome)
			return pubsub.NackDiscard
		}

		log.Println("Move handler rejects and discards (unexpected); outcome:", outcome)
		return pubsub.NackDiscard
	}
}
