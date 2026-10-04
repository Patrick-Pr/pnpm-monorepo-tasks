package model

import (
	"log"
	"os"
	"sync"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/service/pnpm"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/delegate"
	itemgenerator "github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/itemGenerator"
	item "github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/items"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/styles"
)

type Model struct {
	styles        styles.Styles
	darkBG        bool
	width, height int
	once          *sync.Once
	list          list.Model
	itemGenerator *itemgenerator.RandomItemGenerator
	keys          *listKeyMap
	delegateKeys  *delegate.DelegateKeyMap
	workspaceRoot *pnpm.PnpmWorkspaceRoot
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor)
}

func (m *Model) updateListProperties() {
	x, y := m.styles.App.GetFrameSize()
	m.list.SetSize(m.width-x, m.height-y)

	m.styles = styles.NewStyles(m.darkBG)
	m.list.Styles.Title = m.styles.Title
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd
		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil
		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil
		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil
		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		case key.Matches(msg, m.keys.insertItem):
			newItem := m.itemGenerator.Next()
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(m.styles.StatusMessage.Render("Added " + newItem.TitleText))
			return m, tea.Batch(insCmd, statusCmd)
		}
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	view := tea.NewView(m.styles.App.Render(m.list.View()))
	view.AltScreen = true
	return view
}

func InitialModel() Model {
	// Parsing the current directory to find the workspace root.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get the current Directory!")
	}
	workspace, err := pnpm.ParseWorkspace(cwd)
	if err != nil {
		log.Fatalf("Failed to parse workspace: %v", err)
	}

	// Initialize the Model and list.
	m := Model{}
	m.styles = styles.NewStyles(false) // default to dark background styles
	m.workspaceRoot = &workspace

	delegateKeys := delegate.NewDelegateKeyMap()
	listKeys := newListKeyMap()

	// Make initial list of items.
	var itemGenerator itemgenerator.RandomItemGenerator
	const numItems = 24
	items := make([]list.Item, 0, numItems)
	//for i := range numItems {
	//	items[i] = itemGenerator.Next()
	//}
	for _, pkg := range m.workspaceRoot.Packages {
		for name, script := range pkg.Scripts {
			items = append(items, item.Item{TitleText: name, DescriptionText: script, Path: pkg.ContextPath, PackageName: pkg.Name})
		}
	}

	// Setup list.
	delegate := delegate.NewItemDelegate(delegateKeys, &m.styles)

	groceryList := list.New(items, delegate, 0, 0)
	groceryList.KeyMap.Quit = key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	)
	groceryList.Title = "Workspace Packages"
	groceryList.Styles.Title = m.styles.Title
	groceryList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{listKeys.insertItem}
	}
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	m.list = groceryList
	m.keys = listKeys
	m.delegateKeys = delegateKeys
	m.itemGenerator = &itemGenerator

	return m
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleSpinner:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "toggle spinner")),
		toggleTitleBar:   key.NewBinding(key.WithKeys("T"), key.WithHelp("T", "toggle title")),
		toggleStatusBar:  key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "toggle status")),
		togglePagination: key.NewBinding(key.WithKeys("P"), key.WithHelp("P", "toggle pagination")),
		toggleHelpMenu:   key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "toggle help menu")),
		insertItem:       key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add item")),
	}
}
