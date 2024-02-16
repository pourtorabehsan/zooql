package main

import (
	"log"
	"time"

	"github.com/Shopify/zk"
)

type listener struct {
	basePath string
	znodes   *znodeRepository
}

// OnNodeCreated implements zk.TreeCacheListener.
func (l *listener) OnNodeCreated(path string, data []byte, stat *zk.Stat) {
	go l.znodes.create(l.absolutePath(path), string(data))
}

// OnNodeDataChanged implements zk.TreeCacheListener.
func (l *listener) OnNodeDataChanged(path string, data []byte, stat *zk.Stat) {
	go l.znodes.update(l.absolutePath(path), string(data))
}

// OnNodeDeleted implements zk.TreeCacheListener.
func (l *listener) OnNodeDeleted(path string) {
	go l.znodes.delete(l.absolutePath(path))
}

// OnNodeDeleting implements zk.TreeCacheListener.
func (l *listener) OnNodeDeleting(path string, data []byte, stat *zk.Stat) {
	// noop
}

// OnSyncError implements zk.TreeCacheListener.
func (l *listener) OnSyncError(err error) {
	log.Println("sync error:", err)
}

// OnSyncStarted implements zk.TreeCacheListener.
func (l *listener) OnSyncStarted() {
	log.Println("sync started")
}

// OnSyncStopped implements zk.TreeCacheListener.
func (l *listener) OnSyncStopped(err error) {
	log.Println("sync stopped:", err)
}

// OnTreeSynced implements zk.TreeCacheListener.
func (l *listener) OnTreeSynced(elapsed time.Duration) {
	log.Println("tree synced in", elapsed)
}

func (l *listener) absolutePath(path string) string {
	if path == "/" {
		return l.basePath
	}

	return l.basePath + path
}
