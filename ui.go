package main

import (
	"fmt"
	"math"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{
	screens []Screen
	selected int
	info string
	all bool
}

func (m Model) Init() tea.Cmd{
	return nil
}
func (m *Model)Change(increment float64){
	if m.all{
		for i:=range m.screens{
			e:=m.screens[i].Change(increment)
			if e!=nil{
				m.info=e.Error()
			}
		}
		return
	}
	e:=m.screens[m.selected].Change(increment)
	if e!=nil{
		m.info=e.Error()
	}
}
func (m Model) Update(msg tea.Msg)(tea.Model,tea.Cmd){
	switch msg:=msg.(type){
	case tea.KeyMsg:
		switch msg.String(){
		case " ":
			m.all=!m.all
		case "j":
			m.selected=min(m.selected+1,len(m.screens)-1)
		case "k":
			m.selected=max(0,m.selected-1)
		case "h":
			m.Change(-0.1)
		case "H":
			m.Change(-0.01)
		case "l":
			m.Change(0.1)
		case "L":
			m.Change(0.01)
		case "q","esc","ctrl+c":
			return m,tea.Quit
		}
	}
	return m,nil
}
func change(in float64,change float64)float64{
	if change>0{
		return min(round(in+change,2),1)
	}
	return max(round(in+change,2),0)
}
func round(in float64,digit int)float64{
	return math.Round( in*math.Pow10(digit) )/100
}
func (m Model) View()string{
	var out []string
	for i,screen:=range m.screens{
		marker:=" "
		if m.all{
			marker="#"
		}else if i==m.selected{
			marker="*"
		}
		out=append(out,fmt.Sprintf("%v %v: %v",marker,screen.name,screen.brightness))
	}
	return strings.Join(out,"\n")+"\n"+m.info
}

// Model for single screen
type Sodel struct{
	screen Screen
}

func (s Sodel)Init()tea.Cmd{
	return nil
}
func (s Sodel)Update(msg tea.Msg)(tea.Model,tea.Cmd){
	switch msg:=msg.(type){
		case tea.KeyMsg:
			switch msg.String(){
			case "h":
				s.screen.Change(-0.1)
			case "H":
				s.screen.Change(-0.01)
			case "l":
				s.screen.Change(0.1)
			case "L":
				s.screen.Change(0.01)
		case "q","esc","ctrl-c":
				return s,tea.Quit
			}
	}
	return s,nil
}
func (s Sodel)View()string{
	return fmt.Sprintf("%v: %v",s.screen.name,s.screen.brightness)
}

func ui(screens []Screen) {
	if len(screens)==1{
		tea.NewProgram(Sodel{screen:screens[0]}).Run()
	}else{
		tea.NewProgram(Model{screens:screens,all:true}).Run()
	}
}

