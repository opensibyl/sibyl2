package upload

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	sibylUploader := NewUploadCmd()
	b := bytes.NewBufferString("")
	sibylUploader.SetOut(b)
	sibylUploader.SetArgs([]string{"--src", "../../../..", "--dry"})
	err := sibylUploader.Execute()
	assert.Nil(t, err)
}
