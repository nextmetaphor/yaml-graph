/*
 * Copyright 2020 Paul Tatham <paul@nextmetaphor.io>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_loadReportConf(t *testing.T) {
	t.Run("InvalidConfiguration", func(t *testing.T) {
		ms, err := loadReportConf("./_test/invalid-config.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, ms)
	})

	t.Run("ValidConfiguration", func(t *testing.T) {
		ms, err := loadReportConf("./_test/valid-config.yaml")

		assert.Nil(t, err)
		assert.NotNil(t, ms)
	})

	t.Run("MissingConfiguration", func(t *testing.T) {
		ms, err := loadReportConf("./_test/wheres-my-config.yaml")

		assert.NotNil(t, err)
		assert.Nil(t, ms)
	})

}
