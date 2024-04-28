package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/lakrsv/parkour-engine/engine"
	"io"
	"io/fs"
	"math"
	"strings"
	"unicode"
)

func Run(level int) {
	w := engine.NewWorld()

	grid := loadLevel(level, w)
	w.AddSystems(
		&DeferDoorRenderSystem{},
		&InputSystem{},
		&PlayerInputSystem{},
		&CreateSummonSystem{},
		&SummonInputSystem{},
		&MoveSystem{},
		&SummonPickupSystem{},
		&TriggerSystem{},
		&DirectionIndicatorSystem{},
		&RenderSystem{palette: NewRunePalette(
			map[rune]color.Color{
				Floor:           *color.New(color.FgWhite),
				Wall:            *color.New(color.FgWhite),
				Player:          *color.New(color.FgGreen),
				Button:          *color.New(color.FgGreen),
				TriggeredButton: *color.New(color.FgGreen),
				OpenDoor:        *color.New(color.FgWhite),
				DoorHorizontal:  *color.New(color.FgWhite),
				DoorVertical:    *color.New(color.FgWhite),
				UpIndicator:     *color.New(color.FgGreen, color.Bold),
				DownIndicator:   *color.New(color.FgGreen, color.Bold),
				LeftIndicator:   *color.New(color.FgGreen, color.Bold),
				RightIndicator:  *color.New(color.FgGreen, color.Bold),
				Exit:            *color.New(color.FgWhite),
			}),
		},
	)
	_ = w.CreateEntity(
		*grid,
	)
	_ = w.CreateEntity(
		DoorOpenPlayCountComponent{},
	)

	if err := w.Simulate(context.Background()); err != nil {
		panic(err)
	}
}

func loadLevel(level int, w *engine.World) *GridComponent {
	file, e := content.Open(fmt.Sprintf("assets/levels/level_%d.txt", level))
	defer func(file fs.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	if e != nil {
		panic(e)
	}

	// Load config
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "CONFIG" {
			break
		}
		if strings.HasPrefix("//", line) {
			continue
		}
	}
	config := make(map[rune]map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "HEADER" || line == "LEVEL" {
			break
		}
		if strings.HasPrefix("//", line) {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		for _, char := range parts[0] {
			if config[char] == nil {
				config[char] = make(map[string]string)
			}
			for _, modifier := range strings.Split(strings.TrimSpace(parts[1]), ",") {
				if strings.Contains(modifier, ":") {
					parts := strings.Split(modifier, ":")
					config[char][parts[0]] = parts[1]
				} else {
					config[char][modifier] = ""
				}
			}
		}
	}

	// Load header
	if _, err := file.(io.Seeker).Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "HEADER" {
			break
		}
		if strings.HasPrefix("//", line) {
			continue
		}
	}

	var header []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "LEVEL" {
			break
		}
		if strings.HasPrefix("//", line) {
			continue
		}
		header = append(header, line)
	}

	_ = w.CreateEntity(
		LevelComponent{Level: level, Header: header},
	)

	width := 0
	height := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix("//", line) {
			continue
		}
		if width == 0 {
			width = len(line)
		}
		height += 1
	}

	grid := &GridComponent{Width: width, Height: height, BackgroundEntities: make([]uint32, width*height), ForegroundEntities: make([]uint32, width*height), EffectEntities: make([]uint32, width*height)}

	if _, err := file.(io.Seeker).Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(file)
	// Skip until we get to the level block
	for scanner.Scan() {
		line := scanner.Text()
		if line == "LEVEL" {
			break
		}
		if strings.HasPrefix("//", line) {
			continue
		}
	}
	cellOffset := 0
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix("//", line) {
			continue
		}
		runes := bufio.NewReader(strings.NewReader(line))
		for {
			char, _, err := runes.ReadRune()
			if err == io.EOF {
				cellOffset += 1
				idx++
				break
			} else if err != nil {
				panic(err)
			}
			x, y := grid.GetPosition(idx - cellOffset)
			grid.ForegroundEntities[idx-cellOffset] = math.MaxUint32
			grid.EffectEntities[idx-cellOffset] = math.MaxUint32

			var components []any
			components = append(components, getConfigComponents(config[char])...)

			switch char {
			case Floor:
				components = append(components, FloorBlueprint(x, y)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case Wall:
				components = append(components, WallBlueprint(x, y)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case Exit:
				components = append(components, ExitBlueprint(x, y, level+1)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case Player:
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					FloorBlueprint(x, y)...,
				)
				components = append(components, PlayerBlueprint(x, y)...)
				grid.ForegroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case CyanSummon:
				components = append(components, SummonPickupBlueprint(x, y, color.FgCyan)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case RedSummon:
				components = append(components, SummonPickupBlueprint(x, y, color.FgRed)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case YellowSummon:
				components = append(components, SummonPickupBlueprint(x, y, color.FgYellow)...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			default:
				// Button
				if unicode.IsLower(char) {
					components = append(components, ButtonBlueprint(x, y, unicode.ToUpper(char))...)
					grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
						components...,
					)
				} else if unicode.IsUpper(char) {
					if modifiers, ok := config[char]; ok {
						if _, ok := modifiers[OpenDoorModifier]; ok {
							components = append(components, OpenDoorBlueprint(x, y, char)...)
							grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
								components...,
							)
						} else if _, ok := modifiers[ClosedDoorModifier]; ok {
							components = append(components, ClosedDoorBlueprint(x, y, char)...)
							grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
								components...,
							)
						}
					}
				}
			}
			idx++
		}
	}
	return grid
}

func getConfigComponents(modifiers map[string]string) []any {
	var components []any
	if configColor, ok := modifiers[ColorModifier]; ok {
		var cAttr color.Attribute
		switch configColor {
		case "Cyan":
			cAttr = color.FgCyan
		case "Green":
			cAttr = color.FgGreen
		case "Red":
			cAttr = color.FgRed
		case "Yellow":
			cAttr = color.FgYellow
		default:
			cAttr = color.FgWhite
		}

		components = append(components, ColorComponent{color: *color.New(cAttr), colorAttr: cAttr})
	}
	return components
}
