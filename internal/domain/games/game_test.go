package games_test

import (
	"testing"

	"planningpoker/test"
)

func TestSimpleGame(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User1, "XS").
		Then().ShouldHaveVote(test.User1, "XS").
		When().UserJoins(test.User2).
		And().UserVotes(test.User2, "S").
		Then().ShouldHaveVote(test.User2, "S").
		When().UserReveals(test.User1).
		Then().GameShouldBeFinished()
}

func TestFailedToReveal(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserJoins(test.User2).
		And().UserReveals(test.User2).
		Then().ShouldFail("user can not reveal cards").
		And().GameShouldBeRunning()
}

func TestEveryoneCanReveal(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, true)).
		When().UserJoins(test.User1).
		And().UserJoins(test.User2).
		And().UserReveals(test.User2).
		Then().GameShouldBeFinished()
}

func TestUserCanJoinTwice(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, true)).
		When().UserJoins(test.User1).
		And().UserJoins(test.User1).
		Then().ShouldSucceed()
}

func TestRestartGame(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, true)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User1, "XS").
		And().UserReveals(test.User1).
		Then().GameShouldBeFinished().
		And().ShouldHaveVote(test.User1, "XS").
		When().UserRestartsGame(test.User1).
		Then().ShouldHaveNoVote(test.User1).
		And().GameShouldBeRunning()
}

func TestNonPlayerCanNotRestartAGame(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, true)).
		When().UserJoins(test.User1).
		And().UserRestartsGame(test.User2).
		Then().ShouldFail("user is not a player")
}

func TestPlayerWhoLeftCanNotParticipate(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, true)).
		When().UserJoins(test.User1).
		And().UserJoins(test.User2).
		And().UserLeaves(test.User2).
		And().UserVotes(test.User2, "S").
		Then().ShouldFail("user is not a player")
}

func TestCanNotVoteWhenGameFinished(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserReveals(test.User1).
		And().UserVotes(test.User1, "XS").
		Then().ShouldFail("can not vote on ended game")
}

func TestCannotVoteWhenNotAPlayer(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User2, "XS").
		Then().ShouldFail("user is not a player")
}

func TestCannotVoteWithWrongCard(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User1, "foo").
		Then().ShouldFail("unknown card")
}

func TestCanUnVote(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User1, "XS").
		Then().ShouldHaveVote(test.User1, "XS").
		When().UserUnVotes(test.User1).
		Then().ShouldHaveNoVote(test.User1)
}

func TestNonPlayerCanNotUnVote(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserUnVotes(test.User2).
		Then().ShouldFail("user is not a player")
}

func TestCanNotUnVoteOnFinishedGame(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserVotes(test.User1, "XS").
		And().UserReveals(test.User1).
		When().UserUnVotes(test.User1).
		Then().ShouldFail("can not un-vote on ended game")
}

func TestCannotRevealWhenNotAPlayer(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserReveals(test.User2).
		Then().ShouldFail("user is not a player")
}

func TestPlayerCanUpdateGameName(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserUpdatesGameName(test.User1, "new name").
		Then().ShouldHaveGameName("new name")
}

func TestCanNotUpdateGameNameWhenNotAPlayer(t *testing.T) {
	test.NewTestGame(t, test.NewSimpleGame(t, false)).
		When().UserJoins(test.User1).
		And().UserUpdatesGameName(test.User2, "new name").
		Then().ShouldFail("user is not a player")
}
