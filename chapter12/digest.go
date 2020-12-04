package chapter12

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/KeKe-Li/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func main() {
	_, ctx, _ := log.FromContextOrNew(context.Background(), nil)
	m, err := MD5All(os.Args[1])
	if err != nil {
		log.ErrorContext(ctx, "MD5All-failed", "error", err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles(done, root)
	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)
	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path: path, sum: md5.Sum(data), err: err}:
				case <-done:
				}
				wg.Done()
			}()
			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		// walk has returned,so all calls to wg.Add are done. start a goroutine to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()
	return c, errc
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}
