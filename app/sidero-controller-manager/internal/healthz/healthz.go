// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package healthz

import (
	"net/http"
)

func RegisterServer(mux *http.ServeMux) error {
	mux.HandleFunc("/healthz", healthzHandler)

	return nil
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	// do nothing, consider to be healthy always
}
