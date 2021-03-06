package main

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	ctx           = context.Background()
	credential 		azblob.TokenCredential
	tenant			*string
	appId			*string
	subId			*string
	store			*string
	container		*string
)

// Authenticate with the Azure services using file-based authentication
func init() {
	subId = flag.String("subid", "", "Azure SubscriptionId (Required)")
	appId = flag.String("appid", "", "App Registration Id (Required)")
	tenant = flag.String("tenantid", "", "Tenant Id (Required)")
	store = flag.String("store", "", "Storage Acct for Upload (Required)")
	container = flag.String("container", "", "Container for Upload (Required)")
	flag.Parse()

	if *subId == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *appId == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *tenant == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *store == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *container == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Setup a new Device Token configuration
	dfc := auth.NewDeviceFlowConfig(*appId, *tenant)
	
	// Set the resource for the storage endpoint
	dfc.Resource = "https://storage.azure.com/"
	
	// Prompt user for authentication
	storageToken, err := dfc.ServicePrincipalToken()
	if err != nil {
		log.Fatal(err)
	}

	// Setup a storage Token credential to use for upload
	credential = azblob.NewTokenCredential(storageToken.OAuthToken(), nil)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}
}

func main() {
	// Create a file in container of store
	createFileinStorageAcct(*store, *container)

}

func createFileinStorageAcct(storageAcct string, container string) {
	rand.Seed(time.Now().UnixNano())
	randnum := rand.Intn(999999)
	// Create a file to test the upload
	data := []byte("hello from Microsoft. this is a blob " + strconv.Itoa(randnum) + "\n")
	fileName := "file-" + strconv.Itoa(randnum)
	err := ioutil.WriteFile(fileName, data, 0700)
	if err != nil {
		log.Fatal(err)
	}

	pipeline := azblob.NewPipeline(
		credential,
		azblob.PipelineOptions{},
	)

	// Hardcoded for Azure commercial.  Could be adjusted for MAG
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", storageAcct, container))
	if err != nil {
		log.Fatal(err)
	}
	containerURL := azblob.NewContainerURL(
		*URL,
		pipeline,
	)

	blobURL := containerURL.NewBlockBlobURL(fileName)

	file, err := os.Open(fileName)

	fmt.Printf("Uploading file with blob name: %s", fileName)

	_, err = azblob.UploadFileToBlockBlob(
		ctx,
		file,
		blobURL,
		azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16},
	)
	if err != nil {
		log.Fatal(err)
	}
}

