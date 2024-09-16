package main

import (
	"dense"
	"fmt"
	"path/filepath"
	"sort"
)

func main() {
	// Step 1: Define project parameters for MNIST dataset and AI model testing
	projectName := "AIModelTestProject"
	inputSize := 28 * 28 // Example input size for MNIST (28x28 pixel images)
	outputSize := 10     // Example output size for MNIST digits (0-9)
	outputTypes := []string{"softmax"} // Use softmax for classification
	modelLocation := "models"
	methods := []string{"HillClimb"}          // Define the optimization method
	layerTypes := []string{"FFNN", "CNN", "LSTM"} // Define types of layers to test
	numModels := 100                     // Number of models to create and test
	cycleAllMutations := true          // Flag to cycle through all mutations
	topX := 3                          // Number of top models to track
	numGenerations := 5                // Number of generations to run the hill climb

	// Load file path for project state
	loadFilePath := filepath.Join(modelLocation, projectName+"_save_state.json")

	// Step 2: Create the models folder if it doesn't exist
	err := dense.CreateDirectory(modelLocation)
	if err != nil {
		fmt.Printf("Error creating models folder: %v\n", err)
		return
	}

	// Step 3: Initialize the AIModelManager (either load from saved state or create new)
	manager := &dense.AIModelManager{}
	manager.Init(projectName, inputSize, outputSize, outputTypes, modelLocation, methods, layerTypes, numModels, cycleAllMutations, topX, loadFilePath)

	// Check which generation to resume from
	startGeneration := manager.History.CurrentGeneration
	fmt.Printf("Resuming from generation %d...\n", startGeneration)

	// Step 5: Load the MNIST dataset (downloading if necessary)
	err = dense.EnsureMNISTDownloads()
	if err != nil {
		fmt.Printf("Error downloading MNIST data: %v\n", err)
		return
	}
	mnist, err := dense.LoadMNISTOLD() // Loading the MNIST dataset
	if err != nil {
		fmt.Printf("Error loading MNIST: %v\n", err)
		return
	}

	// Split the MNIST data into training and testing sets
	_, testData := splitData(mnist)

	// Step 6: Run the hill-climbing algorithm across multiple generations
	for generation := startGeneration; generation < numGenerations; generation++ {
		fmt.Printf("Running generation %d...\n", generation)
		generationFolder := filepath.Join(modelLocation, fmt.Sprintf("%d", generation))

		// Create the first generation of models if it's the first run
		if generation == 0 {
			err := dense.CreateDirectory(generationFolder)
			if err != nil {
				fmt.Printf("Error creating generation folder: %v\n", err)
				return
			}

			// Create the first generation of models
			manager.CreateFirstGeneration(generationFolder)
			fmt.Println("First generation of models created.")
		}

		// Step 7: Evaluate each model and rank them by accuracy
		modelScores := evaluateAndRankModels(manager, generationFolder, numModels, testData)

		// Step 8: Select the top X models and mutate them to create the next generation
		if generation < numGenerations-1 {
			nextGenerationFolder := filepath.Join(modelLocation, fmt.Sprintf("%d", generation+1))
			err := dense.CreateDirectory(nextGenerationFolder)
			if err != nil {
				fmt.Printf("Error creating next generation folder: %v\n", err)
				return
			}

			createNextGeneration(manager, modelScores, nextGenerationFolder)
		}

		// Save the current generation number
		manager.History.CurrentGeneration = generation + 1
		err = manager.SaveProjectState()
		if err != nil {
			fmt.Printf("Error saving project state: %v\n", err)
		} else {
			fmt.Println("Project state saved successfully.")
		}
	}

	fmt.Println("Hill climbing optimization completed.")
}

// evaluateAndRankModels evaluates the fitness of each model and returns them sorted by accuracy
func evaluateAndRankModels(manager *dense.AIModelManager, generationFolder string, numModels int, testData *dense.MNISTData) []dense.ModelData {
	var modelScores []dense.ModelData

	// Evaluate each model
	for i := 0; i < numModels; i++ {
		modelFile := filepath.Join(generationFolder, fmt.Sprintf("model-%d.json", i+1))
		config, err := dense.LoadNetworkFromFile(modelFile)
		if err != nil {
			fmt.Printf("Error loading model %d: %v\n", i+1, err)
			continue
		}

		// Evaluate the fitness of the model on the test dataset
		fitness := evaluateFitness(config, testData)
		fmt.Printf("Model %d accuracy: %.4f%%\n", i+1, fitness*100)

		// Append the model data with fitness score
		modelScores = append(modelScores, dense.ModelData{
			ModelName: fmt.Sprintf("model-%d", i+1),
			Accuracy:  fitness,
		})
	}

	// Sort models by accuracy in descending order
	sort.Slice(modelScores, func(i, j int) bool {
		return modelScores[i].Accuracy > modelScores[j].Accuracy
	})

	return modelScores
}

// createNextGeneration creates the next generation of models by mutating the top models
func createNextGeneration(manager *dense.AIModelManager, modelScores []dense.ModelData, nextGenerationFolder string) {
    topX := manager.TopX
    learningRate := 0.01
    mutationRate := 20

    // Select the top X models and mutate them
    for i := 0; i < topX; i++ {
        modelName := modelScores[i].ModelName

        // Load the model configuration
        config, err := dense.LoadNetworkFromFile(filepath.Join(nextGenerationFolder, modelName+".json"))
        if err != nil {
            fmt.Printf("Error loading model %s: %v\n", modelName, err)
            continue
        }

        // Mutate the model
        manager.ApplyAllMutations(config, learningRate, mutationRate)

        // Save the mutated model to the next generation
        mutatedModelFile := filepath.Join(nextGenerationFolder, fmt.Sprintf("mutated-%s.json", modelName))
        err = dense.SaveNetworkConfig(config, mutatedModelFile)
        if err != nil {
            fmt.Printf("Error saving mutated model %s: %v\n", modelName, err)
        } else {
            fmt.Printf("Mutated model %s saved to %s\n", modelName, mutatedModelFile)
        }
    }
}



// Function to split the MNIST data into training (80%) and testing (20%)
func splitData(mnist *dense.MNISTData) (trainData, testData *dense.MNISTData) {
	totalImages := len(mnist.Images)
	splitIndex := int(float64(totalImages) * 0.8)

	trainData = &dense.MNISTData{
		Images: mnist.Images[:splitIndex],
		Labels: mnist.Labels[:splitIndex],
	}

	testData = &dense.MNISTData{
		Images: mnist.Images[splitIndex:],
		Labels: mnist.Labels[splitIndex:], // Fix here, use slice for remaining test data
	}

	return trainData, testData
}

// Evaluate the model's performance on the MNIST test dataset
func evaluateFitness(config *dense.NetworkConfig, mnist *dense.MNISTData) float64 {
	correct := 0
	total := len(mnist.Images)

	for i, image := range mnist.Images {
		// Prepare input data
		input := make(map[string]interface{})
		for j, pixel := range image {
			inputKey := fmt.Sprintf("input%d", j)
			input[inputKey] = float64(pixel) / 255.0 // Normalize pixel values
		}

		// Run the model's feedforward function
		outputs := dense.Feedforward(config, input)

		// Interpret the model output (e.g., predicted digit)
		predictedDigit := 0
		highestProb := 0.0
		for k := 0; k < 10; k++ {
			outputKey := fmt.Sprintf("output%d", k)
			if prob, ok := outputs[outputKey]; ok && prob > highestProb {
				highestProb = prob
				predictedDigit = k
			}
		}

		expectedDigit := int(mnist.Labels[i]) // Use the byte value directly
		if predictedDigit == expectedDigit {
			correct++
		}
	}

	accuracy := float64(correct) / float64(total)
	return accuracy
}
