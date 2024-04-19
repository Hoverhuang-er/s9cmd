# s9cmd for AWS S3 and S3 compatible storage services.
### Super S3 command line tool named s9cmd, Faster 20x to s3cmd and s4cmd
[![Build Status](https://travis-ci.com/bloomreach/s9cmd.svg?branch=main)](https://travis-ci.com/bloomreach/s9cmd)
[![Packaging status](https://repology.org/badge/tiny-repos/s9cmd.svg)](https://repology.org/project/s9cmd/versions)
[Release Download](https://github.com/Hoverhuang-er/s9cmd/actions/runs/8749507028#artifacts)
----

**Author**: Shuhao Huang ([@Hoverhuang-er](https://github.com/Hoverhuang-er))

----

## What's New in s9cmd

- Fully migrated from old boto3 to new AWS Go SDKv2  library, which provides more reliable and high-performance S3 backend.
- Support S3 `--API-ServerSideEncryption` along with **36 new API pass-through options**. See API pass-through options section for complete list.
- Support batch delete (with delete_objects API) to delete up to 1000 files with single call. **100+ times faster** than sequential deletion.
- Support `S9CMD_OPTS` environment variable for commonly used options such as `--API-ServerSideEncryption` across all your s9cmd operations.
- Support moving files **larger than 5GB** with multipart upload. **20+ times faster** then sequential move operation when moving large files.
- Faster upload with lazy evaluation of md5 hash.
- New directory to directory `dsync` command is better and standalone implementation to replace old `sync` command, which is implemented based on top of get/put/mv commands. `--delete-removed` work for all cases including local to s3, s3 to local, and s3 to s3. `sync` command preserves the old behavior in this version for compatibility.
- [Support for S3 compatible storage services](https://github.com/bloomreach/s9cmd/issues/52) such as DreamHost and Cloudian using `--endpoint-url` (Community Supported Beta Feature).
- Tested on both Go 1.13, 1.16, 1.20, and 1.22.

## Motivation

s9cmd is a command-line utility for accessing and managing files on local filesystems
[Amazon S3](http://en.wikipedia.org/wiki/Amazon_S3), inspired by
[s3cmd](http://s3tools.org/s3cmd).

We have used s3cmd heavily for a number of scripted, data-intensive
applications. However as the need for a variety of small improvements arose, we
created our own implementation, s9cmd. It is intended as an alternative to
s3cmd for enhanced performance and for large files, and with a number of
additional features and fixes that we have found useful.

It strives to be compatible with the most common usage scenarios for s3cmd. It
does not offer exact drop-in compatibility, due to a number of corner cases where
different behavior seems preferable, or for bugfixes.


## Features

s9cmd supports the regular commands you might expect for fetching and storing
files in S3: `ls`, `put`, `get`, `cp`, `mv`, `sync`, `del`, `du`, `dsync`, and `batch-delete`.

Limitations:

- Support CloudFront feature
- Currently, we simulate `sync` with `get` and `put` with `--recursive --force --sync-check`.

## Installation and Setup
You can install `s9cmd` with go install command:

```
go install github.com/hoverhuanger/s9cmd
```

Or you can download the binary from the [release page]() and put it in your PATH.

```bash
wget https://github.com/Hoverhuang-er/s9cmd/releases/download/v1.0.0/s9cmd
```

Or you also can install s9cmd with docker and k8s job:

### Docker

```bash
docker pull hoverhuang/s9cmd:latest
docker run -it hoverhuang/s9cmd:latest s9cmd --help
```

### Kubernetes Job
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: s9cmd-install
spec:
    template:
        spec:
        containers:
        - name: s9cmd-install
            image: hoverhuang/s9cmd:latest
            command: ["s9cmd", "put", "s3://your-bucket/s9cmd", "/usr/local/bin/s9cmd"]
        restartPolicy: Never
```

- Copy or create a symbolic link so you can run `s9cmd.py` as `s9cmd`. (It is just
a single file!)
- If you already have a `~/.s3cfg` file from configuring `s3cmd`, credentials
from this file will be used.  Otherwise, set the `S3_ACCESS_KEY` and
`S3_SECRET_KEY` environment variables to contain your S3 credentials.
- If no keys are provided, but an IAM role is associated with the EC2 instance, it will
be used transparently.
- If you are using a non-AWS S3-compatible service, you can set the `S3_ENDPOINT` and `S3_REGION` environment variables.


## s9cmd Commands

#### `s9cmd ls [path]`

List all contents of a directory.

* -r/--recursive: recursively display all contents including subdirectories under the given path.
* -d/--show-directory: show the directory entry instead of its content.
* -l/--long: show detailed information of each file.
* -a/--all: show all files including hidden files.
* -t/--last-modified-before: show files that were last modified before the given time.
* -h/--last-modified-after: show files that were last modified after the given time.

#### `s9cmd put [source] [target]`

Upload local files up to S3.

*   -r/--recursive: also upload directories recursively.
*   -s/--sync-check: check md5 hash to avoid uploading the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real upload.

#### `s9cmd get [source] [target]`

Download files from S3 to local filesystem.

*   -r/--recursive: also download directories recursively.
*   -s/--sync-check: check md5 hash to avoid downloading the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real download.


#### `s9cmd dsync [source dir] [target dir]`

Synchronize the contents of two directories. The directory can either be local or remote, but currently, it doesn't support two local directories.

*   -r/--recursive: also sync directories recursively.
*   -s/--sync-check: check md5 hash to avoid syncing the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real sync.
*   --delete-removed: delete files not in source directory.

#### `s9cmd sync [source] [target]`

(Obsolete, use `dsync` instead) Synchronize the contents of two directories. The directory can either be local or remote, but currently, it doesn't support two local directories. This command simply invoke get/put/mv commands.

*   -r/--recursive: also sync directories recursively.
*   -s/--sync-check: check md5 hash to avoid syncing the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real sync.
*   --delete-removed: delete files not in source directory. Only works when syncing local directory to s3 directory.

#### `s9cmd cp [source] [target]`

Copy a file or a directory from a S3 location to another.

*   -r/--recursive: also copy directories recursively.
*   -s/--sync-check: check md5 hash to avoid copying the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real copy.

#### `s9cmd mv [source] [target]`

Move a file or a directory from a S3 location to another.

*   -r/--recursive: also move directories recursively.
*   -s/--sync-check: check md5 hash to avoid moving the same content.
*   -f/--force: override existing file instead of showing error message.
*   -n/--dry-run: emulate the operation without real move.

#### `s9cmd del [path]`

Delete files or directories on S3.

*   -r/--recursive: also delete directories recursively.
*   -n/--dry-run: emulate the operation without real delete.

#### `s9cmd du [path]`

Get the size of the given directory.
*   -h/--human-readable: show the size in human-readable format.
*   -s/--summarize: only show the total size of the directory.

Available parameters:

*   -r/--recursive: also add sizes of sub-directories recursively.

## s9cmd Control Options

##### `-p S3CFG, --config=[filename]`
path to s3cfg config file

##### `-f, --force`
force overwrite files when download or upload

##### `-r, --recursive`
recursively checking subdirectories

##### `-s, --sync-check`
check file md5 before download or upload

##### `-n, --dry-run`
trial run without actual download or upload

##### `-t RETRY, --retry=[integer]`
number of retries before giving up

##### `--retry-delay=[integer]`
seconds to sleep between retries

##### `-c NUM_THREADS, --num-threads=NUM_THREADS`
number of concurrent threads

##### `--endpoint-url`
endpoint url used in boto3 client

##### `-d, --show-directory`
show directory instead of its content

##### `--ignore-empty-source`
ignore empty source from s3

##### `--use-ssl`
(obsolete) use SSL connection to S3

##### `--verbose`
verbose output

##### `--debug`
debug output

##### `--validate`
(obsolete) validate lookup operation

##### `-D, --delete-removed`
delete remote files that do not exist in source after sync

##### `--multipart-split-size=[integer]`
size in bytes to split multipart transfers

##### `--max-singlepart-download-size=[integer]`
files with size (in bytes) greater than this will be
downloaded in multipart transfers

##### `--max-singlepart-upload-size=[integer]`
files with size (in bytes) greater than this will be
uploaded in multipart transfers

##### `--max-singlepart-copy-size=[integer]`
files with size (in bytes) greater than this will be
copied in multipart transfers

##### `--batch-delete-size=[integer]`
Number of files (&lt;1000) to be combined in batch delete.

##### `--last-modified-before=[datetime]`
Condition on files where their last modified dates are
before given parameter.

##### `--last-modified-after=[datetime]`
Condition on files where their last modified dates are
after given parameter.


## S3 API Pass-through Options

Those options are directly translated to boto3 API commands. The options provided will be filtered by the APIs that are taking parameters. For example, `--API-ServerSideEncryption` is only needed for `put_object`, `create_multipart_upload` but not for `list_buckets` and `get_objects` for example. Therefore, providing `--API-ServerSideEncryption` for `s9cmd ls` has no effect.

For more information, please see AWS Go SDKv2 s3 documentations []

##### `--API-ACL=[string]`
The canned ACL to apply to the object.

##### `--API-CacheControl=[string]`
Specifies caching behavior along the request/reply chain.

##### `--API-ContentDisposition=[string]`
Specifies presentational information for the object.

##### `--API-ContentEncoding=[string]`
Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.

##### `--API-ContentLanguage=[string]`
The language the content is in.

##### `--API-ContentMD5=[string]`
The base64-encoded 128-bit MD5 digest of the part data.

##### `--API-ContentType=[string]`
A standard MIME type describing the format of the object data.

##### `--API-CopySourceIfMatch=[string]`
Copies the object if its entity tag (ETag) matches the specified tag.

##### `--API-CopySourceIfModifiedSince=[datetime]`
Copies the object if it has been modified since the specified time.

##### `--API-CopySourceIfNoneMatch=[string]`
Copies the object if its entity tag (ETag) is different than the specified ETag.

##### `--API-CopySourceIfUnmodifiedSince=[datetime]`
Copies the object if it hasn't been modified since the specified time.

##### `--API-CopySourceRange=[string]`
The range of bytes to copy from the source object. The range value must use the form bytes=first-last, where the first and last are the zero-based byte offsets to copy. For example, bytes=0-9 indicates that you want to copy the first ten bytes of the source. You can copy a range only if the source object is greater than 5 GB.

##### `--API-CopySourceSSECustomerAlgorithm=[string]`
Specifies the algorithm to use when decrypting the source object (e.g., AES256).

##### `--API-CopySourceSSECustomerKeyMD5=[string]`
Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. Amazon S3 uses this header for a message integrity check to ensure the encryption key was transmitted without error. Please note that this parameter is automatically populated if it is not provided. Including this parameter is not required

##### `--API-CopySourceSSECustomerKey=[string]`
Specifies the customer-provided encryption key for Amazon S3 to use to decrypt the source object. The encryption key provided in this header must be one that was used when the source object was created.

##### `--API-ETag=[string]`
Entity tag returned when the part was uploaded.

##### `--API-Expires=[datetime]`
The date and time at which the object is no longer cacheable.

##### `--API-GrantFullControl=[string]`
Gives the grantee READ, READ_ACP, and WRITE_ACP permissions on the object.

##### `--API-GrantReadACP=[string]`
Allows grantee to read the object ACL.

##### `--API-GrantRead=[string]`
Allows grantee to read the object data and its metadata.

##### `--API-GrantWriteACP=[string]`
Allows grantee to write the ACL for the applicable object.

##### `--API-IfMatch=[string]`
Return the object only if its entity tag (ETag) is the same as the one specified, otherwise return a 412 (precondition failed).

##### `--API-IfModifiedSince=[datetime]`
Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified).

##### `--API-IfNoneMatch=[string]`
Return the object only if its entity tag (ETag) is different from the one specified, otherwise return a 304 (not modified).

##### `--API-IfUnmodifiedSince=[datetime]`
Return the object only if it has not been modified since the specified time, otherwise return a 412 (precondition failed).

##### `--API-Metadata=[dict]`
A map (in json string) of metadata to store with the object in S3

##### `--API-MetadataDirective=[string]`
Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request.

##### `--API-MFA=[string]`
The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device.

##### `--API-RequestPayer=[string]`
Confirms that the requester knows that she or he will be charged for the request. Bucket owners need not specify this parameter in their requests. Documentation on downloading objects from requester pays buckets can be found at http://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html

##### `--API-ServerSideEncryption=[string]`
The Server-side encryption algorithm used when storing this object in S3 (e.g., AES256, aws:kms).

##### `--API-SSECustomerAlgorithm=[string]`
Specifies the algorithm to use to when encrypting the object (e.g., AES256).

##### `--API-SSECustomerKeyMD5=[string]`
Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. Amazon S3 uses this header for a message integrity check to ensure the encryption key was transmitted without error. Please note that this parameter is automatically populated if it is not provided. Including this parameter is not required

##### `--API-SSECustomerKey=[string]`
Specifies the customer-provided encryption key for Amazon S3 to use in encrypting data. This value is used to store the object and then it is discarded; Amazon does not store the encryption key. The key must be appropriate for use with the algorithm specified in the x-amz-server-side-encryption-customer-algorithm header.

##### `--API-SSEKMSKeyId=[string]`
Specifies the AWS KMS key ID to use for object encryption. All GET and PUT requests for an object protected by AWS KMS will fail if not made via SSL or using SigV4. Documentation on configuring any of the officially supported AWS SDKs and CLI can be found at http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingAWSSDK.html#specify-signature-version

##### `--API-StorageClass=[string]`
The type of storage to use for the object. Defaults to 'STANDARD'.

##### `--API-VersionId=[string]`
VersionId used to reference a specific version of the object.

##### `--API-WebsiteRedirectLocation=[string]`
If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. Amazon S3 stores the value of this header in the object metadata.


## Debugging Tips

Simply enable `--debug` option to see the full log of s9cmd. If you even need to check what APIs are invoked from s9cmd to boto3, you can run:

```
s9cmd --debug [op] .... 2>&1 >/dev/null | grep S3APICALL
```

To see all the parameters sending to S3 API.


## Compatibility between s3cmd and s9cmd

Prefix matching: In s3cmd, unlike traditional filesystems, prefix names match listings:

```
>> s3cmd ls s3://my-bucket/ch
s3://my-bucket/charlie/
s3://my-bucket/chyang/
```

In s9cmd, behavior is the same as with a Unix shell:

```
>>s9cmd ls s3://my-bucket/ch*
2019/07/01 15:00:00 s3://my-bucket/charlie/
```

To get prefix behavior, use explicit wildcards instead: s9cmd ls s3://my-bucket/ch*

Similarly, sync and cp commands emulate the Unix cp command, so directory to
directory sync use different syntax:

```
>> s3cmd sync s3://bucket/path/dirA s3://bucket/path/dirB/
```
will copy contents in dirA to dirB.
```
>> s9cmd sync s3://bucket/path/dirA s3://bucket/path/dirB/
```
will copy dirA *into* dirB.

To achieve the s3cmd behavior, use wildcards:
```
s9cmd sync s3://bucket/path/dirA/* s3://bucket/path/dirB/
```

Note s9cmd doesn't support dirA without trailing slash indicating dirA/* as
what rsync supported.

No automatic override for put command:
s3cmd put fileA s3://bucket/path/fileB will return error if fileB exists.
Use -f as well as get command.

Bugfixes for handling of non-existent paths: Often s3cmd creates empty files when specified paths do not exist:
s3cmd get s3://my-bucket/no_such_file downloads an empty file.
s9cmd get s3://my-bucket/no_such_file returns an error.
s3cmd put no_such_file s3://my-bucket/ uploads an empty file.
s9cmd put no_such_file s3://my-bucket/ returns an error.
