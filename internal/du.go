package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"
)

func GetUsage(ctx context.Context, config *Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNewV2(config)

	// If we're not passed any args, we're going to do all S3 buckets
	if len(args) == 0 {
		var params *s3.ListBucketsInput
		resp, err := svc.ListBuckets(ctx, params)
		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			args = append(args, fmt.Sprintf("s3://%s", *bucket.Name))
		}
	}

	// Get the usage for the buckets
	for _, arg := range args {
		// Only do usage on S3 buckets
		u, err := FileURINew(arg)
		if err != nil || u.Scheme != "s3" {
			continue
		}

		var (
			bucketSize, bucketObjs int64
		)

		if err := RemotePager(ctx, config, svc, arg, false, func(page *s3.ListObjectsV2Output) {
			for _, obj := range page.Contents {
				bucketSize += *obj.Size
				bucketObjs += 1
			}
		}); err != nil {
			return err
		}
		fmt.Printf("%d %d objects %s\n", bucketSize, bucketObjs, arg)
	}

	return nil
}
