package utils

import (
	"encoding/xml"
	"go.uber.org/zap"
	"io"
)

func DecodeXml[T any](buff io.Reader) (T, error) {
	var res T
	b, err := io.ReadAll(buff)
	if err != nil {
		zap.L().Sugar().Error(err)
		return res, err
	}
	err = xml.Unmarshal(b, &res)
	return res, err
}
