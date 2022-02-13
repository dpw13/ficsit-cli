package scenes

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

var _ tea.Model = (*profiles)(nil)

type profiles struct {
	root   components.RootModel
	list   list.Model
	parent tea.Model
}

func NewProfiles(root components.RootModel, parent tea.Model) tea.Model {
	l := list.NewModel(profilesToList(root), utils.NewItemDelegate(), root.Size().Width, root.Size().Height-root.Height())
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetSpinner(spinner.MiniDot)
	l.Title = "Profiles"
	l.Styles = utils.ListStyles
	l.SetSize(l.Width(), l.Height())
	l.KeyMap.Quit.SetHelp("q", "back")
	l.DisableQuitKeybindings()

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("n", "new profile")),
		}
	}

	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithHelp("n", "new profile")),
		}
	}

	return &profiles{
		root:   root,
		list:   l,
		parent: parent,
	}
}

func (m profiles) Init() tea.Cmd {
	return nil
}

func (m profiles) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// List enables its own keybindings when they were previously disabled
	m.list.DisableQuitKeybindings()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.SettingFilter() {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch keypress := msg.String(); keypress {
		case "n":
			newModel := NewNewProfile(m.root, m)
			return newModel, newModel.Init()
		case KeyControlC:
			return m, tea.Quit
		case "q":
			if m.parent != nil {
				m.parent.Update(m.root.Size())
				return m.parent, nil
			}
			return m, tea.Quit
		case KeyEnter:
			i, ok := m.list.SelectedItem().(utils.SimpleItem)
			if ok {
				if i.Activate != nil {
					newModel, cmd := i.Activate(msg, m)
					if newModel != nil || cmd != nil {
						if newModel == nil {
							newModel = m
						}
						return newModel, cmd
					}
					return m, nil
				}
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(m.root.Height(), 2, 0).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	case updateProfileList:
		m.list.ResetSelected()
		cmd := m.list.SetItems(profilesToList(m.root))

		// Done to refresh keymap
		m.list.SetFilteringEnabled(m.list.FilteringEnabled())
		return m, cmd
	case updateProfileNames:
		cmd := m.list.SetItems(profilesToList(m.root))

		// Done to refresh keymap
		m.list.SetFilteringEnabled(m.list.FilteringEnabled())
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m profiles) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}

func profilesToList(root components.RootModel) []list.Item {
	items := make([]list.Item, len(root.GetGlobal().Profiles.Profiles))

	i := 0
	for _, profile := range root.GetGlobal().Profiles.Profiles {
		temp := profile
		items[i] = utils.SimpleItem{
			ItemTitle: temp.Name,
			Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
				newModel := NewEditProfile(root, currentModel, temp)
				return newModel, newModel.Init()
			},
		}
		i++
	}

	return items
}
