package models

import (
    "time"
    
)

// CardModel represents a card in the database
type CardModel struct {
    ID        int
    TypeID    int
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
    IsEquipment  bool
    TargetType   string
    Timing       string
    Continuous   bool
}