package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
	"ynufes-mypage-backend/pkg/setting"
)

func TestObjectRef_IssueDownloadSignedURL(t *testing.T) {
	ctx := context.Background()
	conf := setting.Get()
	certPath := conf.Infrastructure.Firebase.JsonCredentialFile
	c, err := storage.NewClient(ctx, option.WithCredentialsFile(certPath))
	if err != nil {
		t.Fatal(err)
	}
	b := c.Bucket("ynufes-mypage-staging-bucket")
	ref := NewFolderRef(b, "some-folder")
	uploadRef, err := ref.Upload(ctx, "some-file", []byte("some-content"))
	if err != nil {
		t.Fatal(err)
		return
	}
	issueOpt := IssueLinkOptions{
		ExpiresIn: 5 * time.Second,
		AuthHeaders: map[string]string{
			"user": "shion-test",
		},
		AuthMetaHeaders: map[string]string{
			"user": "shion-meta",
		},
		AuthQueries: map[string]string{
			"user-query": "shion-query",
		},
	}
	targetUrl, err := uploadRef.IssueDownloadSignedURL(issueOpt, true)
	if err != nil {
		t.Fatal(err)
	}
	if err := verifySignedURL(targetUrl, issueOpt, true); err != nil {
		t.Fatal(err)
	}
}

//func uploadTestFile(
//	ctx context.Context,
//	bucket *storage.BucketHandle,
//	fileName string,
//	content []byte,
//) (ObjectRef, error) {
//	ref := NewFolderRef(bucket, "some-folder")
//	return ref.Upload(ctx, fileName, content)
//}

func verifySignedURL(
	target SignedObjectLink,
	opt IssueLinkOptions,
	addQuery bool,
) error {
	targetURL, err := url.Parse(string(target))
	if err != nil {
		return fmt.Errorf("failed to parse signed URL: %w", err)
	}
	if addQuery {
		q := targetURL.Query()
		for k, v := range opt.AuthQueries {
			q.Set(k, v)
		}
		targetURL.RawQuery = q.Encode()
	}
	fmt.Println("targetURL: ", targetURL.String())
	req, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range opt.AuthHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range opt.AuthMetaHeaders {
		req.Header.Set("x-goog-meta-"+k, v)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	payload, _ := io.ReadAll(resp.Body)
	if string(payload) != "some-content" {
		return fmt.Errorf("unexpected payload: %s", string(payload))
	}
	return nil
}
