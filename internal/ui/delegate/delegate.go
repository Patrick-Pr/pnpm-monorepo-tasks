package delegate

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/items"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/styles"
	"github.com/charmbracelet/x/ansi"
)

type commandFinishedMsg struct{ err error }

type ItemDelegate struct {
	list.DefaultDelegate
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(items.Item)
	if !ok || m.Width() <= 0 {
		return
	}

	const gap = 2
	totalWidth := m.Width()
	leftWidth := min(32, totalWidth)
	pathWidth := max(0, totalWidth-leftWidth-gap)

	// m is a copy. Narrow it for the default renderer.
	m.SetWidth(leftWidth)

	var left bytes.Buffer
	d.DefaultDelegate.Render(&left, m, index, listItem)

	if pathWidth == 0 {
		fmt.Fprint(w, left.String())
		return
	}

	leftColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Render(left.String())

	var packageName string
	if index == m.Index() {
		packageName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Render(ansi.Truncate(item.PackageName, pathWidth, "…"))
	} else {
		packageName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDD")).
			Render(ansi.Truncate(item.PackageName, pathWidth, "…"))
	}

	fmt.Fprint(w, lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		strings.Repeat(" ", gap),
		packageName,
	))
}

func NewItemDelegate(keys *DelegateKeyMap, styles *styles.Styles) ItemDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.Choose):
				// return m.NewStatusMessage(styles.StatusMessage.Render("You chose " + title))
				item, ok := m.SelectedItem().(items.Item)
				if !ok {
					return nil
				}
				cmd := exec.Command("pnpm", "-F", item.PackageName, "run", item.TitleText)
				cmd.Dir = item.Path
				return tea.Exec(&scriptCommand{Cmd: cmd}, func(err error) tea.Msg {
					return commandFinishedMsg{err: err}
				})
			}
		}

		return nil
	}

	help := []key.Binding{keys.Choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return ItemDelegate{DefaultDelegate: d}
}

type DelegateKeyMap struct {
	Choose key.Binding
}

// ShortHelp Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.Choose,
	}
}

// FullHelp Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.Choose,
		},
	}
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		Choose: key.NewBinding(
			key.WithKeys("enter", "v"),
			key.WithHelp("enter/v", "choose"),
		),
	}
}
