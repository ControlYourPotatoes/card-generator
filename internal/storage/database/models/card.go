package models

import (
	"time"
)

// CardModel represents a card in the database
type CardModel struct {
	ID        int
	TypeID    int
	TypeName  string
	Name      string
	Cost      int
	Effect    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TypeSpecificData holds type-specific card data
type TypeSpecificData struct {
	// Fields for all possible card types
	Attack       int
	Defense      int
	TraitID      int
	TraitName    string
	IsEquipment  bool
	TargetType   string
	Timing       string
	Continuous   bool
}

// KeywordModel represents a keyword in the database
type KeywordModel struct {
	ID          int
	Name        string
	Description string
}

// TraitModel represents a creature trait in the database
type TraitModel struct {
	ID          int
	Name        string
	Description string
}

// CardTypeModel represents a card type in the database
type CardTypeModel struct {
	ID          int
	Name        string
	Description string
}

// CardMetadataModel represents a card metadata entry in the database
type CardMetadataModel struct {
	CardID int
	Key    string
	Value  string
}

// CardImageModel represents a card image in the database
type CardImageModel struct {
	ID        int
	CardID    int
	ImagePath string
	Version   int
	CreatedAt time.Time
}

// CardSetModel represents a card set in the database
type CardSetModel struct {
	ID          int
	Name        string
	Code        string
	ReleaseDate time.Time
	Description string
}

// CardSetCardModel represents a card in a set in the database
type CardSetCardModel struct {
	SetID      int
	CardID     int
	CardNumber string
	Rarity     string
}