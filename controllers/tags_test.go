package controllers

import (
	"net/http"
	"testing"

	"chirp.com/models"
)

func TestTags(t *testing.T) {
	services, router := setUpTests()
	defer services.Close()
	tt := newTweetsTester()
	ut := newUsersTester()
	tt.users = ut.users
	testCases := tt.createTagTestCases()
	runAPITests(t, router, testCases)
}

func (tt *tweetsTester) createTagTestCases() (testCases []apiTestCase) {
	getTweetsWithTag := apiTestCase{
		tag:    "Get array of tweets with tag name",
		method: "GET",
		url:    "/tags/lakers",
		status: http.StatusOK,
		want: []map[string]interface{}{
			toMap(tt.tweetsFromSetup[1001]),
			toMap(tt.tweetsFromSetup[1002]),
			toMap(tt.tweetsFromSetup[1003]),
			toMap(tt.tweetsFromSetup[1005]),
		},
	}

	postTweetWithTags := apiTestCase{
		tag:    "post tweet with tags",
		method: "POST",
		body: TweetForm{
			Post: "Ready for game 7!",
			Tags: []string{"Rockets", "OKC"},
		},
		url:    "/tweets",
		status: http.StatusOK,
		want: toMap(&models.Tweet{
			ID:       1,
			Username: vinceTester,
			Post:     "Ready for game 7!",
			Tags:     []string{"rockets", "okc"},
		}),
		remember: tokenUserRequired,
	}
	updateTags := apiTestCase{
		tag:    "update tags in tweet",
		method: "POST",
		body: TweetForm{
			Post: "wow, lost game 7",
			Tags: []string{"lakers", "cavs", "okc"},
		},
		url:    "/vinceTester/1005/update",
		status: http.StatusOK,
		want: toMap(&models.Tweet{
			ID:       1005,
			Username: vinceTester,
			Post:     "wow, lost game 7",
			Tags:     []string{"lakers", "cavs", "okc"},
		}),
		remember: tokenUserRequired,
	}

	testCases = append(testCases,
		getTweetsWithTag,
		postTweetWithTags,
		updateTags,
	)

	return testCases
}
