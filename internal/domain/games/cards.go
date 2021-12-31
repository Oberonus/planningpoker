package games

import "errors"

// Card represents a playing card.
type Card string

// NewCard creates a card with specific type.
func NewCard(typ string) (*Card, error) {
	if typ == "" {
		return nil, errors.New("card type should be provided")
	}

	if len(typ) > 3 {
		return nil, errors.New("card type should be 1-3 chars long")
	}

	card := Card(typ)

	return &card, nil
}

// NewUnrevealedCard creates a new unrevealed card, which type is not determined.
func NewUnrevealedCard() Card {
	return "*"
}

// Type returns a type of the card.
func (c Card) Type() string {
	return string(c)
}

// CardsDeck represents a deck of cards.
type CardsDeck struct {
	name  string
	cards []Card
}

// NewCardsDeck creates a new named deck of cards.
func NewCardsDeck(name string, cards []Card) (*CardsDeck, error) {
	if name == "" {
		return nil, errors.New("name should be provided")
	}
	if len(cards) == 0 {
		return nil, errors.New("cards should be provided")
	}

	return &CardsDeck{
		name:  name,
		cards: cards,
	}, nil
}

// Name returns cards deck name.
func (d CardsDeck) Name() string {
	return d.name
}

// Cards return all cards from the deck.
func (d CardsDeck) Cards() []Card {
	return d.cards
}

// IsInDeck checks if specific card exists in the deck.
func (d CardsDeck) IsInDeck(card Card) bool {
	for _, c := range d.cards {
		if c == card {
			return true
		}
	}

	return false
}
