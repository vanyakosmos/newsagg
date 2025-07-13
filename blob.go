package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type BlobTracker struct {
	client       *minio.Client
	fileAgeLimit time.Duration
	bucketName   string
}

func NewBlobTracker(ctx context.Context, endpoint string, accessKey string, secretKey string, region string, bucketName string) Tracker {
	client := mustInitMinioClient(endpoint, accessKey, secretKey, region)
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{ForceCreate: false})
	if err != nil {
		if !strings.Contains(err.Error(), "request to create the named bucket succeeded") {
			log.Fatalln("Bucket creation error:", err)
		}
	}
	return &BlobTracker{
		client:       client,
		fileAgeLimit: 5 * time.Minute,
		bucketName:   bucketName,
	}
}

func mustInitMinioClient(endpoint string, accessKey string, secretKey string, region string) *minio.Client {
	secure := !strings.HasPrefix(endpoint, "localhost")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
		Region: region,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v\n", err)
	}
	return client
}

func (t *BlobTracker) IsTracked(ctx context.Context, articleID string) bool {
	filename := t.getFilename(articleID)

	_, err := t.client.StatObject(ctx, t.bucketName, filename, minio.StatObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return false
		}
		log.Fatalln("Getting tracker error:", err)
	}
	return true
}

func (t *BlobTracker) MarkAsTracked(ctx context.Context, articleID string) {
	filename := t.getFilename(articleID)

	emptyContent := strings.NewReader("")
	_, err := t.client.PutObject(ctx, t.bucketName, filename, emptyContent, 0, minio.PutObjectOptions{
		ContentType: "plain/text",
	})
	if err != nil {
		log.Fatalln("Marking tracker error:", err)
	}
}

func (t *BlobTracker) CleanupOldTrackers(ctx context.Context) {
	opts := minio.ListObjectsOptions{
		Prefix:    "trackers/",
		Recursive: false,
	}

	for object := range t.client.ListObjects(ctx, t.bucketName, opts) {
		if object.Err != nil {
			log.Println("Cleanup error:", object.Err)
			return
		}
		if time.Since(object.LastModified) > t.fileAgeLimit {
			err := t.client.RemoveObject(ctx, t.bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Printf("Failed to remove object %s: %v", object.Key, err)
			} else {
				log.Printf("Deleted old object: %s", object.Key)
			}
		}
	}
}

func (t *BlobTracker) getFilename(articleID string) string {
	return fmt.Sprintf("trackers/hn_%s.txt", articleID)
}
