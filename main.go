package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/url"
	"os"
)

func main() {

	type Tweet struct {
		TwitterID string `gorm:"unique_index"`
		Tweet     string `gorm:"type:varchar(512)"`
	}

	db, err := gorm.Open("sqlite3", "tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	twitterApi := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	v := url.Values{}
	v.Set("screen_name", os.Getenv("TWITTER_TARGET_SCREEN_NAME"))
	v.Set("count", "200")
	tweets, err := twitterApi.GetUserTimeline(v)
	if err != nil {
		panic(err)
	}
	tx := db.Begin()
	for _, tweet := range tweets {
		tweetDb := Tweet{
			Tweet:     tweet.Text,
			TwitterID: fmt.Sprint(tweet.Id),
		}
		tx.Create(&tweetDb)
	}
	tx.Commit()

}
