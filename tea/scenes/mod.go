package scenes

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/ficsit-cli/tea/components"
	"github.com/satisfactorymodding/ficsit-cli/tea/utils"
)

var _ tea.Model = (*modMenu)(nil)

type modMenu struct {
	root   components.RootModel
	list   list.Model
	parent tea.Model
}

func NewModMenu(root components.RootModel, parent tea.Model, mod utils.Mod) tea.Model {
	model := modMenu{
		root:   root,
		parent: parent,
	}

	var items []list.Item
	if root.GetCurrentProfile().HasMod(mod.Reference) {
		items = []list.Item{
			utils.SimpleItem{
				ItemTitle: "Remove Mod",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					root.GetCurrentProfile().RemoveMod(mod.Reference)
					return currentModel.(modMenu).parent, nil
				},
			},
			utils.SimpleItem{
				ItemTitle: "Change Version",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					newModel := NewModVersion(root, currentModel.(modMenu).parent, mod)
					return newModel, newModel.Init()
				},
			},
		}
	} else {
		items = []list.Item{
			utils.SimpleItem{
				ItemTitle: "Install Mod",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					err := root.GetCurrentProfile().AddMod(mod, ">=0.0.0")
					if err != nil {
						menu := currentModel.(exitMenu)
						log.Error().Err(err).Msg(ErrorFailedAddMod)
						cmd := menu.list.NewStatusMessage(ErrorFailedAddMod)
						return currentModel, cmd
					}
					return currentModel.(modMenu).parent, nil
				},
			},
			utils.SimpleItem{
				ItemTitle: "Install Mod with specific version",
				Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
					newModel := NewModVersion(root, currentModel.(modMenu).parent, mod)
					return newModel, newModel.Init()
				},
			},
		}
	}

	items = append(items, utils.SimpleItem{
		ItemTitle: "View Mod info",
		Activate: func(msg tea.Msg, currentModel tea.Model) (tea.Model, tea.Cmd) {
			newModel := NewModInfo(root, currentModel, mod)
			return newModel, newModel.Init()
		},
	})

	model.list = list.NewModel(items, utils.NewItemDelegate(), root.Size().Width, root.Size().Height-root.Height())
	model.list.SetShowStatusBar(false)
	model.list.SetFilteringEnabled(false)
	model.list.Title = mod.Name
	model.list.Styles = utils.ListStyles
	model.list.SetSize(model.list.Width(), model.list.Height())
	model.list.KeyMap.Quit.SetHelp("q", "back")
	model.list.StatusMessageLifetime = time.Second * 3
	model.list.DisableQuitKeybindings()

	return model
}

func (m modMenu) Init() tea.Cmd {
	return nil
}

func (m modMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
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
							newModel.Update(m.root.Size())
							newModel = m
						}
						return newModel, cmd
					}
					return m, nil
				}
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(2, 2).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	}

	return m, nil
}

func (m modMenu) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}
