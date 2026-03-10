package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wtran29/go-ds-regression/model"
)

func main() {
	// parse command-line arguments
	config := parseCommandLineArgs()

	// set up a logger
	rootPath, _ := os.Getwd()
	dataLogger := slog.New(slog.Default().Handler())
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.StringValue(time.Now().Format("2006/01/02 03:04:05PM"))

			case slog.SourceKey:
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					// source.File = path.Base(source.File)
					source.File = strings.TrimPrefix(source.File, filepath.ToSlash(rootPath)+"/")

				}
			}
			return a
		},
	}))

	logger.Info("Parsed command line flags:", "features", config.FeatureVars)

	// Either load or train a model
	dataModel, dataContext, err := getOrTrainModel(config, dataLogger)
	if err != nil {
		logger.Error(fmt.Sprintf("Model error: %v", err))
	}

	// save model if requested
	if config.SaveModelPath != "" {
		if err := model.SaveModelToJSON(dataModel, config.SaveModelPath, config.ModelDesc, dataContext.Data.Nrow()); err != nil {
			logger.Error(fmt.Sprintf("Error saving model: %v", err))
		}
	}
}
