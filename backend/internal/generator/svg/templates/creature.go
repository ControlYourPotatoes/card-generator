// svg/templates/creature.go
package templates

import (
	"image"
	"os"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/metadata"
)

// CreatureTemplate implements SVGTemplate for creature cards
type CreatureTemplate struct {
	templateDir string
	loader      *Loader
	svgContent  string
}

// GetFrame implements base.Template interface for compatibility
func (t *CreatureTemplate) GetFrame(data *card.CardDTO) (image.Image, error) {
	// SVG templates don't use image frames, but we need this for interface compatibility
	// Return a placeholder 1x1 image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	return img, nil
}

// GetTextBounds implements base.Template interface for compatibility
func (t *CreatureTemplate) GetTextBounds(data *card.CardDTO) map[string]image.Rectangle {
	bounds := make(map[string]image.Rectangle)
	
	// Standard creature card text bounds (based on 1500x2100 card size)
	bounds["name"] = image.Rect(125, 90, 1375, 170)
	bounds["cost"] = image.Rect(1320, 90, 1420, 190)
	bounds["type"] = image.Rect(125, 1885, 1375, 1955)
	bounds["effect"] = image.Rect(160, 1250, 1340, 1750)
	bounds["attack"] = image.Rect(130, 1820, 230, 1900)
	bounds["defense"] = image.Rect(1270, 1820, 1370, 1900)
	bounds["collector"] = image.Rect(110, 2010, 750, 2090)
	
	return bounds
}

// GetArtBounds implements base.Template interface for compatibility
func (t *CreatureTemplate) GetArtBounds() image.Rectangle {
	// Standard creature art bounds
	return image.Rect(170, 240, 1330, 1000)
}

// GetSVGTemplate implements SVGTemplate interface
func (t *CreatureTemplate) GetSVGTemplate() string {
	if t.svgContent == "" {
		t.loadSVGContent()
	}
	return t.svgContent
}

// GetInteractiveZones implements SVGTemplate interface
func (t *CreatureTemplate) GetInteractiveZones() map[string]metadata.InteractiveZone {
	zones := make(map[string]metadata.InteractiveZone)
	
	// Main card tap zone
	zones["card-tap"] = metadata.InteractiveZone{
		ID:      "card-tap",
		Bounds:  image.Rect(0, 0, 1500, 2100),
		Action:  "tap",
		Trigger: "click",
		Data: map[string]interface{}{
			"zone_type": "main_card",
		},
	}
	
	// Card inspection zone (hover for details)
	zones["card-inspect"] = metadata.InteractiveZone{
		ID:      "card-inspect",
		Bounds:  image.Rect(0, 0, 1500, 2100),
		Action:  "inspect",
		Trigger: "hover",
		Data: map[string]interface{}{
			"zone_type": "inspect",
		},
	}
	
	// Stats zone for power/toughness targeting
	zones["stats-target"] = metadata.InteractiveZone{
		ID:      "stats-target",
		Bounds:  image.Rect(100, 1800, 1400, 1950),
		Action:  "target_stats",
		Trigger: "click",
		Data: map[string]interface{}{
			"zone_type": "stats",
			"targetable": "power_toughness",
		},
	}
	
	return zones
}

// GetAnimationTargets implements SVGTemplate interface
func (t *CreatureTemplate) GetAnimationTargets() []metadata.AnimationTarget {
	targets := []metadata.AnimationTarget{
		{
			ElementID:     "card-frame",
			AnimationType: "glow",
			Properties: map[string]interface{}{
				"color":     "#00ff00",
				"intensity": 0.8,
			},
			Duration: "0.3s",
			Trigger:  "hover",
		},
		{
			ElementID:     "stats-group",
			AnimationType: "pulse",
			Properties: map[string]interface{}{
				"scale_min": 1.0,
				"scale_max": 1.1,
			},
			Duration: "0.5s",
			Trigger:  "target",
		},
		{
			ElementID:     "card-frame",
			AnimationType: "shake",
			Properties: map[string]interface{}{
				"intensity": 2.0,
				"direction": "horizontal",
			},
			Duration: "0.2s",
			Trigger:  "damage",
		},
	}
	return targets
}

// loadSVGContent loads the SVG template from file or generates it
func (t *CreatureTemplate) loadSVGContent() {
	// Try to load from file first
	svgPath := filepath.Join(t.templateDir, "creature.svg")
	if content, err := os.ReadFile(svgPath); err == nil {
		t.svgContent = string(content)
		return
	}
	
	// Generate default SVG template if file doesn't exist
	t.svgContent = t.generateDefaultSVG()
}

// generateDefaultSVG creates a default creature SVG template
func (t *CreatureTemplate) generateDefaultSVG() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<svg width="1500" height="2100" viewBox="0 0 1500 2100" xmlns="http://www.w3.org/2000/svg">
  <!-- Card Frame -->
  <g id="card-frame" class="card-frame">
    <!-- Background -->
    <rect x="0" y="0" width="1500" height="2100" fill="#1a1a1a" stroke="#333" stroke-width="2" rx="60"/>
    
    <!-- Art Frame -->
    <rect id="art-frame" x="170" y="240" width="1160" height="760" fill="#333" stroke="#666" stroke-width="2" rx="20"/>
    
    <!-- Name Background -->
    <rect id="name-bg" x="125" y="90" width="1250" height="80" fill="#2a2a2a" stroke="#555" stroke-width="1" rx="10"/>
    
    <!-- Type Line Background -->
    <rect id="type-bg" x="125" y="1885" width="1250" height="70" fill="#2a2a2a" stroke="#555" stroke-width="1" rx="10"/>
    
    <!-- Effect Text Background -->
    <rect id="effect-bg" x="160" y="1250" width="1180" height="500" fill="#1e1e1e" stroke="#444" stroke-width="1" rx="15"/>
    
    <!-- Stats Background -->
    <g id="stats-group" class="stats-group">
      <circle id="attack-bg" cx="180" cy="1860" r="50" fill="#8b0000" stroke="#ff0000" stroke-width="2"/>
      <circle id="defense-bg" cx="1320" cy="1860" r="50" fill="#006400" stroke="#00ff00" stroke-width="2"/>
    </g>
  </g>
  
  <!-- Text Placeholders -->
  <g id="text-elements" class="text-elements">
    <!-- Card Name -->
    <text id="card-name" x="750" y="140" text-anchor="middle" font-family="serif" font-size="60" font-weight="bold" fill="white">
      {{.Name}}
    </text>
    
    <!-- Mana Cost -->
    <text id="mana-cost" x="1370" y="140" text-anchor="middle" font-family="serif" font-size="50" font-weight="bold" fill="white">
      {{.Cost}}
    </text>
    
    <!-- Type Line -->
    <text id="type-line" x="750" y="1935" text-anchor="middle" font-family="serif" font-size="40" fill="white">
      Creature{{if .Trait}} - {{.Trait}}{{end}}
    </text>
    
    <!-- Effect Text -->
    <foreignObject id="effect-text" x="170" y="1260" width="1160" height="480">
      <div xmlns="http://www.w3.org/1999/xhtml" style="color: white; font-family: serif; font-size: 32px; line-height: 1.2; padding: 10px;">
        {{.Effect}}
      </div>
    </foreignObject>
    
    <!-- Stats -->
    <text id="attack-text" x="180" y="1870" text-anchor="middle" font-family="serif" font-size="48" font-weight="bold" fill="white">
      {{.Attack}}
    </text>
    <text id="defense-text" x="1320" y="1870" text-anchor="middle" font-family="serif" font-size="48" font-weight="bold" fill="white">
      {{.Defense}}
    </text>
  </g>
  
  <!-- Interactive Zones (invisible overlay) -->
  <g id="interactive-zones" class="interactive-zones">
    <rect id="tap-zone" x="0" y="0" width="1500" height="2100" fill="transparent" data-action="tap" data-trigger="click"/>
    <rect id="inspect-zone" x="0" y="0" width="1500" height="2100" fill="transparent" data-action="inspect" data-trigger="hover"/>
    <rect id="stats-zone" x="100" y="1800" width="1300" height="150" fill="transparent" data-action="target_stats" data-trigger="click"/>
  </g>
  
  <!-- CSS Styles -->
  <style>
    <![CDATA[
      .card-frame:hover { filter: drop-shadow(0 0 10px rgba(255,255,255,0.3)); }
      .stats-group:hover { transform: scale(1.05); transition: transform 0.2s ease; }
      .interactive-zones rect { cursor: pointer; }
      .text-elements text { user-select: none; }
    ]]>
  </style>
</svg>`
} 