// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package siderolink provides server-side implementation of the SideroLink API.
package siderolink

// SecretName is the name of the Secret Sidero stores information about siderolink installation.
//
// Secret holds private Sidero Wireguard key and installation ID.
const SecretName = "siderolink"

// LogReceiverPort is the port of the log receiver container.
//
// LogReceiverPort is working only over Wireguard.
const LogReceiverPort = 4001

// Cfg is a default global instance of the SideroLink configuration.
//
// Cfg should be initialized first with `LoadOrCreate`.
var Cfg Config
