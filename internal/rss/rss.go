package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"

	log "github.com/sirupsen/logrus"
)

type storage interface {
	WriteNews([]*database.Post) error
}

type NewsParser struct {
	db            storage
	rssLinks      []string
	requestPeriod time.Duration

	errorChan chan error
	postChan  chan []*database.Post
}

func NewNewsParser(cfg config.RSS, db storage) *NewsParser {
	return &NewsParser{
		db:            db,
		rssLinks:      cfg.Links,
		requestPeriod: cfg.RequestPeriod,
		errorChan:     make(chan error),
		postChan:      make(chan []*database.Post),
	}
}

func (p *NewsParser) Run() {
	for {
		posts := p.readAllRSS()
		err := p.db.WriteNews(posts)
		if err != nil {
			log.WithError(err).Error("fail to write data to database")
		}

		<-time.After(p.requestPeriod)
	}
}

func (p *NewsParser) readAllRSS() []*database.Post {
	rssCount := len(p.rssLinks)

	for _, link := range p.rssLinks {
		go p.readRSS(link)
	}

	var posts []*database.Post

	for rssCount > 0 {
		select {
		case err := <-p.errorChan:
			log.WithError(err).Error("failed to read rss")

		case i := <-p.postChan:
			posts = append(posts, i...)
		}

		rssCount--
	}

	return posts
}

func (p *NewsParser) readRSS(link string) {

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		p.errorChan <- err
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		p.errorChan <- err
		return
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)

	var rss RSS

	err = xml.Unmarshal(text, &rss)
	if err != nil {
		p.errorChan <- err
		return
	}

	posts, err := p.convertDataModel(rss.Channel.Items)
	if err != nil {
		p.errorChan <- err
		return
	}

	p.postChan <- posts
}

func (p *NewsParser) convertDataModel(items []Item) ([]*database.Post, error) {
	posts := make([]*database.Post, 0, len(items))

	for _, item := range items {
		var post database.Post

		post.Title = item.Title
		post.Content = item.Content
		post.Link = item.Link

		dateLayout := "Mon, 2 Jan 2006 15:04:05 MST"
		pubTime, err := time.Parse(dateLayout, item.PubTime)
		if err != nil {
			return nil, err
		}
		post.PubTime = pubTime.Unix()

		posts = append(posts, &post)
	}
	return posts, nil
}
