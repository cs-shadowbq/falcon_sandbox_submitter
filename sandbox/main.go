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
	FalconClientId string

	FalconClientSecret string

	ClientCloud string

	Filename string

	SandboxEnvId int32

	NetworkSettings string

	ActionScript string
}

func (sub CmdSubmission) SubmitFile(verbose bool) {

	client, err := falcon.NewClient(&falcon.ApiConfig{
		ClientId:     sub.FalconClientId,
		ClientSecret: sub.FalconClientSecret,
		Cloud:        falcon.Cloud(sub.ClientCloud),
		Context:      context.Background(),
	})
	if err != nil {
		panic(err)
	}

	// get full path of the file
	full_filename, err := filepath.Abs(sub.Filename)

	// get file name from the full_filename
	var filename = filepath.Base(full_filename)

	if err != nil {
		panic(err)
	}

	// open the file
	file_handler, err := os.Open(full_filename)
	if err != nil {
		panic(err)
	}
	// read the file contents of fileh into fileContents
	fileReadCloser := runtime.NamedReader(filename, file_handler)

	if verbose {
		fmt.Printf("Uploading file %s \n", full_filename)
		fmt.Printf("Client ID: %s\n", sub.FalconClientId)
		fmt.Printf("Client Secret: %s\n", sub.FalconClientSecret)
		fmt.Printf("Cloud: %s\n", sub.ClientCloud)
		fmt.Printf("Full filename: %s\n", full_filename)
	}
	//type UploadSampleV2Params struct

	var submissionParams = falconx_sandbox.UploadSampleV2Params{
		Context:  context.Background(),
		FileName: filename,
		Sample:   fileReadCloser,
	}

	if verbose {
		fmt.Printf("Building payload\n")
		fmt.Printf("File name: %s\n", submissionParams.FileName)
	}

	upload, err := client.FalconxSandbox.UploadSampleV2(&submissionParams)

	if err != nil {
		panic(err)
	}

	// Print uploaded
	if verbose {
		fmt.Printf("Uploaded file %s \n", full_filename)
	}
	payload := upload.GetPayload()

	// fmt.Printf("Payload: %v\n", payload)
	if err = falcon.AssertNoError(payload.Errors); err != nil {
		panic(err)
	}
	if verbose {
		fmt.Printf("Uploaded file %s with ID %s\n", filename, *payload.Resources[0].Sha256)
		fmt.Printf("Submitting file %s for analysis to env %d\n", filename, sub.SandboxEnvId)
	}

	var submitParams = falconx_sandbox.SubmitParams{
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
		panic(err)
	}

	submitPayload := submit.GetPayload()
	// fmt.Printf("Payload: %v\n", submitPayload)
	if err = falcon.AssertNoError(submitPayload.Errors); err != nil {
		panic(err)
	}

	//fmt.Printf("Submitted file %s with ID %s\n", filename, *submitPayload.Resources)
	// Print submitted

}
