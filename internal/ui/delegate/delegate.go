package delegate

import (
	"fmt"
	"os/exec"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/items"
	"github.com/Patrick-Pr/pnpm-monorepo-tasks/internal/ui/styles"
)

type commandFinishedMsg struct{ err error }

func NewItemDelegate(keys *DelegateKeyMap, styles *styles.Styles) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(items.Item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.Choose):
				// return m.NewStatusMessage(styles.StatusMessage.Render("You chose " + title))
				item, ok := m.SelectedItem().(items.Item)
				if !ok {
					return nil
				}
				script := fmt.Sprintf("pnpm -F %s %s", item.PackageName, item.TitleText)
				cmd := exec.Command("sh", "-c", script)
				return tea.ExecProcess(cmd, func(err error) tea.Msg {
					return commandFinishedMsg{err: err}
				})

			case key.Matches(msg, keys.Remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.Remove.SetEnabled(false)
				}
				return m.NewStatusMessage(styles.StatusMessage.Render("Deleted " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.Choose, keys.Remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type DelegateKeyMap struct {
	Choose key.Binding
	Remove key.Binding
}

// ShortHelp Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.Choose,
		d.Remove,
	}
}

// FullHelp Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.Choose,
			d.Remove,
		},
	}
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		Choose: key.NewBinding(
			key.WithKeys("enter", "v"),
			key.WithHelp("enter/v", "choose"),
		),
		Remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
