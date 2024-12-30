package service

import (
	"fmt"
	"gemini-poc/app/adapter"
	"gemini-poc/app/dto"
	"gemini-poc/utils/config"
	"gemini-poc/utils/custom"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MirrorService struct {
	da   *adapter.DestinationAdapter
	pool custom.WorkerPool

	mirrorWhitelist map[string]*dto.MethodWhitelist
	pathTrie        *Trie

	conf []config.MirrorsConfig
	log  *zap.Logger
}

func NewMirrorService(
	da *adapter.DestinationAdapter,
	pool custom.WorkerPool,

	conf []config.MirrorsConfig,

	log *zap.Logger,
) *MirrorService {
	return &MirrorService{
		da:              da,
		pool:            pool,
		mirrorWhitelist: NewMirrorWhitelist(conf),
		pathTrie:        NewAndInitTrie(conf),

		conf: conf,
		log:  log,
	}
}

func NewMirrorWhitelist(conf []config.MirrorsConfig) map[string]*dto.MethodWhitelist {
	mirrorWhitelist := make(map[string]*dto.MethodWhitelist)
	for _, mirror := range conf {
		methods := dto.NewMethodWhitelist(mirror.Methods)

		for _, path := range mirror.Endpoints {
			if _, ok := mirrorWhitelist[path]; ok {
				mirrorWhitelist[path].Append(mirror.Methods)
				continue
			}

			mirrorWhitelist[path] = &methods
		}
	}

	return mirrorWhitelist
}

func NewAndInitTrie(conf []config.MirrorsConfig) *Trie {
	trie := newTrie()

	mapConf := make(map[string]bool)
	orderedConf := make([]string, 0)

	for _, mirror := range conf {
		for _, path := range mirror.Endpoints {
			if _, exists := mapConf[path]; !exists {
				mapConf[path] = true
				orderedConf = append(orderedConf, path)
			}
		}
	}

	for _, path := range orderedConf {
		trie.insert(path)
	}

	return trie
}

type TrieNode struct {
	children map[string]*TrieNode
	isEnd    bool
	isParam  bool
	paramKey string
}

func newTrieNode() *TrieNode {
	return &TrieNode{children: make(map[string]*TrieNode)}
}

type Trie struct {
	root *TrieNode
}

func newTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

func (t *Trie) insert(pattern string) {
	node := t.root
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		isParam := strings.HasPrefix(part, ":")
		key := part
		if isParam {
			key = ":"
		}
		if _, exists := node.children[key]; !exists {
			node.children[key] = newTrieNode()
		}
		node.children[key].isParam = isParam
		if isParam {
			node.children[key].paramKey = part
		}
		node = node.children[key]
	}
	node.isEnd = true
}

func (t *Trie) match(input string) (bool, string, map[string]string) {
	node := t.root
	parts := strings.Split(input, "/")
	var matchedParts []string
	var params = make(map[string]string)

	for _, part := range parts {
		// Match exact part
		if next, exists := node.children[part]; exists {
			matchedParts = append(matchedParts, part)
			node = next

			// Match path variable
		} else if next, exists := node.children[":"]; exists {
			matchedParts = append(matchedParts, next.paramKey)
			params[next.paramKey] = part
			node = next
		} else {
			return false, "", nil
		}
	}

	if node.isEnd {
		return true, strings.Join(matchedParts, "/"), params
	}
	return false, "", nil
}

func (m *MirrorService) MirrorRequest(
	path string,
	method string,
	queries map[string]string,
	requestHeaders map[string][]string,
	body []byte,
	isSuccess bool,
	responseHeaders map[string][]string,
) {
	if !isSuccess {
		return
	}

	// check if the request path and ethod is in the trie whitelist
	matched, pattern, _ := m.pathTrie.match(path)
	if !matched {
		return
	}

	// check in the whitelist map
	allowedMethods, ok := m.mirrorWhitelist[pattern]
	if !ok {
		return
	}

	// check if the method is allowed
	if method == fiber.MethodGet && !allowedMethods.GET {
		return
	}
	if method == fiber.MethodPost && !allowedMethods.POST {
		return
	}
	if method == fiber.MethodPut && !allowedMethods.PUT {
		return
	}
	if method == fiber.MethodDelete && !allowedMethods.DELETE {
		return
	}
	if method == fiber.MethodPatch && !allowedMethods.PATCH {
		return
	}
	if method == fiber.MethodOptions && !allowedMethods.OPTIONS {
		return
	}
	if method == fiber.MethodHead && !allowedMethods.HEAD {
		return
	}
	if method == fiber.MethodTrace && !allowedMethods.TRACE {
		return
	}
	if method == fiber.MethodConnect && !allowedMethods.CONNECT {
		return
	}

	traceparent := uuid.New().String()
	if traceparentHeader, ok := responseHeaders["Traceparent"]; ok {
		traceparent = traceparentHeader[0]
	}

	// add the task to the worker pool
	m.pool.AddTask(func() {
		m.log.Info(fmt.Sprintf("Mirror %s request", method), zap.String("path", path), zap.String("traceparent", traceparent))
		res, err := m.da.Do(path, method, queries, requestHeaders, body)
		if err != nil {
			m.log.Error(fmt.Sprintf("Failed to mirror %s request", method), zap.String("path", path), zap.String("traceparent", traceparent), zap.Error(err))
		}
		m.log.Info(fmt.Sprintf("Mirror %s response", method), zap.ByteString("response", res), zap.String("traceparent", traceparent))
	})

}