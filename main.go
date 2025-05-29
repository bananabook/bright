package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"flag"
)

type Screen struct{
	name string
	brightness float64
}
func (s *Screen)Change(increment float64)error{
	s.brightness=change(s.brightness,increment)
	return exec.Command("xrandr","--output",s.name,"--brightness",fmt.Sprint(s.brightness)).Run()
}

func main() {
	flag.Usage=func(){
		fmt.Println(
			"interactive brightness control using xrandr\n"+
			" k,j:     if multiple screens: move selector up and down\n"+
			" <space>: if multiple screens: select all screens, to apply actions to all\n"+
			" l,h:     move brightness up and down\n"+
			" L,H:     move brightness up and down in smaller increments\n"+
			"any keypress, that has no action will close the program",
		)
	}
	flag.Parse()
	if e:=run();e!=nil{
		fmt.Println("ERROR:",e)
	}
}

func run() error{
	screens,e:=getScreens()
	if e != nil{
		return e
	}
	
	ui( screens )
	return nil
}
func getScreens() ([]Screen,error){
	cmd:=exec.Command("xrandr","--verbose")
	out,e:=cmd.Output()
	if e!=nil{
		return []Screen{},e
	}
	connectedRE:=regexp.MustCompile(`\n(\w+)\sconnected`)
	findings:=connectedRE.FindAllSubmatch(out,-1)
	screenNames:=[]string{}
	for _,finding:=range findings{
		screenNames=append(screenNames,string(finding[1]))
	}
	
	splits:=connectedRE.Split(string(out),-1)
	brghtnessRE:=regexp.MustCompile(`Brightness: ([\d\.]*)`)
	screens:=[]Screen{}
	for i,v:=range splits{
		if i==0{
			continue
		}
		findings:=brghtnessRE.FindAllStringSubmatch(v,-1)
		brightness,e:=strconv.ParseFloat(string(findings[0][1]),64)
		if e != nil{
			return []Screen{},e
		}
		name:=screenNames[i-1]
		screens=append(screens,Screen{
			name,
			brightness,
		})
		
	}
	return screens,nil
}
