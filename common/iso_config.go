//go:generate struct-markdown

package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/packer/template/interpolate"
)

// By default, Packer will symlink, download or copy image files to the Packer
// cache into a "`hash($iso_url+$iso_checksum).$iso_target_extension`" file.
// Packer uses [hashicorp/go-getter](https://github.com/hashicorp/go-getter) in
// file mode in order to perform a download.
//
// go-getter supports the following protocols:
//
// * Local files
// * Git
// * Mercurial
// * HTTP
// * Amazon S3
//
//
// \~&gt; On windows - when referencing a local iso - if packer is running
// without symlinking rights, the iso will be copied to the cache folder. Read
// [Symlinks in Windows 10
// !](https://blogs.windows.com/buildingapps/2016/12/02/symlinks-windows-10/)
// for more info.
//
// Examples:
// go-getter can guess the checksum type based on `iso_checksum` len.
//
// ``` json
// {
//   "iso_checksum": "946a6077af6f5f95a51f82fdc44051c7aa19f9cfc5f737954845a6050543d7c2",
//   "iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"
// }
// ```
//
// ``` json
// {
//   "iso_checksum": "file:ubuntu.org/..../ubuntu-14.04.1-server-amd64.iso.sum",
//   "iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"
// }
// ```
//
// ``` json
// {
//   "iso_checksum": "file://./shasums.txt",
//   "iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"
// }
// ```
//
// ``` json
// {
//   "iso_checksum": "file:./shasums.txt",
//   "iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"
// }
// ```
//
type ISOConfig struct {
	// The checksum for the ISO file or virtual hard drive file. The algorithm
	// to use when computing the checksum will be determined automatically
	// based on `iso_checksum` length. `iso_checksum` can be also be a file or
	// an URL, in which case iso_checksum must be prefixed with `file:`; the
	// go-getter will download it and use the first hash found.
	//
	// `iso_checksum` can be set to `"none"` if you want no checksumming
	// operation to be run.
	ISOChecksum string `mapstructure:"iso_checksum" required:"true"`
	// A URL to the ISO containing the installation image or virtual hard drive
	// (VHD or VHDX) file to clone.
	RawSingleISOUrl string `mapstructure:"iso_url" required:"true"`
	// Multiple URLs for the ISO to download. Packer will try these in order.
	// If anything goes wrong attempting to download or while downloading a
	// single URL, it will move on to the next. All URLs must point to the same
	// file (same checksum). By default this is empty and `iso_url` is used.
	// Only one of `iso_url` or `iso_urls` can be specified.
	ISOUrls []string `mapstructure:"iso_urls"`
	// The path where the iso should be saved after download. By default will
	// go in the packer cache, with a hash of the original filename and
	// checksum as its name.
	TargetPath string `mapstructure:"iso_target_path"`
	// The extension of the iso file after download. This defaults to `iso`.
	TargetExtension string `mapstructure:"iso_target_extension"`
}

func (c *ISOConfig) Prepare(ctx *interpolate.Context) (warnings []string, errs []error) {
	if len(c.ISOUrls) != 0 && c.RawSingleISOUrl != "" {
		errs = append(
			errs, errors.New("Only one of iso_url or iso_urls must be specified"))
		return
	}

	if c.RawSingleISOUrl != "" {
		// make sure only array is set
		c.ISOUrls = append([]string{c.RawSingleISOUrl}, c.ISOUrls...)
		c.RawSingleISOUrl = ""
	}

	if len(c.ISOUrls) == 0 {
		errs = append(
			errs, errors.New("One of iso_url or iso_urls must be specified"))
		return
	}
	if c.TargetExtension == "" {
		c.TargetExtension = "iso"
	}
	c.TargetExtension = strings.ToLower(c.TargetExtension)

	// Warnings
	if c.ISOChecksum == "none" {
		warnings = append(warnings,
			"A checksum of 'none' was specified. Since ISO files are so big,\n"+
				"a checksum is highly recommended.")
		return warnings, errs
	} else if c.ISOChecksum == "" {
		errs = append(errs, fmt.Errorf("A checksum must be specified"))
	}

	if strings.HasSuffix(strings.ToLower(c.ISOChecksum), ".iso") {
		errs = append(errs, fmt.Errorf("Error parsing checksum:"+
			" .iso is not a valid checksum ending"))
	}

	return warnings, errs
}
