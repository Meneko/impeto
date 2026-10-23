package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type ComponentType string

const (
	ComponentButton     ComponentType = "Button"
	ComponentSlider     ComponentType = "Slider"
	ComponentBackground ComponentType = "Background"
	ComponentSwitch     ComponentType = "Switch"
	ComponentTitle      ComponentType = "Title"
	ComponentText       ComponentType = "Text"
)

type ComponentTrigger string

const (
	TriggerOnClick     ComponentTrigger = "onClick"
	TriggerOnMouseDown ComponentTrigger = "onMouseDown"
	TriggerOnLoad      ComponentTrigger = "onLoad"
)

type ComponentAction struct {
	Type   string   `json:"type"`            // eg 'navigate', 'pause'
	Params []string `json:"params,omitzero"` // eg 'url', 'id'
}

type Component struct {
	Type    ComponentType    `json:"type"`
	Trigger ComponentTrigger `json:"trigger,omitempty"`
	Action  ComponentAction  `json:"action,omitzero"`
}

type Image struct {
	Source string `json:"source"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type Layout struct {
	Position Position `json:"position"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
}

type UIElement struct {
	Component Component `json:"component"`
	Image     Image     `json:"image,omitzero"`
	Layout    Layout    `json:"layout,omitzero"`
}

type App struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Icon  string      `json:"icon"`
	Modal []UIElement `json:"modal"`
}

type AppsService struct {
	Apps []App
	mu   sync.RWMutex
}

func NewAppsService() *AppsService {
	apps := []App{}

	err := os.MkdirAll("cache", 0755)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("cache")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := os.ReadFile("cache/" + file.Name())

		if err != nil {
			log.Fatal(err)
		}

		if len(data) == 0 {
			continue
		}

		app := App{}

		fmt.Println(string(data))

		err = json.Unmarshal(data, &app)

		if err != nil {
			log.Fatal(err)
		}

		apps = append(apps, app)
	}

	return &AppsService{
		Apps: apps,
	}
}

func (a *AppsService) GetAll() []App {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Apps
}
