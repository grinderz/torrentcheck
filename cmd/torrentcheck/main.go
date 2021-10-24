package main

import (
	"encoding/json"
	"fmt"
	"github.com/grinderz/go-torrent-parser"
	"github.com/grinderz/grgo/logging"
	"github.com/grinderz/grgo/osutils"
	"github.com/grinderz/torrentcheck"
	"github.com/grinderz/torrentcheck/internal/config"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Config    string `long:"config" default:"/boot/config/config.user/torrentcheck/config.yml" required:"true" description:"config file"`
	LogLevel  string `long:"log-level" default:"info" description:"logs level"`
	LogCaller bool   `long:"log-caller" description:"report log caller"`
}

func torrentCheck(cfg *config.Config) error {
	var torrentsCnt int
	rootTorrents := make(map[string]interface{})

	files, err := ioutil.ReadDir(cfg.SessionDir)
	if err != nil {
		return fmt.Errorf("read session dir: %v", err)
	}

	var size int64
	var found bool
	var foundSubdir string

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".torrent") {
			continue
		}
		inFile := filepath.Join(cfg.SessionDir, file.Name())
		torrent, err := gotorrentparser.ParseFromFile(inFile)
		if err != nil {
			return fmt.Errorf("parse file [%s]: %v", inFile, err)
		}

		torrentsCnt++
		for _, t := range torrent.Files {
			rootTorrents[t.Path[0]] = nil
			pth := filepath.Join(t.Path...)

			found = false
			for _, sd := range cfg.TorrentsSubdirs {
				fPath := filepath.Join(cfg.RootTorrentsDir, sd, pth)
				isExists, err := osutils.IsExists(fPath)
				if err != nil {
					return fmt.Errorf("open path [%s]: %v", fPath, err)
				}
				if isExists {
					fi, err := os.Stat(fPath)
					if err != nil {
						return fmt.Errorf("stat path [%s]: %v", fPath, err)
					}
					size = fi.Size()
					found = true
					foundSubdir = sd
					break
				}
			}
			if !found {
				return fmt.Errorf("not found %s [ %s ]", pth, file.Name())
			}
			if t.Length != size {
				return fmt.Errorf("size mismatch %s", filepath.Join(foundSubdir, pth))
			}
		}
	}

	logging.Log.Infof("torrents count: %d\n", torrentsCnt)
	logging.Log.Infof("root torrents: %d\n", len(rootTorrents))

	var inTorrentCheckCnt int
	for _, sd := range cfg.TorrentsSubdirs {
		fPath := filepath.Join(cfg.RootTorrentsDir, sd)
		files, err := os.ReadDir(fPath)
		if err != nil {
			return fmt.Errorf("read subdir error [%s]: %v", sd, err)
		}
		for _, f := range files {
			inTorrentCheckCnt++
			_, exist := rootTorrents[f.Name()]
			if !exist {
				return fmt.Errorf("path not exists in torrent %s", filepath.Join(sd, f.Name()))
			}
		}
	}
	logging.Log.Infof("in torrent check count: %d\n", torrentsCnt)
	return nil
}

func main() {
	var opts Options
	var parser = flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	logging.Configure(opts.LogLevel, opts.LogCaller)

	logging.Log.Infof("torrentcheck-v%s [ %s ]", torrentcheck.VersionAndBuild, torrentcheck.BuildTimestamp)

	optionsBytes, err := json.MarshalIndent(opts, "", "  ")
	if err != nil {
		logging.Log.Fatalf("marshal options error: %v\n", err)
	}
	logging.Log.Infof("options:\n%s", string(optionsBytes))

	cfg, err := config.Load(opts.Config)
	if err != nil {
		logging.Log.Fatalf("%v\n", err)
	}

	cfgBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logging.Log.Fatalf("marshal config error: %v\n", err)
	}
	logging.Log.Infof("config:\n%s", string(cfgBytes))

	if err := torrentCheck(cfg); err != nil {
		logging.Log.Fatalf("torrent check %v\n", err)
	}

}
