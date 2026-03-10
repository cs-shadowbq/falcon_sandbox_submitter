/*
Copyright © 2024-2026 CrowdStrike - Scott MacGregor scott.macgregor@crowdstrike.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package sandbox

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client/falconx_sandbox"
	"github.com/crowdstrike/gofalcon/falcon/models"
	"github.com/go-openapi/runtime"
)

type CmdSubmission struct {
	FalconClientId     string
	FalconClientSecret string
	ClientCloud        string
	Filename           string
	SandboxEnvId       int32
	NetworkSettings    string
	ActionScript       string
}

// maskSecret returns the first 4 characters of a secret followed by asterisks, or
// all asterisks if the secret is shorter than 4 characters.
func maskSecret(s string) string {
	if len(s) <= 4 {
		return "********"
	}
	return s[:4] + "********"
}

func (sub CmdSubmission) SubmitFile(verbose bool) error {
	cloud, err := falcon.CloudValidate(sub.ClientCloud)
	if err != nil {
		return fmt.Errorf("invalid cloud region %q: %w", sub.ClientCloud, err)
	}

	client, err := falcon.NewClient(&falcon.ApiConfig{
		ClientId:     sub.FalconClientId,
		ClientSecret: sub.FalconClientSecret,
		Cloud:        cloud,
		Context:      context.Background(),
	})
	if err != nil {
		return fmt.Errorf("failed to create Falcon client: %w", err)
	}

	fullFilename, err := filepath.Abs(sub.Filename)
	if err != nil {
		return fmt.Errorf("failed to resolve file path %q: %w", sub.Filename, err)
	}
	filename := filepath.Base(fullFilename)

	fileHandler, err := os.Open(fullFilename)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", fullFilename, err)
	}
	defer fileHandler.Close()

	fileReadCloser := runtime.NamedReader(filename, fileHandler)

	if verbose {
		fmt.Printf("Uploading file %s\n", fullFilename)
		fmt.Printf("Client ID:     %s\n", sub.FalconClientId)
		fmt.Printf("Client Secret: %s\n", maskSecret(sub.FalconClientSecret))
		fmt.Printf("Cloud:         %s\n", cloud.String())
	}

	submissionParams := falconx_sandbox.UploadSampleV2Params{
		Context:  context.Background(),
		FileName: filename,
		Sample:   fileReadCloser,
	}

	if verbose {
		fmt.Printf("Building payload for file: %s\n", submissionParams.FileName)
	}

	upload, err := client.FalconxSandbox.UploadSampleV2(&submissionParams)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	payload := upload.GetPayload()
	if err = falcon.AssertNoError(payload.Errors); err != nil {
		return fmt.Errorf("upload API error: %w", err)
	}

	if verbose {
		fmt.Printf("Uploaded file %s with SHA256 %s\n", filename, *payload.Resources[0].Sha256)
		fmt.Printf("Submitting file %s for analysis to env %d\n", filename, sub.SandboxEnvId)
	}

	submitParams := falconx_sandbox.SubmitParams{
		Context: context.Background(),
		Body: &models.FalconxSubmissionParametersV1{
			Sandbox: []*models.FalconxSandboxParametersV1{
				{
					Sha256:          *payload.Resources[0].Sha256,
					EnvironmentID:   sub.SandboxEnvId,
					SubmitName:      filename,
					ActionScript:    sub.ActionScript,
					NetworkSettings: sub.NetworkSettings,
				},
			},
		},
	}

	submit, err := client.FalconxSandbox.Submit(&submitParams)
	if err != nil {
		return fmt.Errorf("submission failed: %w", err)
	}

	submitPayload := submit.GetPayload()
	if err = falcon.AssertNoError(submitPayload.Errors); err != nil {
		return fmt.Errorf("submission API error: %w", err)
	}

	if verbose {
		fmt.Printf("Successfully submitted %s (env %d)\n", filename, sub.SandboxEnvId)
	}

	return nil
}
