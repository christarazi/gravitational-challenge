package manager

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"

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

	if !m.IsAJob(id) {
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

	return m.stop(m.jobs[id-1])
}

func (m *Manager) stop(j *models.Job) error {
	err := j.Process.Process.Kill()
	if err != nil {
		j.Status = "Failed to kill"
		return fmt.Errorf("failed to kill job %d: %v", j.ID, err)
	}

	// TODO: We still want to wait on the process to retrieve its correct exit
	// code. Need to figure out how to best do this without blocking the
	// request.
	// err = j.Process.Wait()
	// if err != nil {
	// 	j.Status = "Failed to wait"
	// 	return err
	// }

	// TODO: For now the exit code will always be -1 until we do the above
	// todo.
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
