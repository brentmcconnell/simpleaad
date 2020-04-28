package main

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

	dfc := auth.NewDeviceFlowConfig(*appId, *tenant)
	dfc.Resource = "https://storage.azure.com/"
	storageToken, err := dfc.ServicePrincipalToken()
	if err != nil {
		log.Fatal(err)
	}

	credential = azblob.NewTokenCredential(storageToken.OAuthToken(), nil)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}
}

func main() {
	// Create a file in container1 of storageAcct1
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

