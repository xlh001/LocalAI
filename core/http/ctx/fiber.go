package fiberContext

import (
	"fmt"
	"strings"

	"github.com/go-skynet/LocalAI/pkg/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ModelFromContext returns the model from the context
// If no model is specified, it will take the first available
// Takes a model string as input which should be the one received from the user request.
// It returns the model name resolved from the context and an error if any.
func ModelFromContext(ctx *fiber.Ctx, loader *model.ModelLoader, modelInput string, firstModel bool) (string, error) {
	if ctx.Params("model") != "" {
		modelInput = ctx.Params("model")
	}

	// Set model from bearer token, if available
	bearer := strings.TrimLeft(ctx.Get("authorization"), "Bearer ")
	bearerExists := bearer != "" && loader.ExistsInModelPath(bearer)

	// If no model was specified, take the first available
	if modelInput == "" && !bearerExists && firstModel {
		models, _ := loader.ListModels()
		if len(models) > 0 {
			modelInput = models[0]
			log.Debug().Msgf("No model specified, using: %s", modelInput)
		} else {
			log.Debug().Msgf("No model specified, returning error")
			return "", fmt.Errorf("no model specified")
		}
	}

	// If a model is found in bearer token takes precedence
	if bearerExists {
		log.Debug().Msgf("Using model from bearer token: %s", bearer)
		modelInput = bearer
	}
	return modelInput, nil
}
