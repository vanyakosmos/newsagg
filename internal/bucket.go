package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type bucketTracker struct {
	client       *minio.Client
	fileAgeLimit time.Duration
	bucketName   string
}

func NewBucketTracker(ctx context.Context, endpoint string, accessKey string, secretKey string, region string, bucketName string) tracker {
	client := mustInitMinioClient(endpoint, accessKey, secretKey, region)
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{ForceCreate: false})
	if err != nil {
		if !strings.Contains(err.Error(), "request to create the named bucket succeeded") {
			log.Printf("Bucket %s creation error: %s\n", bucketName, err)
		}
	}
	return &bucketTracker{
		client:       client,
		fileAgeLimit: 7 * 24 * time.Hour,
		bucketName:   bucketName,
	}
}

func mustInitMinioClient(endpoint string, accessKey string, secretKey string, region string) *minio.Client {
	secure := !strings.HasPrefix(endpoint, "minio")

	var creds *credentials.Credentials
	if accessKey == "" && secretKey == "" {
		creds = credentials.NewIAM(endpoint)
	} else {
		creds = credentials.NewStaticV4(accessKey, secretKey, "")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: secure,
		Region: region,
	})
	if err != nil {
		log.Panicf("Failed to create MinIO client: %v\n", err)
	}
	return client
}

func (t *bucketTracker) IsTracked(ctx context.Context, article newsArticle) bool {
	filename := t.getFilename(article)

	_, err := t.client.StatObject(ctx, t.bucketName, filename, minio.StatObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return false
		}
		log.Panicln("Getting tracker error:", err)
	}
	return true
}

func (t *bucketTracker) MarkAsTracked(ctx context.Context, article newsArticle) {
	filename := t.getFilename(article)

	emptyContent := strings.NewReader("")
	_, err := t.client.PutObject(ctx, t.bucketName, filename, emptyContent, 0, minio.PutObjectOptions{
		ContentType: "plain/text",
	})
	if err != nil {
		log.Panicln("Marking tracker error:", err)
	}
}

func (t *bucketTracker) CleanupOldTrackers(ctx context.Context) {
	log.Printf("Attepting cleanup...")
	opts := minio.ListObjectsOptions{
		Prefix:    "trackers/",
		Recursive: false,
	}

	for object := range t.client.ListObjects(ctx, t.bucketName, opts) {
		if object.Err != nil {
			log.Panicln("Cleanup error:", object.Err)
			return
		}
		if time.Since(object.LastModified) > t.fileAgeLimit {
			err := t.client.RemoveObject(ctx, t.bucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Printf("Failed to delete object %s: %v", object.Key, err)
			} else {
				log.Printf("Deleted old object: %s", object.Key)
			}
		}
	}
}

func (t *bucketTracker) getFilename(article newsArticle) string {
	return fmt.Sprintf("trackers/%s_%s.txt", article.Source, article.ID)
}
