package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "Card-Generator/internal/parser"
)

func main() {
    inputFile := flag.String("input", "", "Input CSV file")
    outputDir := flag.String("output", "output", "Output directory for card images")
    flag.Parse()

    if *inputFile == "" {
        log.Fatal("Input file is required")
    }

    file, err := os.Open(*inputFile)
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer file.Close()

    p := parser.NewParser(file)
    cards, err := p.Parse()
    if err != nil {
        log.Fatalf("Failed to parse cards: %v", err)
    }

    fmt.Printf("Successfully parsed %d cards\n", len(cards))
    
    // TODO: Add image generation logic
}