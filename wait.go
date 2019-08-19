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

// boundedwait provides a wait group that uses a semaphore to control concurrent execution and will
// block during the add call.
package boundedwait

import "sync"

type Group struct {
	wg    sync.WaitGroup
	limit chan struct{}
}

func NewGroup(limit int) *Group {
	return &Group{limit: make(chan struct{}, limit)}
}

func (g *Group) Add(delta int) {
	g.wg.Add(delta)
	for i := 0; i < delta; i++ {
		g.limit <- struct{}{}
	}
}

func (g *Group) Done() {
	<-g.limit
	g.wg.Done()
}

func (g *Group) Wait() {
	g.wg.Wait()
}
