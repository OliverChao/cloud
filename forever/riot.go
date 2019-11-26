package forever

import (
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/sirupsen/logrus"
)

var searcher = &riot.Engine{}

func RiotRegister() {
	searcher.Init(types.EngineOpts{
		Using:         3,
		StopTokenFile: "search/stop_tokens.txt",
		// IDOnly:        true,
		//GseDict: "./dictionary.txt",
	})
}

func RiotUnregister() {
	searcher.Close()
	logrus.Info("search engine")
}

func GetGlobalSearchEngine() *riot.Engine {
	return searcher
}
