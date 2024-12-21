package types

import (

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
	"github.com/ControlYourPotatoes/card-generator/internal/core/common"
    
)

type Artifact struct {
    card.BaseCard
    IsEquipment bool
}

func (a *Artifact) Validate() *common.ValidationError {
    
    baseValidator := validation.BaseValidator{
        Name: a.Name,
        Cost: a.Cost,
        Effect: a.Effect,
    } 
    
    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    // Validate artifact-specific properties
    if err := validation.ValidateArtifact(a.IsEquipment, a.Effect); err != nil {
        return err
    }


    return nil
}

func (a *Artifact) ToData() *card.CardData {
    data := a.BaseCard.ToData()
    data.IsEquipment = a.IsEquipment
    return data
}