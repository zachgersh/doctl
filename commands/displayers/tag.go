/*
Copyright 2018 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package displayers

import (
	"github.com/digitalocean/doctl/do"
)

type Tag struct {
	Tags do.Tags
}

var _ Displayable = &Tag{}

func (t *Tag) JSON() (string, error) {
	return writeJSON(t.Tags)
}

func (t *Tag) Cols() []string {
	return []string{"Name", "DropletCount"}
}

func (t *Tag) ColMap() map[string]string {
	return map[string]string{
		"Name":         "Name",
		"DropletCount": "Droplet Count",
	}
}

func (t *Tag) KV() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, x := range t.Tags {
		dropletCount := x.Resources.Droplets.Count
		o := map[string]interface{}{
			"Name":         x.Name,
			"DropletCount": dropletCount,
		}
		out = append(out, o)
	}

	return out
}
