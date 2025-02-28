package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
)

type FolderRef struct {
	// basePath should be the prefix of object name
	// e.g. basePath = "foo/bar/", object name = "foo/bar/baz"
	// e.g. basePath = "ynufes/", object name = "ynufes/shion"
	basePath string
	bucket   *storage.BucketHandle
}

func NewFolderRef(bucket *storage.BucketHandle, basePath string) FolderRef {
	if basePath != "" && basePath[len(basePath)-1] != '/' {
		basePath += "/"
	}
	return FolderRef{
		basePath: basePath,
		bucket:   bucket,
	}
}

func (f FolderRef) Object(name string) ObjectRef {
	return ObjectRef{
		objName: f.basePath + name,
		bucket:  f.bucket,
	}
}

func (f FolderRef) Folder(name string) FolderRef {
	return FolderRef{
		basePath: f.basePath + name + "/",
		bucket:   f.bucket,
	}
}

func (f FolderRef) List(ctx context.Context) ([]FolderRef, []ObjectRef, error) {
	results := f.bucket.Objects(ctx, &storage.Query{
		Prefix:    f.basePath,
		Delimiter: "/",
	})
	var folders []FolderRef
	var objects []ObjectRef
	for {
		next, err := results.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		if next.Prefix != "" {
			folders = append(folders, FolderRef{
				basePath: next.Prefix,
				bucket:   f.bucket,
			})
		} else {
			objects = append(objects, ObjectRef{
				objName: next.Name,
				bucket:  f.bucket,
			})
		}
	}
	return folders, objects, nil
}

func (f FolderRef) Upload(ctx context.Context, name string, data []byte) (*ObjectRef, error) {
	w := f.bucket.Object(f.basePath + name).NewWriter(ctx)
	defer w.Close()
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	return &ObjectRef{
		objName: f.basePath + name,
		bucket:  f.bucket,
	}, nil
}
