package local

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// LocalDataRepository handles OS-level data operations
type LocalDataRepository struct{}

// NewLocalDataRepository creates a new LocalDataRepository
func NewLocalDataRepository() *LocalDataRepository {
	return &LocalDataRepository{}
}

// GetSystemInfo retrieves basic system information
func (r *LocalDataRepository) GetSystemInfo() (*model.SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &model.SystemInfo{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		CPUCount:     runtime.NumCPU(),
		Timestamp:    time.Now(),
	}, nil
}

// WriteDataToFile writes data to local file
func (r *LocalDataRepository) WriteDataToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0o644)
}

// ReadDataFromFile reads data from local file
func (r *LocalDataRepository) ReadDataFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
