package manager

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/christarazi/gravitational-challenge/server/models"
)

type Manager struct {
	sync.Mutex
	jobs []*models.Job
}

func NewManager() *Manager {
	return &Manager{
		jobs: []*models.Job{},
	}
}

func (m *Manager) IsAJob(id uint64) bool {
	m.Lock()
	defer m.Unlock()

	// The reason we are subtracting one here is because we want to make sure
	// that (id - 1) is an index within the length of the Jobs list. Because
	// jobs aren't removed from the list even when they've been stopped, the ID
	// is monotonically increasing.
	return (id - 1) < uint64(len(m.jobs))
}

// TODO: Should Job-specifc functions in here go in a dedicated separate file?

func (m *Manager) GetJobs() []*models.Job {
	m.Lock()
	defer m.Unlock()

	return m.jobs
}

func (m *Manager) GetJobByID(id uint64) *models.Job {
	m.Lock()
	defer m.Unlock()

	return m.jobs[id-1]
}

func (m *Manager) SetJobStatus(j *models.Job, status string) {
	m.Lock()
	defer m.Unlock()

	j.Status = status
}

func (m *Manager) AddAndStartJob(j *models.Job) (uint64, error) {
	m.Lock()
	defer m.Unlock()

	m.jobs = append(m.jobs, j)
	j.ID = uint64(len(m.jobs))
	j.Process = exec.Command(j.Command, j.Args...)
	j.Status = "Running"

	return j.ID, j.Process.Start()
}

func (m *Manager) StopJobByID(id uint64) error {
	m.Lock()
	defer m.Unlock()

	j := m.jobs[id-1]

	err := j.Process.Process.Kill()
	if err != nil {
		j.Status = "Failed to kill"
		return err
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
