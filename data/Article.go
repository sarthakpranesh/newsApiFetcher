package data

import (
	"encoding/json"
	"io"
)

type Src struct {
	Id		string		`json:"id"`
	Name	string		`json:"name"`
}

type Article struct {
	Source		Src		`json:"source"`
	Author		string	`json:"author"`
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	Url			string	`json:"url"`
	UrlToImage	string	`json:"urlToImage"`
	publishedAt	string	`json:"-"`
	content		string	`json:"-"`
}

type Articles [] *Article

var articleList Articles

type ApiRequestFmt struct {
	Status			string		`json:"-"`
	TotalResults	int			`json:"-"`
	AllArticles		Articles	`json:"articles"`
}

func (arf *ApiRequestFmt) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(arf)
}

func (a *Articles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func SetArticles(a Articles) {
	articleList = a
}

func GetArticles() Articles{
	return articleList
}
