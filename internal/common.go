package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// DefaultRegion to use for S3 credential creation
const defaultRegion = "ap-southeast-1"
const DATE_FMT = "2006-01-02 15:04"

var (
	Threads  int  = 10
	IsGCR    bool = false
	LogLevel int  = 3
)

type FileURI struct {
	Scheme string
	Bucket string
	Path   string
}
type FileObject struct {
	Source   int64 // used by sync
	Name     string
	Size     int64
	Checksum string
}

func FileURINew(path string) (*FileURI, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "" && u.Scheme != "s3" && u.Scheme != "file" {
		return nil, fmt.Errorf("Invalid URI scheme must be one of file/s3/NONE")
	}

	uri := FileURI{
		Scheme: u.Scheme,
		Bucket: u.Host,
		Path:   u.Path,
	}

	if uri.Scheme == "" {
		uri.Scheme = "file"
	}
	if uri.Scheme == "s3" && uri.Path != "" {
		uri.Path = uri.Path[1:]
	}
	if uri.Path == "" && uri.Scheme == "s3" {
		uri.Path = "/"
	}

	return &uri, nil
}
func RemotePager(ctx context.Context, config *Config, svc *s3.Client, uri string, delim bool, pager func(page *s3.ListObjectsV2Output)) error {
	u, err := FileURINew(uri)
	if err != nil || u.Scheme != "s3" {
		return fmt.Errorf("requires buckets to be prefixed with s3://")
	}

	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(u.Bucket), // Required
		MaxKeys: aws.Int32(1000),
	}
	if u.Path != "" && u.Path != "/" {
		params.Prefix = u.Key()
	}
	if delim {
		params.Delimiter = aws.String("/")
	}

	if svc == nil {
		svc = SessionNewV2(config)
	}

	bsvc, err := SessionForBucket(ctx, config, u.Bucket)
	if err != nil {
		return err
	}
	res, err := bsvc.ListObjectsV2(ctx, params)
	if err != nil {
		return err
	}
	pager(res)
	return nil
}

func RemoteList(ctx context.Context, config *Config, svc *s3.Client, args []string) ([]FileObject, error) {
	result := make([]FileObject, 0)

	for _, arg := range args {
		pager := func(page *s3.ListObjectsV2Output) {
			for _, obj := range page.Contents {
				result = append(result, FileObject{
					Name:     *obj.Key,
					Size:     *obj.Size,
					Checksum: *obj.ETag,
				})
			}
		}

		go func() {
			if err := RemotePager(ctx, config, svc, arg, false, pager); err != nil {
				slog.Error(err.Error())
			}
		}()
	}

	return result, nil
}

// Return the path as a valid S3 bucket key
func (uri *FileURI) Key() *string {
	if uri.Path[0] == '/' {
		s := uri.Path[1:]
		return &s
	}
	return &uri.Path
}

// Return a string version of the path
func (uri *FileURI) String() string {
	if uri.Scheme == "s3" {
		return fmt.Sprintf("s3://%s/%s", uri.Bucket, *uri.Key())
	} else {
		return fmt.Sprintf("file://%s", uri.Path)
	}
}

// Do a path.Join() style operation on this FileURI to generate a new one
func (uri *FileURI) Join(elem string) *FileURI {
	nuri := FileURI{
		Scheme: uri.Scheme,
		Bucket: uri.Bucket,
	}

	if elem == "" {
		nuri.Path = uri.Path
	} else if elem[0] == '/' {
		nuri.Path = elem
	} else {
		// TODO: https://golang.org/pkg/net/url/#URL.ResolveReference
		nuri.Path = path.Join(filepath.Dir(uri.Path), elem)
		if elem[len(elem)-1] == '/' {
			nuri.Path += "/"
		}
	}

	return &nuri
}

func (uri *FileURI) SetPath(elem string) *FileURI {
	nuri := FileURI{
		Scheme: uri.Scheme,
		Bucket: uri.Bucket,
		Path:   elem,
	}
	if uri.Path == "" && uri.Scheme == "s3" {
		uri.Path = "/"
	}

	return &nuri
}

// CamelToSnake converts a given string to snake case
func CamelToSnake(s string) string {
	result := ""
	words := make([]string, 0)
	lastPos := 0
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}
func GetEnv(name string) *string {
	name = name + "="
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, name) {
			value := env[len(name):len(env)]
			return &value
		}
	}
	return nil
}

// buildSessionConfigV2 - Build a session configuration based on the provided config
func buildSessionConfigV2(config *Config) aws.Config {
	// By default make sure a region is specified, this is required for S3 operations
	sessionConfig := aws.Config{Region: defaultRegion}

	if config.AccessKey != "" && config.SecretKey != "" {
		sessionConfig.Credentials = credentials.NewStaticCredentialsProvider(config.AccessKey, config.SecretKey, "")
	}

	return sessionConfig
}

func buildEndpointResolver(hostname string) endpoints.Resolver {
	defaultResolver := endpoints.DefaultResolver()

	fixedHost := hostname
	if !strings.HasPrefix(hostname, "http") {
		fixedHost = "https://" + hostname
	}

	return endpoints.ResolverFunc(func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL: fixedHost,
			}, nil
		}

		return defaultResolver.EndpointFor(service, region, optFns...)
	})
}

// SessionNew - Read the config for default credentials, if not provided use environment based variables
func SessionNewV2(config *Config) *s3.Client {
	sessionConfig := buildSessionConfigV2(config)

	if config.HostBase != "" && config.HostBase != "s3.amazon.com" {
		sessionConfig.BaseEndpoint = &config.HostBase
	}
	return s3.NewFromConfig(sessionConfig)
}

// SessionForBucket - For a given S3 bucket, create an approprate session that references the region
// that this bucket is located in
func SessionForBucket(ctx context.Context, config *Config, bucket string) (*s3.Client, error) {
	sessionConfig := buildSessionConfigV2(config)

	if config.HostBucket == "" || config.HostBucket == "%(bucket)s.s3.amazonaws.com" {
		svc := SessionNewV2(config)

		if loc, err := svc.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: &bucket}); err != nil {
			return nil, err
		} else if loc.LocationConstraint == "" {
			// Use default service
			return svc, nil
		}
	}

	return s3.NewFromConfig(sessionConfig), nil
}
