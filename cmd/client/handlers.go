package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

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

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(move gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)

		if outcome == gamelogic.MoveOutcomeMakeWar {
			war := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
			}
			pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+move.Player.Username, war)
			return pubsub.NackRequeue
		} else if outcome == gamelogic.MoveOutComeSafe {
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

func handlerWar(gs *gamelogic.GameState) func(recognition gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(recognition gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(recognition)

		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		} else if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		} else if outcome == gamelogic.WarOutcomeOpponentWon {
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeYouWon {
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeDraw {
			return pubsub.Ack
		} else {
			fmt.Println("Error acknowledging outcome")
			return pubsub.NackDiscard
		}
	}
}
