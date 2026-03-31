package uconfig

import (
	"context"

	"github.com/omeid/uconfig/plugins"
)

func (c *config[C]) Watch(ctx context.Context, fn func(ctx context.Context, c *C) error) error {
	conf, err := c.Parse()
	if err != nil {
		return err
	}

	// Collect updaters.
	var updaters []plugins.Updater
	for _, p := range c.plugins {
		if u, ok := p.(plugins.Updater); ok {
			updaters = append(updaters, u)
		}
	}

	// No updaters: call fn once.
	if len(updaters) == 0 {
		return fn(ctx, conf)
	}

	// Start persistent update watchers. Each goroutine loops calling
	// Updated, sending to the shared channel on each change. All
	// goroutines exit when watchCtx is cancelled.
	watchCtx, watchCancel := context.WithCancel(ctx)
	defer watchCancel()
	changed := startUpdaters(watchCtx, updaters)

	for {
		// Run fn with a cancellable sub-context.
		runCtx, runCancel := context.WithCancel(ctx)
		fnDone := make(chan error, 1)
		go func() { fnDone <- fn(runCtx, conf) }()

		select {
		case <-changed:
			// Source changed: stop fn, re-parse.
			runCancel()
			<-fnDone

		case err := <-fnDone:
			// fn returned on its own — exit Watch.
			runCancel()
			return err

		case <-ctx.Done():
			runCancel()
			<-fnDone
			return ctx.Err()
		}

		// Re-parse config.
		newConf, err := c.Parse()
		if err != nil {
			// Bad config — wait for next change and retry.
			select {
			case <-changed:
			case <-ctx.Done():
				return ctx.Err()
			}
			newConf, err = c.Parse()
			if err != nil {
				continue // keep retrying
			}
		}
		conf = newConf
	}
}

// startUpdaters launches a goroutine per Updater that loops calling
// Updated and fans results into a single channel. The channel has
// capacity 1 so rapid changes coalesce.
func startUpdaters(ctx context.Context, updaters []plugins.Updater) <-chan struct{} {
	changed := make(chan struct{}, 1)
	for _, u := range updaters {
		go func(u plugins.Updater) {
			for {
				if !u.Updated(ctx) {
					return // context cancelled
				}
				select {
				case changed <- struct{}{}:
				default:
				}
			}
		}(u)
	}
	return changed
}
