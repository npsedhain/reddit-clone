package main

import (
	"fmt"
	"reddit/messages"
	"reddit/simulation"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ActionType string

func main() {
    system := actor.NewActorSystem()

    // Create controller (single instance)
    controller := simulation.NewSimulationController(system)
    controllerProps := actor.PropsFromProducer(func() actor.Actor {
        return controller
    })
    system.Root.Spawn(controllerProps)

    // Run for specific duration
    simulationDuration := 3 * time.Second
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
