/*
Copyright 2020 Talos Systems, Inc.
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

package tftp

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pin/tftp"
	"github.com/talos-systems/talos/pkg/safepath"
)

// readHandler is called when client starts file download from server.
func readHandler(filename string, rf io.ReaderFrom) error {
	filename = filepath.Join("/var/lib/sidero/tftp", safepath.CleanPath(filename))

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%v", err)

		return err
	}

	n, err := rf.ReadFrom(file)
	if err != nil {
		log.Printf("%v", err)

		return err
	}

	log.Printf("%d bytes sent", n)

	return nil
}

func ServeTFTP() error {
	if err := os.MkdirAll("/var/lib/sidero/tftp", 0o777); err != nil {
		return err
	}

	s := tftp.NewServer(readHandler, nil)

	// A standard TFTP server implementation receives requests on port 69 and
	// allocates a new high port (over 1024) dedicated to that request. In single
	// port mode, the same port is used for transmit and receive. If the server
	// is started on port 69, all communication will be done on port 69.
	// This option is required since the Kubernetes service definition defines a
	// single port.
	s.EnableSinglePort()
	s.SetTimeout(5 * time.Second)

	if err := s.ListenAndServe(":69"); err != nil {
		return err
	}

	return nil
}
