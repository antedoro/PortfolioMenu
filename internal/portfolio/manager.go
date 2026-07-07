package portfolio

import (
	"sync"

	"github.com/antedoro/PortfolioMenu/internal/models"
)

type Manager struct {
	mu sync.RWMutex

	Portfolio models.Portfolio
}

func NewManager() *Manager {

	return &Manager{}

}

func (m *Manager) GetPortfolio() models.Portfolio {

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Portfolio

}

func (m *Manager) SetPortfolio(p models.Portfolio) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Portfolio = p

}
