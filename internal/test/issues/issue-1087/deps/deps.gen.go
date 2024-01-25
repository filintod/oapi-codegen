// Package deps provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0-00010101000000-000000000000 DO NOT EDIT.
package deps

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// BaseError defines model for BaseError.
type BaseError struct {
	// Code The underlying http status code
	Code int32 `json:"code"`

	// Domain The domain where the error is originating from as defined by the service
	Domain string `json:"domain"`

	// Message A simple message in english describing the error and can be returned to the consumer
	Message string `json:"message"`

	// Metadata Any additional details to be conveyed as determined by the service. If present, will return map of key value pairs
	Metadata *map[string]string `json:"metadata,omitempty"`
}

// Error defines model for Error.
type Error = BaseError

// N401 defines model for 401.
type N401 = Error

// N403 defines model for 403.
type N403 = Error

// N410 defines model for 410.
type N410 = Error

// DefaultError defines model for DefaultError.
type DefaultError = Error

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7SV34/jNBDH/5WR4TFq0+siobwVDtAiAacDnk6rlRtPmjmcsc+edLes+r8jO2nTH6tb",
	"7qR9ahPPeL7zmR95UrXrvGNkiap6UgGjdxwxP9yUi/RTOxZkSX+195ZqLeR4/jE6Tu/wUXfe4mBpUFU3",
	"5aJQxnWaWFVq1UurCtVhjHqDqlK3vNWWDOheWmQZr1OFCqjzlQeL+9W5xb5QsW6x0ynUtwEbValv5pP+",
	"+XAa5z+F4ILa7wsl+Chzb7OSpxPvo2b164MARcBHTwGNKpTsfHofJRBv1D7dYjDWgXwWUam/OSl3gf5F",
	"M4NbNkkfRpBWC/jgtmTQQB3QJO3axnQ/O4GcVMriplx+Hdfl57iu6hpjhLfIlBM54hwO7seDV6F4GftF",
	"iD+7sCZjkK8IBvzUYxSoNSdoa4QJd4a3KL8K3qI8gXcG7hfHeMorP78KpjHSi3TeY+2CAXZgHW8wgN5q",
	"snpts6632OjeyhD4ZRRflsaVljEaYDKAw36AxgXQvBteRyAGaRFW724zivHWFPQHHfEo1QfnMQgN+2Wo",
	"zNNFwL9ahJ4NBrsj3kAr4iGKlj5CdigmoN+VZaEaFzotqlLEsnwz4SUW3GBIxA51fy7UcAYPLQbMOQyJ",
	"UgQXaEOsJalogutARzDYEKOB9S7bRgxbqs80qesCnzTbpYIVREp+MFokkMgbS7GFwXKdwk+6NJs0Gmku",
	"Akofkhhx2aB2HPsOw5ma1QZPRsmmKZVWMyy+f16naKMld4s2hpJKbd+dVe3K6SIj3sHkCgZFk41J4zpL",
	"3OIOzYBSMHTP0JzBbQM+YESWAh7I2jFV6LQH18A/uEvLtEfwmkI8zffYYrvfdZdknj6moqb1kjf9/pi+",
	"W3/EWnLfHk+rD2pstrF3phreXTkW6tjg2to/GlV9+PywTTOxvysuhuKwhq47ZTjJQwDRY00N1Yfaj+hO",
	"26OPQ2tQ/g41A2J81HX64MUeZ/Bn63prsi3Tpx7hgaQlBg3HpKdGep+j3/84ULlcYf8P3XHJXjNMVxA3",
	"LncYSQ75mzNoU3m3GOJA4c2snJWJuPPI2pOq1HJWzpaqUF5Lmwim9YMhueQ69MGqSs3V/m7/XwAAAP//",
	"TAbhS+0IAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
