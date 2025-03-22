package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"math"
	"strings"
	"unicode"

	"github.com/lakrsv/parkour-engine/engine"
)

func Run(level int) {
	// TODO: Reuse world so we don't get multiple windows
	w := engine.GetInstance()
	w.InitWindow("Colormancer", 800, 480)

	grid := loadLevel(level, w)
	w.AddSystems(
		&DeferDoorRenderSystem{},
		&PlayerInputSystem{},
		&CreateSummonSystem{},
		&SummonInputSystem{},
		&MoveSystem{},
		&SummonPickupSystem{},
		&TriggerSystem{},
		&DirectionIndicatorSystem{},
		&RenderSystem{palette: NewRunePalette(
			map[rune]Color{
				Floor:           Color{R: 255, G: 255, B: 255},
				Wall:            Color{R: 255, G: 255, B: 255},
				Player:          Color{R: 0, G: 255, B: 0},
				Button:          Color{R: 0, G: 255, B: 0},
				TriggeredButton: Color{R: 0, G: 255, B: 0},
				OpenDoor:        Color{R: 255, G: 255, B: 255},
				DoorHorizontal:  Color{R: 255, G: 255, B: 255},
				DoorVertical:    Color{R: 255, G: 255, B: 255},
				UpIndicator:     Color{R: 0, G: 255, B: 0},
				DownIndicator:   Color{R: 0, G: 255, B: 0},
				LeftIndicator:   Color{R: 0, G: 255, B: 0},
				RightIndicator:  Color{R: 0, G: 255, B: 0},
				Exit:            Color{R: 255, G: 255, B: 255},
			}),
		},
	)
	_ = w.CreateEntity(
		*grid,
	)
	_ = w.CreateEntity(
		DoorOpenPlayCountComponent{},
	)

	if err := w.Simulate(); err != nil {
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
		if strings.HasPrefix(line, "//") {
			continue
		}
	}
	config := make(map[rune]map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "HEADER" || line == "LEVEL" {
			break
		}
		if strings.HasPrefix(line, "//") {
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
		if strings.HasPrefix(line, "//") {
			continue
		}
	}

	var header []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "LEVEL" {
			break
		}
		if strings.HasPrefix(line, "//") {
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
		if strings.HasPrefix(line, "//") {
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
		if strings.HasPrefix(line, "//") {
			continue
		}
	}
	cellOffset := 0
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") {
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
				components = append(components, SummonPickupBlueprint(x, y, Color{R: 0, G: 255, B: 255})...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case RedSummon:
				components = append(components, SummonPickupBlueprint(x, y, Color{R: 255, G: 0, B: 0})...)
				grid.BackgroundEntities[idx-cellOffset] = w.CreateEntity(
					components...,
				)
			case YellowSummon:
				components = append(components, SummonPickupBlueprint(x, y, Color{R: 255, G: 255, B: 0})...)
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
		var color Color
		switch configColor {
		case "Cyan":
			color = Color{R: 0, G: 255, B: 255}
		case "Green":
			color = Color{R: 0, G: 255, B: 0}
		case "Red":
			color = Color{R: 255, G: 0, B: 0}
		case "Yellow":
			color = Color{R: 255, G: 255, B: 0}
		default:
			color = Color{R: 255, G: 255, B: 255}
		}

		components = append(components, ColorComponent{color: color})
	}
	return components
}
