package lcu

import "sync"

// PhaseBroadcaster manages subscribers of game phase updates.
type PhaseBroadcaster struct {
	mu   sync.RWMutex
	subs map[chan string]struct{}
}

// NewPhaseBroadcaster creates a new broadcaster.
func NewPhaseBroadcaster() *PhaseBroadcaster {
	return &PhaseBroadcaster{subs: make(map[chan string]struct{})}
}

// Subscribe returns a channel to receive game phase updates.
func (b *PhaseBroadcaster) Subscribe() chan string {
	ch := make(chan string, 1)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes a subscriber and closes the channel.
func (b *PhaseBroadcaster) Unsubscribe(ch chan string) {
	b.mu.Lock()
	if _, ok := b.subs[ch]; ok {
		delete(b.subs, ch)
		close(ch)
	}
	b.mu.Unlock()
}

// Broadcast sends an update to all subscribers.
func (b *PhaseBroadcaster) Broadcast(phase string) {
	b.mu.RLock()
	for ch := range b.subs {
		select {
		case ch <- phase:
		default:
		}
	}
	b.mu.RUnlock()
}
