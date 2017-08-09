/*
   GOLANG REST API Skeleton

   Copyright (C) 2017 Carsten Seeger

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

   @author Carsten Seeger
   @copyright Copyright (C) 2017 Carsten Seeger
   @license http://www.gnu.org/licenses/gpl-3.0 GNU General Public License 3
   @link https://github.com/cseeger-epages/rest-api-go-skeleton
*/

package main

import (
	"crypto/sha256"
	_ "encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

/*
Default Handler Template
func Handler(w http.ResponseWriter, r*http.Request) {
	// caching stuff is handler specific
	w.Header().Set("Cache-Control", "no-store")

	qs := ParseQueryStrings(r)
	msg := HelpMsg{Message: "im a default Handler"}
	EncodeAndSend(w, r, qs, msg)
}
*/

// root handler giving basic API information
func Index(w http.ResponseWriter, r *http.Request) {
	// dont know what should happen here
	w.Header().Set("Cache-Control", "no-store")

	qs := ParseQueryStrings(r)
	message := fmt.Sprintf("Welcome to GOLANG REST API SKELETON please take a look at https://%s/help", r.Host)
	msg := HelpMsg{Message: message}
	EncodeAndSend(w, r, qs, msg)
}

// help reference for all routes
func Help(w http.ResponseWriter, r *http.Request) {
	// never cache help commands
	w.Header().Set("Cache-Control", "no-store")

	qs := ParseQueryStrings(r)
	message := fmt.Sprintf("http://%s/help/[cmd]", r.Host)
	msg := HelpMsg{Message: message}
	EncodeAndSend(w, r, qs, msg)
}

// specific route help reference
func HelpCmd(w http.ResponseWriter, r *http.Request) {
	// never cache help commands
	w.Header().Set("Cache-Control", "no-store")

	qs := ParseQueryStrings(r)
	vars := mux.Vars(r)
	cmd := vars["cmd"]
	message := fmt.Sprintf("help show: %s", cmd)
	msg := HelpMsg{Message: message}
	EncodeAndSend(w, r, qs, msg)
}

// Handler returning a list of all Projects
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	// 30 days cache since project did not change often
	w.Header().Set("Cache-Control", "max-age=2592000")

	qs := ParseQueryStrings(r)
	projects, err := GetProjects()
	Error(err)
	EncodeAndSend(w, r, qs, projects)
}

// Project specific information
func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	// 1 day cache since data only change once a day
	w.Header().Set("Cache-Control", "max-age=86400")

	qs := ParseQueryStrings(r)
	vars := mux.Vars(r)
	project := vars["project"]

	var msg interface{}
	var projects Projects
	var perr error = nil

	if pid, err := strconv.Atoi(project); err == nil {
		projects, perr = GetProject(pid)
	} else {
		project = strings.ToLower(project)
		projects, perr = GetProject(project)
	}
	if perr != nil {
		msg = ErrorMsg{"project does not exists"}
	} else {
		msg = projects
	}
	EncodeAndSend(w, r, qs, msg)
}

// Parse filter functions
func ParseQueryStrings(r *http.Request) QueryStrings {
	vals := r.URL.Query()

	// set defaults
	qs := QueryStrings{false}

	// Parse
	_, ok := vals["prettify"]
	if ok {
		qs.prettify = true
	}

	return qs
}

// Handles some filters and does what the name says
func EncodeAndSend(w http.ResponseWriter, r *http.Request, qs QueryStrings, msg interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// i need to encode the data twice for checking etag
	// and for sending with/without prettyfy maybe there
	// is a better way
	etagdata, _ := json.Marshal(msg)
	etagsha := sha256.Sum256([]byte(etagdata))
	etag := fmt.Sprintf("%x", etagsha)
	w.Header().Set("ETag", etag)

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	var err error

	if qs.prettify {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		err = encoder.Encode(msg)
	} else {
		err = json.NewEncoder(w).Encode(msg)
	}
	Error(err)
}
