package tui

import (
	"fmt"
	"sort"
	"strings"
	"wlpv/xmlparser"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().Padding(0, 1)
	infoStyle  = lipgloss.NewStyle().Padding(0, 1)
)

type view uint8

const (
	listView view = iota
	pagerView
)

type item struct {
	protocol     xmlparser.Protocol
	namespace    string
	pagerYOffset int
}

func newItem(protocol xmlparser.Protocol, namespace string) item {
	return item{
		protocol:     protocol,
		namespace:    namespace,
		pagerYOffset: 0,
	}
}
func (i item) Title() string       { return i.protocol.Name }
func (i item) Description() string { return i.namespace }
func (i item) FilterValue() string { return i.protocol.Name }

type model struct {
	list              list.Model
	viewport          viewport.Model
	ready             bool
	pending           view
	current           view
	selectedItemIndex int
	items             []item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc", "q", "h":
			if m.current == pagerView {
				m.exitPagerView()
			}

		case "enter", "l":
			if m.current == listView {
				{
					selectedItem := m.list.SelectedItem().(item)
					for index, item := range m.items {
						if item.protocol.Name == selectedItem.protocol.Name {
							m.viewport.SetYOffset(item.pagerYOffset)
							m.selectedItemIndex = index
							break
						}
					}
				}

				m.pending = pagerView

				selectedItem := m.items[m.selectedItemIndex]
				m.viewport.SetContent(selectedItem.protocol.Render())
			}

		case "g":
			if m.current == pagerView {
				m.viewport.GotoTop()
			}

		case "G":
			if m.current == pagerView {
				m.viewport.GotoBottom()
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		footerHeight := lipgloss.Height(m.footerView())

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
			if m.selectedItemIndex != -1 {
				m.viewport.SetContent(m.items[m.selectedItemIndex].protocol.Render())
			}
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - footerHeight
		}
	}

	switch m.current {
	case pagerView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	case listView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.current = m.pending

	return m, tea.Batch(cmds...)
}

func (m *model) exitPagerView() {
	m.pending = listView
	selectedItem := &m.items[m.selectedItemIndex]
	selectedItem.pagerYOffset = m.viewport.YOffset
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	var selectedTitle string
	if m.selectedItemIndex != -1 {
		selectedTitle = fmt.Sprintf("%s ", m.items[m.selectedItemIndex].protocol.Name)
	} else {
		selectedTitle = ""
	}

	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)-lipgloss.Width(selectedTitle)))

	return lipgloss.JoinHorizontal(lipgloss.Center, selectedTitle, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	var v string

	switch m.current {
	case pagerView:
		if !m.ready {
			v = "\n  Initializing..."
		} else {
			v = fmt.Sprintf("%s\n%s", m.viewport.View(), m.footerView())
		}

	case listView:
		v = docStyle.Render(m.list.View())
	}

	return v
}

func listShortHelpCallback() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func Run(protocolToOpen string, protocols map[string][]xmlparser.Protocol) error {
	var items []list.Item
	var mItems []item
	selectedIndex := -1

	{
		namespaces := []string{
			"core",
			"stable",
			"staging",
			"unstable",
			"wlroots",
			"weston",
			"kde",
		}
		count := 0
		for _, namespace := range namespaces {
			sort.Slice(protocols[namespace], func(i, j int) bool {
				a := protocols[namespace][i]
				b := protocols[namespace][j]

				return a.Name < b.Name
			})

			for _, protocol := range protocols[namespace] {
				if protocol.Name == protocolToOpen {
					selectedIndex = count
				}

				item := newItem(protocol, namespace)
				items = append(items, item)
				mItems = append(mItems, item)

				count++
			}
		}
	}

	var currentView view
	if protocolToOpen == "" {
		currentView = listView
		selectedIndex = -1
	} else {
		currentView = pagerView
	}

	defaultDelegate := list.NewDefaultDelegate()
	defaultDelegate.ShortHelpFunc = listShortHelpCallback

	m := model{
		list:              list.New(items, defaultDelegate, 0, 0),
		current:           currentView,
		items:             mItems,
		selectedItemIndex: selectedIndex,
	}

	if m.current == pagerView {
		m.pending = pagerView
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
