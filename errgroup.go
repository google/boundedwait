// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
//     Unless required by applicable law or agreed to in writing, software
//     distributed under the License is distributed on an "AS IS" BASIS,
//     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//     See the License for the specific language governing permissions and
//     limitations under the License

package boundedwait

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// ErrGroup provides a errgroup with a bounded amount of concurrency.
type ErrGroup struct {
	eg    *errgroup.Group
	limit chan struct{}
}

func WithContext(ctx context.Context, limit int) (*ErrGroup, context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	return &ErrGroup{limit: make(chan struct{}, limit), eg: eg}, ctx
}

func (g *ErrGroup) Go(f func() error) {
	g.limit <- struct{}{}
	g.eg.Go(func() error {
		defer func() {
			<-g.limit
		}()
		return f()
	})
}

func (g *ErrGroup) Wait() error {
	return g.eg.Wait()
}
