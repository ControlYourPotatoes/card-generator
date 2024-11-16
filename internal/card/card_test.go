package card

import (
    "testing"
)

func TestCreature(t *testing.T) {
    tests := []struct {
        name    string
        card    Creature
        wantErr bool
    }{
        {
            name: "valid creature",
            card: Creature{
                BaseCard: BaseCard{
                    Name:   "Demon Pup",
                    Cost:   1,
                    Effect: "Each time you OFFER; gain +1/1.",
                    Type:   TypeCreature,
                },
                Attack:  1,
                Defense: 1,
                SubType: "Demon",
            },
            wantErr: false,
        },
        {
            name: "invalid attack",
            card: Creature{
                BaseCard: BaseCard{
                    Name:   "Invalid Creature",
                    Cost:   1,
                    Effect: "Test",
                    Type:   TypeCreature,
                },
                Attack:  -1,
                Defense: 1,
                SubType: "Demon",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.card.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Creature.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestArtifact(t *testing.T) {
    tests := []struct {
        name    string
        card    Artifact
        wantErr bool
    }{
        {
            name: "valid equipment artifact",
            card: Artifact{
                BaseCard: BaseCard{
                    Name:   "Enchanted Mace",
                    Cost:   2,
                    Effect: "pay 1 - Equip to target creature. Equipped creature gains CRITICAL.",
                    Type:   TypeArtifact,
                },
                IsEquipment: true,
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.card.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Artifact.Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

