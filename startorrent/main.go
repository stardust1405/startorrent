package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Home Model

type homeModel struct {
	cursor           int
	menuItems        []string
	selectedMenuItem string
}

func initHomeModel() homeModel {
	return homeModel{
		menuItems: []string{"List Torrents", "Add Torrent", "Exit"},
	}
}

func (m homeModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Startorrent")
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "enter", "right":
			m.selectedMenuItem = m.menuItems[m.cursor]
			switch m.selectedMenuItem {
			case "Exit":
				return m, tea.Quit
			case "List Torrents":
				return initListTorrentModel(), nil
			case "Add Torrent":
				return initAddTorrentModel(), nil
			}
		}
	}

	return m, nil
}

func (m homeModel) View() string {
	s := "\n******** Startorrent welcomes you! **********\n\n"

	for i, item := range m.menuItems {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, item)
	}

	s += "\nSelected: " + m.selectedMenuItem + "\n"

	if m.selectedMenuItem == "Exit" {
		s += "\nBye!\n"
	}

	return s
}

// Add Torrent Model

type addTorrentModel struct {
	magnetLink string
}

func initAddTorrentModel() addTorrentModel {
	return addTorrentModel{
		magnetLink: "starmagnet 🧲\n",
	}
}

func (a addTorrentModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Add Torrent")
}

func (a addTorrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m", "left":
			return initHomeModel(), nil
		}
	}
	return a, nil
}

func (a addTorrentModel) View() string {
	s := "\nAdd Torrent\n\n"

	s += "Magnet Link: " + a.magnetLink + "\n"

	s += "\nPress m or < to go back to main menu\n"

	return s
}

// List Torrent Model

type listTorrentModel struct {
	cursor                 int
	torrents               []string
	torrentDownloadPercent map[int]int
}

func initListTorrentModel() listTorrentModel {
	list := []string{"Avengers Endgame", "Linux Kernel", "Star Wars", "Star Trek", "Forza Horizon 5"}
	percent := make(map[int]int)
	for i := range list {
		percent[i] = rand.Intn(101)
	}
	return listTorrentModel{
		torrents:               list,
		torrentDownloadPercent: percent,
	}
}

func (l listTorrentModel) Init() tea.Cmd {
	return tea.SetWindowTitle("List Torrents")
}

func (l listTorrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m", "left":
			return initHomeModel(), nil
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
			}
		case "down", "j":
			if l.cursor < len(l.torrents)-1 {
				l.cursor++
			}
		case "enter", "right":
			return initViewTorrentModel(l.torrents[l.cursor], l.torrentDownloadPercent[l.cursor], "/home/stardust/Downloads"), nil
		}
	}
	return l, nil
}

func (l listTorrentModel) View() string {
	s := "\nList Torrents\n\n"

	for i, torrent := range l.torrents {
		cursor := " "
		if l.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %d. %s (%d%%)\n", cursor, i+1, torrent, l.torrentDownloadPercent[i])

		status := "Downloading"
		if l.torrentDownloadPercent[i] == 100 {
			status = "Downloaded"
		}

		s += "     Progress: "

		prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
		s += prog.ViewAs(float64(l.torrentDownloadPercent[i]))
		s += "\n"
		s += fmt.Sprintf("     Status: %s\n\n", status)
	}

	s += "\nPress m or < to go back to main menu\n"

	return s
}

// View Torrent Model

type viewTorrentModel struct {
	name                   string
	torrentDownloadPercent int
	downloadPath           string
}

func initViewTorrentModel(name string, percent int, downloadPath string) viewTorrentModel {
	return viewTorrentModel{
		name:                   name,
		torrentDownloadPercent: percent,
		downloadPath:           downloadPath,
	}
}

func (v viewTorrentModel) Init() tea.Cmd {
	return tea.SetWindowTitle("View Torrent")
}

func (v viewTorrentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m", "left":
			return initListTorrentModel(), nil
		}
	}
	return v, nil
}

func (v viewTorrentModel) View() string {
	s := "\nView Torrent\n\n"

	s += "Torrent Name: " + v.name + "\n"
	s += "Download Path: " + v.downloadPath + "\n"
	s += "Download Percent: " + fmt.Sprintf("%d%%", v.torrentDownloadPercent) + "\n"
	s += "Status: "
	if v.torrentDownloadPercent == 100 {
		s += "Downloaded\n"
	} else {
		s += "Downloading\n"
	}
	s += "Progress: "

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	s += prog.ViewAs(float64(v.torrentDownloadPercent))

	s += "\nPress m or < to go back to list torrents\n"

	return s
}

func main() {
	p := tea.NewProgram(initHomeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error encountered, terminating: %v", err)
		os.Exit(1)
	}
}
