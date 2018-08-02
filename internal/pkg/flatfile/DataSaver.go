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

package flatfile

import (
	"context"
	"encoding/json"
	"github.com/zlepper/welp/internal/pkg/models"
	"os"
	"time"
)

type dataSaveable interface {
	// Should put the data in readonly mode to avoid mutations during save
	// Will be invoked before data is saved
	Lock()
	// Should free the data from readonly mode
	// Will be invoked after data has been saved
	Unlock()
	// Should return the data to save
	GetData() interface{}
	// Should return true if the data has changed
	HasChanged() bool
	// Should register that that changes has been saved
	SetChanged(changed bool)
}

type DataSaverArgs struct {
	filename     string
	saveInterval time.Duration
	saveable     dataSaveable
	logger       models.Logger
}

func NewDataSaver(args DataSaverArgs) *DataSaver {
	return &DataSaver{
		filename:     args.filename,
		saveable:     args.saveable,
		logger:       args.logger,
		saveInterval: args.saveInterval,
	}
}

type DataSaver struct {
	filename     string
	saveable     dataSaveable
	logger       models.Logger
	saveInterval time.Duration
}

func (d *DataSaver) StartSaveCycle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			d.logger.Info("Shutting down file data storage")
			d.saveChanges()
			// We are completely done saving
			return
		case <-time.After(d.saveInterval):
			err := d.saveChanges()
			if err != nil {
				d.logger.Error(err)
			}
			// And just start waiting again
		}
	}
}

func (d *DataSaver) saveChanges() error {
	if !d.saveable.HasChanged() {
		return nil
	}

	d.logger.Info("Data has changed. Saving...")

	d.saveable.Lock()
	defer d.saveable.Unlock()

	tempName := d.filename + ".temp"

	// Save to temp file first to ensure no errors
	file, err := os.Create(tempName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(d.saveable.GetData())
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempName, d.filename)
	if err != nil {
		return err
	}

	// Register that there are no new changes to be saved
	d.saveable.SetChanged(false)

	return nil
}

// Loads data from the underlying data file into the given pointer
func (d *DataSaver) LoadData(data interface{}) error {
	file, err := os.Open(d.filename)
	if err != nil {
		// If the file doesn't exist, then it's simply because we haven't saved anything yet.
		if os.IsNotExist(err) {
			d.logger.Infof("Data file '%s' doesn't exist. Assuming it's because it hasn't been created yet.", d.filename)
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(data)
}
