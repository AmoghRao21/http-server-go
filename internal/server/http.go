package server

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Req struct {
	Method string
	Path   string
	Ver    string
	Hdr    map[string]string
	Body   []byte
	Query  map[string]string
}

func rdReq(r io.Reader) (*Req, error) {
	reader := bufio.NewReader(r)

	reqLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	reqLine = strings.TrimSpace(reqLine)

	lineParts := strings.SplitN(reqLine, " ", 3)
	if len(lineParts) < 3 {
		return nil, io.ErrUnexpectedEOF
	}

	hdrMap := map[string]string{}
	for {
		rawHdr, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		rawHdr = strings.TrimSpace(rawHdr)
		if rawHdr == "" {
			break
		}

		sepIndex := strings.IndexByte(rawHdr, ':')
		if sepIndex < 0 {
			continue
		}

		hdrKey := strings.TrimSpace(rawHdr[:sepIndex])
		hdrVal := strings.TrimSpace(rawHdr[sepIndex+1:])
		hdrMap[strings.ToLower(hdrKey)] = hdrVal
	}

	var bodyBuf []byte
	contentLenStr, hasLen := hdrMap["content-length"]
	if hasLen {
		size, _ := strconv.Atoi(contentLenStr)
		if size > 0 {
			bodyBuf = make([]byte, size)
			io.ReadFull(reader, bodyBuf)
		}
	}

	qp := map[string]string{}
	fullPath := lineParts[1]

	pathOnly := fullPath
	if idx := strings.IndexByte(fullPath, '?'); idx >= 0 {
		raw := fullPath[idx+1:]
		pathOnly = fullPath[:idx]
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

	return &Req{
		Method: lineParts[0],
		Path:   pathOnly,
		Ver:    lineParts[2],
		Hdr:    hdrMap,
		Body:   bodyBuf,
		Query:  qp,
	}, nil

}
