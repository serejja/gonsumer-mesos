package framework

import "errors"

var ErrEmptyZPath = errors.New("Specified blank path")

var ErrUnsupportedStorage = errors.New("Unsupported storage")

var ErrStorageUninitialized = errors.New("Storage is uninitialized")
