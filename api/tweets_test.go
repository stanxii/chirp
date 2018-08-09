package api

import (
	"net/http"
	"testing"

	"chirp.com/models"
)

func TestTweets(t *testing.T) {
	services, router := setUpTests()
	defer services.Close()

	tt := newTweetsTester()
	ut := newUsersTester()
	tt.users = ut.users
	testCases := tt.createTestCases()
	runAPITests(t, router, testCases)
}

type tweetsTester struct {
	users                map[string]*models.User
	tweetsFromSetup      map[uint]*models.Tweet
	tweetsCountFromTests uint
	tweetsFromTests      map[uint]*models.Tweet
}

func newTweetsTester() *tweetsTester {
	tt := &tweetsTester{}
	tt.createTweets()
	return tt
}

func (tt *tweetsTester) createTweets() {
	const io = iota
	tt.tweetsFromSetup = make(map[uint]*models.Tweet, 0)
	tt.tweetsFromSetup[1001] = &models.Tweet{
		Username: duasings,
		Post:     "Hey, this is my first tweet!",
	}
	tt.tweetsFromSetup[1002] = &models.Tweet{
		Username: duasings,
		Post:     "Second tweet! Let's go!",
	}
	tt.tweetsFromSetup[1003] = &models.Tweet{
		Username: bobbyd,
		Post:     "I love playing the guitar.",
	}
	tt.tweetsFromSetup[1004] = &models.Tweet{
		Username: vinceTester,
		Post:     "this tweet will be deleted...",
	}
	tt.tweetsFromSetup[1005] = &models.Tweet{
		Username: vinceTester,
		Post:     "this tweet will be updated...",
	}
	tt.tweetsFromSetup[1006] = &models.Tweet{
		Username: kanye_west,
		Post:     "amazing tweet by kanye",
	}

	tt.tweetsFromTests = make(map[uint]*models.Tweet, 0)
	tt.tweetsFromTests[1] = &models.Tweet{
		Username: vinceTester,
		Post:     "new tweet in testing!",
	}
	tt.tweetsFromTests[2] = &models.Tweet{
		Username:  vinceTester,
		Retweet:   tt.tweetsFromSetup[1006],
		RetweetID: 1006,
	}

	for k, tweet := range tt.tweetsFromSetup {
		tweet.ID = k
	}
	for k, tweet := range tt.tweetsFromTests {
		tweet.ID = k
	}
}

func (tt *tweetsTester) getTweetsByUsername(username string) (ret []map[string]interface{}) {
	var tweets []*models.Tweet
	for _, t := range tt.tweetsFromSetup {
		tweets = append(tweets, t)
	}
	for _, t := range tt.tweetsFromTests {
		tweets = append(tweets, t)
	}
	for _, t := range tweets {
		if t.Username == username {
			ret = append(ret, toMap(t))
		}
	}
	return ret
}

func (tt *tweetsTester) createTestCases() (testCases []apiTestCase) {
	getTweet := apiTestCase{
		tag:    "Get tweet info",
		method: "GET",
		url:    "/duasings/1001",
		status: http.StatusOK,
		want:   toMap(tt.tweetsFromSetup[1001]),
	}
	// getSelfTweets := apiTestCase{
	// 	tag:      "get all self tweets",
	// 	method:   "GET",
	// 	url:      "/i/tweets",
	// 	status:   http.StatusOK,
	// 	want:     tt.getTweetsByUsername(vinceTester),
	// 	remember: tokenUserRequired,
	// }
	getUserTweets := apiTestCase{
		tag:      "get all user tweets",
		method:   "GET",
		url:      "/duasings/tweets",
		status:   http.StatusOK,
		want:     tt.getTweetsByUsername(duasings),
		remember: tokenUserRequired,
	}
	postTweet := apiTestCase{
		tag:    "post tweet",
		method: "POST",
		body: TweetForm{
			Post: "new tweet in testing!",
		},
		url:      "/tweets",
		status:   http.StatusOK,
		want:     toMap(tt.tweetsFromTests[1]),
		remember: tokenUserRequired,
	}
	deleteTweet := apiTestCase{
		tag:    "delete tweet",
		method: "POST",
		url:    "/tweets/vinceTester/1004/delete",
		status: http.StatusOK,
		want: toMap(&models.Tweet{
			ID: 1004,
		}),
		remember: tokenUserRequired,
	}
	updateTweet := apiTestCase{
		tag:    "update tweeet",
		method: "POST",
		body: TweetForm{
			Post: "just updated this tweet",
		},
		url:    "/vinceTester/1005/update",
		status: http.StatusOK,
		want: toMap(&models.Tweet{
			ID:       1005,
			Username: vinceTester,
			Post:     "just updated this tweet",
		}),
		remember: tokenUserRequired,
	}

	likedTweet := *tt.tweetsFromSetup[1002]
	likedTweet.LikesCount = 1
	likeTweet := apiTestCase{
		tag:      "like tweet",
		method:   "POST",
		url:      "/vinceTester/1002/like",
		status:   http.StatusOK,
		want:     toMap(likedTweet),
		remember: tokenUserRequired,
	}
	deleteLike := apiTestCase{
		tag:      "remove like on tweet",
		method:   "POST",
		url:      "/kanye_west/1006/like/delete",
		status:   http.StatusOK,
		want:     toMap(tt.tweetsFromSetup[1006]),
		remember: tokenUserRequired,
	}
	getUsersWhoLiked := apiTestCase{
		tag:    "get users who liked the tweet",
		method: "GET",
		url:    "/bobbyd/1003/liked",
		status: http.StatusOK,
		want: []map[string]interface{}{
			toMap(tt.users[samsmith], "email"),
			toMap(tt.users[kanye_west], "email"),
			toMap(tt.users[duasings], "email"),
		},
		remember: tokenUserRequired,
	}

	createRetweet := apiTestCase{
		tag:      "retweet tweet",
		method:   "POST",
		url:      "/kanye_west/1006/retweet",
		status:   http.StatusOK,
		want:     toMap(tt.tweetsFromTests[2]),
		remember: tokenUserRequired,
	}

	testCases = append(testCases,
		getTweet,
		postTweet,
		getUserTweets,
		deleteTweet,
		updateTweet,
		likeTweet,
		deleteLike,
		getUsersWhoLiked,
		createRetweet,
	)
	return testCases
}
