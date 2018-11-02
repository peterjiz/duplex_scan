// Copyright © 2017 Peter El Jiz <peter.eljiz@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

/*
folder batch 1:
get today's date:
set i = 0
loop over files
rename file YYYY-MM-DD-i,
i += 2
exit loop

folder batch 1:
get today's date:
get list of files
set i = 1
reverse iterate through list of files
rename file YYYY-MM-DD-i,
i += 2
exit loop
 */


// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Runs the duplex scan postprocessor",
	Long: `Run it by typing:
go run main.go process -f /path/to/batch1 -s /path/to/batch2`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("processing batches")
		processBatches()
	},
}

var prefix string
var batch1Directory string
var batch2Directory string

func processBatches() {
	var wg sync.WaitGroup

	wg.Add(2)
	go processEvenBatch(&wg)
	go processOddBatch(&wg)
	wg.Wait()

}

func processEvenBatch(wg *sync.WaitGroup){

	defer wg.Done()

	//Batch 1 - Even
	batch1DirInfo, err := ioutil.ReadDir(batch1Directory)
	if err != nil {
		fmt.Errorf("Batch1 directory invalid")
		return
	}

	fileIndex := 0
	for _, file := range batch1DirInfo {
		if strings.Contains(file.Name(), "DS_Store") {
			continue
		}

		oldfilepath := fmt.Sprintf("%v/%v", batch1Directory, file.Name())
		newfilename := fmt.Sprintf("%v/%v-%04v.%v", batch1Directory, prefix, fileIndex, filepath.Ext(oldfilepath))
		fmt.Printf("%v\n", newfilename)
		os.Rename(oldfilepath, newfilename)
		fileIndex += 2
	}
}

func processOddBatch(wg *sync.WaitGroup){

	defer wg.Done()

	//Batch 2 - Odd
	batch2DirInfo, err := ioutil.ReadDir(batch2Directory)
	if err != nil {
		fmt.Errorf("Batch2 directory invalid")
		return
	}

	fileIndex := 1
	for i := len(batch2DirInfo)-1; i >= 0; i-- {
		file := batch2DirInfo[i]

		if strings.Contains(file.Name(), "DS_Store") {
			continue
		}

		oldfilepath := fmt.Sprintf("%v/%v", batch2Directory, file.Name())
		newfilename := fmt.Sprintf("%v/%v-%04v.%v", batch2Directory, prefix, fileIndex, filepath.Ext(oldfilepath))
		fmt.Printf("%v\n", newfilename)
		os.Rename(oldfilepath, newfilename)
		fileIndex += 2
	}
}


func init() {
	RootCmd.AddCommand(processCmd)
	processCmd.Flags().StringVarP(&prefix, "prefix", "p", time.Now().Local().Format("2006-01-02"), "Filenames Prefix")
	processCmd.Flags().StringVarP(&batch1Directory, "batch1", "f", "./batch1", "First Batch")
	processCmd.Flags().StringVarP(&batch2Directory, "batch2", "s", "./batch2", "Second Batch")
}
