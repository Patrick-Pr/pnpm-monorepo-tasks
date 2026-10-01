package main

import (
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/service/pnpm"
)

func main() {
	pnpm.RecognizeWorkspace()

	// p := tea.NewProgram(model.InitialModel())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatalf("Alas, theres ben an error: %v", err)
	// }
}
