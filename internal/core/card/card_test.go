package card

import (
    "testing"
)

func TestBaseCardValidation(t *testing.T) {
    tests := []struct {
        name    string
        card    BaseCard
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid base card",
            card: BaseCard{
                Name:   "Test Card",
                Cost:   1,
                Effect: "Test effect",
                Type:   TypeCreature,
            },
            wantErr: false,
        },
        {
            name: "empty name",
            card: BaseCard{
                Name:   "",
                Cost:   1,
                Effect: "Test effect",
                Type:   TypeCreature,
            },
            wantErr: true,
            errMsg:  "name cannot be empty",
        },
        {
            name: "negative cost",
            card: BaseCard{
                Name:   "Test Card",
                Cost:   -1,
                Effect: "Test effect",
                Type:   TypeCreature,
            },
            wantErr: true,
            errMsg:  "cost cannot be negative",
        },
        {
            name: "empty effect",
            card: BaseCard{
                Name:   "Test Card",
                Cost:   1,
                Effect: "",
                Type:   TypeCreature,
            },
            wantErr: true,
            errMsg:  "effect cannot be empty",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.card.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("BaseCard.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err != nil && err.Error() != tt.errMsg {
                t.Errorf("BaseCard.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
            }
        })
    }
}

func TestIncantationValidation(t *testing.T) {
    tests := []struct {
        name    string
        card    Incantation
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid on clash incantation",
            card: Incantation{
                BaseCard: BaseCard{
                    Name:   "Quick Strike",
                    Cost:   1,
                    Effect: "Deal 2 damage to target creature",
                    Type:   TypeIncantation,
                },
                Timing: "ON ANY CLASH",
            },
            wantErr: false,
        },
        {
            name: "valid on attack incantation",
            card: Incantation{
                BaseCard: BaseCard{
                    Name:   "Battle Fury",
                    Cost:   2,
                    Effect: "Target creature gets +2/+0",
                    Type:   TypeIncantation,
                },
                Timing: "ON ATTACK",
            },
            wantErr: false,
        },
        {
            name: "invalid timing",
            card: Incantation{
                BaseCard: BaseCard{
                    Name:   "Invalid Timing",
                    Cost:   1,
                    Effect: "Test effect",
                    Type:   TypeIncantation,
                },
                Timing: "INVALID TIMING",
            },
            wantErr: true,
            errMsg:  "invalid timing: INVALID TIMING",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.card.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Incantation.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err != nil && err.Error() != tt.errMsg {
                t.Errorf("Incantation.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
            }
        })
    }
}

func TestAnthemValidation(t *testing.T) {
    tests := []struct {
        name    string
        card    Anthem
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid anthem",
            card: Anthem{
                BaseCard: BaseCard{
                    Name:   "Glory of Battle",
                    Cost:   3,
                    Effect: "All creatures you control get +1/+1",
                    Type:   TypeAnthem,
                },
                Continuous: true,
            },
            wantErr: false,
        },
        {
            name: "invalid non-continuous anthem",
            card: Anthem{
                BaseCard: BaseCard{
                    Name:   "Invalid Anthem",
                    Cost:   3,
                    Effect: "Test effect",
                    Type:   TypeAnthem,
                },
                Continuous: false,
            },
            wantErr: true,
            errMsg:  "anthem must be continuous",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.card.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Anthem.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err != nil && err.Error() != tt.errMsg {
                t.Errorf("Anthem.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
            }
        })
    }
}

func TestCardDataConversion(t *testing.T) {
    tests := []struct {
        name     string
        card     Card
        validate func(*testing.T, *CardData)
    }{
        {
            name: "creature conversion",
            card: &Creature{
                BaseCard: BaseCard{
                    Name:   "Test Creature",
                    Cost:   2,
                    Effect: "Test effect",
                    Type:   TypeCreature,
                },
                Attack:  2,
                Defense: 3,
                Trait: "Beast",
            },
            validate: func(t *testing.T, data *CardData) {
                if data.Attack != 2 || data.Defense != 3 || data.Trait != "Beast" {
                    t.Error("Creature data not converted correctly")
                }
            },
        },
        {
            name: "artifact conversion",
            card: &Artifact{
                BaseCard: BaseCard{
                    Name:   "Test Artifact",
                    Cost:   2,
                    Effect: "Equip - Test effect",
                    Type:   TypeArtifact,
                },
                IsEquipment: true,
            },
            validate: func(t *testing.T, data *CardData) {
                if !data.IsEquipment {
                    t.Error("Artifact equipment status not converted correctly")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            data := ToData(tt.card)
            if data.Name != tt.card.GetName() || data.Cost != tt.card.GetCost() || 
               data.Effect != tt.card.GetEffect() || data.Type != tt.card.GetType() {
                t.Error("Base card data not converted correctly")
            }
            tt.validate(t, data)
        })
    }
}

// Mock store for testing
type mockStore struct{}

func (m *mockStore) Save(card Card) (string, error)    { return "test-id", nil }
func (m *mockStore) Load(id string) (Card, error)      { return nil, nil }
func (m *mockStore) List() ([]Card, error)            { return nil, nil }
func (m *mockStore) Delete(id string) error           { return nil }