package search

import (
	"cloud/forever"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestInitRiot(t *testing.T) {
	RiotRegister()
	InitRiot_Test()
	AndSearch("IMAX")

	logrus.Info(searcher.NumIndexed())
	AddDoc("test_kind", "战争IMAX")
	logrus.Info(searcher.NumIndexed())
	AndSearch("IMAX")
	AddMultiDocs("test_kind", "IMAX1", "IMAX2", "IMAX3")
	AndSearch("IMAX")

	logrus.Info(searcher.NumIndexed())
	RiotUnregister()
}

func TestOrSearch(t *testing.T) {
	RiotRegister()
	InitRiot()

	AddDoc("test_kind", "战争IMAX")
	AddMultiDocs("test_kind", "IMAX1", "IMAX2", "IMAX3")

	logrus.Info(searcher.NumIndexed())
	OrSearch("战争影院")
	RiotUnregister()
}

func TestInitRiot1(t *testing.T) {
	forever.RedisRegister()
	RiotRegister()
	InitRiot()
	OrSearch("t")
	RiotUnregister()
	forever.RedisUnRegister()
}
