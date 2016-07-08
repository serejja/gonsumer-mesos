package framework

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Storage interface {
	Save([]byte) error
	Load() ([]byte, error)
}

type FileStorage struct {
	file string
}

func NewFileStorage(file string) *FileStorage {
	return &FileStorage{
		file: file,
	}
}

func (fs *FileStorage) Save(contents []byte) error {
	return ioutil.WriteFile(fs.file, contents, 0644)
}

func (fs *FileStorage) Load() ([]byte, error) {
	data, err := ioutil.ReadFile(fs.file)
	if os.IsNotExist(err) {
		return nil, ErrStorageUninitialized
	}

	return data, err
}

func (fs *FileStorage) String() string {
	return fmt.Sprintf("%s", fs.file)
}

type ZKStorage struct {
	zkConnect string
	zPath     string
}

func NewZKStorage(zk string) (*ZKStorage, error) {
	zkConnect := zk
	path := "/"
	chrootIdx := strings.Index(zk, "/")
	if chrootIdx != -1 {
		zkConnect = zk[:chrootIdx]
		path = zk[chrootIdx:]
	}

	storage := &ZKStorage{
		zkConnect: zkConnect,
		zPath:     path,
	}

	err := storage.createChrootIfRequired()
	return storage, err
}

func (zs *ZKStorage) Save(contents []byte) error {
	conn, err := zs.newZkClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, stat, _ := conn.Get(zs.zPath)
	_, err = conn.Set(zs.zPath, contents, stat.Version)
	return err
}

func (zs *ZKStorage) Load() ([]byte, error) {
	conn, err := zs.newZkClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	contents, _, err := conn.Get(zs.zPath)
	if len(contents) == 0 {
		return nil, ErrStorageUninitialized
	}

	return contents, err
}

func (zs *ZKStorage) String() string {
	return fmt.Sprintf("%s%s", zs.zkConnect, zs.zPath)
}

func (zs *ZKStorage) createChrootIfRequired() error {
	if zs.zPath != "" {
		conn, err := zs.newZkClient()
		if err != nil {
			return err
		}
		defer conn.Close()

		err = zs.createZPath(conn, zs.zPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (zs *ZKStorage) newZkClient() (*zk.Conn, error) {
	conn, _, err := zk.Connect([]string{zs.zkConnect}, 30*time.Second)
	return conn, err
}

func (zs *ZKStorage) createZPath(conn *zk.Conn, zpath string) error {
	_, err := conn.Create(zpath, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		if zk.ErrNodeExists == err {
			return nil
		} else {
			parent, _ := path.Split(zpath)
			if len(parent) == 0 {
				return ErrEmptyZPath
			}
			err = zs.createZPath(conn, parent[:len(parent)-1])
			if err != nil {
				return err
			}

			_, err = conn.Create(zpath, nil, 0, zk.WorldACL(zk.PermAll))
			if err == zk.ErrNodeExists {
				err = nil
			}
		}
	}

	if zk.ErrNodeExists == err {
		return nil
	} else {
		return err
	}
}

func NewStorage(storage string) (Storage, error) {
	storageTokens := strings.SplitN(storage, ":", 2)
	if len(storageTokens) != 2 {
		return nil, ErrUnsupportedStorage
	}

	switch storageTokens[0] {
	case "file":
		return NewFileStorage(storageTokens[1]), nil
	case "zk":
		return NewZKStorage(storageTokens[1])
	default:
		return nil, ErrUnsupportedStorage
	}
}
