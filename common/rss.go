package common

type Channel struct {
	Title string `xml:title`
	Items []Item `xml:channel>item`
}

type Item struct {
	Title string `xml:title`
	Link  string `xml:link`
}
