package model

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ModelMetadata struct {
	SavedAt     time.Time `json:"saved_at"`
	Description string    `json:"description"`
	NumSamples  int       `json:"num_samples,omitempty"`
	Version     string    `json:"version,omitempty"`
}

func SaveModelToJSON(model *LinearRegression, filePath string, description string, numSamples int) error {
	// add metadata fields directly to the model for saving
	model.SavedAt = time.Now()
	model.Description = description
	model.NumSamples = numSamples
	model.Version = "1.0"

	// marshal the model directly to JSON
	modelJSON, err := json.MarshalIndent(model, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling model to JSON %v:", err)
	}

	// write json to file
	err = os.WriteFile(filePath, modelJSON, 0644)
	if err != nil {
		return fmt.Errorf("error saving model to JSON: %v", err)
	}

	fmt.Printf("Model saved to: %s\n", filePath)
	fmt.Printf("  - target: %s\n", model.Target)
	fmt.Printf("  - Features: %s\n", model.Features)
	fmt.Printf("  - R-Squared: %.4f\n", model.RSquared)

	return nil
}

func LoadModelFromJSON(filePath string) (*LinearRegression, error) {
	// read the json file
	modelJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading most file: %v", err)
	}

	// unmarshall
	var loadedModel LinearRegression
	err = json.Unmarshal(modelJSON, &loadedModel)
	if err != nil {
		return nil, fmt.Errorf("invalid mdoel format in file: %v", err)
	}

	// print some model information
	fmt.Printf("Model sucessfully loaded from %s\n", filePath)
	if loadedModel.Version != "" {
		fmt.Printf("Model Versions: %s\n", loadedModel.Version)
	}
	if loadedModel.Description != "" {
		fmt.Printf("Description: %s\n", loadedModel.Description)
	}

	fmt.Printf("loaded model information: \n")
	fmt.Printf("- Target: %s\n", loadedModel.Target)
	fmt.Printf("- Feature: %s\n", loadedModel.Features)
	fmt.Printf("- Intercept: %.4f\n", loadedModel.Intercept)
	fmt.Printf("- Coefficients: %.v\n", loadedModel.Coefficients)
	fmt.Printf("- R-squared: %.4f\n", loadedModel.RSquared)

	return &loadedModel, nil

}
