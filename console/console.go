package console

import (
	"fmt"
	"strings"

	"github.com/kreuzwerker/asc/database"

	prompt "github.com/c-bata/go-prompt"
)

type Console struct {
	db      *database.Database
	p       *prompt.Prompt
	region  string
	service string
}

func New(d *database.Database) *Console {

	c := &Console{
		db: d,
	}

	c.p = prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionLivePrefix(c.prefix),
		prompt.OptionTitle("bla"),
	)

	return c

}

func (c *Console) completer(d prompt.Document) []prompt.Suggest {

	var (
		// a = strings.Split(d.TextBeforeCursor(), " ")
		s []prompt.Suggest
		w = d.GetWordBeforeCursor()
	)

	if w == "region" {

		regions, err := c.db.Terms("Other.Location")

		if err != nil {
			panic(err)
		}

		for _, region := range regions { // export this as region

			s = append(s, prompt.Suggest{
				Text:        region,
				Description: "none",
			})

		}

		return prompt.FilterContains(s, w, true)

	}

	s = []prompt.Suggest{
		{Text: "region", Description: "Region string with wildcard support"},
		{Text: "service", Description: "Service description"},
	}

	return prompt.FilterHasPrefix(s, w, true)

}

func (c *Console) executor(in string) {

	in = strings.TrimSpace(in)

	fmt.Println("ok!")

}

func (c *Console) prefix() (string, bool) {

	var prefix []string

	if c.region != "" {
		prefix = append(prefix, fmt.Sprintf("region=%s", c.region))
	}

	if c.service != "" {
		prefix = append(prefix, fmt.Sprintf("service=%s", c.service))
	}

	return strings.Join(prefix, ","), false
}

func (c *Console) Run() {
	c.p.Run()
}
