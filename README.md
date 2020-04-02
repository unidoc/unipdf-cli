# unipdf-cli

[![Build Status](https://travis-ci.org/unidoc/unipdf-cli.svg?branch=master)](https://travis-ci.org/unidoc/unipdf-cli)
[![GoDoc](https://godoc.org/github.com/unidoc/unipdf-cli?status.svg)](https://godoc.org/github.com/unidoc/unipdf-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/unidoc/unipdf-cli)](https://goreportcard.com/report/github.com/unidoc/unipdf-cli)

unipdf-cli is a CLI tool which makes working with PDF files very easy. It supports
the most common PDF operations. The application is written in Golang and is
powered by the [UniPDF](https://github.com/unidoc/unipdf-cli) PDF library.

## Features

- [Merge PDF files](#merge)
- [Split PDF files](#split)
- [Explode PDF files](#explode)
- [Encrypt PDF files](#encrypt)
- [Decrypt PDF files](#decrypt)
- [Change user/owner password](#passwd)
- [Optimize PDF files](#optimize)
- [Rotate PDF pages](#rotate)
- [Add watermark images to PDF files](#watermark)
- [Convert PDF files to grayscale](#grayscale)
- [Validate and print PDF file information](#info)
- [Extract text from PDF files](#extract-text)
- [Extract images from PDF files](#extract-images)
- [Search text in PDF files](#search)
- [Export PDF form fields as JSON](#form-export)
- [Fill PDF form fields from JSON file](#form-fill)
- [Fill PDF form fields from FDF file](#fdf-merge)
- [Flatten PDF form fields](#form-flatten)
- [Render PDF pages to images](#render)

## Short demo

[![asciicast](https://i.imgur.com/nQZq6T7.png)](https://asciinema.org/a/220314)

## Installation

Minimum required Go version: 1.11

```
git clone git@github.com:unidoc/unipdf-cli.git
cd unipdf-cli/cmd/unipdf
go build
```

In Go 1.11 modules are disabled by default in GOPATH/src (`GO111MODULE=auto`).
Newer versions will have Go modules enabled by default. If you choose to clone
the project somewhere in this location, you must explicitly enable Go modules.

```
git clone git@github.com:unidoc/unipdf-cli.git
cd unipdf-cli/cmd/unipdf
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
unipdf merge OUTPUT_FILE INPUT_FILE...

Examples:
unipdf merge output_file.pdf input_file1.pdf input_file2.pdf
```

#### Split

Extract one or more page ranges from PDF file and save the result as a
single output file.

```
unipdf split [FLAG]... INPUT_FILE OUTPUT_FILE [PAGES]

Flags:
-p, --password string   PDF file password

Examples:
unipdf split input_file.pdf output_file.pdf 1-2
unipdf split -p pass input_file.pd output_file.pdf 1-2,4

PAGES argument example: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be present in the output file,
while page number 5 is skipped.
```

#### Explode

Splits the input file into separate single page PDF files and saves the result
as a ZIP archive.

```
Usage:
unipdf explode [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages to extract from the input file
-p, --password string      Input file password

Examples:
unipdf explode input_file.pdf
unipdf explode -o pages.zip input_file.pdf
unipdf explode -o pages.zip -P 1-3 input_file.pdf
unipdf explode -o pages.zip -P 1-3 -p pass input_file.pdf

Pages flag example: 1-3,4,6-7
Pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be extracted, while page
number 5 is skipped.
```

#### Encrypt

Add password protection to PDF files. Owner and user passwords can be
specified, along with a set of user permissions. The encryption algorithm
used for protecting the file is configurable.

```
unipdf encrypt [FLAG]... INPUT_FILE OWNER_PASSWORD [USER_PASSWORD]

Flags:
-m, --mode string         Algorithm to use for encrypting the file (default "rc4")
-o, --output-file string  Output file
-P, --perms string        User permissions (default "all")

Examples:
unipdf encrypt input_file.pdf owner_pass
unipdf encrypt input_file.pdf owner_pass user_pass
unipdf encrypt -o output_file.pdf -m aes256 input_file.pdf owner_pass user_pass
unipdf encrypt -o output_file.pdf -P none -m aes256 input_file.pdf owner_pass user_pass
unipdf encrypt -o output_file.pdf -P modify,annotate -m aes256 input_file.pdf owner_pass user

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
unipdf decrypt [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-p, --password string      PDF file password

Examples:
unipdf decrypt -p pass input_file.pdf
unipdf decrypt -p pass -o output_file.pdf input_file.pdf
```

#### Passwd

Change protected PDF user/owner password.

```
unipdf passwd [FLAG]... INPUT_FILE NEW_OWNER_PASSWORD [NEW_USER_PASSWORD]

Flags:
-o, --output-file string   Output file
-p, --password string      PDF file password

Examples:
unipdf passwd -p pass input_file.pdf new_owner_pass
unipdf passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass
unipdf passwd -p pass -o output_file.pdf input_file.pdf new_owner_pass new_user_pass
```

#### Optimize

Optimize PDF files by optimizing structure, compression and image quality.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_optimized" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the optimized output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

The quality of the images in the output files can be configured through
the --image-quality flag (default 90).
The resolution of the output images can be controlled using the --image-ppi flag.
Common pixels per inch values are 100 (screen), 150-300 (print), 600 (art). If
not specified, the PPI of the output images is 100.

```
unipdf optimize [FLAG]... INPUT_FILES...

Flags:
-P, --image-ppi float     output images pixels per inch (default 100)
-q, --image-quality int   output JPEG image quality (default 90)
-O, --overwrite           overwrite input files
-p, --password string     file password
-r, --recursive           search PDF files in subdirectories
-t, --target-dir string   output directory

Examples:
unipdf optimize file_1.pdf file_n.pdf
unipdf optimize -O file_1.pdf file_n.pdf
unipdf optimize -O -r file_1.pdf file_n.pdf dir_1 dir_n
unipdf optimize -t out_dir file_1.pdf file_n.pdf dir_1 dir_n
unipdf optimize -t out_dir -r file_1.pdf file_n.pdf dir_1 dir_n
unipdf optimize -t out_dir -r -q 75 file_1.pdf file_n.pdf dir_1 dir_n
unipdf optimize -t out_dir -r -q 75 -P 100 file_1.pdf file_n.pdf dir_1 dir_n
unipdf optimize -t out_dir -r -q 75 -P 100 -p pass file_1.pdf file_n.pdf dir_1 dir_n
```

#### Rotate

Rotate PDF file pages by a specified angle. The angle argument is specified in
degrees and it must be a multiple of 90.

```
unipdf rotate [FLAG]... INPUT_FILE ANGLE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages to rotate
-p, --password string      PDF file password

Examples:
unipdf rotate input_file.pdf 90
unipdf rotate -- input_file.pdf -270
unipdf rotate -o output_file.pdf input_file.pdf 90
unipdf rotate -o output_file.pdf -P 1-3 input_file.pdf 90
unipdf rotate -o output_file.pdf -P 1-3 -p pass input_file.pdf 90

Pages flag example: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be rotated, while
page number 5 is skipped.
```

#### Watermark

Add watermark images to PDF files.

```
unipdf watermark [FLAG]... INPUT_FILE WATERMARK_IMAGE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages on which to add watermark
-p, --password string      PDF file password

Examples:
unipdf watermark input_file.pdf watermark.png
unipdf watermark -o output file.png input_file.pdf watermark.png
unipdf watermark -o output file.png -P 1-3 input_file.pdf watermark.png
unipdf watermark -o output file.png -P 1-3 -p pass input_file.pdf watermark.png

Pages flag example: 1-3,4,6-7
Watermark will only be applied to pages 1,2,3 (1-3), 4 and 6,7 (6-7), while
page number 5 is skipped.
```

#### Grayscale

Convert PDF files to grayscale.

```
unipdf grayscale [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   Output file
-P, --pages string         Pages to convert to grayscale
-p, --password string      PDF file password

Examples:
unipdf grayscale input_file.pdf
unipdf grayscale -o output_file input_file.pdf
unipdf grayscale -o output_file -P 1-3 input_file.pdf
unipdf grayscale -o output_file -P 1-3 -p pass input_file.pdf

Pages flag example: 1-3,4,6-7
Only pages 1,2,3 (1-3), 4 and 6,7 (6-7) will be converted to grayscale, while
page number 5 is skipped.
```

#### Info

Outputs file information. Also does some basic validation.

```
unipdf info [FLAG]... INPUT_FILE

Flags:
-p, --password string   PDF file password

Examples:
unipdf info input_file.pdf
unipdf info -p pass input_file.pdf
```

#### Extract text

Extracts PDF text. The extracted text is always printed to STDOUT.

```
unipdf extract text [FLAG]... INPUT_FILE

Flags:
-P, --pages string           Pages to extract text from
-p, --user-password string   Input file password

Examples:
unipdf extract text input_file.pdf
unipdf extract text -P 1-3 input_file.pdf
unipdf extract text -P 1-3 -p pass input_file.pdf

Pages flag example: 1-3,4,6-7
Text will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7), while
page number 5 is skipped.
```

#### Extract images

Extracts PDF images. The images are extracted in a ZIP file and saved at the
destination specified by the --output-file parameter. If no output file is
specified, the ZIP archive is saved in the same directory as the input file.

```
unipdf extract [FLAG]... INPUT_FILE

Flags:
-S, --include-inline-stencil-masks   Include inline stencil masks
-o, --output-file string             Output file
-P, --pages string                   Pages to extract images from
-p, --password string                Input file password

Examples:
unipdf extract images input_file.pdf
unipdf extract images -o images.zip input_file.pdf
unipdf extract images -P 1-3 -p pass -o images.zip input_file.pdf

Pages flag example: 1-3,4,6-7
Images will only be extracted from pages 1,2,3 (1-3), 4 and 6,7 (6-7), while
page number 5 is skipped.
```

#### Search

Search text in PDF files.

```
unipdf search [FLAG]... INPUT_FILE TEXT

Flags:
-p, --password string   PDF file password

Examples:
unipdf search input_file.pdf text_to_search
unipdf search -p pass input_file.pdf text_to_search
```

#### Form Export

Export JSON representation of form fields.

By default, the resulting JSON content is printed to STDOUT. The output can be
saved to a file by using the --output-file flag.

```
unipdf form export [FLAG]... INPUT_FILE

Flags:
-o, --output-file string   output file

Examples:
unipdf form export in_file.pdf
unipdf form export in_file.pdf > out_file.json
unipdf form export -o out_file.json in_file.pdf
```

#### Form Fill

Fill form fields from JSON file.

The field values specified in the JSON file template are used to fill the form
fields in the input PDF files. In addition, the output file form fields can be
flattened by using the --flatten flag. The flattening process makes the form
fields of the output files read-only by appending the form field annotation
XObject Form data to the page content stream, thus making it part of the page
contents.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_filled" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the filled output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

```
unipdf form fill [FLAG]... JSON_FILE INPUT_FILES...

Flags:
-f, --flatten             flatten form annotations
-O, --overwrite           overwrite input files
-p, --password string     input file password
-r, --recursive           search PDF files in subdirectories
-t, --target-dir string   output directory

Examples:
unipdf form fill fields.json file_1.pdf file_n.pdf
unipdf form fill -O fields.json file_1.pdf file_n.pdf
unipdf form fill -O -r -f fields.json file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fill -t out_dir fields.json file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fill -t out_dir -r fields.json file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fill -t out_dir -r -p pass fields.json file_1.pdf file_n.pdf dir_1 dir_n
```
#### FDF Merge

Fill form fields from FDF file.

The field values specified in the FDF file template are used to fill the form
fields in the input PDF files. In addition, the output file form fields can be
flattened by using the --flatten flag. The flattening process makes the form
fields of the output files read-only by appending the form field annotation
XObject Form data to the page content stream, thus making it part of the page
contents.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_filled" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the filled output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

```
Usage:
unipdf form fdfmerge [FLAG]... FDF_FILE INPUT_FILES...

Flags:
-f, --flatten             flatten form annotations
-O, --overwrite           overwrite input files
-p, --password string     input file password
-r, --recursive           search PDF files in subdirectories
-t, --target-dir string   output directory

Examples:
unipdf form fdfmerge fields.fdf file_1.pdf file_n.pdf
unipdf form fdfmerge -O fields.fdf file_1.pdf file_n.pdf
unipdf form fdfmerge -O -r -f fields.fdf file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fdfmerge -t out_dir fields.fdf file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fdfmerge -t out_dir -r fields.fdf file_1.pdf file_n.pdf dir_1 dir_n
unipdf form fdfmerge -t out_dir -r -p pass fields.fdf file_1.pdf file_n.pdf dir_1 dir_n
```

#### Form Flatten

Flatten PDF file form annotations.

The flattening process makes the form fields of the output files read-only by
appending the form field annotation XObject Form data to the page content
stream, thus making it part of the page contents.

The command can take multiple files and directories as input parameters.
By default, each PDF file is saved in the same location as the original file,
appending the "_flattened" suffix to the file name. Use the --overwrite flag
to overwrite the original files.
In addition, the flattened output files can be saved to a different directory
by using the --target-dir flag.
The command can search for PDF files inside the subdirectories of the
specified input directories by using the --recursive flag.

```
unipdf form flatten [FLAG]... INPUT_FILES...

Flags:
-O, --overwrite           overwrite input files
-p, --password string     input file password
-r, --recursive           search PDF files in subdirectories
-t, --target-dir string   output directory

Examples:
unipdf form flatten file_1.pdf file_n.pdf
unipdf form flatten -O file_1.pdf file_n.pdf
unipdf form flatten -O -r file_1.pdf file_n.pdf dir_1 dir_n
unipdf form flatten -t out_dir file_1.pdf file_n.pdf dir_1 dir_n
unipdf form flatten -t out_dir -r file_1.pdf file_n.pdf dir_1 dir_n
unipdf form flatten -t out_dir -r -p pass file_1.pdf file_n.pdf dir_1 dir_n
```

#### Render

Render PDF pages to image targets.

The rendered image files are saved in a ZIP file, at the location specified
by the --output-file parameter. If no output file is specified, the ZIP file
is saved in the same directory as the input file.

The format of the rendered image files can be specified using
the --image-format flag (default jpeg). The quality of the image files can be
configured through the --image-quality flag (default 100, only applies to
JPEG images).
```
unipdf render [FLAG]... INPUT_FILE

Flags:
-f, --image-format string   format of the output images (default "jpeg")
-q, --image-quality int     quality of the output images (default 100)
-o, --output-file string    output file
-P, --pages string          pages to render from the input file
-p, --password string       input file password

Examples:
unipdf render in_file.pdf
unipdf render -o images.zip in_file.pdf
unipdf render -o images.zip -P 1-3 in_file.pdf
unipdf render -o images.zip -P 1-3 -p pass in_file.pdf
unipdf render -o images.zip -P 1-3 -p pass -f jpeg -q 100 in_file.pdf

Pages flag example: 1-3,4,6-7
Images will only be rendered for pages 1,2,3 (1-3), 4 and 6,7 (6-7), while
page number 5 is skipped.

Supported image formats:
  - jpeg (default)
  - png
```

## License

The application is licensed under the same conditions as the
[UniPDF](https://github.com/unidoc/unipdf) library.
It has a dual license, a commercial one suitable for closed source projects
and an AGPL license that can be used in open source software.

Depending on your needs, you can choose one of them and follow its policies.
A detail of the policies and agreements for each license type are available in
the [LICENSE.COMMERCIAL](LICENSE.COMMERCIAL) and [LICENSE.AGPL](LICENSE.AGPL)
files.

Please see [pricing](http://unidoc.io/pricing) to purchase a commercial
[UniPDF](https://github.com/unidoc/unipdf) license or contact sales at
sales@unidoc.io for more info.

If you have a license for [UniPDF](https://github.com/unidoc/unipdf), you can
set it through the UNIDOC_LICENSE_FILE and UNIDOC_LICENSE_CUSTOMER environment
variables.

```
export UNIDOC_LICENSE_FILE="PATH_TO_LICENSE_FILE"
export UNIDOC_LICENSE_CUSTOMER="CUSTOMER_NAME"
```
