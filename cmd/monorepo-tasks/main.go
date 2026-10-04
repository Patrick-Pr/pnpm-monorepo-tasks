package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/model"
)

func main() {
	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, theres ben an error: %v", err)
	}
}
