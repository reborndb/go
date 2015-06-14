// Copyright 2015 Reborndb Org. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package time2

import "time"

type BackoffDelay struct {
	min   time.Duration
	max   time.Duration
	delay time.Duration
}

func NewBackoffDelay(min, max time.Duration) *BackoffDelay {
	return &BackoffDelay{min, max, min}
}

func (bd *BackoffDelay) NextDelay() time.Duration {
	delay := bd.delay
	bd.delay = 2 * bd.delay
	if bd.delay > bd.max {
		bd.delay = bd.max
	}
	return delay
}

func (bd *BackoffDelay) Reset() {
	bd.delay = bd.min
}
