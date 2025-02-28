package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"net/url"
	"time"
)

type (
	ObjectRef struct {
		objName string
		bucket  *storage.BucketHandle
	}
	SignedObjectLink string
)

func (r ObjectRef) Name() string {
	return r.objName
}

func (r ObjectRef) Download(ctx context.Context) ([]byte, error) {
	reader, err := r.bucket.Object(r.objName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// IssueLinkOptions
// ExpiresIn is the duration the link is valid for.
// AuthHeaders are the headers to sign.
// AuthMetaHeaders are the metadata headers to sign. The key will be prefixed with "x-goog-meta-".
// headers starting with "x-goog-meta-" are considered metadata headers, and it will be set as metadata in the object.
// Both AuthHeaders and AuthMetaHeaders are required when using the URL.
// AuthQueries are the query parameters to sign.
type IssueLinkOptions struct {
	ExpiresIn       time.Duration
	AuthHeaders     map[string]string
	AuthMetaHeaders map[string]string
	AuthQueries     map[string]string
}

// IssueDownloadSignedURL
// removeAuthQueries will remove the query parameters from the signed URL.
// It will be useful for preventing CSRF attacks.
func (r ObjectRef) IssueDownloadSignedURL(ops IssueLinkOptions, removeAuthQueries bool) (SignedObjectLink, error) {
	headers := make([]string, 0, len(ops.AuthHeaders))
	for k, v := range ops.AuthHeaders {
		headers = append(headers, k+":"+v)
	}
	for k, v := range ops.AuthMetaHeaders {
		headers = append(headers, "x-goog-meta-"+k+":"+v)
	}
	queries := url.Values{}
	for k, v := range ops.AuthQueries {
		queries.Set(k, v)
	}
	signedUrl, err := r.bucket.SignedURL(r.objName, &storage.SignedURLOptions{
		Method:          "GET",
		Expires:         time.Now().Add(ops.ExpiresIn),
		Headers:         headers,
		QueryParameters: queries,
		Scheme:          storage.SigningSchemeV4,
	})
	if err != nil {
		return "", err
	}
	if removeAuthQueries {
		u, err := url.Parse(signedUrl)
		if err != nil {
			return "", err
		}
		q := u.Query()
		for k := range ops.AuthQueries {
			q.Del(k)
		}
		u.RawQuery = q.Encode()
		signedUrl = u.String()
	}
	return SignedObjectLink(signedUrl), nil
}
