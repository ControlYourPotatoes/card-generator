// metadata_builder.go builds metadata for parsed templates
package ingestion

import (
	"time"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// MetadataBuilder generates metadata for parsed templates
type MetadataBuilder struct {
	// Configuration for metadata building
}

// NewMetadataBuilder creates a new metadata builder
func NewMetadataBuilder() *MetadataBuilder {
	return &MetadataBuilder{}
}

// BuildMetadata generates comprehensive metadata for a parsed template
func (mb *MetadataBuilder) BuildMetadata(
	sourceFile string,
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
	positioning *TransparencyPositioning,
) *TemplateMetadata {
	
	metadata := &TemplateMetadata{
		SourceFile:     sourceFile,
		ParsedAt:       time.Now().Format(time.RFC3339),
		ObjectCount:    len(objects),
		BoundaryCount:  len(boundaries),
		SupportedTypes: mb.determineSupportedCardTypes(objects, boundaries),
		Version:        "1.0.0", // Phase 3a initial version
		Validation:     mb.validateTemplate(objects, boundaries, positioning),
	}
	
	return metadata
}

// determineSupportedCardTypes analyzes objects and boundaries to determine supported card types
func (mb *MetadataBuilder) determineSupportedCardTypes(
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
) []card.CardType {
	
	var supportedTypes []card.CardType
	
	// Check each card type for compatibility
	allCardTypes := []card.CardType{
		card.TypeCreature,
		card.TypeArtifact,
		card.TypeSpell,
		card.TypeAnthem,
	}
	
	for _, cardType := range allCardTypes {
		if mb.isCardTypeSupported(cardType, objects, boundaries) {
			supportedTypes = append(supportedTypes, cardType)
		}
	}
	
	return supportedTypes
}

// isCardTypeSupported checks if a card type is supported by the template
func (mb *MetadataBuilder) isCardTypeSupported(
	cardType card.CardType,
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
) bool {
	
	// Check required objects
	requiredObjects := GetRequiredObjectsForCardType(cardType)
	for _, required := range requiredObjects {
		if _, exists := objects[required]; !exists {
			return false
		}
	}
	
	// Check required boundaries
	requiredBoundaries := GetRequiredBoundariesForCardType(cardType)
	for _, required := range requiredBoundaries {
		if _, exists := boundaries[required]; !exists {
			return false
		}
	}
	
	return true
}

// validateTemplate performs comprehensive validation of the template
func (mb *MetadataBuilder) validateTemplate(
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
	positioning *TransparencyPositioning,
) ValidationResult {
	
	var errors []string
	var warnings []string
	
	// Validate objects
	objectErrors, objectWarnings := mb.validateObjects(objects)
	errors = append(errors, objectErrors...)
	warnings = append(warnings, objectWarnings...)
	
	// Validate boundaries
	boundaryErrors, boundaryWarnings := mb.validateBoundaries(boundaries)
	errors = append(errors, boundaryErrors...)
	warnings = append(warnings, boundaryWarnings...)
	
	// Validate positioning
	positioningErrors, positioningWarnings := mb.validatePositioning(positioning)
	errors = append(errors, positioningErrors...)
	warnings = append(warnings, positioningWarnings...)
	
	// Cross-validate objects and boundaries
	crossErrors, crossWarnings := mb.validateObjectBoundaryCompatibility(objects, boundaries)
	errors = append(errors, crossErrors...)
	warnings = append(warnings, crossWarnings...)
	
	return ValidationResult{
		IsValid:  len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}

// validateObjects validates the object collection
func (mb *MetadataBuilder) validateObjects(objects map[ObjectType]*CardObject) ([]string, []string) {
	var errors []string
	var warnings []string
	
	if len(objects) == 0 {
		errors = append(errors, "template has no objects - visual elements are required")
		return errors, warnings
	}
	
	// Check for essential objects
	hasFrameBase := false
	hasNameTitle := false
	hasEffectBody := false
	
	for objectType := range objects {
		switch objectType {
		case ObjectFrameBase:
			hasFrameBase = true
		case ObjectNameTitle:
			hasNameTitle = true
		case ObjectEffectBody:
			hasEffectBody = true
		}
	}
	
	if !hasFrameBase {
		errors = append(errors, "template missing ObjectFrameBase - card frame is required")
	}
	
	if !hasNameTitle {
		warnings = append(warnings, "template missing ObjectNameTitle - card name styling recommended")
	}
	
	if !hasEffectBody {
		warnings = append(warnings, "template missing ObjectEffectBody - effect text styling recommended")
	}
	
	// Check for conflicting frame types
	frameTypes := []ObjectType{ObjectFrameCreature, ObjectFrameAnthem, ObjectFrameArtifact, ObjectFrameSpell}
	frameCount := 0
	for _, frameType := range frameTypes {
		if _, exists := objects[frameType]; exists {
			frameCount++
		}
	}
	
	if frameCount > 1 {
		errors = append(errors, "template has multiple frame types - only one frame type allowed")
	}
	
	if frameCount == 0 {
		warnings = append(warnings, "template has no specific frame type - consider adding one for better card type support")
	}
	
	return errors, warnings
}

// validateBoundaries validates the boundary collection
func (mb *MetadataBuilder) validateBoundaries(boundaries map[BoundaryType]*TextBoundary) ([]string, []string) {
	var errors []string
	var warnings []string
	
	if len(boundaries) == 0 {
		errors = append(errors, "template has no boundaries - text areas are required")
		return errors, warnings
	}
	
	// Check for essential boundaries
	hasNameText := false
	hasEffectText := false
	hasCostSymbols := false
	
	for boundaryType := range boundaries {
		switch boundaryType {
		case BoundaryNameText:
			hasNameText = true
		case BoundaryEffectText:
			hasEffectText = true
		case BoundaryCostSymbols:
			hasCostSymbols = true
		}
	}
	
	if !hasNameText {
		errors = append(errors, "template missing BoundaryNameText - card name area is required")
	}
	
	if !hasEffectText {
		warnings = append(warnings, "template missing BoundaryEffectText - effect text area recommended")
	}
	
	if !hasCostSymbols {
		warnings = append(warnings, "template missing BoundaryCostSymbols - cost display area recommended")
	}
	
	// Validate individual boundaries
	for boundaryType, boundary := range boundaries {
		boundaryErrors, boundaryWarnings := mb.validateSingleBoundary(boundaryType, boundary)
		errors = append(errors, boundaryErrors...)
		warnings = append(warnings, boundaryWarnings...)
	}
	
	return errors, warnings
}

// validateSingleBoundary validates an individual boundary
func (mb *MetadataBuilder) validateSingleBoundary(boundaryType BoundaryType, boundary *TextBoundary) ([]string, []string) {
	var errors []string
	var warnings []string
	
	// Check safe zone
	if boundary.SafeZone.Empty() {
		errors = append(errors, string(boundaryType)+" has empty safe zone")
	}
	
	// Check preferred zone
	if boundary.PreferredZone.Empty() {
		warnings = append(warnings, string(boundaryType)+" has empty preferred zone")
	}
	
	// Check font constraints
	if boundary.FontConstraints.MinSize <= 0 {
		errors = append(errors, string(boundaryType)+" has invalid minimum font size")
	}
	
	if boundary.FontConstraints.MaxSize <= boundary.FontConstraints.MinSize {
		errors = append(errors, string(boundaryType)+" has invalid font size range")
	}
	
	if boundary.FontConstraints.PreferredSize < boundary.FontConstraints.MinSize ||
		boundary.FontConstraints.PreferredSize > boundary.FontConstraints.MaxSize {
		warnings = append(warnings, string(boundaryType)+" preferred font size is outside min/max range")
	}
	
	// Check character limits
	if boundary.MaxCharacters <= 0 {
		warnings = append(warnings, string(boundaryType)+" has no character limit - consider setting one")
	}
	
	return errors, warnings
}

// validatePositioning validates the positioning system
func (mb *MetadataBuilder) validatePositioning(positioning *TransparencyPositioning) ([]string, []string) {
	var errors []string
	var warnings []string
	
	if positioning == nil {
		errors = append(errors, "template has no positioning information")
		return errors, warnings
	}
	
	if len(positioning.Layers) == 0 {
		warnings = append(warnings, "template has no transparency layers")
	}
	
	if len(positioning.OpacityMaps) == 0 {
		warnings = append(warnings, "template has no opacity maps")
	}
	
	// Validate that layer count matches opacity map count
	if len(positioning.Layers) != len(positioning.OpacityMaps) {
		warnings = append(warnings, "layer count doesn't match opacity map count")
	}
	
	return errors, warnings
}

// validateObjectBoundaryCompatibility validates that objects and boundaries work together
func (mb *MetadataBuilder) validateObjectBoundaryCompatibility(
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
) ([]string, []string) {
	var errors []string
	var warnings []string
	
	// Check that text styling objects have corresponding boundaries
	textStyleObjects := []ObjectType{ObjectNameTitle, ObjectEffectBody, ObjectStatsText}
	correspondingBoundaries := []BoundaryType{BoundaryNameText, BoundaryEffectText, BoundaryStatsText}
	
	for i, textObj := range textStyleObjects {
		if _, hasObject := objects[textObj]; hasObject {
			if _, hasBoundary := boundaries[correspondingBoundaries[i]]; !hasBoundary {
				warnings = append(warnings, 
					string(textObj)+" object exists but corresponding boundary "+
					string(correspondingBoundaries[i])+" is missing")
			}
		}
	}
	
	// Check that boundaries have styling support
	for i, boundary := range correspondingBoundaries {
		if _, hasBoundary := boundaries[boundary]; hasBoundary {
			if _, hasObject := objects[textStyleObjects[i]]; !hasObject {
				warnings = append(warnings,
					string(boundary)+" boundary exists but corresponding text style object "+
					string(textStyleObjects[i])+" is missing")
			}
		}
	}
	
	return errors, warnings
}

// GetValidationSummary creates a human-readable validation summary
func (mb *MetadataBuilder) GetValidationSummary(validation ValidationResult) map[string]interface{} {
	summary := map[string]interface{}{
		"is_valid":      validation.IsValid,
		"error_count":   len(validation.Errors),
		"warning_count": len(validation.Warnings),
		"status":        "unknown",
	}
	
	if validation.IsValid {
		if len(validation.Warnings) == 0 {
			summary["status"] = "perfect"
		} else {
			summary["status"] = "valid_with_warnings"
		}
	} else {
		summary["status"] = "invalid"
	}
	
	summary["errors"] = validation.Errors
	summary["warnings"] = validation.Warnings
	
	return summary
}

// CreateMetadataSummary creates a summary of template metadata
func (mb *MetadataBuilder) CreateMetadataSummary(metadata *TemplateMetadata) map[string]interface{} {
	return map[string]interface{}{
		"source_file":     metadata.SourceFile,
		"parsed_at":       metadata.ParsedAt,
		"object_count":    metadata.ObjectCount,
		"boundary_count":  metadata.BoundaryCount,
		"supported_types": metadata.SupportedTypes,
		"version":         metadata.Version,
		"is_valid":        metadata.Validation.IsValid,
		"validation":      mb.GetValidationSummary(metadata.Validation),
	}
} 