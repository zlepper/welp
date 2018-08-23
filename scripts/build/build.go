/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/zlepper/welp/internal/pkg/consts"
	"log"
	"os"
	"os/exec"
	"sync"
)

type configuration struct {
	OS        string
	Arch      string
	Extension string
}

var (
	configurations = []configuration{
		{
			OS:        "windows",
			Arch:      "386",
			Extension: "windows-x32.exe",
		},
		{
			OS:        "windows",
			Arch:      "amd64",
			Extension: "windows-x64.exe",
		},
		{
			OS:        "darwin",
			Arch:      "386",
			Extension: "osx-x32",
		},
		{
			OS:        "darwin",
			Arch:      "amd64",
			Extension: "osx-x64",
		},
		{
			OS:        "linux",
			Arch:      "386",
			Extension: "linux-x32",
		},
		{
			OS:        "linux",
			Arch:      "amd64",
			Extension: "linux-x64",
		},
	}
)

func main() {
	log.Println("Starting build")
	goBinary, err := exec.LookPath("go")
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(configurations))
	for _, conf := range configurations {
		conf := conf
		go func() {
			defer wg.Done()
			log.Printf("building binary for '%s'\n", conf.Extension)
			cmd := exec.Cmd{
				Path: goBinary,
				Args: []string{
					goBinary,
					"build",
					"-o",
					fmt.Sprintf("build/welp-%s-%s", consts.Version, conf.Extension),
					"github.com/zlepper/welp",
				},
				Env: append(
					os.Environ(),
					fmt.Sprintf("GOOS=%s", conf.OS),
					fmt.Sprintf("GOARCH=%s", conf.Arch),
				),
			}
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Println("build args", cmd.Args)
				log.Fatalln("Error when building", err, "\n", string(output))
				return
			}
			log.Printf("Successfully build binary for '%s'\n", conf.Extension)
		}()
	}
	wg.Wait()
	log.Println("Finished building all configurations.")
}
