package main

import (
	"fmt"
	"os"

	"github.com/amoghe/go-upstart"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	jobName = kingpin.Flag("job-name", "Name of job.").String()
	action  = kingpin.Flag("action", "Action to perform").String()
)

func main() {
	var err error

	kingpin.Parse()

	switch *action {
	case "start":
		err = start(*jobName)
	case "stop":
		err = stop(*jobName)
	case "list":
		err = list(*jobName)
	case "restart":
		err = restart(*jobName)
	case "detect":
		detect()
	default:
		err = fmt.Errorf("Invalid action specified")
	}

	if err != nil {
		fmt.Printf("Error: %s (%T)\n", err, err)
		os.Exit(1)
	}
}

func start(name string) error {
	fmt.Println("Starting job", name)
	return upstart.StartJob(name)
}

func stop(name string) error {
	fmt.Println("Stopping job", name)
	return upstart.StopJob(name)
}

func restart(name string) error {
	fmt.Println("Restarting job", name)
	return upstart.RestartJob(name)
}

func list(name string) error {
	fmt.Println("Listing instances for job", name)

	instances, err := upstart.ListJobInstances(name)
	if err != nil {
		return err
	}

	for i, instance := range instances {
		fmt.Println("Instance", i, ":", instance)
	}

	return nil
}

func detect() {
	fmt.Printf("Upstart detected: %t\n", upstart.Detect())
}
