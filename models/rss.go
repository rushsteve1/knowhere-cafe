package models

import (
	"encoding/xml"
	"io"
	"time"

	"gorm.io/gorm"
)

// https://medium.com/parallel-thinking/building-a-personal-rss-email-digest-service-in-go-7c8b71ac5b89

type Feed struct {
	gorm.Model
	URL     string   `gorm:"unique"` // TODO xml parse this
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
	Updated atomTime `xml:"updated"`
}

func parseFeed(r io.Reader) (*Feed, error) {
	feed := Feed{}
	err := xml.NewDecoder(r).Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

type Entry struct {
	gorm.Model
	FeedID  uint
	URL     string   `gorm:"unique"`
	Title   string   `xml:"title"`
	Summary string   `xml:"summary"`
	Body    string   `xml:"content"`
	Updated atomTime `xml:"updated"`
	Author  Author   `xml:"author" gorm:"embedded"`
	Archive Archive  `gorm:"many2many:entry_archives"`
}

type Author struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type atomTime time.Time

func (a *atomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	*a = atomTime(parsed)
	return nil
}

func (a *atomTime) Time() time.Time {
	return time.Time(*a)
}

func (a *atomTime) LocalString() string {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return a.Time().String()
	}
	return a.Time().In(loc).Format("Jan 02, 2006 3:04PM")
}
