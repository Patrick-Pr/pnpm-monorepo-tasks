package items

type Item struct {
	TitleText       string
	DescriptionText string
}

func (i Item) Title() string {
	return i.TitleText
}

func (i Item) Description() string {
	return i.DescriptionText
}

func (i Item) FilterValue() string {
	return i.TitleText
}
