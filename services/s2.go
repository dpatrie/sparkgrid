package services

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
Should read all CSV files from a path given in a config file and for each csv entry,
push the data to S1
Should be able to scale to thousands of files if needed
*/

type S2 struct{}

const workers = 10

func (s *S2) ProcessDir(toProcess, processedDir string) (err error) {
	wg := sync.WaitGroup{}
	workChan := make(chan string, workers)

	for i := 0; i < workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			s.work(workChan, processedDir)
		}()
	}

	err = filepath.Walk(toProcess, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("failed walking path:", err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".csv") {
			return nil
		}

		workChan <- path

		return nil
	})

	close(workChan)
	wg.Wait()

	return nil
}

func (s *S2) work(c chan string, processedDir string) {
	for path := range c {
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()

		req, err := http.NewRequest("PUT", "http://localhost:5555/api/records", f)
		if err != nil {
			log.Println(err)
			continue
		}
		req.Header.Set("Content-type", "text/csv")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("service returned an http %d", resp.StatusCode)
			continue
		}

		if err := os.Rename(path, filepath.Join(processedDir, filepath.Base(path))); err != nil {
			log.Println(err)
			continue
		}
	}
	return
}
