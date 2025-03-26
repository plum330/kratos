package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/plum330/kratos/v2/encoding"
	"github.com/plum330/kratos/v2/errors"
	"github.com/plum330/kratos/v2/internal/httputil"
	"github.com/plum330/kratos/v2/transport/http/binding"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

type requestVarType string

const requestVars requestVarType = "request-vars"

// Redirector replies to the request with a redirect to url
// which may be a path relative to the request path.
type Redirector interface {
	Redirect() (string, int)
}

// Request type net/http.
type Request = http.Request

// ResponseWriter type net/http.
type ResponseWriter = http.ResponseWriter

// Flusher type net/http
type Flusher = http.Flusher

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, any) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, any) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// DefaultRequestVars decodes the request vars to object.
func DefaultRequestVars(r *http.Request, v any) error {
	params, ok := r.Context().Value(requestVars).(map[string]string)
	if !ok {
		return nil
	}
	vars := make(url.Values, len(params))
	for key, value := range params {
		vars[key] = []string{value}
	}
	return binding.BindQuery(vars, v)
}

// DefaultRequestQuery decodes the request vars to object.
func DefaultRequestQuery(r *http.Request, v any) error {
	return binding.BindQuery(r.URL.Query(), v)
}

// DefaultRequestDecoder decodes the request body to object.
func DefaultRequestDecoder(r *http.Request, v any) error {
	codec, ok := CodecForRequest(r, "Content-Type")
	if !ok {
		return errors.BadRequest("CODEC", fmt.Sprintf("unregister Content-Type: %s", r.Header.Get("Content-Type")))
	}
	data, err := io.ReadAll(r.Body)

	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}

type Head struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Response struct {
	Head
	Data any `json:"data"`
}

// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(Redirector); ok {
		u, code := rd.Redirect()
		http.Redirect(w, r, u, code)
		return nil
	}

	rsp := &Response{
		Head: Head{
			Msg:  "�ɹ�",
			Code: http.StatusOK,
		},
		Data: v,
	}
	codec, _ := CodecForRequest(r, "Accept")
	hb, err := codec.Marshal(rsp.Head)
	if err != nil {
		return err
	}
	pb, err := codec.Marshal(rsp.Data)
	if err != nil {
		return err
	}
	data := make([]byte, 0, len(hb)+len(pb)+8)
	data = append(data, hb[:len(hb)-1]...)
	data = append(data, []byte(`,"data":`)...)
	data = append(data, pb...)
	data = append(data, '}')
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	_, err = w.Write(data)
	return err
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	rsp := &Response{
		Head: Head{
			Code: int(se.Code),
			Msg:  se.Message,
		},
	}
	codec, _ := CodecForRequest(r, "Accept")
	body, err := codec.Marshal(rsp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(httputil.ContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec("json"), false
}
