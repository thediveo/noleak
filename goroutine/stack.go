// Copyright 2022 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package goroutine

import "runtime"

const startStackBufferSize = 64 * 1024 // 64kB

// stacks returns stack trace information for either all goroutines or only the
// current goroutine. It is a convenience wrapper around runtime.Stack, hiding
// the result allocation.
func stacks(all bool) []byte {
	for size := startStackBufferSize; ; size *= 2 {
		buffer := make([]byte, size)
		if n := runtime.Stack(buffer, all); n < size {
			return buffer[:n]
		}
	}
}
