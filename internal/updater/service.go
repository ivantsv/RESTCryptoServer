package updater

import (
	"RESTCryptoServer/internal/crypto"
	"log"
	"sync"
	"time"
)

type Updater struct {
	UpdateTime    time.Duration
	CryptoService *crypto.CryptoService
	mu            sync.Mutex
	StopChan      chan struct{}
	LastUpdate    time.Time 
	Enabled       bool
}

func NewUpdater(cs *crypto.CryptoService, t time.Duration) (*Updater) {
	if t == 0 {
		t = 30 * time.Second
	}

	return &Updater{
		UpdateTime: t,
		CryptoService: cs,
		StopChan: make(chan struct{}),
		Enabled: false,
	}
}

func (u *Updater) StartUpdating() {
	u.mu.Lock()
	defer u.mu.Unlock()

	if !u.Enabled {
		return
	}

	log.Printf("Updater started with interval: %s", u.UpdateTime)

	go func() {
		ticker := time.NewTicker(u.UpdateTime)
		defer ticker.Stop()

		for {
			select {
			case <- ticker.C:
				u.CryptoService.UpdateAllCryptos()
				u.LastUpdate = time.Now()
			case <- u.StopChan:
				return
			}
		}
	}()
}

func (u *Updater) EndUpdating() {
	if !u.Enabled {
		return
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	select {
	case <-u.StopChan:
		log.Println("Updater: stop requested, but already stopped")
	default:
		close(u.StopChan)
		log.Println("Updater: stop requested, stopping updater")
	}
}

func (u *Updater) RestartUpdating(seconds int) {
	if !u.Enabled {
		return
	}

	u.EndUpdating()
	u.mu.Lock()
	defer u.mu.Unlock()

	u.UpdateTime = time.Duration(seconds) * time.Second
	u.StopChan = make(chan struct{})

	log.Printf("Updater: restarting with new interval %s", u.UpdateTime)
	u.StartUpdating()
}

func (u *Updater) GetUpdateTime() int {
	u.mu.Lock()
	defer u.mu.Unlock()

	return int(u.UpdateTime)
}

func (u *Updater) GetLastUpdate() time.Time {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.LastUpdate
}

func (u *Updater) IsEnabled() bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.Enabled
}

func (u *Updater) Update() (int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	cnt, err := u.CryptoService.UpdateAllCryptos() 
	return cnt, err
}