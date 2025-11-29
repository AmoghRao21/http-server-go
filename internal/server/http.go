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

	return &Req{
		Method: lineParts[0],
		Path:   lineParts[1],
		Ver:    lineParts[2],
		Hdr:    hdrMap,
		Body:   bodyBuf,
	}, nil
}
