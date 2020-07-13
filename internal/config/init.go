/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package config

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"regexp"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/memory"
	"git.maxset.io/web/knaxim/internal/database/mongo"
	"git.maxset.io/web/knaxim/internal/handlers/spa"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"github.com/google/go-tika/tika"
	yaml "gopkg.in/yaml.v3"
)

// V => Configuration of Knaxim
// populate with ParseConfig
var V Configuration

// DB => database instance specified by configuration
var DB database.Database

// T => Tika server and connection information
// if Server is nil, path should hold the http path to the tika server
var T struct {
	Path   string
	Server *tika.Server
}

// StaticHandler is a http.Handler generated from configuration to
// handle reuests from server for static files. also will respond with
// index.html if the path does not map to any static file, on the
// assumption that the path is relevant to a single page application
var StaticHandler spa.Handler

var extensionRegex = regexp.MustCompile("\\.[^.]*$")

// ParseConfig loads configuration file and populates global vars
func ParseConfig(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	var isYAML bool
	switch extensionRegex.FindString(path) {
	case ".yml":
		fallthrough
	case ".yaml":
		isYAML = true
		err = yaml.NewDecoder(fp).Decode(&V)
	default:
		err = json.NewDecoder(fp).Decode(&V)
	}
	if err != nil {
		return err
	}
	fp.Close()
	if V.FileLimit == 0 {
		V.FileLimit = 5 * 1024 * 1024
	} else if V.FileLimit < 0 {
		V.FileLimit = math.MaxInt64
	}
	StaticHandler.StaticPath = V.StaticPath
	StaticHandler.IndexPath = V.IndexPath
	switch V.DatabaseType {
	case "mongo":
		DB = new(mongo.Database)
	case "memory":
		DB = new(memory.Database)
	default:
		return errors.New("Unrecognized config database type")
	}
	if isYAML {
		err = V.Database.YAML.Decode(DB)
	} else {
		err = json.Unmarshal(V.Database.JSON, DB)
	}
	if err != nil {
		return err
	}
	if V.Tika.Type == "local" {
		T.Server, err = tika.NewServer(V.Tika.Path, V.Tika.Port)
		if err != nil {
			return err
		}
		T.Server.ChildMode(&tika.ChildOptions{
			MaxFiles:          V.Tika.MaxFiles,
			TaskPulseMillis:   V.Tika.TaskPulse,
			TaskTimeoutMillis: V.Tika.TaskTimeout,
			PingPulseMillis:   V.Tika.PingPulse,
			PingTimeoutMillis: V.Tika.PingTimeout,
		})
		T.Path = T.Server.URL()
	} else if V.Tika.Type == "external" {
		if V.Tika.Port == "" {
			V.Tika.Port = "9998"
		}
		T.Path = V.Tika.Path + ":" + V.Tika.Port
	} else {
		return errors.New("unrecognized tika config type")
	}
	if V.ActiveFileProcessing > 0 {
		resources = make(chan struct{}, V.ActiveFileProcessing)
		for i := 0; i < V.ActiveFileProcessing; i++ {
			resources <- struct{}{}
		}
	}
	srverror.LogPath = V.LogPath
	return nil
}
