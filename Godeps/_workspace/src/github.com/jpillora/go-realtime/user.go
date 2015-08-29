package realtime

import (
	"encoding/json"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type versions map[key]int64

type User struct {
	mut       sync.Mutex //protects all user fields
	Connected bool
	ID        string
	uptime    time.Time
	conn      *websocket.Conn
	versions  versions
	pending   []*update
}

func (u *User) sendPending() {
	u.mut.Lock()
	if len(u.pending) > 0 {
		b, _ := json.Marshal(u.pending)
		u.conn.Write(b)
		u.pending = nil
	}
	u.mut.Unlock()
}
