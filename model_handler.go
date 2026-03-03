package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/wtran29/go-ds-regression/model"
)

func getOrTrainModel(config Config, logger *slog.Logger) (*model.LinearRegression, *DataContext, error) {
	var dataModel *model.LinearRegression
	var err error
	var dataContext DataContext

	// application can either load a saved model or train a new one
	if config.LoadModelPath != "" {
		// TODO

	}

	// training a new model from a csv file
	if config.CsvFilePath == "" {
		flag.Usage()
		return nil, nil, fmt.Errorf("please provide a path to the csv file using the -file flag")
	}

	// load and prepare training data
	dataContext, err = loadAndPrepareData(config, logger)
	if err != nil {
		return nil, nil, err
	}

	// train linear regression model using the data
	logger.Info("Training model:", "features", config.FeatureVars)
	dataModel, err = model.TrainLinearRegression(dataContext.Data, config.FeatureVars, config.TargetVariable, config.Normalize)
	if err != nil {
		return nil, nil, fmt.Errorf("error training model: %v", err)
	}
	dataModel.PrintModelSummary()

	return dataModel, &dataContext, nil
}
