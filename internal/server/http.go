package server

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const maxBodySize = 1 << 20
const maxStartLineSize = 8 * 1024
const maxHeaderLineSize = 8 * 1024
const maxHeaderTotalSize = 1 << 20

var errTooLarge = errors.New("body too large")

type Req struct {
	Method string
	Path   string
	Ver    string
	Hdr    map[string]string
	Body   []byte
	Query  map[string]string
}

func readLimitedLine(br *bufio.Reader, limit int) (string, error) {
	if limit <= 0 {
		return "", errTooLarge
	}

	var (
		buf     []byte
		prefix  bool
		chunk   []byte
		readErr error
	)

	for {
		chunk, prefix, readErr = br.ReadLine()
		if readErr != nil {
			return "", readErr
		}

		if len(buf)+len(chunk) > limit {
			return "", errTooLarge
		}

		buf = append(buf, chunk...)

		if !prefix {
			break
		}
	}

	return string(buf), nil
}

func rdReq(br *bufio.Reader) (*Req, error) {
	startLine, err := readLimitedLine(br, maxStartLineSize)
	if err != nil {
		return nil, err
	}

	for startLine == "" {
		startLine, err = readLimitedLine(br, maxStartLineSize)
		if err != nil {
			return nil, err
		}
		if startLine != "" {
			break
		}
	}

	parts := strings.SplitN(startLine, " ", 3)
	if len(parts) < 3 {
		return nil, io.ErrUnexpectedEOF
	}
	method := parts[0]
	fullPath := parts[1]
	ver := parts[2]

	hdrMap := make(map[string]string, 8)

	totalHdrBytes := 0
	for {
		line, err := readLimitedLine(br, maxHeaderLineSize)
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}

		totalHdrBytes += len(line) + 2
		if totalHdrBytes > maxHeaderTotalSize {
			return nil, errTooLarge
		}

		sepIndex := strings.IndexByte(line, ':')
		if sepIndex < 0 {
			continue
		}

		hdrKey := strings.TrimSpace(line[:sepIndex])
		hdrVal := strings.TrimSpace(line[sepIndex+1:])
		if hdrKey == "" {
			continue
		}
		hdrMap[hdrKey] = hdrVal

	}

	var qp map[string]string
	pathOnly := fullPath
	if idx := strings.IndexByte(fullPath, '?'); idx >= 0 {
		raw := fullPath[idx+1:]
		pathOnly = fullPath[:idx]

		qp = make(map[string]string, 4)
		for _, part := range strings.Split(raw, "&") {
			if part == "" {
				continue
			}
			kv := strings.SplitN(part, "=", 2)
			k := kv[0]
			v := ""
			if len(kv) > 1 {
				v = kv[1]
			}
			qp[k] = v
		}
	}

	var bodyBuf []byte

	contentLenStr, hasLen := hdrMap["Content-Length"]
	if hasLen {
		size, convErr := strconv.Atoi(contentLenStr)
		if convErr != nil || size < 0 {
			return nil, io.ErrUnexpectedEOF
		}

		if size > maxBodySize {
			remain := size
			tmp := make([]byte, 4096)
			for remain > 0 {
				chunk := 4096
				if remain < chunk {
					chunk = remain
				}
				_, readErr := io.ReadFull(br, tmp[:chunk])
				if readErr != nil {
					break
				}
				remain -= chunk
			}
			return nil, errTooLarge
		}

		if size > 0 {
			bodyBuf = make([]byte, size)
			if _, err := io.ReadFull(br, bodyBuf); err != nil {
				return nil, err
			}
		}
	}

	return &Req{
		Method: method,
		Path:   pathOnly,
		Ver:    ver,
		Hdr:    hdrMap,
		Body:   bodyBuf,
		Query:  qp,
	}, nil
}
