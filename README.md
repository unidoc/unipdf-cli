# unicli

[![GoDoc](https://godoc.org/github.com/unidoc/unicli?status.svg)](https://godoc.org/github.com/unidoc/unicli)
[![Go Report Card](https://goreportcard.com/badge/github.com/unidoc/unicli)](https://goreportcard.com/report/github.com/unidoc/unicli)

unicli is a CLI tool which makes working with PDF files very easy. It supports
the most common PDF operations. The application is written in Golang and is
powered by the [UniDoc](https://github.com/unidoc/unidoc) PDF library.

## Features

- [Merge PDF files](#merge)
- [Split PDF files](#split)
- [Explode PDF files](#explode)
- [Encrypt PDF files](#encrypt)
- [Decrypt PDF files](#decrypt)
- [Change user/owner password](#passwd)
- [Optimize PDF files](#optimize)
- [Add watermark images to PDF files](#watermark)
- [Convert PDF files to grayscale](#grayscale)
- [Validate and print PDF file information](#info)
- [Extract text from PDF files](#extract)
- [Extract images from PDF files](#extract)
- [Search text in PDF files](#search)

## Short demo

[![asciicast](https://i.imgur.com/nQZq6T7.png)](https://asciinema.org/a/220314)

## Installation

Minimum required Go version: 1.11

```
git clone git@github.com:unidoc/unicli.git
cd unicli
go build
```

Go modules are disabled by default in GOPATH/src. If you choose to clone the
project somewhere in this location, you must explicitly enable Go modules.

```
git clone git@github.com:unidoc/unicli.git
cd unicli
export GO111MODULE=on
go build
```

## Showcase

#### Grayscale conversion

![encrypt example](https://i.imgur.com/9QgXWUc.png)

#### Add watermark

![watermark example](https://i.imgur.com/GIRsTnT.png)

## Usage

#### Merge

Merge multiple PDF files into a single output file.

```
unicli merge OUTPUT_FILE INPUT_FILE...

Examples:
unicli merge output_file.pdf input_file1.pdf input_file2.pdf
```

#### Split

Extract one or more page ranges from PDF file and save the result as a
single output file.

```
unicli split [FLAG]... INPUT_FILE OUTPUT_FILE [PAGES]

Flags:
-p, --password string   PDF file password

Examples:
unicli split input_file.pdf output_file.pdf 1-2
unicli split -p pass input_file.pd output_file.pdf 1-2,4

PAGES parameter example: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be present in the output file,
while page number 5 is skipped.
```

#### Explode

Splits the input file into separate single page PDF files and saves the result
as a ZIP archive.

```
Usage:
unicli explode [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages to extract from the input file
-p, --password string      Input file password

Examples:
unicli explode input_file.pdf
unicli explode -o pages.zip input_file.pdf
unicli explode -o pages.zip -P 1-3 input_file.pdf
unicli explode -o pages.zip -P 1-3 -p pass input_file.pdf

Pages parameter example: 1-3,4,6-7
Pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be extracted, while page
number 5 is skipped.
```

#### Encrypt

Add password protection to PDF files. Owner and user passwords can be
specified, along with a set of user permissions. The encryption algorithm
used for protecting the file is configurable.

```
unicli encrypt [FLAG]... INPUT_FILE OWNER_PASSWORD [USER_PASSWORD]

Flags:
-m, --mode string         Algorithm to use for encrypting the file (default "rc4")
-o, --output-file string  Output file
-P, --perms string        User permissions (default "all")

Examples:
unicli encrypt input_file.pdf owner_pass
unicli encrypt input_file.pdf owner_pass user_pass
unicli encrypt -o output_file.pdf -m aes256 input_file.pdf owner_pass user_pass
unicli encrypt -o output_file.pdf -P none -m aes256 input_file.pdf owner_pass user_pass
unicli encrypt -o output_file.pdf -P modify,annotate -m aes256 input_file.pdf owner_pass user

Supported encryption algorithms:
- rc4 (default)
- aes128
- aes256

Supported user permissions:
- all (default)
- none
- print-low-res
- print-high-res
- modify
- extract
- extract-graphics
- annotate
- fill-forms
- rotate
```

#### Decrypt

Remove password protection from PDF files.

```
unicli decrypt [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-p, --password string      PDF file password

Examples:
unicli decrypt -p pass input_file.pdf
unicli decrypt -p pass -o output_file.pdf input_file.pdf
```

#### Passwd

Change protected PDF user/owner password.

```
unicli passwd [FLAG]... INPUT_FILE NEW_OWNER_PASSWORD [NEW_USER_PASSWORD]

Flags:
-o, --output-file string   Output file
-p, --password string      PDF file password

Examples:
unicli passwd -p pass input_file.pdf new_owner_pass
unicli passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass
unicli passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass new_user_pass
```

#### Optimize

Optimize PDF files by removing redundant objects. The quality of the
contained images can also be configured.

```
unicli optimize [FLAG]... INPUT_FILE

Flags:
-q, --image-quality int    Optimized image quality (default 100)
-o, --output-file string   Output file
-p, --password string      File password

Examples:
unicli optimize input_file.pdf
unicli optimize -o output_file input_file.pdf
unicli optimize -o output_file -i 75 input_file.pdf
unicli optimize -o output_file -i 75 -p pass input_file.pdf
```

#### Watermark

Add watermark images to PDF files.

```
unicli watermark [FLAG]... INPUT_FILE WATERMARK_IMAGE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages on which to add watermark
-p, --password string      PDF file password

Examples:
unicli watermark input_file.pdf watermark.png
unicli watermark -o output file.png input_file.pdf watermark.png
unicli watermark -o output file.png -P 1-3 input_file.pdf watermark.png
unicli watermark -o output file.png -P 1-3 -p pass input_file.pdf watermark.png

Pages parameter example: 1-3,4,6-7
Watermark will only be applied to pages 1,2,3 (1-3), 4 and 6,7 (6-7), while
page number 5 is skipped.
```

#### Grayscale

Convert PDF files to grayscale.

```
unicli grayscale [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages to convert to grayscale
-p, --password string      PDF file password

Examples:
unicli grayscale input_file.pdf
unicli grayscale -o output_file input_file.pdf
unicli grayscale -o output_file -P 1-3 input_file.pdf
unicli grayscale -o output_file -P 1-3 -p pass input_file.pdf

Pages parameter example: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be converted to grayscale, while
page number 5 is skipped.
```

#### Info

Outputs file information. Also does some basic validation.

```
unicli info [FLAG]... INPUT_FILE

Flags:
-p, --password string   PDF file password

Examples:
unicli info input_file.pdf
unicli info -p pass input_file.pdf
```

#### Extract

Extract resources (text, images) from PDF files.

```
unicli extract [FLAG]... INPUT_FILE

Flags:
-o, --output-file string     Output file
-P, --pages string           Pages to extract resources from
-r, --resource string        Resource to extract
-p, --user-password string   Input file password

Examples:
unicli extract -r text input_file.pdf
unicli extract -r text -P 1-3 input_file.pdf
unicli extract -r text -P 1-3 -p pass input_file.pdf
unicli extract -r images input_file.pdf
unicli extract -r images -o images.zip input_file.pdf
unicli extract -r images -P 1-3 -p pass -o images.zip input_file.pdf

Pages parameter example: 1-3,4,6-7
Resources will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7),
while page number 5 is skipped.

Supported resources:
- text
- images
```

#### Search

Search text in PDF files.

```
unicli search [FLAG]... INPUT_FILE TEXT

Flags:
-p, --password string   PDF file password

Examples:
unicli search input_file.pdf text_to_search
unicli search -p pass input_file.pdf text_to_search
```

## License

The application is licensed under the same conditions as the
[UniDoc](https://github.com/unidoc/unidoc) library.
Is has a dual license, a commercial one suitable for closed source projects
and an AGPL license that can be used in open source software.

Depending on your needs, you can choose one of them and follow its policies.
A detail of the policies and agreements for each license type are available in
the [LICENSE.COMMERCIAL](LICENSE.COMMERCIAL) and [LICENSE.AGPL](LICENSE.AGPL)
files.

Please see [pricing](http://unidoc.io/pricing) to purchase a commercial
[UniDoc](https://github.com/unidoc/unidoc) license or contact sales at
sales@unidoc.io for more info.

If you have a license for [UniDoc](https://github.com/unidoc/unidoc), you can
set it through the UNIDOC_LICENSE_FILE and UNIDOC_LICENSE_CUSTOMER environment
variables.

```
export UNIDOC_LICENSE_FILE="PATH_TO_LICENSE_FILE"
export UNIDOC_LICENSE_CUSTOMER="CUSTOMER_NAME"
```
