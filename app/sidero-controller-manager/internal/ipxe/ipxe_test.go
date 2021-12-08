// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/ipxe"
)

func TestEmbeddedLength(t *testing.T) {
	var buf bytes.Buffer

	assert.NoError(t, ipxe.BootTemplate.Execute(&buf, struct {
		Endpoint string
		Port     string
	}{ // use bigger values here to get maximum length of the script
		Endpoint: "[2001:470:6d:30e:e5b8:903e:3701:7332]",
		Port:     "12345",
	}))

	// iPXE script should fit length of the reserved space in the iPXE binary
	assert.Less(t, buf.Len(), 1000)
}
