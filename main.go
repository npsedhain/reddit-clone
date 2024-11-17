package main

import (
	"fmt"
	"log"
	"time"

	"reddit/actors"
	"reddit/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
    system := actor.NewActorSystem()

    // Create the root user actor
    props := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewUserActor()
    })

    userPID := system.Root.Spawn(props)

    // Test registration
    registerMsg := &messages.RegisterUser{
        Username: "testuser",
        Password: "password123",
    }

    result, err := system.Root.RequestFuture(userPID, registerMsg, 5*time.Second).Result()
    if err != nil {
        log.Fatal(err)
    }

    if response, ok := result.(*messages.RegisterUserResponse); ok {
        if response.Success {
            fmt.Println("User registered successfully!")
        } else {
            fmt.Printf("Registration failed: %s\n", response.Error)
        }
    }

    // Test login
    loginMsg := &messages.LoginUser{
        Username: "testuser",
        Password: "password123",
    }

    result, err = system.Root.RequestFuture(userPID, loginMsg, 5*time.Second).Result()
    if err != nil {
        log.Fatal(err)
    }

    if response, ok := result.(*messages.LoginUserResponse); ok {
        if response.Success {
            fmt.Println("User logged in successfully!")
            fmt.Printf("Token: %s\n", response.Token)
        } else {
            fmt.Printf("Login failed: %s\n", response.Error)
        }
    }

    // Keep the program running for a moment to see the results
    time.Sleep(1 * time.Second)
}
