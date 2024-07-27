package statistic

import (
	"github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/headless"
	"time"

	"github.com/metacubex/mihomo/common/atomic"

	"github.com/puzpuzpuz/xsync/v3"
)

var DefaultManager *Manager
var ProxyManager *Manager
var DirectManager *Manager
var RejectManager *Manager

func init() {
	DefaultManager = &Manager{
		connections:   xsync.NewMapOf[string, Tracker](),
		uploadTemp:    atomic.NewInt64(0),
		downloadTemp:  atomic.NewInt64(0),
		uploadBlip:    atomic.NewInt64(0),
		downloadBlip:  atomic.NewInt64(0),
		uploadTotal:   atomic.NewInt64(0),
		downloadTotal: atomic.NewInt64(0),
	}

	ProxyManager = &Manager{
		connections:   xsync.NewMapOf[string, Tracker](),
		uploadTemp:    atomic.NewInt64(0),
		downloadTemp:  atomic.NewInt64(0),
		uploadBlip:    atomic.NewInt64(0),
		downloadBlip:  atomic.NewInt64(0),
		uploadTotal:   atomic.NewInt64(0),
		downloadTotal: atomic.NewInt64(0),
	}

	DirectManager = &Manager{
		connections:   xsync.NewMapOf[string, Tracker](),
		uploadTemp:    atomic.NewInt64(0),
		downloadTemp:  atomic.NewInt64(0),
		uploadBlip:    atomic.NewInt64(0),
		downloadBlip:  atomic.NewInt64(0),
		uploadTotal:   atomic.NewInt64(0),
		downloadTotal: atomic.NewInt64(0),
	}

	RejectManager = &Manager{
		connections:   xsync.NewMapOf[string, Tracker](),
		uploadTemp:    atomic.NewInt64(0),
		downloadTemp:  atomic.NewInt64(0),
		uploadBlip:    atomic.NewInt64(0),
		downloadBlip:  atomic.NewInt64(0),
		uploadTotal:   atomic.NewInt64(0),
		downloadTotal: atomic.NewInt64(0),
	}

	go DefaultManager.handle()
	go ProxyManager.handle()
	go DirectManager.handle()
	go RejectManager.handle()
}

type Manager struct {
	connections   *xsync.MapOf[string, Tracker]
	uploadTemp    atomic.Int64
	downloadTemp  atomic.Int64
	uploadBlip    atomic.Int64
	downloadBlip  atomic.Int64
	uploadTotal   atomic.Int64
	downloadTotal atomic.Int64
}

type Managers []*Manager

func (m *Manager) Join(c Tracker) {
	m.connections.Store(c.ID(), c)
}

func (m *Manager) Leave(c Tracker) {
	m.connections.Delete(c.ID())
}

func (m *Manager) Get(id string) (c Tracker) {
	if value, ok := m.connections.Load(id); ok {
		c = value
	}
	return
}

func (m *Manager) Range(f func(c Tracker) bool) {
	m.connections.Range(func(key string, value Tracker) bool {
		return f(value)
	})
}

func (m *Manager) PushUploaded(size int64) {
	m.uploadTemp.Add(size)
	m.uploadTotal.Add(size)
}

func (m *Manager) PushDownloaded(size int64) {
	m.downloadTemp.Add(size)
	m.downloadTotal.Add(size)
}

func (m *Manager) Now() (up int64, down int64) {
	return m.uploadBlip.Load(), m.downloadBlip.Load()
}

func (m *Manager) Snapshot() *Snapshot {
	var connections []*TrackerInfo
	m.Range(func(c Tracker) bool {
		connections = append(connections, c.Info())
		return true
	})
	return &Snapshot{
		UploadTotal:   m.uploadTotal.Load(),
		DownloadTotal: m.downloadTotal.Load(),
		Connections:   connections,
	}
}

func (m *Manager) ResetStatistic() {
	m.uploadTemp.Store(0)
	m.uploadBlip.Store(0)
	m.uploadTotal.Store(0)
	m.downloadTemp.Store(0)
	m.downloadBlip.Store(0)
	m.downloadTotal.Store(0)
}

func (m *Manager) handle() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			m.uploadBlip.Store(m.uploadTemp.Load())
			m.uploadTemp.Store(0)
			m.downloadBlip.Store(m.downloadTemp.Load())
			m.downloadTemp.Store(0)
		case <-headless.Register():
			return
		}
	}
}

type Snapshot struct {
	DownloadTotal int64          `json:"downloadTotal"`
	UploadTotal   int64          `json:"uploadTotal"`
	Connections   []*TrackerInfo `json:"connections"`
	Memory        uint64         `json:"memory"`
}

func (m *Managers) Join(c Tracker) {
	for _, m := range *m {
		m.connections.Store(c.ID(), c)
	}
}

func (m *Managers) Leave(c Tracker) {
	for _, m := range *m {
		m.connections.Delete(c.ID())
	}
}
func (m *Managers) PushUploaded(size int64) {
	for _, m := range *m {
		m.uploadTemp.Add(size)
		m.uploadTotal.Add(size)
	}
}

func (m *Managers) PushDownloaded(size int64) {
	for _, m := range *m {
		m.downloadTemp.Add(size)
		m.downloadTotal.Add(size)
	}
}

func GetManagers(c constant.Chain) *Managers {
	switch c.First() {
	case "PROXY":
		return &Managers{DefaultManager, ProxyManager}
	case "REJECT":
		return &Managers{DefaultManager, RejectManager}
	case "DIRECT":
		return &Managers{DefaultManager, DirectManager}
	}
	return &Managers{DefaultManager, ProxyManager}
}
