package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/93lykevin/go-twit-backend/internal/store"

	"github.com/gin-gonic/gin"
)

// createTweet is the Gin handler function to create a new tweet
// it will call addTweet function from /store/tweets.go to communicate with the
// db and actually create the tweet object in memory
func createTweet(ctx *gin.Context) {
	// Define a new tweet object
	tweet := new(store.Tweet)
	// BIND the context to this newly created 'tweet' object
	// TODO 03-29-2023: Understand ctx.Bind() from GIN framework. I'm not sure how it works.
	if err := ctx.Bind(tweet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := store.AddTweet(user, tweet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "Tweet created successfully",
		"tweets": tweet,
	})
}

// fetch current user tweets
func getCurrentUserTweets(ctx *gin.Context) {
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := store.GetCurrentUserTweets(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "getCurrentUserTweets fetched successfully",
		"tweets": user.Tweets,
	})
}

func getTweetById(ctx *gin.Context) {
	param := ctx.Param("tweet_id")
	tweetId, err := strconv.Atoi(param)

	tweet, err := store.FetchTweet(tweetId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "fetch tweet by id",
		"tweet": tweet,
	})
}

/*
get ALL tweets
TODO: Add in pagination
*/
func getAllTweets(ctx *gin.Context) {
	tweets, err := store.GetAllTweets()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "booger getTweets fetched successfully",
		"tweets": tweets,
	})
}

// fetch user tweets by id param
func getTweetsByUserId(ctx *gin.Context) {
	param := ctx.Param("user_id")
	userId, err := strconv.Atoi(param)

	tweets, err := store.GetTweetsByUserId(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "fetchTweetsByUserId fetched successfully",
		"tweets": tweets,
	})
}

func updateTweet(ctx *gin.Context) {
	jsonTweet := new(store.Tweet)
	if err := ctx.Bind(jsonTweet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dbTweet, err := store.FetchTweet(jsonTweet.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.ID != dbTweet.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not authorized."})
		return
	}
	jsonTweet.UpdatedAt = time.Now()
	if err := store.UpdateTweet(jsonTweet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "Tweet updated successfully.",
		"tweets": jsonTweet,
	})
}

func deleteTweet(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not valid ID."})
		return
	}
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tweet, err := store.FetchTweet(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.ID != tweet.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not authorized."})
		return
	}
	if err := store.DeleteTweet(tweet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Tweet deleted successfully."})
}
