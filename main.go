package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"net/url"
	"os"
	"os/exec"
)

func main() {

	type Tweet struct {
		TwitterID string `gorm:"unique_index"`
		Tweet     string `gorm:"type:varchar(512)"`
	}
	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})

	s3file, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String("/tweets.tar.xz"),
	})
	if err != nil {
		panic(err)
	}
	defer s3file.Body.Close()

	file, err := os.Create("/tmp/tweets.tar.xz")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, s3file.Body); err != nil {
		panic(err)
	}

	if err := exec.Command("sh", "-c", "cd /tmp && tar Jxvf tweets.tar.xz").Run(); err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", "/tmp/tweets.db")
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
		var tweetDb Tweet
		if db.Where("twitter_id = ?", fmt.Sprint(tweet.Id)).First(&tweetDb).RecordNotFound() {
			tweetDb := Tweet{
				Tweet:     tweet.Text,
				TwitterID: fmt.Sprint(tweet.Id),
			}
			tx.Create(&tweetDb)
		}
	}
	tx.Commit()

	if err := exec.Command("sh", "-c", "cd /tmp && tar Jcvf tweets.tar.xz tweets.db").Run(); err != nil {
		panic(err)
	}

	file2, err := os.Open("/tmp/tweets.tar.xz")

	if _, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String("/tweets.tar.xz"),
		Body:   file2,
	}); err != nil {
		panic(err)
	}

	if err := exec.Command("sh", "-c", "rm -rf /tmp/*").Run(); err != nil {
		panic(err)
	}
}
