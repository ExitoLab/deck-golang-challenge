package controllers

import (
	"context"
	"fmt"
	"golang-execrise/database"
	"golang-execrise/models"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Card represents a standard playing card
type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}

// Deck represents a standard 52-card deck of playing cards
type Deck struct {
	Deck_ID   string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

var validate = validator.New()
var deckCollection *mongo.Collection = database.OpenCollection(database.Client, "decks")

//

// ranks represents the standard ranks of playing cards
var ranks = []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

// suits represents the standard suits of playing cards
var suits = []string{"Spades", "Diamonds", "Clubs", "Hearts"}

func isValidRank(rank string) bool {
	validRanks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	for _, r := range validRanks {
		if r == rank {
			return true
		}
	}
	return false
}

func isValidSuit(suit string) bool {
	validSuits := []string{"S", "D", "C", "H"}
	for _, s := range validSuits {
		if s == suit {
			return true
		}
	}
	return false
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func convertToModelsCards(cards []Card) []models.Card {
	var modelsCards []models.Card
	for _, card := range cards {
		modelsCard := models.Card{
			Suit: card.Suit,
			Rank: card.Rank,
			Code: card.Code,
		}
		modelsCards = append(modelsCards, modelsCard)
	}
	return modelsCards
}

func lookupSuits(suit string) string {
	m := map[string]string{
		"K": "KING",
		"S": "SPADES",
		"D": "DIAMONDS",
		"C": "CLUBS",
	}

	val, ok := m[suit]
	if ok {
		return val
	}
	return suit
}

// CreateNewDeckHandler creates a new deck of playing cards
func CreateNewDeckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var decks models.Decks
		defer cancel()

		validationErr := validate.Struct(decks)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":         "",
				"responseCode": http.StatusBadRequest,
				"message":      validationErr.Error(),
			})
			return
		}

		// Parse query parameters
		shuffled, _ := strconv.ParseBool(c.DefaultQuery("shuffled", "false"))
		cardList := c.Query("cards")

		// Create new deck and set ID
		deck := Deck{Deck_ID: uuid.New().String(), Shuffled: shuffled}

		// Add cards to deck based on query parameters
		if len(cardList) > 0 {
			// Split card list into slice
			cards := strings.Split(cardList, ",")
			for _, cardStr := range cards {
				// Parse card rank and suit
				rank := strings.TrimSpace(cardStr[:len(cardStr)-1])
				suitParts := strings.TrimSpace(cardStr[len(cardStr)-1:])
				suit := lookupSuits(suitParts)

				if !isValidRank(rank) || !isValidSuit(suitParts) {
					c.JSON(http.StatusNotFound, gin.H{"Message": "You supplied an invalid card"})
					return
				}

				// Add card to deck if valid rank and suit
				if isValidRank(rank) && isValidSuit(suitParts) {
					code := rank[0:1] + suit[0:1]
					card := Card{Rank: rank, Suit: suit, Code: code}
					deck.Cards = append(deck.Cards, card)
					deck.Remaining = len(cards)
				}
			}
		} else {
			// Add all standard cards to deck
			for _, rank := range ranks {
				for _, suit := range suits {
					code := rank[0:1] + suit[0:1]
					card := Card{Rank: rank, Suit: suit, Code: code}
					deck.Cards = append(deck.Cards, card)
					deck.Remaining = 52
				}
			}
		}

		// Shuffle deck if requested
		if shuffled {
			deck.shuffle()
		}

		decks.ID = primitive.NewObjectID()
		decks.Remaining = deck.Remaining
		decks.Shuffled = deck.Shuffled
		decks.Deck_ID = deck.Deck_ID
		decks.Cards = convertToModelsCards(deck.Cards)

		_, insertErr := deckCollection.InsertOne(ctx, decks)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": insertErr.Error()})
			return
		}
		c.JSON(http.StatusOK, deck)
	}
}

// FindDeckByDeckIDHandler opens an existing deck of playing cards
func FindDeckByDeckIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var decks models.Decks
		defer cancel()

		if c.Param("deck_id") == "" {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Deck id is not supplied!"})
			return
		}

		// Get deck ID from path parameter
		deckID := c.Param("deck_id")
		deckIDDetails := deckCollection.FindOne(ctx, bson.M{"deck_id": deckID}).Decode(&decks)
		if deckIDDetails != nil {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Deck is not found!"})
			return
		}

		// Return JSON
		c.JSON(http.StatusOK, decks)
	}
}

func Draw(d *models.Decks, count int) ([]models.Card, error) {
	if count > len(d.Cards) {
		return nil, fmt.Errorf("not enough cards in deck")
	}

	cards := make([]models.Card, count)
	for i := 0; i < count; i++ {
		cards[i] = d.Cards[i]
	}

	d.Cards = d.Cards[count:]

	return cards, nil
}

func DrawAndUpdateDeck(d *models.Decks, count int) (*models.Decks, []models.Card, error) {
	cards, err := Draw(d, count)
	if err != nil {
		return nil, nil, err
	}

	return d, cards, nil
}

func DrawCardDeckByDeckIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var decks models.Decks
		defer cancel()

		if c.Param("deck_id") == "" {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Deck id is not supplied!"})
			return
		}

		// Get deck ID from path parameter
		deck_id := c.Param("deck_id")
		deckIDDetails := deckCollection.FindOne(ctx, bson.M{"deck_id": deck_id}).Decode(&decks)
		if deckIDDetails != nil {
			c.JSON(http.StatusNotFound, gin.H{"Message": "Deck is not found!"})
			return
		}

		// Parse count parameter
		countStr := c.DefaultQuery("count", "1")
		countInt, err := strconv.Atoi(countStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		if countInt > 52 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can not supply card greater than 52"})
			return
		}

		var drawnCards []models.Card
		var updatedDecks *models.Decks

		updatedDecks, drawnCards, err = DrawAndUpdateDeck(&decks, countInt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentRemaining := decks.Remaining - countInt
		fmt.Println(currentRemaining)

		//Update deck with the remaining cards and remaining
		updateObj := bson.M{
			"$set": bson.M{
				"remaining": currentRemaining,
				"cards":     updatedDecks.Cards,
			},
		}

		filter := bson.M{"deck_id": deck_id}

		result, err := deckCollection.UpdateOne(
			ctx,
			filter,
			updateObj,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result != nil {
			resultDrawCards := make(map[string]interface{})
			resultDrawCards["Cards"] = drawnCards
			c.JSON(http.StatusOK, resultDrawCards)
		}
	}
}
