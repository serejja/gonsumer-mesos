package framework

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFileStorage(t *testing.T) {
	file := "tmp_file_storage.txt"
	contents := "hello world"
	defer func() {
		os.Remove(file)
	}()

	storage := NewFileStorage(file)
	_, err := storage.Load()
	assert.Equal(t, ErrStorageUninitialized, err)
	err = storage.Save([]byte(contents))
	require.Nil(t, err)
	loadedContents, err := storage.Load()
	require.Nil(t, err)
	assert.Equal(t, contents, string(loadedContents))
	assert.Equal(t, file, storage.String())
}

func TestZKStorage(t *testing.T) {
	zkConnect := "localhost:2181"
	zpath := "/tmp/zk/storage"
	contents := "hello world"

	conn, _, err := zk.Connect([]string{zkConnect}, 30*time.Second)
	_, _, err = conn.Exists("/tmp") // check if zk is alive
	if err != nil {
		t.Skipf("localhost:2181 is not responding (error %s). To run this test please spin up ZK on localhost:2181", err)
	}
	defer conn.Close()

	storage, err := NewZKStorage(fmt.Sprintf("%s%s", zkConnect, zpath))
	require.Nil(t, err)

	_, err = storage.Load()
	assert.Equal(t, ErrStorageUninitialized, err)

	err = storage.Save([]byte(contents))
	require.Nil(t, err)

	loadedContents, err := storage.Load()
	require.Nil(t, err)
	assert.Equal(t, contents, string(loadedContents))

	err = zkDelete(conn, zpath)
	require.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s%s", zkConnect, zpath), storage.String())
}

func zkDelete(conn *zk.Conn, zpath string) error {
	if zpath != "" {
		_, stat, _ := conn.Get(zpath)
		err := conn.Delete(zpath, stat.Version)
		if err != nil {
			return err
		}

		index := strings.LastIndex(zpath, "/")
		if index != -1 {
			return zkDelete(conn, zpath[:index])
		} else {
			return nil
		}
	} else {
		return nil
	}
}
