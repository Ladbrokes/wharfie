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
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"github.com/vharitonsky/iniflags"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	config       = &configuration{}
	printVersion = false
)

var log = logging.MustGetLogger("wharfie")

func init() {
	cwd, _ := os.Getwd()
	flag.StringVar(&config.BindHost, "host", "172.17.42.1", "IP Address to bind to")
	flag.UintVar(&config.BindPort, "port", 2864, "Port to listen on")
	flag.StringVar(&config.BuildPath, "path", cwd, "Path to build from")
	flag.BoolVar(&printVersion, "version", printVersion, "Display the current version")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	log.Warning("exists? %s", err)

	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	log := logging.MustGetLogger("request")
	log.Info("Handling %s", r.URL.Path)

	pathParts := strings.Split(r.URL.Path[1:], "/")
	image := pathParts[0]
	buildRoot := path.Join(config.BuildPath, image)

	if !exists(buildRoot) {
		log.Notice("404 %s Not found", buildRoot)
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 && pathParts[1] == "bundle" {
		log.Info("Handling archive request %s", r.URL.Path)
		provisionPath := path.Join(buildRoot, "provision")
		if !exists(provisionPath) {
			http.Error(w, "Missing provisioning directory", http.StatusNotImplemented)
			return
		}
		archive(w, provisionPath)
	} else if len(pathParts) == 1 {
		if r.Host == "" {
			log.Warning("Bad request, no host header")
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		fmt.Fprintf(w, "#!/bin/bash\n# Build script for %[1]s\nmkdir -p /tmp/provision\necho -e \"GET /%[1]s/bundle HTTP/1.0\\r\\nhost: %[2]s\\r\\n\\r\\n\" | nc %[2]s %[3]d | tail -n+5 | tar -x -C /tmp/provision; exec /tmp/provision/deploy.sh", image, r.Host, config.BindPort)
	} else {
		log.Warning("Not sure what to do with %s", r.URL.Path)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func main() {
	iniflags.Parse()

	if printVersion {
		fmt.Printf("wharfie %s\n", Version)
		os.Exit(0)
	}

	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	syslogBackend, err := logging.NewSyslogBackend("")
	if err != nil {
		log.Fatal(err)
	}

	logging.SetBackend(consoleBackend, syslogBackend)
	logging.SetFormatter(logging.MustStringFormatter("%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"))

	log := logging.MustGetLogger("main")

	if !exists(config.BuildPath) {
		log.Fatalf("Build path %s does not exist", config.BuildPath)
	}

	http.HandleFunc("/", handler)

	log.Info("Opening http socket on %s:%d resolving to builds in %s", config.BindHost, config.BindPort, config.BuildPath)

	if http.ListenAndServe(fmt.Sprintf("%s:%d", config.BindHost, config.BindPort), nil) != nil {
		log.Fatalf("Failed to start serving: %s", err)
	}
}
