package types

// CardElement identifies different text elements on a card
type CardElement string

const (
    ElementTitle    CardElement = "title"
    ElementCost     CardElement = "cost"
    ElementEffect   CardElement = "effect"
    ElementStats    CardElement = "stats"
    ElementKeyword  CardElement = "keyword"
    ElementType     CardElement = "type"
    ElementSubtype  CardElement = "subtype"
)

// ElementAttributes defines properties for each card element
type ElementAttributes struct {
    MaxLength    int
    MinFontSize  float64
    MaxFontSize  float64
    AllowWrap    bool
    MaxLines     int
    Required     bool
}

// Default attributes for each element
var DefaultAttributes = map[CardElement]ElementAttributes{
    ElementTitle: {
        MaxLength:    40,
        MinFontSize:  48,
        MaxFontSize:  72,
        AllowWrap:    false,
        MaxLines:     1,
        Required:     true,
    },
    ElementEffect: {
        MaxLength:    500,
        MinFontSize:  36,
        MaxFontSize:  48,
        AllowWrap:    true,
        MaxLines:     12,
        Required:     true,
    },
    ElementStats: {
        MaxLength:    5,
        MinFontSize:  48,
        MaxFontSize:  72,
        AllowWrap:    false,
        MaxLines:     1,
        Required:     false,
    },
    // Add other elements...
}