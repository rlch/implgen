package shrek

import "context"

// Embedded types and method sets

type Ogre struct {
	Name   string
	Layers int
	Swamp  string
}

type Princess struct {
	Name      string
	Kingdom   string
	Cursed    bool
	TrueForm  string
}

type Donkey struct {
	Name        string
	TalkingRate int // words per minute
}

// Embedded interfaces
type Roarer interface {
	Roar() string
}

type Storyteller interface {
	TellStory(ctx context.Context) (string, error)
}

// Combined interface
type SwampDweller interface {
	Roarer
	Storyteller
	GetTerritory() string
}

// Union types simulation
type FairyTaleCharacter interface {
	GetCharacterType() string
}

func (o Ogre) GetCharacterType() string     { return "ogre" }
func (p Princess) GetCharacterType() string { return "princess" }
func (d Donkey) GetCharacterType() string   { return "talking_animal" }