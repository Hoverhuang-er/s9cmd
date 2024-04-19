package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v2"
)

func ListBucket(ctx context.Context, config *Config, c *cli.Context) error {
	args := c.Args().Slice()

	svc := SessionNewV2(config)

	if len(args) == 0 {
		var params *s3.ListBucketsInput
		resp, err := svc.ListBuckets(ctx, params)
		if err != nil {
			return err
		}

		for _, bucket := range resp.Buckets {
			fmt.Printf("%s  s3://%s\n", bucket.CreationDate.Format(DATE_FMT), *bucket.Name)
		}
		return nil
	}

	return listBucket(ctx, config, svc, args)
}

func listBucket(ctx context.Context, config *Config, svc *s3.Client, args []string) error {
	for _, arg := range args {
		u, err := FileURINew(arg)
		if err != nil || u.Scheme != "s3" {
			return fmt.Errorf("ls requires buckets to be prefixed with s3://")
		}

		//_, err = SessionForBucket(config, u.Bucket)
		//if err != nil {
		//	return err
		//}

		todo := []string{arg}

		for len(todo) != 0 {
			var item string
			item, todo = todo[0], todo[1:]

			if err := RemotePager(ctx, config, svc, item, !config.Recursive, func(page *s3.ListObjectsV2Output) {
				for _, item := range page.CommonPrefixes {
					uri := fmt.Sprintf("s3://%s/%s", u.Bucket, *item.Prefix)

					if config.Recursive {
						todo = append(todo, uri)
					} else {
						fmt.Printf("%16s %9s   %s\n", "", "DIR", uri)
					}
				}
				if page.Contents != nil {
					for _, item := range page.Contents {
						fmt.Printf("%16s %9d   s3://%s/%s\n", item.LastModified.Format(DATE_FMT), *item.Size, u.Bucket, *item.Key)
					}
				}
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
