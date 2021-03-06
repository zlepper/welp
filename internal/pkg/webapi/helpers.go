/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package webapi

import (
	"bytes"
	"context"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"strings"
)

const (
	HeaderAccept       = echo.HeaderAccept
	HeaderCacheControl = "Cache-Control"
	MIMEHTML           = "text/html"
	MIMEJSON           = "application/json"
	MIMEXML            = "application/xml"
)

// Attempts to normalize the accept header
// If no type is known, return MIMEHTML
func GetResponseType(r *http.Request) string {
	accept := r.Header.Get(HeaderAccept)

	options := strings.Split(accept, ",")

	for _, option := range options {
		switch {
		case strings.HasPrefix(option, echo.MIMEApplicationJSON):
			return MIMEJSON
		case strings.HasPrefix(option, echo.MIMETextHTML):
			return MIMEHTML
		case strings.HasPrefix(option, echo.MIMEApplicationXML), strings.HasPrefix(option, echo.MIMETextXML):
			return MIMEXML
		}
	}

	return MIMEHTML
}

func GetContext(r *http.Request) context.Context {
	return r.Context()
}

// Attempts to detect the content type of the response
// As this function consumes the first 512 bytes of the response
// a new reader is returned, that also contains the start bytes
func DetectContentType(r io.Reader) (string, io.Reader, error) {

	data := make([]byte, 512)
	read, err := r.Read(data)
	if err != nil && err != io.EOF {
		return "", nil, err
	}

	data = data[0:read]

	contentType := http.DetectContentType(data)

	// Get a reader that contains the initial content too
	newReader := io.MultiReader(bytes.NewReader(data), r)

	return contentType, newReader, nil
}
