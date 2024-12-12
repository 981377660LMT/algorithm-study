package main

type Book struct {
	title string // unexported field
}

// Getter
func (b *Book) Title() string {
	if len(b.title) == 0 {
		return "No title"
	}
	return b.title
}

// Setter
func (b *Book) SetTitle(title string) {
	b.title = title
}
