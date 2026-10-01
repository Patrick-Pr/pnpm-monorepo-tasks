package itemgenerator

import (
	"math/rand/v2"
	"sync"

	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/items"
)

type item = items.Item

type RandomItemGenerator struct {
	Titles     []string
	Descs      []string
	TitleIndex int
	DescIndex  int
	Mtx        *sync.Mutex
	Shuffle    *sync.Once
}

func (r *RandomItemGenerator) reset() {
	r.Mtx = &sync.Mutex{}
	r.Shuffle = &sync.Once{}

	r.Titles = []string{
		"Artichoke",
		"Baking Flour",
		"Bananas",
		"Barley",
		"Bean Sprouts",
		"Bitter Melon",
		"Black Cod",
		"Blood Orange",
		"Brown Sugar",
		"Cashew Apple",
		"Cashews",
		"Cat Food",
		"Coconut Milk",
		"Cucumber",
		"Curry Paste",
		"Currywurst",
		"Dill",
		"Dragonfruit",
		"Dried Shrimp",
		"Eggs",
		"Fish Cake",
		"Furikake",
		"Garlic",
		"Gherkin",
		"Ginger",
		"Granulated Sugar",
		"Grapefruit",
		"Green Onion",
		"Hazelnuts",
		"Heavy whipping cream",
		"Honey Dew",
		"Horseradish",
		"Jicama",
		"Kohlrabi",
		"Leeks",
		"Lentils",
		"Licorice Root",
		"Meyer Lemons",
		"Milk",
		"Molasses",
		"Muesli",
		"Nectarine",
		"Niagamo Root",
		"Nopal",
		"Nutella",
		"Oat Milk",
		"Oatmeal",
		"Olives",
		"Papaya",
		"Party Gherkin",
		"Peppers",
		"Persian Lemons",
		"Pickle",
		"Pineapple",
		"Plantains",
		"Pocky",
		"Powdered Sugar",
		"Quince",
		"Radish",
		"Ramps",
		"Star Anise",
		"Sweet Potato",
		"Tamarind",
		"Unsalted Butter",
		"Watermelon",
		"Weißwurst",
		"Yams",
		"Yeast",
		"Yuzu",
		"Snow Peas",
	}

	r.Descs = []string{
		"A little weird",
		"Bold flavor",
		"Can’t get enough",
		"Delectable",
		"Expensive",
		"Expired",
		"Exquisite",
		"Fresh",
		"Gimme",
		"In season",
		"Kind of spicy",
		"Looks fresh",
		"Looks good to me",
		"Maybe not",
		"My favorite",
		"Oh my",
		"On sale",
		"Organic",
		"Questionable",
		"Really fresh",
		"Refreshing",
		"Salty",
		"Scrumptious",
		"Delectable",
		"Slightly sweet",
		"Smells great",
		"Tasty",
		"Too ripe",
		"At last",
		"What?",
		"Wow",
		"Yum",
		"Maybe",
		"Sure, why not?",
	}

	r.Shuffle.Do(func() {
		shuf := func(x []string) {
			rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
		}
		shuf(r.Titles)
		shuf(r.Descs)
	})
}

func (r *RandomItemGenerator) Next() item {
	if r.Mtx == nil {
		r.reset()
	}

	r.Mtx.Lock()
	defer r.Mtx.Unlock()

	i := item{
		TitleText:       r.Titles[r.TitleIndex],
		DescriptionText: r.Descs[r.DescIndex],
	}

	r.TitleIndex++
	if r.TitleIndex >= len(r.Titles) {
		r.TitleIndex = 0
	}

	r.DescIndex++
	if r.DescIndex >= len(r.Descs) {
		r.DescIndex = 0
	}

	return i
}
