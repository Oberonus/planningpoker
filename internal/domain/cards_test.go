package domain_test

import (
	"planningpoker/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCard(t *testing.T) {
	sucCard := domain.Card("xs")

	testCases := map[string]struct {
		typ     string
		expCard *domain.Card
		expErr  string
	}{
		"success": {
			typ:     sucCard.Type(),
			expCard: &sucCard,
		},
		"fail on no type": {
			typ:    "",
			expErr: "card type should be provided",
		},
		"fail on too long type": {
			typ:    "qwert",
			expErr: "card type should be 1-3 chars long",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			card, err := domain.NewCard(tt.typ)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expCard, card)
		})
	}
}

func TestNewUnrevealedCard(t *testing.T) {
	card := domain.NewUnrevealedCard()
	assert.Equal(t, "*", card.Type())
}

func TestNewCardsDeck(t *testing.T) {
	card, err := domain.NewCard("xs")
	require.NoError(t, err)

	testCases := map[string]struct {
		name   string
		cards  []domain.Card
		expErr string
	}{
		"success": {
			name:  "name",
			cards: []domain.Card{*card},
		},
		"fail on empty name": {
			name:   "",
			cards:  []domain.Card{*card},
			expErr: "name should be provided",
		},
		"fail on empty cards": {
			name:   "name",
			cards:  nil,
			expErr: "cards should be provided",
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			deck, err := domain.NewCardsDeck(tt.name, tt.cards)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Nil(t, deck)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, deck)
				assert.Equal(t, tt.name, deck.Name())
				assert.Equal(t, tt.cards, deck.Cards())
			}
		})
	}
}

func TestCardsDeck_IsInDeck(t *testing.T) {
	card1, err := domain.NewCard("XS")
	require.NoError(t, err)
	card2, err := domain.NewCard("S")
	require.NoError(t, err)
	deck, err := domain.NewCardsDeck("foo", []domain.Card{*card1, *card2})
	require.NoError(t, err)

	testCases := map[string]struct {
		deck     domain.CardsDeck
		card     string
		expFound bool
	}{
		"success found": {
			card:     "XS",
			expFound: true,
		},
		"success not found": {
			card:     "L",
			expFound: false,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			card, err := domain.NewCard(tt.card)
			require.NoError(t, err)
			found := deck.IsInDeck(*card)
			assert.Equal(t, tt.expFound, found)
		})
	}
}
