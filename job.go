package upstart

import (
	"fmt"

	"github.com/godbus/dbus"
)

const (
	// Upstart is available on dbus at this "well known" address
	upstartServiceDBusPath = "com.ubuntu.Upstart"

	// the manager object where all our queries begin
	upstartManagerObject = "/com/ubuntu/Upstart"
)

// job represents an Upstart Job.
type job struct {
	Name string
}

// StartWithEnv starts this job (an instance, really) with the specified env.
func (j *job) StartWithEnv(env map[string]string) error {
	wait := true // TODO

	// connect to the system bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	jobpath, err := j.dbusPath(conn)
	if err != nil {
		return err
	}

	// manipulate env map into env arr
	envarr := []string{}
	for k, v := range env {
		envarr = append(envarr, fmt.Sprintf("%s=%s", k, v))
	}

	// start a job instance based on the env
	err = conn.
		Object(upstartServiceDBusPath, jobpath).
		Call("com.ubuntu.Upstart0_6.Job.Start", 0, envarr, wait).
		Store(&jobpath)
	if err != nil {
		return err
	}

	return nil
}

// Stop this job by stopping all its instances
func (j *job) Stop() error {
	wait := true // TODO

	doStop := func(conn *dbus.Conn, inst dbus.ObjectPath) error {
		err := conn.
			Object(upstartServiceDBusPath, inst).
			Call("com.ubuntu.Upstart0_6.Instance.Stop", 0, wait).
			Err
		if err != nil {
			return fmt.Errorf("Failed to stop %s instance %s: %s", j.Name, inst, err)
		}
		return nil
	}

	return j.foreachInstance(doStop)
}

// Restart all instances of this job
func (j *job) Restart() error {
	wait := true // TODO

	doRestart := func(conn *dbus.Conn, inst dbus.ObjectPath) error {
		err := conn.
			Object(upstartServiceDBusPath, inst).
			Call("com.ubuntu.Upstart0_6.Instance.Restart", 0, wait).
			Err
		if err != nil {
			return fmt.Errorf("Failed to restart inst %s: %s", inst, err)
		}
		return nil
	}

	return j.foreachInstance(doRestart)
}

//
// Helpers
//

// dbusPath returns the dbus path of the job object.
func (j *job) dbusPath(conn *dbus.Conn) (dbus.ObjectPath, error) {

	var jobpath dbus.ObjectPath

	// get the job path
	err := conn.
		Object(upstartServiceDBusPath, upstartManagerObject).
		Call("com.ubuntu.Upstart0_6.GetJobByName", 0, j.Name).
		Store(&jobpath)
	if err != nil {
		return jobpath, err
	}

	return jobpath, nil
}

// execute function 'f' for each instance of this job.
func (j *job) foreachInstance(f func(*dbus.Conn, dbus.ObjectPath) error) error {

	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	jobpath, err := j.dbusPath(conn)
	if err != nil {
		return err
	}

	// list the instances
	var instpaths []dbus.ObjectPath
	err = conn.
		Object(upstartServiceDBusPath, jobpath).
		Call("com.ubuntu.Upstart0_6.Job.GetAllInstances", 0).
		Store(&instpaths)
	if err != nil {
		return err
	}

	for _, inst := range instpaths {
		err = f(conn, inst)
		if err != nil {
			return err
		}
	}

	return nil
}
