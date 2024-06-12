// utilities for managing xplane features

package feature

import "sync"

type Flag string

type Flags struct {
	enabled map[Flag]bool
	m       sync.RWMutex
}

func (fs *Flags) Enable(f Flag) {
	fs.m.Lock()
	if fs.enabled == nil {
		fs.enabled = make(map[Flag]bool)
	}
	fs.enabled[f] = true
	fs.m.Unlock()
}

func (fs *Flags) Enabled(f Flag) bool {
	if fs == nil {
		return false
	}
	fs.m.RLock()
	defer fs.m.RUnlock()
	return fs.enabled[f]
}
