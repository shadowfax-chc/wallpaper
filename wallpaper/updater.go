package wallpaper

import (
	"context"
	"time"
)

// UpdaterConfig is used to configure an Updater.
type UpdaterConfig struct {
	Mode       Mode
	Repository Repository
	Frequency  time.Duration
}

// ReloadConfig is used to reload the Updater.
type ReloadConfig struct {
	Mode      Mode
	Location  string
	Shuffle   bool
	Frequency time.Duration
}

// Updater updates the background using images from the repository.
type Updater struct {
	Background *Background
	Repository Repository
	Frequency  time.Duration

	next chan struct{} // used to trigger the next image manually
}

// NewUpdater creates an Updater.
func NewUpdater(c *UpdaterConfig) *Updater {
	return &Updater{
		Background: &Background{
			Mode: c.Mode,
		},
		Repository: c.Repository,
		Frequency:  c.Frequency,
		next:       make(chan struct{}, 1),
	}
}

// Next tells the updater to set the next image available.
// This is used to update the image on demand instead of via the timer.
func (u *Updater) Next() {
	u.next <- struct{}{}
}

// Run starts up the updater.
func (u *Updater) Run(ctx context.Context) error {
	timer := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-u.next:
			timer.Stop()
			u.Background.Set(u.Repository.Next())
		case <-timer.C:
			timer.Stop()
			u.Background.Set(u.Repository.Next())
		}
		timer = time.NewTimer(u.Frequency)
	}
}

// Reload refreshes the repository and sets a new wallpaper
func (u *Updater) Reload(c *ReloadConfig) {
	if c != nil {
		u.Background.Mode = c.Mode
		u.Repository.SetLocation(c.Location)
		u.Repository.SetShuffle(c.Shuffle)
		u.Frequency = c.Frequency
	}
	u.Repository.Reload()
	u.Next()
}
