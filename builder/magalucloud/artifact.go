// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

package magalucloud

import "fmt"

type Artifact struct {
	ID        string
	Region    Region
	StateData map[string]any
}

func (*Artifact) BuilderId() string {
	return BuilderId
}

func (a *Artifact) Files() []string {
	return []string{}
}

func (a *Artifact) Id() string {
	return fmt.Sprintf("%s:%s", a.Region, a.ID)
}

func (a *Artifact) String() string {
	return fmt.Sprintf("A snapshot was created in the '%s' region: %s", a.Region, a.ID)
}

func (a *Artifact) State(name string) any {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	return nil
}
