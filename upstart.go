// Package upstart provides functions to control upstart jobs.
// Note that this is only usable on systems where upstart is the init system.
package upstart

import (
	"github.com/guelfey/go.dbus"
)

// StartJob starts the specified job.
func StartJob(name string) error {
	job := &job{Name: name}
	return job.StartWithEnv(map[string]string{})
}

// StopJob stops the specified job.
func StopJob(name string) error {
	job := &job{Name: name}
	return job.Stop()
}

// RestartJob restarts the specified job.
func RestartJob(name string) error {
	job := &job{Name: name}
	return job.Restart()
}

// ListJobInstances lists the instances of the specified job. Unless configured to run
// as multiple instances, a job has a single instance.
func ListJobInstances(name string) ([]string, error) {
	ret := []string{}
	job := &job{Name: name}

	printInst := func(_ *dbus.Conn, inst dbus.ObjectPath) error {
		ret = append(ret, string(inst))
		return nil
	}

	err := job.foreachInstance(printInst)
	return ret, err

}
