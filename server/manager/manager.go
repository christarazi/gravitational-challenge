package manager

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/christarazi/gravitational-challenge/server/models"
)

type Manager struct {
	Mutex *sync.Mutex
	Jobs  []*models.Job
}

func NewManager() *Manager {
	return &Manager{
		Mutex: &sync.Mutex{},
		Jobs:  []*models.Job{},
	}
}

func (m *Manager) IsAJob(id uint64) bool {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if (id - 1) >= uint64(len(m.Jobs)) {
		return false
	}

	return true
}

// TODO: Should Job-specifc functions in here go in a dedicated separate file?

func (m *Manager) GetJobByID(id uint64) *models.Job {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	return m.Jobs[id-1]
}

func (m *Manager) SetJobStatus(j *models.Job, status string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	j.Status = status
}

func (m *Manager) AddAndStartJob(j *models.Job) (uint64, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Jobs = append(m.Jobs, j)
	j.ID = uint64(len(m.Jobs))
	j.Process = exec.Command(j.Command, j.Args...)
	j.Status = "Running"

	return j.ID, j.Process.Start()
}

func (m *Manager) StopJobByID(id uint64) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	j := m.Jobs[id-1]

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
