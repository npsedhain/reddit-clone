package main

import (
	"flag"
	"fmt"
	"reddit/messages"
	"reddit/simulation"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ActionType string

type SimulationConfig struct {
	NumClients    int
	NumEngines    int
	DurationSecs  int
}

func main() {
    // Define command-line flags
    config := SimulationConfig{}

    flag.IntVar(&config.NumClients, "clients", 1000, "number of client actors")
    flag.IntVar(&config.NumEngines, "engines", 10, "number of engine actors")
    flag.IntVar(&config.DurationSecs, "duration", 3, "simulation duration in seconds")

    // Parse command-line arguments
    flag.Parse()

    // Validate inputs
    if config.NumClients <= 0 || config.NumEngines <= 0 || config.DurationSecs <= 0 {
        fmt.Println("Error: All parameters must be positive numbers")
        flag.Usage()
        return
    }

    system := actor.NewActorSystem()

    // Create controller with configuration
    controller := simulation.NewSimulationController(system, config.NumEngines, config.NumClients)
    controllerProps := actor.PropsFromProducer(func() actor.Actor {
        return controller
    })
    system.Root.Spawn(controllerProps)

    // Run for specified duration
    simulationDuration := time.Duration(config.DurationSecs) * time.Second
    time.Sleep(simulationDuration)

    // Get and print metrics
    metrics := controller.GetMetrics()
    printMetrics(metrics)
}

func printMetrics(metrics map[string]interface{}) {
    fmt.Println("\nSimulation Results:")
    fmt.Printf("Total Requests: %d\n", metrics["total_requests"])
    fmt.Printf("Successful Requests: %d\n", metrics["successful_requests"])
    fmt.Printf("Failed Requests: %d\n", metrics["failed_requests"])

    if avgResponseTimes, ok := metrics["avg_response_times"].(map[messages.ActionType]time.Duration); ok {
        fmt.Println("\nAverage Response Times:")
        for action, duration := range avgResponseTimes {
            fmt.Printf("%s: %v\n", action, duration)
        }
    }

    // Try different type assertions
    if actionCounts, ok := metrics["action_counts"].(map[messages.ActionType]int64); ok {
        fmt.Println("\nAction Counts:")
        for action, count := range actionCounts {
            fmt.Printf("%s: %d\n", action, count)
        }
    }

    if errorCounts, ok := metrics["error_counts"].(map[messages.ActionType]int64); ok {
        fmt.Println("\nError Counts:")
        for action, count := range errorCounts {
            fmt.Printf("%s: %d\n", action, count)
        }
    }
}
