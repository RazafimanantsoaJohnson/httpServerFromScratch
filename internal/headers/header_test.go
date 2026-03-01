package headers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// test with correct value
	headers := Header{}
	data := []byte("Host: localhost:42069  \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done) //??

	// test with incorrect values
	data = []byte("     Host    : loclahost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test correct values with 2 headers
	data = []byte("Host: localhost:42069\r\naccept: application/json\r\n\r\n")
	firstn, done, err := headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, strings.Index(string(data), "\r\n"), firstn)
	assert.False(t, done)
	n, done, err = headers.Parse(data[firstn+2:])
	assert.NoError(t, err)
	assert.Equal(t, strings.Index(string(data[firstn+2:]), "\r\n"), n)
	assert.False(t, done)

	// test parse done
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 0, n)

	// test with invalid key name
	data = []byte("MoneyH|@st: 50029\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

	// test with many values for same key
	data = []byte("key1: value1\r\nkey1: value2\r\nKEy1: value3\r\n\r\n")
	firstn, done, err = headers.Parse(data)
	secondn, done, err := headers.Parse(data[firstn+2:])
	n, done, err = headers.Parse(data[firstn+secondn+4:])
	assert.NoError(t, err)
	assert.Equal(t, "value1,value2,value3", headers["key1"])
}
