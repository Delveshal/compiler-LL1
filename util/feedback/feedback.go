package feedback

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type feedBack struct {
	DistWriter http.ResponseWriter `json:"-"`
	FbCode     int                 `json:"-"`
	FbMsg      string              `json:"msg,omitempty"`
	FbData     interface{}         `json:"data,omitempty"`
}

type FbBuilder interface {
	Dist(w http.ResponseWriter) FbBuilder
	Code(code int) FbBuilder
	Msg(msg string) FbBuilder

	Data(data interface{}) FbBuilder
	Response() error
	Clear()
}

func NewFeedBack(w http.ResponseWriter) FbBuilder {
	return &feedBack{DistWriter: w}
}

func (f *feedBack) Dist(w http.ResponseWriter) FbBuilder {
	f.DistWriter = w
	return f
}

func (f *feedBack) Code(code int) FbBuilder {
	f.FbCode = code
	return f
}

func (f *feedBack) Msg(msg string) FbBuilder {
	f.FbMsg = msg
	return f
}

func (f *feedBack) Data(data interface{}) FbBuilder {
	f.FbData = data
	return f
}

func (f *feedBack) Response() (err error) {
	if f.DistWriter == nil {
		return errors.New("DistWriter is empty")
	}
	if f.FbCode == 0{
		f.DistWriter.WriteHeader(http.StatusOK)
	}else{
		f.DistWriter.WriteHeader(f.FbCode)
	}
	buf, _ := json.Marshal(f)
	fmt.Fprint(f.DistWriter, string(buf))
	f.Clear()
	return nil
}

func (f *feedBack) Clear() {
	f.FbData = nil
	f.FbMsg = ""
	f.FbCode = 0
}
