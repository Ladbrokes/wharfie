/*
 *   Wharfie - Companion for docker
 *   Copyright (c) 2014 Shannon Wynter, Ladbrokes Digital Australia Pty Ltd.
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *   Author: Shannon Wynter <http://fremnet.net/contact>
 */

package main

import (
	"archive/tar"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func archive(w http.ResponseWriter, root string) {
	log := logging.MustGetLogger("archive")
	log.Info("Distributing bundle from %s", root)

	w.Header().Set("Content-Type", "application/octet-stream")

	fw := &flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}

	tw := tar.NewWriter(fw)

	pathTrim := len(root) + 1

	visit := func(path string, f os.FileInfo, err error) error {
		if path == root {
			return nil
		}

		relpath := path[pathTrim:]
		target := relpath
		if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			t, err := os.Readlink(path)
			if err != nil {
				log.Error("Unable to resolve symlink %s: %#v", relpath, err)
				return nil
			}
			target = t
		}

		log.Debug("Adding %s>%s to archive", relpath, target)

		hdr, err := tar.FileInfoHeader(f, target)
		if err != nil {
			log.Error("Unable to create header from fileinfo for %s: %s", relpath, err)
			return nil
		}
		hdr.Name = relpath
		if f.IsDir() {
			hdr.Name = relpath + "/"
			if err := tw.WriteHeader(hdr); err != nil {
				log.Error("Unable to write tar header for %s: %s", relpath, err)
				return nil
			}
		} else if f.Mode()&os.ModeSymlink == os.ModeSymlink {
			if err := tw.WriteHeader(hdr); err != nil {
				log.Error("Unable to write tar header for %s: %s", relpath, err)
				return nil
			}
		} else {
			inf, err := ioutil.ReadFile(path)
			if err != nil {
				log.Error("Unable to read %s: %s", relpath, err)
				return nil
			}

			if err := tw.WriteHeader(hdr); err != nil {
				log.Error("Unable to write tar header for %s: %s", relpath, err)
				return nil
			}

			if _, err := tw.Write(inf); err != nil {
				log.Error("Unable to write tar body for %s: %s", relpath, err)
				return nil
			}
		}

		return nil
	}

	filepath.Walk(root, visit)

	if err := tw.Close(); err != nil {
		log.Error("Error while closing %s", err)
	}
}
