package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/wtran29/go-ds-regression/utils"
)

type DataContext struct {
	Data         dataframe.DataFrame
	FeatureData  map[string][]float64
	TargetValues []float64
}

func loadAndPrepareData(config Config, logger *slog.Logger) (DataContext, error) {
	var dataContext DataContext

	// read data from csv
	dataFile, err := os.Open(config.CsvFilePath)
	if err != nil {
		return dataContext, fmt.Errorf("could not open file: %v", err)
	}

	dataContext.Data = dataframe.ReadCSV(dataFile)

	// display a summary of the data for the user
	printDataSummary(dataContext.Data, logger, "before outlier removal")

	dataContext.Data, err = utils.ValidateData(dataContext.Data, config.FeatureVars, config.TargetVariable, config.OutlierLowerBound, config.OutlierUpperBound)
	if err != nil {
		return dataContext, fmt.Errorf("data validation error: %v", err)
	}

	if dataContext.Data.Nrow() > 0 {
		printDataSummary(dataContext.Data, logger, "after outlier removal")
	}

	if err := utils.CheckDatasetSize(dataContext.Data.Nrow(), len(config.FeatureVars)); err != nil {
		logger.Warn(err.Error())
	}

	dataContext.FeatureData = make(map[string][]float64)
	for _, feature := range config.FeatureVars {
		featureCol := dataContext.Data.Col(feature)
		featureValues := make([]float64, featureCol.Len())
		for i := range featureCol.Len() {
			featureValues[i] = featureCol.Elem(i).Float()
		}

		dataContext.FeatureData[feature] = featureValues
	}

	targetCol := dataContext.Data.Col(config.TargetVariable)
	dataContext.TargetValues = make([]float64, targetCol.Len())
	for i := range targetCol.Len() {
		dataContext.TargetValues[i] = targetCol.Elem(i).Float()
	}

	return dataContext, nil

}

func printDataSummary(df dataframe.DataFrame, logger *slog.Logger, stage string) {
	logger.Info(fmt.Sprintf("Data Preview (%s):\n", stage))
	logger.Info(df.Describe().String())
	logger.Info(fmt.Sprintf("Columns in dataset: %v\n", df.Names()))
	logger.Info(fmt.Sprintf("Row count: %v\n", df.Nrow()))

	// show sample rows
	if df.Nrow() > 0 {
		// get the minimum of 3 or the number of columns available
		numCols := min(df.Ncol(), 3)

		// get minimum of 5 or the number of rows available
		numRows := min(df.Nrow(), 5)

		columnsToShow := make([]int, numCols)
		for i := range numCols {
			columnsToShow[i] = i
		}

		logger.Info(df.Select(columnsToShow).Subset(numRows).String())
	}
}
