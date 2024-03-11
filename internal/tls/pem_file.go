package tls

import (
	"bytes"
	"encoding/pem"
	"os"

	"github.com/0loff/go_link_shortener/internal/logger"
	"go.uber.org/zap"
)

// PEMFileCreate - the pem encrypted file creation method
func PEMFileCreate(filepath, cypherType string, cypher []byte) error {
	var buf bytes.Buffer
	var f *os.File

	err := pem.Encode(&buf, &pem.Block{
		Type:  cypherType,
		Bytes: cypher,
	})
	if err != nil {
		logger.Log.Error("Unable to PEM encode", zap.Error(err))
		return err
	}

	f, err = os.Create(filepath)
	if err != nil {
		logger.Log.Error("Unable to PEM file create", zap.Error(err))
		return err
	}

	_, err = buf.WriteTo(f)
	if err != nil {
		logger.Log.Error("Unable write cypher to file", zap.Error(err))
		return err
	}

	err = f.Close()
	if err != nil {
		logger.Log.Error("Error during write cypher to file", zap.Error(err))
		return err
	}

	return nil
}
