package gxjpath

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type GXJValueType int8
type GXJContainerType int8

var ErrNotFound = errors.New("Path not found")
var ErrWrongPath = errors.New("Wrong path")
var ErrInvalidIndex = errors.New("Invalid array index")

const (
	GXJMapContainer GXJContainerType = iota
	GXJSliceContainer
)

type GXJPathItem struct {
	key           string
	index         int
	containerType GXJContainerType
}

func (s *GXJPathItem) String() string {
	var t string
	var c string

	switch s.containerType {
	case GXJMapContainer:
		c = "map[string]" + t
	case GXJSliceContainer:
		c = "[]" + t
	}
	return fmt.Sprintf("Key: %s, Index: %d, container: %s", s.key, s.index, t, c)
}

type GXJPath []GXJPathItem

var cacheLock sync.Mutex
var compiledCache map[string]GXJPath = make(map[string]GXJPath)

func (x GXJPath) String() string {
	buf := make([]string, 0)
	for _, v := range x {
		buf = append(buf, v.String())
	}
	return strings.Join(buf, "\n")
}

func getIndex(path string) (int, error) {
	if len(path) == 0 {
		return 0, ErrInvalidIndex
	}

	switch path {
	case "first":
		return 0, nil
	case "last":
		return -1, nil
	default:
		idx, err := strconv.Atoi(path)
		if err != nil {
			return idx, ErrInvalidIndex
		}
		return idx, nil
	}
}

// cutSegment cuts a path segment taking into account escaped "." characters.
func cutSegment(path string) (string, string) {
	normChar := true
	var idx int

	for idx = 0; idx < len(path); idx++ {
		if normChar {
			if path[idx] == '\\' {
				normChar = false
			}
			if path[idx] == '.' {
				leftover := path[idx+1:]
				if len(leftover) == 0 {
					// Will cause an error later. Used on purpose.
					leftover = "."
				}
				return path[:idx], leftover
			}
		} else {
			normChar = true
		}
	}
	return path, ""
}

// unescape removed escape characters.
func unescape(segment string) string {
	buf := make([]byte, 0, len(segment))
	escChar := false
	for idx := 0; idx < len(segment); idx++ {
		b := segment[idx]
		if escChar {
			escChar = false
			buf = append(buf, b)
			continue
		}
		if segment[idx] == '\\' {
			escChar = true
			continue
		}
		buf = append(buf, b)
	}
	if len(segment) == len(buf) {
		return segment
	}
	return string(buf)
}

// compileSegment compiles a single segment.
func compileSegment(segment string) (GXJPathItem, error) {
	var idx int
	var err error
	var containerType GXJContainerType = GXJMapContainer

	if len(segment) == 0 {
		return GXJPathItem{}, ErrWrongPath
	}
	if segment[0] == '@' {
		containerType = GXJSliceContainer
		segment = segment[1:]
		idx, err = getIndex(segment)
		if err != nil {
			return GXJPathItem{}, err
		}
	}

	return GXJPathItem{
		index:         idx,
		key:           unescape(segment),
		containerType: containerType,
	}, nil
}

// CompileXJPath creates a slice of parsed path segments
func CompilePath(path string) (GXJPath, error) {
	output := make(GXJPath, 0, 4)
	var segment string
	for len(path) > 0 {
		segment, path = cutSegment(path)
		if len(segment) == 0 {
			return nil, ErrWrongPath
		}
		compiled, err := compileSegment(segment)
		if err != nil {
			return nil, err
		}
		output = append(output, compiled)
	}
	return output, nil
}

func lookupSegment(pathItem *GXJPathItem, data interface{}) (interface{}, error) {
	if pathItem.containerType == GXJSliceContainer {
		v, ok := data.([]interface{})
		if !ok {
			return nil, ErrNotFound
		}
		idx := pathItem.index
		if idx < 0 {
			idx = len(v) + idx
		}
		if idx < 0 && idx >= len(v) {
			return nil, ErrNotFound
		}
		return v[idx], nil
	}
	if pathItem.containerType == GXJMapContainer {
		v, ok := data.(map[string]interface{})
		if !ok {
			return nil, ErrNotFound
		}

		mapValue, ok := v[pathItem.key]
		if !ok {
			return nil, ErrNotFound
		}
		return mapValue, nil
	}
	panic("Unknown container type!")
	return nil, nil
}

// LookupCompiledPath looks up a value using pre-compiled path that works much faster.
func LookupCompiledPath(path GXJPath, data interface{}) (interface{}, error) {
	var err error
	for _, x := range path {
		data, err = lookupSegment(&x, data)
		if err != nil {
			return nil, err
		}
	}
	return data, err
}

// LookupRawPath looks up a value using raw string path.
// Before actual lookup value is compiled, so if path used frequently,
// it would be a good idea to compiled it first.
func LookupRawPath(path string, data interface{}) (interface{}, error) {
	cp, err := CompilePath(path)
	if err != nil {
		return nil, err
	}
	return LookupCompiledPath(cp, data)
}

// CachedLookup compiles and caches compiled path in the module level dictionary.
// Lookup result is equal to the previous functions.
func CachedLookup(path string, data interface{}) (interface{}, error) {
	cacheLock.Lock()
	compiled, ok := compiledCache[path]
	cacheLock.Unlock()
	if !ok {
		cacheLock.Lock()
		defer cacheLock.Unlock()
		c, err := CompilePath(path)
		if err != nil {
			return nil, err
		}
		compiledCache[path] = c
		compiled = c
	}
	return LookupCompiledPath(compiled, data)
}
