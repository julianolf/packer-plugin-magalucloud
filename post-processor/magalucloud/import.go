// Copyright (c) Juliano Fernandes 2026
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package magalucloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

const (
	BuilderId    = "julianolf.post-processor.magalucloud-import"
	ImageHeader  = "QFI\xfb"
	ImageFormat  = "qcow2"
	ImageFileExt = "." + ImageFormat
	WaitInterval = 15 * time.Second
)

type Region string

type URL struct {
	API client.MgcUrl
	OBJ objectstorage.Endpoint
}

var Regions = map[Region]URL{
	"br-se1": {API: client.BrSe1, OBJ: objectstorage.BrSe1},
	"br-ne1": {API: client.BrNe1, OBJ: objectstorage.BrNe1},
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	APIKey              string                 `mapstructure:"api_key"`
	AccessKey           string                 `mapstructure:"access_key"`
	SecretKey           string                 `mapstructure:"secret_key"`
	Region              Region                 `mapstructure:"region"`
	Bucket              string                 `mapstructure:"bucket"`
	Filename            string                 `mapstructure:"filename"`
	Expires             time.Duration          `mapstructure:"expires"`
	URL                 client.MgcUrl          `mapstructure:"url"`
	Endpoint            objectstorage.Endpoint `mapstructure:"endpoint"`
	ImageName           string                 `mapstructure:"image_name"`
	Platform            compute.Platform       `mapstructure:"platform"`
	Architecture        compute.Architecture   `mapstructure:"architecture"`
	License             compute.License        `mapstructure:"license"`
	Version             string                 `mapstructure:"version"`
	Description         string                 `mapstructure:"description"`
	UEFI                bool                   `mapstructure:"uefi"`
	SkipCleanup         bool                   `mapstructure:"skip_cleanup"`

	ctx interpolate.Context
}

type Importer struct {
	config        Config
	compute       *compute.VirtualMachineClient
	objectstorage *objectstorage.ObjectStorageClient
}

func (i *Importer) ConfigSpec() hcldec.ObjectSpec {
	return i.config.FlatMapstructure().HCL2Spec()
}

func (i *Importer) Configure(raws ...any) error {
	err := config.Decode(
		&i.config,
		&config.DecodeOpts{
			PluginType:         BuilderId,
			Interpolate:        true,
			InterpolateContext: &i.config.ctx,
			InterpolateFilter:  &interpolate.RenderFilter{Exclude: []string{}},
		},
		raws...,
	)
	if err != nil {
		return err
	}

	url, ok := Regions[i.config.Region]
	if !ok {
		return fmt.Errorf("Invalid region: %s", i.config.Region)
	}
	if i.config.URL == "" {
		i.config.URL = url.API
	}
	if i.config.Endpoint == "" {
		i.config.Endpoint = url.OBJ
	}

	if i.config.Expires == time.Duration(0) {
		i.config.Expires = time.Hour
	}

	cli := client.NewMgcClient(
		client.WithAPIKey(i.config.APIKey),
		client.WithBaseURL(i.config.URL),
	)
	obj, err := objectstorage.New(
		cli,
		i.config.AccessKey,
		i.config.SecretKey,
		objectstorage.WithEndpoint(i.config.Endpoint),
	)
	if err != nil {
		return err
	}

	i.objectstorage = obj
	i.compute = compute.New(cli)

	return nil
}

func (i *Importer) PostProcess(ctx context.Context, ui packersdk.Ui, artifact packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	sURL, err := i.uploadImage(ctx, ui, artifact)
	if err != nil {
		return nil, false, false, err
	}

	_, err = i.importImage(ctx, ui, sURL)
	if err != nil {
		return nil, false, false, err
	}

	err = i.cleanup(ctx, ui)
	if err != nil {
		return nil, false, false, err
	}

	return artifact, true, true, nil
}

func findImage(artifact packersdk.Artifact) (string, error) {
	source := ""
	for _, path := range artifact.Files() {
		if strings.HasSuffix(path, ImageFileExt) {
			log.Printf("[DEBUG] Image file %s found in artifact from builder", path)
			source = path
			break
		}
	}
	if source == "" {
		return "", fmt.Errorf("No %s image file found in artifact from builder", ImageFormat)
	}

	return source, nil
}

func validateImage(file *os.File) error {
	header := make([]byte, 4)
	_, err := file.Read(header)
	if err != nil {
		return err
	}
	if string(header) != ImageHeader {
		return fmt.Errorf("Invalid %s image header %s", ImageFormat, header)
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}

func (i *Importer) uploadImage(ctx context.Context, ui packersdk.Ui, artifact packersdk.Artifact) (string, error) {
	filename, err := findImage(artifact)
	if err != nil {
		return "", err
	}

	if i.config.Filename == "" {
		i.config.Filename = filename
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = validateImage(file)
	if err != nil {
		return "", err
	}

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := stat.Size()

	ui.Sayf("Uploading %s to %s bucket", i.config.Filename, i.config.Bucket)

	err = i.objectstorage.Objects().Upload(
		ctx,
		i.config.Bucket,
		i.config.Filename,
		file,
		size,
		"application/octet-stream",
	)
	if err != nil {
		return "", err
	}

	ui.Sayf("Generating presigned URL for %s/%s/%s", i.config.Endpoint, i.config.Bucket, i.config.Filename)

	sURL, err := i.objectstorage.Presigner().GeneratePresignedURL(
		ctx,
		http.MethodGet,
		i.config.Bucket,
		i.config.Filename,
		i.config.Expires,
		url.Values{},
	)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] Presigned URL %s", sURL)
	return sURL.String(), nil
}

func (i *Importer) importImage(ctx context.Context, ui packersdk.Ui, sURL string) (string, error) {
	req := compute.CreateCustomImageRequest{
		Name:         i.config.ImageName,
		Platform:     i.config.Platform,
		Architecture: i.config.Architecture,
		License:      i.config.License,
		URL:          sURL,
		Version:      helpers.StrPtr(i.config.Version),
		Description:  helpers.StrPtr(i.config.Description),
		UEFI:         helpers.BoolPtr(i.config.UEFI),
	}

	data, _ := json.Marshal(req)
	log.Printf("[DEBUG] Create custom image data %s", string(data))

	ui.Sayf("Importing image %s using presigned URL", i.config.ImageName)

	id, err := i.compute.CustomImages().Create(ctx, req)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] Custom image ID %s", id)

	ticker := time.NewTicker(WaitInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			image, err := i.compute.CustomImages().Get(ctx, id)
			if err != nil {
				return "", err
			}

			switch image.Status {
			case compute.ImageStatusActive:
				ui.Sayf("Finished importing image with ID %s", id)
				return id, nil
			case compute.ImageStatusCreating, compute.ImageStatusPending, compute.ImageStatusImporting:
			default:
				return "", fmt.Errorf("Invalid image status: %s", image.Status)
			}
		}
	}
}

func (i *Importer) cleanup(ctx context.Context, ui packersdk.Ui) error {
	if !i.config.SkipCleanup {
		ui.Sayf("Deleting %s from %s bucket", i.config.Filename, i.config.Bucket)

		return i.objectstorage.Objects().Delete(
			ctx,
			i.config.Bucket,
			i.config.Filename,
			&objectstorage.DeleteOptions{},
		)
	}

	return nil
}
