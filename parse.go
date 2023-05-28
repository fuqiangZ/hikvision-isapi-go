package hikvision

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type MultipartReader struct {
	boundary  string
	bufReader *bufio.Reader
	ctx       context.Context
}

const (
	ContentT           = "Content-Type:"
	ContentL           = "Content-Length:"
	ContentDisposition = "Content-Disposition:"
	End                = "\r\n\r\n"
	TYPE_XML           = "xml"
	TYPE_JSON          = "json"
	TYPE_IMAGE         = "image"
)

func NewMultipart(ctx context.Context, rc io.Reader, boundary string) *MultipartReader {

	return &MultipartReader{ctx: ctx, bufReader: bufio.NewReaderSize(rc, 4096), boundary: boundary}
}

// 根据类型解析
func (m *MultipartReader) NextPart() (error, *Content) {
	for {
		line, err := m.bufReader.ReadSlice('\n')
		if err == io.EOF {
			return io.EOF, nil
		}
		if err != nil {
			return err, nil
		}
		//新行的分隔符
		if m.isBoundaryDelimiterLine(line) {
			//解析第一条消息，如果是
			return m.readPart()
		}
	}

}

func (m *MultipartReader) readPart() (error, *Content) {
	b := &Content{}
	b.Header = make(HeaderType)
	for {
		line, err := m.bufReader.ReadSlice('\n')
		// fmt.Println(string(line))
		if err == io.EOF {
			return io.EOF, nil
		}
		if strings.EqualFold(string(line), End) {
			fmt.Println("empty line")
			continue
		}

		if bytes.HasPrefix(line, []byte(ContentDisposition)) {
			// Content-Disposition: form-data;name="licensePlatePicture.jpg";filename="licensePlatePicture.jpg"
			s := strings.Split(string(bytes.TrimSuffix(line, []byte(End))), ":")[1]
			pattern := regexp.MustCompile(`name="(?P<name>.*?)";filename="(?P<filename>.*?)"`)
			matches := pattern.FindStringSubmatch(s)
			if len(matches) == 3 {
				b.Header["name"] = matches[1]
				b.Header["filename"] = matches[2]
			}
		}
		if bytes.HasPrefix(line, []byte(ContentT)) {
			if bytes.Contains(line, []byte(TYPE_XML)) {
				// b.ContentType = TYPE_XML
				b.Header[ContentT[:len(ContentT)-1]] = TYPE_XML
			} else if bytes.Contains(line, []byte(TYPE_JSON)) {
				// b.ContentType = TYPE_XML
				b.Header[ContentT[:len(ContentT)-1]] = TYPE_JSON
			} else if bytes.Contains(line, []byte(TYPE_IMAGE)) {
				// b.ContentType = TYPE_IMAGE
				b.Header[ContentT[:len(ContentT)-1]] = TYPE_IMAGE

			} else {
				t := strings.Split(string(line), ":")[1]
				if len(t) > 0 {
					t = strings.Split(t, ";")[0]
				}
				// b.ContentType = t
				b.Header[ContentT[:len(ContentT)-1]] = t

			}
		}
		if bytes.HasPrefix(line, []byte(ContentL)) {
			// fmt.Println(strings.Split(string(line), ":"))
			lenStr := strings.Split(string(bytes.TrimSuffix(line, []byte(End))), ":")[1]
			length, err := strconv.Atoi(strings.TrimSpace(lenStr))
			if err != nil {
				fmt.Println(err)
				return err, nil
			}
			b.Header[ContentL[:len(ContentL)-1]] = length
			//从length之后读取len个字节
			b.Body = make([]byte, length)
			var p int = 0
			for {
				//read full
				if n, err := m.bufReader.Read(b.Body[p:]); err == nil {
					p += n
				}
				if p >= length {
					b.Body = bytes.TrimLeftFunc(b.Body, func(r rune) bool {
						return r == '\r' || r == '\n' || r == ' '
					})
					return nil, b
				}
			}

		}

	}
}

// 起始行
func (m *MultipartReader) isBoundaryDelimiterLine(line []byte) bool {
	if !bytes.Contains(line, []byte(m.boundary)) {
		return false
	}
	return true
}

// 一个part的结束,需要结合length来判断
func (m *MultipartReader) isFinalPart(line []byte) bool {
	if bytes.Equal(line, []byte(End)) {
		return true
	}
	return false
}
