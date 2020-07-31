// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateFields performs validation on each field of the object
func (c *StdOutFilterConfig) ValidateFields(path *field.Path) field.ErrorList {
	// Can we validate the output_format field, do we have an exhaustive list of
	// all the values it can be?
	return nil
}