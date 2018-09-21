package hou_test

import (
	"github.com/ethereum/go-ethereum/log"
	"testing"
)

func TestLenMap(t *testing.T) {
	m := make(map[string]string)
	m["1"] = "1"
	log.Info("123", len(m))
}
