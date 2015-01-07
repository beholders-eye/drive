// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drive

import (
	"fmt"
)

const (
	Barely = iota
	AlmostExceeded
	HalfwayExceeded
	Exceeded
	Unknown
)

func (g *Commands) Quota() (err error) {
	about, err := g.rem.About()
	if err != nil {
		return err
	}

	freeBytes := about.QuotaBytesTotal - about.QuotaBytesUsed
	fmt.Printf(
		"Account type:\t%s\nBytes Used:\t%-20d (%s)\n"+
			"Bytes Free:\t%-20d (%s)\nTotal Bytes:\t%-20d (%s)\n",
		about.QuotaType,
		about.QuotaBytesUsed, prettyBytes(about.QuotaBytesUsed),
		freeBytes, prettyBytes(freeBytes),
		about.QuotaBytesTotal, prettyBytes(about.QuotaBytesTotal))
	return nil
}

func (g *Commands) QuotaStatus(query int64) (status int, err error) {
	if query < 0 {
		return Unknown, err
	}

	about, err := g.rem.About()
	if err != nil {
		return Unknown, err
	}

	// Sanity check
	if about.QuotaBytesTotal < 1 {
		return Unknown, fmt.Errorf("QuotaBytesTotal < 1")
	}

	toBeUsed := query + about.QuotaBytesUsed
	if toBeUsed >= about.QuotaBytesTotal {
		return Exceeded, nil
	}

	percentage := float64(toBeUsed) / float64(about.QuotaBytesTotal)
	if percentage < 0.5 {
		return Barely, nil
	}
	if percentage < 0.8 {
		return HalfwayExceeded, nil
	}
	return AlmostExceeded, nil
}
