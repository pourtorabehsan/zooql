package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/zk"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting server on port :8000")

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var zookeepersStr string
	var basePath string
	flag.StringVar(&zookeepersStr, "zookeepers", "localhost:2181", "zookeeper servers (comma separated), eg: localhost:2181,localhost:2182,localhost:2183 (default: localhost:2181)")
	flag.StringVar(&basePath, "base-path", "/", "zookeeper base path, eg: /myapp (default: /)")
	flag.Parse()

	log.Println("Starting server with zookeepers:", zookeepersStr, "and base path:", basePath)

	zookeepers := strings.Split(zookeepersStr, ",")
	log.Println("Connecting to zookeepers:", zookeepers)
	
	zkConn, _, err := zk.Connect(zookeepers, 10*time.Second, zk.WithHostProvider(zk.NewRefreshDNSHostProvider()))
	if err != nil {
		return fmt.Errorf("failed to connect to zookeeper: %w", err)
	}
	defer zkConn.Close()

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return fmt.Errorf("failed to open sqlite3 database: %w", err)
	}

	err = initDB(db)
	if err != nil {
		return fmt.Errorf("failed to initialize sqlite3 database: %w", err)
	}

	repo := &znodeRepository{db: db}

	tcl := &listener{basePath: basePath, znodes: repo}
	treeCache := zk.NewTreeCache(zkConn, basePath, zk.WithTreeCacheIncludeData(true), zk.WithTreeCacheListener(tcl), zk.WithTreeCacheReservoirLimit(65536), zk.WithTreeCacheBatchSize(1024), zk.WithTreeCacheAbsolutePaths(true))

	ctx := context.Background()
	go func() {
		_ = treeCache.Sync(ctx)
	}()

	syncCtx, syncCtxCancel := context.WithTimeout(ctx, 20*time.Second)
	defer syncCtxCancel()

	log.Println("Waiting for initial sync of zk tree cache")
	if err := treeCache.WaitForInitialSync(syncCtx); err != nil {
		log.Println("Failed to sync zk tree cache:", err)
		return err
	}
	log.Println("Initial sync of zk tree cache complete")

	err = initLoad(basePath, treeCache, repo)
	if err != nil {
		return fmt.Errorf("failed to load initial znodes: %w", err)
	}

	server := newServer(repo, zookeepersStr, basePath)
	return server.Start()
}

func initDB(db *sql.DB) error {
	createTable := "CREATE TABLE IF NOT EXISTS znodes (path TEXT PRIMARY KEY, data TEXT)"
	_, err := db.Exec(createTable)
	if err != nil {
		return err
	}

	createIndex := "CREATE INDEX IF NOT EXISTS idx_znodes_data_path ON znodes (data, path)"
	_, err = db.Exec(createIndex)
	if err != nil {
		return err
	}

	return nil
}

func initLoad(basePath string, treeCache *zk.TreeCache, znodes *znodeRepository) error {
	log.Println("Loading initial znodes")
	w := treeCache.Walker(basePath, zk.BreadthFirstOrder)
	err := w.Walk(func(path string, stat *zk.Stat) error {
		data, _, err := treeCache.Get(path)
		if err != nil {
			return fmt.Errorf("failed to get data for path %s: %w", path, err)
		}

		if err := znodes.create(path, string(data)); err != nil {
			return fmt.Errorf("failed to create znode for path %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk znodes: %w", err)
	}

	log.Println("Initial znodes loaded")
	return nil
}
