package progress

import "time"

const DefaultPeriod = time.Second * 3

type ProgressTicker struct {
	UpdateAt time.Time
}

func (p *ProgressTicker) ShouldUpdate() bool {
	now := time.Now()
	if now.After(p.UpdateAt) {
		p.UpdateAt = now.Add(DefaultPeriod)
		return true
	}
	return false
}
