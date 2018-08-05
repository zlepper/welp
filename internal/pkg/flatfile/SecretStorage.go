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
	"crypto/rand"
	"encoding/json"
	"github.com/zlepper/welp/internal/pkg/models"
	"os"
	"path"
)

type SecretStorageArgs struct {
	Filename string
	Logger   models.Logger
}

func NewSecretStorage(args SecretStorageArgs) (models.SecretService, error) {
	storage := &secretStorage{
		SecretStorageArgs: args,
	}

	err := storage.prepare()
	if err != nil {
		return nil, err
	}

	return storage, err

}

type secretStorage struct {
	SecretStorageArgs
	secretData
}

type secretData struct {
	signingSecret []byte
}

func (s *secretStorage) prepare() error {

	// Ensure the directory exists
	dir := path.Dir(s.Filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Open(s.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return s.generate()
		}
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&s.secretData)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *secretStorage) generate() error {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return err
	}

	s.signingSecret = secret

	file, err := os.Create(s.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(s.secretData)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *secretStorage) GetSigningSecret(ctx context.Context) ([]byte, error) {
	return s.signingSecret, nil
}
