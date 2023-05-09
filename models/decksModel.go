package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Help is the model that inserts a help objects retrived or inserted into the DB
type Decks struct {
	ID        primitive.ObjectID `bson:"_id"`
	Deck_ID   string             `json:"deck_id"`
	Cards     []Card             `json:"cards"`
	Shuffled  bool               `json:"shuffled"`
	Remaining int                `json:"remaining"`
}

type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}
