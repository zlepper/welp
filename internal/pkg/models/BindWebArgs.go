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

package models

import "time"

type BindWebArgs struct {
	// If the web server should automatically fetch an https certificate and use that
	// If enabled, will ignore the port argument, and always bind on port 80 and 433
	// Will also configure HSTS
	UseHttps bool
	// The folder to cache https certificates in
	CertificateCacheFolder string

	// The port to bind to
	Port int

	// Where to save the uploaded files
	FolderPath string

	// The name of the folder where the database files should be stored
	// when using flat-file storage
	DatabaseFolderName string
	// How often the files should be saved when using flatfiles
	SaveInterval time.Duration

	// An optional key for sendgrid
	SendGridApiKey string

	// How long autho
	TokenDuration time.Duration

	EmailSenderName, EmailSenderAddress string
}
