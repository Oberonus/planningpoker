package games_test

import (
	"testing"

	"planningpoker/internal/domain/games"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCard(t *testing.T) {
	t.Parallel()
	sucCard := games.Card("xs")

	testCases := map[string]struct {
		typ     string
		expCard *games.Card
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
			card, err := games.NewCard(tt.typ)
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
	t.Parallel()
	card := games.NewUnrevealedCard()
	assert.Equal(t, "*", card.Type())
}

func TestNewCardsDeck(t *testing.T) {
	t.Parallel()
	card, err := games.NewCard("xs")
	require.NoError(t, err)

	testCases := map[string]struct {
		name   string
		cards  []games.Card
		expErr string
	}{
		"success": {
			name:  "name",
			cards: []games.Card{*card},
		},
		"fail on empty name": {
			name:   "",
			cards:  []games.Card{*card},
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
			deck, err := games.NewCardsDeck(tt.name, tt.cards)
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
	t.Parallel()
	card1, err := games.NewCard("XS")
	require.NoError(t, err)
	card2, err := games.NewCard("S")
	require.NoError(t, err)
	deck, err := games.NewCardsDeck("foo", []games.Card{*card1, *card2})
	require.NoError(t, err)

	testCases := map[string]struct {
		deck     games.CardsDeck
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
			card, err := games.NewCard(tt.card)
			require.NoError(t, err)
			found := deck.IsInDeck(*card)
			assert.Equal(t, tt.expFound, found)
		})
	}
}
