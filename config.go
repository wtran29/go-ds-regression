package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	CsvFilePath       string
	SaveModelPath     string
	LoadModelPath     string
	TargetVariable    string
	ModelDesc         string
	Normalize         bool
	FeatureVars       []string // list of feature variables
	DataToPredict     string   // key value pairs of feature variables
	OutlierLowerBound float64
	OutlierUpperBound float64
	Plot              bool
	PlotURI           string
}

func parseCommandLineArgs() Config {
	var config Config
	var features string
	var target string

	flag.StringVar(&config.CsvFilePath, "file", "house_data.csv", "Path to the CSV file containing data")
	flag.StringVar(&config.SaveModelPath, "save", "", "Path to save the trained model")
	flag.StringVar(&config.LoadModelPath, "load", "", "Path to load a previously trained model")
	flag.StringVar(&config.ModelDesc, "desc", "", "Description of the model (used when saving)")
	flag.BoolVar(&config.Normalize, "normalize", true, "Normalize features (default: true)")
	flag.BoolVar(&config.Plot, "plot", false, "Generate a plot (default: false)")
	flag.StringVar(&features, "features", "", "Comma separated list of features that you want to use")
	flag.StringVar(&target, "target", "", "Name of the target column")
	flag.Float64Var(&config.OutlierLowerBound, "lower-bound", 1.5, "Lower bound multiplier for outlier detection (default: 1.5)")
	flag.Float64Var(&config.OutlierUpperBound, "upper-bound", 1.5, "Upper bound multiplier for outlier detection (default: 1.5)")
	flag.StringVar(&config.DataToPredict, "predict", "", "Enter key value pairs for prediction, e.g. square_footage=1000,bedrooms=2")
	flag.StringVar(&config.PlotURI, "plot-uri", "http://localhost:8000", "URI for the plot app")

	flag.Parse()

	if (len(features) == 0 || len(target) == 0) && config.LoadModelPath == "" {
		fmt.Println("You must specify at least one feature and one target.")
		os.Exit(1)
	}

	if config.OutlierLowerBound < 0 || config.OutlierUpperBound < 0 {
		fmt.Println("Outlier bound mutipliers must be positive values.")
		os.Exit(1)
	}

	var featureList []string

	for _, f := range strings.Split(features, ",") {
		featureList = append(featureList, strings.TrimSpace(f))
	}

	config.FeatureVars = featureList

	if len(target) > 0 {
		config.TargetVariable = target
	}

	return config

}
