package runtime

import (
	cronlib "github.com/hjwalt/platform/cron"
)

type cron struct {
	cron *cronlib.Cron
}

func (r *cron) Start() error {
	r.cron.Start()
	return nil
}

func (r *cron) Stop() {
	r.cron.Stop()
}
