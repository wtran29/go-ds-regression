package model

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/wtran29/go-ds-regression/utils"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type LinearRegression struct {
	Coefficients    []float64 `json:"coefficients,omitempty"`
	Intercept       float64   `json:"intercept,omitempty"`
	Features        []string  `json:"features,omitempty"`
	Target          string    `json:"target"`
	RSquared        float64   `json:"r_squared,omitempty"`
	FeatureMeans    []float64 `json:"feature_means,omitempty"`
	FeaturesStdDevs []float64 `json:"feature_std_devs,omitempty"`
	IsNormalized    bool      `json:"is_normalized,omitempty"`
	SavedAt         time.Time `json:"-"`
	Description     string    `json:"description,omitempty"`
	NumSamples      int       `json:"num_samples,omitempty"`
	Version         string    `json:"version,omitempty"`
}

func TrainLinearRegression(dataFrame dataframe.DataFrame, featureNames []string, targetName string, normalize bool) (*LinearRegression, error) {
	// check if all feature columns and target column exist
	columnNames := dataFrame.Names()
	for _, name := range featureNames {
		if !slices.Contains(columnNames, name) {
			return nil, fmt.Errorf("feature column '%s' not found in the dataset", name)
		}
	}

	if !slices.Contains(columnNames, targetName) {
		return nil, fmt.Errorf("target column '%s' not found in the dataset", targetName)
	}

	// extract features and target from the dataframe
	var featureMatrix [][]float64 // Will hold all feature values (X)
	var targetValues []float64    // Will hold all target values (y)

	// get feature columns as float slices
	featureColumns := make([]series.Series, len(featureNames))
	for i, name := range featureNames {
		featureColumns[i] = dataFrame.Col(name)
	}

	// get target column as a float slice
	targetColumn := dataFrame.Col(targetName)

	numSamples := dataFrame.Nrow()
	featureMatrix = make([][]float64, numSamples)
	targetValues = make([]float64, numSamples)

	// fill feature matrix (X) and target vector (Y) with values from the dataframe
	for rowIndex := range numSamples {
		featureMatrix[rowIndex] = make([]float64, len(featureNames))
		for colIndex, column := range featureColumns {
			featureMatrix[rowIndex][colIndex] = column.Elem(rowIndex).Float()
		}
		targetValues[rowIndex] = targetColumn.Elem(rowIndex).Float()
	}

	// variables to store normalization parameters
	var normalizedFeatures [][]float64
	var featureMeans []float64
	var featureStdDevs []float64

	// Normalize features if requested
	// Normalization helps when features have different scales (like sq ft vs bedroom count)
	// It transforms features to have mean 0 and standard deviation 1, which:
	// - Makes different features comparable (e.g., sq.ft. vs. bedroom count)
	// - Can improve numerical stability during training
	// - Helps interpret relative importance of features
	if normalize {
		normalizedFeatures, featureMeans, featureStdDevs = utils.NormalizeFeatures(featureMatrix)
		featureMatrix = normalizedFeatures
	}

	// create a design matrix
	numFeatures := len(featureNames)

	// Create matrices for gonum linear algebra operations
	// In matrix form, we want to solve: y = Xβ
	// where:
	// - y is the target vector
	// - X is the design matrix (feature values with a column of 1s)
	// - β is the coefficient vector (what we're solving for)
	// Create matrices for gonum linear algebra operations
	designMatrix := mat.NewDense(numSamples, numFeatures+1, nil)
	targetVector := mat.NewVecDense(numSamples, nil)

	for rowIndex := range numSamples {
		designMatrix.Set(rowIndex, 0, 1.0)
		for featureIndex := range numFeatures {
			designMatrix.Set(rowIndex, featureIndex+1, featureMatrix[rowIndex][featureIndex])
		}
		targetVector.SetVec(rowIndex, targetValues[rowIndex])
	}

	// step 1. calculate X^T x (transpose of X multiplied by X)
	var transposeTimesDesign mat.Dense
	transposeTimesDesign.Mul(designMatrix.T(), designMatrix)

	// step 2. calculate (X^T X)^(-1)
	var inverseMatrix mat.Dense
	if err := inverseMatrix.Inverse(&transposeTimesDesign); err != nil {
		return nil, fmt.Errorf("failed to compute inverse: %v - matrix may be singular; try adding more data or removing highly correlated features", err)
	}

	// step 3. calculate X^T y
	var transposeTimesTarget mat.Dense
	transposeTimesTarget.Mul(designMatrix.T(), targetVector)

	// step 4. calculate the optimal coefficients
	var coefficientMatrix mat.Dense
	coefficientMatrix.Mul(&inverseMatrix, &transposeTimesTarget)

	// extract coefficients
	interceptAndCoefficients := make([]float64, numFeatures+1)
	for i := range numFeatures + 1 {
		interceptAndCoefficients[i] = coefficientMatrix.At(i, 0)
	}

	// calculate predictions using the trained model
	predictedValues := make([]float64, numSamples)
	for i := range numSamples {
		// start with intercept
		predictedValues[i] = interceptAndCoefficients[0]

		// add contribution of each feature
		for j := range numFeatures {
			predictedValues[i] += interceptAndCoefficients[j+1] * featureMatrix[i][j]
		}
	}
	// calculate R-squared
	targetMean := stat.Mean(targetValues, nil)

	var totalSumOfSquares, sumOfSquarResiduals float64
	for i := range numSamples {
		totalSumOfSquares += math.Pow(targetValues[i]-targetMean, 2)
		sumOfSquarResiduals += math.Pow(targetValues[i]-predictedValues[i], 2)
	}

	rSquared := 1 - (sumOfSquarResiduals / totalSumOfSquares)

	return &LinearRegression{
		Intercept:       interceptAndCoefficients[0],
		Coefficients:    interceptAndCoefficients[1:],
		Features:        featureNames,
		Target:          targetName,
		RSquared:        rSquared,
		FeatureMeans:    featureMeans,
		FeaturesStdDevs: featureStdDevs,
		IsNormalized:    normalize,
	}, nil
}

func (lr *LinearRegression) PrintModelSummary() {
	fmt.Println("\n==== Model Summary ====\n")
	fmt.Printf("Regression Equation: %s - %.4f", lr.Target, lr.Intercept)
	for i, feature := range lr.Features {
		if lr.Coefficients[i] >= 0 {
			fmt.Printf(" + %.4f x %s", lr.Coefficients[i], feature)
		} else {
			fmt.Printf(" - %.4f x %s", -lr.Coefficients[i], feature)
		}
	}
	fmt.Println()

	// display model fit statistics
	fmt.Printf("\nModel Performance:\n")
	fmt.Printf("- R-squared: %.4f\n", lr.RSquared)
	fmt.Printf("- Interpretation: %.2f%% of variance in %s is explained by theis model\n", lr.RSquared*100, lr.Target)

	fmt.Printf("\nCoefficient Interpretation:\n")
	fmt.Printf("- Intercept (%.4f): The base %s when all features are zero\n", lr.Intercept, lr.Target)

	for i, feature := range lr.Features {
		fmt.Printf("- %s Coefficient (%.4f): For each additional unit of %s, %s changesby %.4f units\n", feature, lr.Coefficients[i], feature, lr.Target, lr.Coefficients[i])
	}

	if lr.IsNormalized {
		fmt.Printf("\nNote: This model was trained on normalized data. Predictions on new data will automatically be normalized.\n")
	}
}
