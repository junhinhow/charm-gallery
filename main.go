package main

import (
	"flag"
	"fmt"
	"image"
	ic "image/color"
	"os"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	a11y "github.com/junhinhow/charm-a11y"
	banner "github.com/junhinhow/charm-banner"
	charts "github.com/junhinhow/charm-charts"
	charmcolor "github.com/junhinhow/charm-color"
	command "github.com/junhinhow/charm-command"
	charmdiff "github.com/junhinhow/charm-diff"
	fuzzy "github.com/junhinhow/charm-fuzzy"
	i18n "github.com/junhinhow/charm-i18n"
	charmimage "github.com/junhinhow/charm-image"
	charminput "github.com/junhinhow/charm-input"
	charmjson "github.com/junhinhow/charm-json"
	keys "github.com/junhinhow/charm-keys"
	layout "github.com/junhinhow/charm-layout"
	logs "github.com/junhinhow/charm-logs"
	nav "github.com/junhinhow/charm-nav"
	palette "github.com/junhinhow/charm-palette"
	panes "github.com/junhinhow/charm-panes"
	progress "github.com/junhinhow/charm-progress"
	spinner "github.com/junhinhow/charm-spinner"
	themes "github.com/junhinhow/charm-themes"
	toast "github.com/junhinhow/charm-toast"
	wizard "github.com/junhinhow/charm-wizard"
)

type demo func(w, h int) string

var demos = map[string]demo{
	"a11y": func(w, h int) string {
		p := a11y.DetectFrom(func(string) string { return "" })
		return p.Describe() + "\n" + a11y.StripANSI("\x1b[31mvermelho\x1b[0m") + "\n" + p.Symbol("▸", ">")
	},
	"banner": func(w, h int) string {
		return banner.Render("CHARM", banner.WithFont(banner.Medium))
	},
	"charts": func(w, h int) string {
		s := charts.NewSparkline([]float64{4, 2, 7, 3, 9, 5, 8, 1, 6, 4}).Width(w - 4).Height(4)
		s = s.Style(lipgloss.NewStyle().Foreground(lipgloss.Color("#a3be8c")))
		return s.Render()
	},
	"color": func(w, h int) string {
		c := charmcolor.Darken(charmcolor.FromHex("#4c566a"), 0.3)
		return "darken(#4c566a, 0.3) = " + charmcolor.ToHex(c) + "\n" +
			"complementar: " + charmcolor.ToHex(charmcolor.Complementary(charmcolor.FromHex("#4c566a")))
	},
	"command": func(w, h int) string {
		p := command.New(
			command.Command{ID: "save", Title: "Salvar", Description: "Salva o arquivo", Keys: []string{"ctrl+s"}},
			command.Command{ID: "quit", Title: "Sair", Description: "Encerra", Keys: []string{"q"}},
		)
		p.SetSize(w, h)
		p.Open()
		return p.View().Content
	},
	"diff": func(w, h int) string {
		return charmdiff.UnifiedDiff("a\nb\nc\n", "a\nX\nc\n")
	},
	"fuzzy": func(w, h int) string {
		matches := fuzzy.Rank("cha", []string{"charm", "charts", "color", "command", "config"})
		var b strings.Builder
		for _, m := range matches {
			b.WriteString(fuzzy.Highlight(m.Value, m.Positions, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ebcb8b"))))
			b.WriteString(fmt.Sprintf("  (score %d)\n", m.Score))
		}
		return b.String()
	},
	"i18n": func(w, h int) string {
		return "locale: pt-BR\n" +
			"RTL(ar): " + fmt.Sprint(i18n.IsRTL("ar")) + "\n" +
			"Reverse(Hello שלום) = " + i18n.Reverse("Hello שלום")
	},
	"image": func(w, h int) string {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				img.Set(x, y, ic.RGBA{R: uint8(x * 60), G: uint8(y * 60), B: 180, A: 255})
			}
		}
		s, _ := charmimage.Render(img, charmimage.WithWidth(12), charmimage.WithProtocol(charmimage.ProtocolHalfBlock))
		return s
	},
	"input": func(w, h int) string {
		r := charminput.NewRange(0, 100, 42)
		r.SetSize(w, h)
		return r.View().Content
	},
	"json": func(w, h int) string {
		node, err := charmjson.ParseJSON([]byte(`{"name":"Charm","version":"1.0","items":[1,2,3]}`))
		if err != nil {
			return err.Error()
		}
		m := charmjson.New(node)
		m.SetSize(w, h)
		m.Update(tea.WindowSizeMsg{Width: w, Height: h})
		return m.View().Content
	},
	"keys": func(w, h int) string {
		r := keys.New()
		r.Global(keys.Binding{Action: "quit", Keys: []string{"q"}, Help: "sair"})
		r.Define("home", keys.Binding{Action: "save", Keys: []string{"ctrl+s"}, Help: "salvar"})
		r.Push("home")
		return r.HelpWidth(w)
	},
	"layout": func(w, h int) string {
		f := layout.NewFlex(layout.Row).SetSize(w, 5).SetGap(1)
		f.AddChild(0.4, "Sidebar")
		f.AddChild(0.6, "Content")
		return f.Render()
	},
	"logs": func(w, h int) string {
		m := logs.New(logs.WithMaxEntries(5))
		m.Append(logs.Entry{Level: logs.Info, Source: "app", Message: "iniciando"})
		m.Append(logs.Entry{Level: logs.Warn, Source: "net", Message: "latencia alta"})
		m.Append(logs.Entry{Level: logs.Error, Source: "db", Message: "conexao perdida"})
		m.SetSize(w, h)
		return m.View().Content
	},
	"nav": func(w, h int) string {
		s := nav.NewSidebar("Dashboard", "Users", "Settings")
		s.Width = 20
		return s.Render()
	},
	"palette": func(w, h int) string {
		return palette.ShowPalette([]ic.Color{charmcolor.FromHex("#bf616a"), charmcolor.FromHex("#a3be8c"), charmcolor.FromHex("#81a1c1")}, w)
	},
	"panes": func(w, h int) string {
		l := panes.New(panes.Horizontal)
		l.Add(panes.Pane{ID: "left", Title: "Left", Content: "conteudo esquerdo", MinSize: 10})
		l.Add(panes.Pane{ID: "right", Title: "Right", Content: "conteudo direito", MinSize: 10})
		l.SetSize(w, h)
		return l.View()
	},
	"progress": func(w, h int) string {
		d := progress.New()
		d.Add(progress.Task{ID: "1", Label: "baixando", Total: 100, Current: 60, State: progress.Running})
		d.Add(progress.Task{ID: "2", Label: "processando", Total: 50, Current: 50, State: progress.Done})
		d.Add(progress.Task{ID: "3", Label: "publicando", Total: 30, Current: 5, State: progress.Running})
		return d.View(w)
	},
	"spinner": func(w, h int) string {
		all := spinner.All()
		s := all["Dots"]
		return s.Render(0) + " " + s.Render(3) + "   " + all["Line"].Render(2)
	},
	"store": func(w, h int) string {
		return "store: estado observavel com undo/redo e snapshots"
	},
	"themes": func(w, h int) string {
		t := themes.Nord
		s := t.Lipgloss()
		return s.Title.Render("Nord") + "\n" + s.Base.Render("texto de exemplo")
	},
	"toast": func(w, h int) string {
		m := toast.NewToastManager()
		m.Add(toast.NewToast("salvo", 0).WithType(toast.Success))
		m.Add(toast.NewToast("dica", 0).WithType(toast.Info))
		return m.Render(w)
	},
	"wizard": func(w, h int) string {
		wz := wizard.New().WithTitle("demo")
		wz.AddStep("primeiro", "primeira etapa", func(map[string]any) error { return nil })
		wz.AddStep("segundo", "segunda etapa", func(map[string]any) error { return nil })
		return wizard.RenderProgress(wz.Steps(), wz.CurrentStep())
	},
}

func componentNames() []string {
	out := make([]string, 0, len(demos))
	for k := range demos {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func runRender(name string, w, h int) error {
	fn, ok := demos[name]
	if !ok {
		var suggestions []string
		for _, n := range componentNames() {
			if strings.Contains(n, name) || strings.Contains(name, n) {
				suggestions = append(suggestions, n)
			}
		}
		if len(suggestions) == 0 {
			return fmt.Errorf("componente %q nao encontrado (use -list)", name)
		}
		return fmt.Errorf("componente %q nao encontrado; talvez voce quis dizer: %s", name, strings.Join(suggestions, ", "))
	}
	out := fn(w, h)
	if w > 0 {
		lines := strings.Split(out, "\n")
		for i, l := range lines {
			if lipgloss.Width(l) > w {
				lines[i] = lipgloss.NewStyle().MaxWidth(w).Render(l)
			}
		}
		out = strings.Join(lines, "\n")
	}
	fmt.Print(out)
	if !strings.HasSuffix(out, "\n") {
		fmt.Println()
	}
	return nil
}

func runList() {
	for _, n := range componentNames() {
		fmt.Println(n)
	}
}

type model struct {
	names     []string
	cursor    int
	width     int
	height    int
	themeName string
}

func newModel(w, h int) model {
	return model{names: componentNames(), width: w, height: h, themeName: "nord"}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.names)-1 {
				m.cursor++
			}
		case "t":
			themes := []string{"nord", "gruvbox", "dracula", "catppuccin", "tokyo-night"}
			for i, th := range themes {
				if th == m.themeName {
					m.themeName = themes[(i+1)%len(themes)]
					break
				}
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 {
		m.width = 100
	}
	if m.height == 0 {
		m.height = 30
	}
	sidebarW := 22
	mainW := m.width - sidebarW - 1
	if mainW < 10 {
		mainW = 10
	}

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#88c0d0")).Render("Charm Suite"))
	sb.WriteString("\n")
	for i, n := range m.names {
		line := " " + n
		if i == m.cursor {
			line = lipgloss.NewStyle().Bold(true).Reverse(true).Render(line)
		}
		sb.WriteString(line + "\n")
	}
	sidebar := lipgloss.NewStyle().Width(sidebarW).Height(m.height - 2).Render(sb.String())

	name := m.names[m.cursor]
	main := lipgloss.NewStyle().Width(mainW).Height(m.height - 2).Render(demos[name](mainW, m.height-2))

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("#7b8794")).Render("up/down navegar  t tema  q sair   tema: " + m.themeName)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, "│", main)
	return tea.NewView(body + "\n" + footer)
}

func main() {
	list := flag.Bool("list", false, "lista os componentes")
	render := flag.String("render", "", "renderiza um componente e sai")
	theme := flag.String("theme", "nord", "tema inicial")
	w := flag.Int("w", 100, "largura")
	h := flag.Int("h", 30, "altura")
	flag.Parse()

	if *list {
		runList()
		return
	}
	if *render != "" {
		if err := runRender(*render, *w, *h); err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		return
	}

	m := newModel(*w, *h)
	m.themeName = *theme
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}
