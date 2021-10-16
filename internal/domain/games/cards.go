package games

import "errors"

type Card string

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

func NewUnrevealedCard() Card {
	return "*"
}

func (c Card) Type() string {
	return string(c)
}

type CardsDeck struct {
	name  string
	cards []Card
}

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

func (d CardsDeck) Name() string {
	return d.name
}

func (d CardsDeck) Cards() []Card {
	return d.cards
}

func (d CardsDeck) IsInDeck(card Card) bool {
	for _, c := range d.cards {
		if c == card {
			return true
		}
	}

	return false
}
