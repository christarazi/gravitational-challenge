package manager

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/christarazi/gravitational-challenge/models"
)

// Manager is a struct that stores all the jobs on the server. The `jobs` list
// is protected by a mutex which is used when accessing and manipulating the
// list.
//
// Potentially several goroutines can be serving requests at any point in time,
// hence the need for the mutex.
//
// There are methods-on-struct defined for Manager to provide access to the
// `jobs`.
type Manager struct {
	sync.Mutex
	jobs []*models.Job
}

// NewManager creates an instance of the Manager struct.
func NewManager() *Manager {
	return &Manager{
		jobs: []*models.Job{},
	}
}

// Jobs returns a copy of the `jobs` list.
func (m *Manager) Jobs() []*models.Job {
	m.Lock()
	defer m.Unlock()

	return m.jobs
}

// JobStatus retrieves the status of a job with the requested ID.
func (m *Manager) JobStatus(rawID string) (*models.Job, error) {
	m.Lock()
	defer m.Unlock()

	id, err := convertIDToUint(rawID)
	if err != nil {
		return nil, err
	}

	if !m.isAJob(id) {
		return nil, fmt.Errorf("job with id %v does not exist", id)
	}

	return m.job(id), nil
}

// StartJob adds a given job to the list and starts the underlying process.
func (m *Manager) StartJob(j *models.Job) (uint64, error) {
	m.Lock()
	defer m.Unlock()

	m.jobs = append(m.jobs, j)
	j.ID = uint64(len(m.jobs))
	j.Process = exec.Command(j.Command, j.Args...)

	err := j.Process.Start()
	if err != nil {
		err = fmt.Errorf("failed to start job %d: %v", j.ID, err)
		j.Status = "Errored"
	} else {
		j.Status = "Running"
	}

	return j.ID, err
}

// StopJob stops a job by the given ID.
func (m *Manager) StopJob(id uint64) error {
	m.Lock()
	defer m.Unlock()

	if !m.isAJob(id) {
		return fmt.Errorf("job id %d does not exist", id)
	}

	j := m.job(id)
	if strings.Contains(j.Status, "Stopped") || strings.Contains(j.Status, "Killed") {
		log.Printf("manager-stop: job %d already stopped", id)
		return nil
	}

	return m.stop(j)
}

func (m *Manager) stop(j *models.Job) error {
	var (
		wg   sync.WaitGroup
		done chan error
		err  error
	)

	// Send a SIGTERM to the process so that it has a chance to clean up and
	// try to wait 5 seconds before it we forcefully kill it via SIGKILL.

	wg.Add(1)
	done = make(chan error, 1)
	go func() {
		if err := j.Process.Process.Signal(syscall.SIGTERM); err != nil {
			// We are not going to handle the error here because we can rely on
			// the timeout to eventually send a SIGKILL if the process is still
			// running. Just log a message for tracing.
			log.Println("manager-stop: failed to send SIGTERM")
		}
		done <- j.Process.Wait()
		wg.Done()
	}()

	select {
	// This value is hardcoded for now. In the future, this can be a
	// configurable value.
	case <-time.After(5 * time.Second):
		log.Println("manager-stop: timeout reached, sending SIGKILL")

		if err = j.Process.Process.Kill(); err != nil {
			log.Printf("manager-stop: failed to kill process: %v", err)

			j.Status = "Failed to kill"
			return fmt.Errorf("failed to kill process: %v", err)
		}

		log.Println("manager-stop: job killed forcefully")

		j.Status = "Killed"
	case err = <-done:
		log.Printf("manager-stop: job terminated: error: %v", err)
	}

	wg.Wait()

	j.Status = fmt.Sprintf("Stopped (ec: %d)",
		j.Process.ProcessState.ExitCode())

	return nil
}

// isAJob checks whether the given `id` is within the list.
func (m *Manager) isAJob(id uint64) bool {
	// The reason we are subtracting one here is because we want to make sure
	// that (id - 1) is an index into the `jobs` list. Because jobs aren't
	// removed from the list even when they've been stopped, the ID is
	// monotonically increasing.
	return (id - 1) < uint64(len(m.jobs))
}

func (m *Manager) job(id uint64) *models.Job {
	return m.jobs[id-1]
}

func convertIDToUint(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
