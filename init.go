package dense

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// AIModelManager manages the lifecycle and growth of AI models.
type AIModelManager struct {
	EnvType          EnvType
	ProjectName      string
	InputSize        int
	OutputSize       int
	OutputTypes      []string
	ModelLocation    string
	Methods          []string // Array of methods like NEAT, HillClimb, DNAS, NAS, etc.
	LayerTypes       []string // Array of layer types like FFNN, LSTM, CNN
	NumModels        int      // Number of models to start with
	CycleAllMutations bool    // Flag to cycle through all mutations at the start
	Config           *NetworkConfig
	History          ProjectHistory
	TopX             int // Number of top models to track per generation
}

// GenerationData holds information about the best models in each generation.
type GenerationData struct {
	Generation int        `json:"generation"`
	Models     []ModelData `json:"models"`
}

// ModelData holds information about individual models in a generation.
type ModelData struct {
	ModelName      string  `json:"model_name"`
	Accuracy       float64 `json:"accuracy"`
	TrainingLoss   float64 `json:"training_loss,omitempty"`
	ValidationLoss float64 `json:"validation_loss,omitempty"`
}

// ProjectHistory keeps track of all generations and can act as a save point.
type ProjectHistory struct {
	ProjectName     string           `json:"project_name"`
	NumModels       int              `json:"num_models"`         // How many models to start with
	CycleAllMutations bool           `json:"cycle_all_mutations"` // Should we cycle through all mutations at the start?
	History         []GenerationData `json:"history"`
	TotalGenerations int             `json:"total_generations"`
	ModelConfig     *NetworkConfig   `json:"model_config"` // Store the latest network configuration
}

// Init initializes the manager with the project-specific parameters or loads from a save point.
func (mgr *AIModelManager) Init(projectName string, inputSize int, outputSize int, outputTypes []string, modelLocation string, methods []string, layerTypes []string, numModels int, cycleAllMutations bool, topX int, loadFilePath string) {
	mgr.ProjectName = projectName
	mgr.InputSize = inputSize
	mgr.OutputSize = outputSize
	mgr.OutputTypes = outputTypes
	mgr.ModelLocation = modelLocation
	mgr.Methods = methods
	mgr.LayerTypes = layerTypes
	mgr.NumModels = numModels
	mgr.CycleAllMutations = cycleAllMutations
	mgr.TopX = topX

	mgr.EnvType = GetEnvType()

	// Check if we should load from a saved state
	if loadFilePath != "" {
		err := mgr.loadProjectState(loadFilePath)
		if err != nil {
			fmt.Printf("Failed to load project state: %v\n", err)
			// Fall back to initializing a new project
			mgr.initialize()
		} else {
			fmt.Printf("Resuming project from save point: %s\n", loadFilePath)
			return
		}
	} else {
		// If no save file is passed, initialize a new project
		mgr.initialize()
	}
}

// initialize creates the neural network configuration dynamically based on the provided layer types.
func (mgr *AIModelManager) initialize() {
	fmt.Printf("Initializing AI Model Manager for project '%s' using methods %v and layers %v...\n", mgr.ProjectName, mgr.Methods, mgr.LayerTypes)
	fmt.Printf("Starting with %d models, cycling through all mutations: %v\n", mgr.NumModels, mgr.CycleAllMutations)

	// Random seed for reproducibility
	rand.Seed(time.Now().UnixNano())

	// Create a random network configuration based on the input and output sizes
	mgr.Config = CreateRandomNetworkConfig(mgr.InputSize, mgr.OutputSize, mgr.OutputTypes, "model-1", mgr.ProjectName)

	// Adjust the network configuration based on the layer types
	for _, layerType := range mgr.LayerTypes {
		switch layerType {
		case "FFNN":
			fmt.Println("Adding FFNN layer...")
			// Add FFNN layer configuration logic
		case "LSTM":
			fmt.Println("Adding LSTM layer...")
			// Add LSTM layer configuration logic
		case "CNN":
			fmt.Println("Adding CNN layer...")
			// Add CNN layer configuration logic
		default:
			fmt.Printf("Unknown layer type: %s\n", layerType)
		}
	}

	// Optionally cycle through all mutations
	if mgr.CycleAllMutations {
		//mgr.cycleThroughMutations()
	}

	// Initialize project history
	mgr.History = ProjectHistory{
		ProjectName:     mgr.ProjectName,
		NumModels:       mgr.NumModels,
		CycleAllMutations: mgr.CycleAllMutations,
		History:         []GenerationData{},
	}

	fmt.Println("Initialization complete. Neural network configuration created.")
	fmt.Printf("Model will be saved at: %s\n", mgr.ModelLocation)
}

/*
// cycleThroughMutations cycles through all the mutations at the start for each model.
func (mgr *AIModelManager) cycleThroughMutations() {
	fmt.Println("Cycling through all mutations for each model at the start...")
	// Implement the mutation cycling logic for each of the initial models
	for i := 0; i < mgr.NumModels; i++ {
		// Example mutation cycle:
		fmt.Printf("Applying mutations to model %d...\n", i+1)
		mgr.MutateNetwork(0.01) // Apply mutation (you can adjust this)
	}
}


// TrainModel executes the training based on the array of methods specified.
func (mgr *AIModelManager) TrainModel() {
	for _, method := range mgr.Methods {
		switch method {
		case "NEAT":
			fmt.Println("Training using NEAT")
			// Implement NEAT logic with species-level Hill Climb here
			mgr.trainNEATWithHillClimb()
		case "HillClimb":
			fmt.Println("Training using Hill Climbing")
			mgr.trainHillClimb()
		case "DNAS":
			fmt.Println("Training using DNAS")
			mgr.trainDNAS()
		case "NAS":
			fmt.Println("Training using NAS")
			mgr.trainNAS()
		default:
			fmt.Printf("Unknown training method: %s\n", method)
		}
	}
}*/

// trainHillClimb implements the Hill Climb logic.
func (mgr *AIModelManager) trainHillClimb() {
	fmt.Println("Performing Hill Climb optimization...")

	// Assume models have been trained and evaluated here
	// For example:
	models := []ModelData{
		{ModelName: "model1", Accuracy: 0.95, TrainingLoss: 0.1, ValidationLoss: 0.2},
		{ModelName: "model2", Accuracy: 0.93, TrainingLoss: 0.12, ValidationLoss: 0.22},
		// Add more models as necessary
	}

	// Update the history for this generation (e.g., generation 1)
	mgr.updateHistory(1, models)

	// Optionally export the history after each generation
	mgr.exportHistory()
}

// trainNEATWithHillClimb integrates NEAT with Hill Climb per species.
/*func (mgr *AIModelManager) trainNEATWithHillClimb() {
	fmt.Println("Performing NEAT with Hill Climb on species...")

	// Implement NEAT logic here, and call Hill Climb within each species
	for speciesID, species := range mgr.Config.Layers.Hidden {
		fmt.Printf("Evolving species %d using Hill Climb...\n", speciesID)
		// Apply Hill Climb to evolve species
		mgr.trainHillClimb() // Integrate Hill Climb within NEAT
	}
}*/

// trainDNAS implements DNAS logic, potentially across the entire population.
func (mgr *AIModelManager) trainDNAS() {
	fmt.Println("Performing DNAS optimization...")
	// Add DNAS logic here
}

// trainNAS implements NAS logic, focusing on architectural search.
func (mgr *AIModelManager) trainNAS() {
	fmt.Println("Performing NAS optimization...")
	// Add NAS logic here
}

// updateHistory records the best models of the current generation in the project history.
func (mgr *AIModelManager) updateHistory(generation int, models []ModelData) {
	// Sort models by accuracy (optional)
	// sort.Slice(models, func(i, j int) bool {
	// 	return models[i].Accuracy > models[j].Accuracy
	// })

	// Limit to top X models
	if len(models) > mgr.TopX {
		models = models[:mgr.TopX]
	}

	// Create generation data
	genData := GenerationData{
		Generation: generation,
		Models:     models,
	}

	// Append to history
	mgr.History.History = append(mgr.History.History, genData)
	mgr.History.TotalGenerations++
}

// exportHistory exports the project history to a JSON file.
func (mgr *AIModelManager) exportHistory() error {
	filePath := fmt.Sprintf("%s_project_history.json", mgr.ProjectName)
	data, err := json.MarshalIndent(mgr.History, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history to JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	fmt.Printf("Project history exported successfully to %s\n", filePath)
	return nil
}

// saveProjectState saves the current state of the project, including model configuration and history.
func (mgr *AIModelManager) saveProjectState() error {
	mgr.History.ModelConfig = mgr.Config // Save the current network configuration

	filePath := fmt.Sprintf("%s_save_state.json", mgr.ProjectName)
	data, err := json.MarshalIndent(mgr.History, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project state to JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write project state file: %w", err)
	}

	fmt.Printf("Project state saved successfully to %s\n", filePath)
	return nil
}

// loadProjectState loads the saved state of the project from a JSON file.
func (mgr *AIModelManager) loadProjectState(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read project state file: %w", err)
	}

	var loadedHistory ProjectHistory
	err = json.Unmarshal(data, &loadedHistory)
	if err != nil {
		return fmt.Errorf("failed to unmarshal project state: %w", err)
	}

	mgr.History = loadedHistory
	mgr.Config = mgr.History.ModelConfig // Restore the saved network configuration

	fmt.Printf("Project state loaded successfully from %s\n", filePath)
	return nil
}
