package main

import (
	"fmt"
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
    simulationDuration := 5 * time.Second
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

    // fmt.Println("\nAverage Response Times:")
    // for action, duration := range metrics["avg_response_times"].(map[ActionType]time.Duration) {
    //     fmt.Printf("%s: %v\n", action, duration)
    // }

    // fmt.Println("\nAction Counts:")
    // for action, count := range metrics["action_counts"].(map[ActionType]int64) {
    //     fmt.Printf("%s: %d\n", action, count)
    // }

    // fmt.Println("\nError Counts:")
    // for action, count := range metrics["error_counts"].(map[ActionType]int64) {
    //     fmt.Printf("%s: %d\n", action, count)
    // }
}
