package telemetry

import (
	"io"
	"sync"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"
)

type ContractEventCollector struct {
	baseCollector
	filePath string
	nodes    map[string]*WrappedContractCollector
}

func NewContractEventCollector(path string) *ContractEventCollector {
	return &ContractEventCollector{
		baseCollector: baseCollector{
			t:        NodeLogType,
			io:       []io.WriteCloser{},
			ioLookup: make(map[string]int),
		},
		filePath: path,
		nodes:    make(map[string]*WrappedContractCollector),
	}
}

func (c *ContractEventCollector) ContractEventCollectorNode(node string) *WrappedContractCollector {
	wc, ok := c.nodes[node]
	if !ok {
		panic("node not available")
	}

	return wc
}

func (c *ContractEventCollector) Data() (map[string]int, map[string][]string) {
	// keyChecks is per id per block
	allKeyChecks := make(map[string]int)
	// idLookup is id checked in blocks
	allKeyIDLookup := make(map[string][]string)

	for _, node := range c.nodes {
		for key, value := range node.keyChecks {
			v, ok := allKeyChecks[key]
			if !ok {
				allKeyChecks[key] = value
			} else {
				allKeyChecks[key] += v
			}
		}

		for key, lookup := range node.keyIDLookup {
			v, ok := allKeyIDLookup[key]
			if !ok {
				allKeyIDLookup[key] = lookup
			} else {
				for _, ls := range lookup {
					found := false
					for _, ex := range v {
						if ls == ex {
							found = true
							break
						}
					}

					if !found {
						v = append(v, ls)
					}
				}
				allKeyIDLookup[key] = v
			}
		}
	}

	return allKeyChecks, allKeyIDLookup
}

func (c *ContractEventCollector) AddNode(node string) error {
	wc := &WrappedContractCollector{
		keyChecks:   make(map[string]int),
		keyIDLookup: make(map[string][]string),
	}

	c.nodes[node] = wc

	return nil
}

type WrappedContractCollector struct {
	mu          sync.Mutex
	keyChecks   map[string]int
	keyIDLookup map[string][]string
}

func (wc *WrappedContractCollector) CheckKey(key types.UpkeepKey) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	k := key.String()

	blockKey, upkeepID, _ := key.BlockKeyAndUpkeepID()

	_, ok := wc.keyChecks[k]
	if !ok {
		wc.keyChecks[k] = 0
	}
	wc.keyChecks[k]++

	val, ok := wc.keyIDLookup[string(upkeepID)]
	if !ok {
		wc.keyIDLookup[string(upkeepID)] = []string{blockKey.String()}
	} else {
		var found bool
		for _, v := range val {
			if v == blockKey.String() {
				found = true
			}
		}

		if !found {
			wc.keyIDLookup[string(upkeepID)] = append(val, blockKey.String())
		}
	}
}
