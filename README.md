# BioIDWebService-PhotoVerify-gRPC-Sample-in-Go

This repository contains a very simple example command line implementation of the **[PhotoVerify gRPC API of the BioID WebService](https://developer.bioid.com/bws/grpc/photoverify)** in Go.

One reference image and one or two live images are required to perform a photoverify request.

## Requirements

Before you can use the PhotoVerify API, you need to create a BWS Client ID and key in the BWS Portal.
You can request trial access on https://bwsportal.bioid.com/register

## Usage

1. Build or download prebuild executable
2. Execute following command to perform a photo verify with two images:

```
 .\BioIDWebService-PhotoVerify-gRPC-Sample-In-Go.exe -BWSClientID <BWSClientID> -BWSKey <BWSKey> -photo example_images\photo.jpg -image1 example_images\testimage1.jpg -image2 example_images\testimage2.jpg
```

Example Output:
Token generation took: 0s
Client creation took: 519.8µs
File reading took: 536.2µs
Photo verification took: 921.5775ms
Verification Status: SUCCEEDED
Verification Errors: [error_code:"RejectedByPassiveLiveDetection" message:"At least one of the live images seem not to be recorded from a live person." error_code:"RejectedByPassiveLiveDetection" message:"At least one of the live images seem not to be recorded from a live person."]
Verification ImageProperties: [faces:{left_eye:{x:716.6898478956991 y:295.90328987385425} right_eye:{x:583.7671688344084 y:295.6642189519129} texture_liveness_score:0.36677824126349556} faces:{left_eye:{x:493.984682522074 y:253.9819667801256} right_eye:{x:382.02615043197886 y:250.22897896148734} texture_liveness_score:0.46766824192470974}]
Verification PhotoProperties: <nil>
Verification Level: NOT_RECOGNIZED
Verification Score: 0
Verification Live: false
Verification LivenessScore: 0.46766824192470974
Total execution time: 923.6721ms

### Available command line parameter

```
./BioIDWebService-PhotoVerify-gRPC-Sample-In-Go --help
  -BWSClientID string
        BioIDWebService ClientID
  -BWSKey string
        BioIDWebService Key
  -image1 string
        1st live image
  -image2 string
        2nd live image (optional)
  -photo string
        reference photo image
```

## Clone and build the project

```
$ git clone https://github.com/danielchristianschroeter/BioIDWebService-PhotoVerify-gRPC-Sample-In-Go
$ cd BioIDWebService-PhotoVerify-gRPC-Sample-In-Go
$ go build .
```

## Generate or Update Go Code from Protocol Buffers file (`.proto`) for BioID Web Service

## Prerequisites & Installation

1. Download Protocol Buffers
   Download the Protocol Buffers binary from the Protocol Buffers GitHub releases page https://github.com/protocolbuffers/protobuf/releases.

2. Install Go Plugins for Protocol Buffers
   Run the following commands to install the necessary Go plugins:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

3. Generating Go Code
   To generate the Go code (bws_grpc.pb.go and bws.pb.go) from your bws.proto file, execute the following command with the Protocol Buffers executable:

```
protoc --proto_path=<path_to_your_proto_folder> \
--go_out=<path_to_your_output_folder> \
--go_opt=paths=source_relative \
--go-grpc_out=<path_to_your_output_folder> \
--go-grpc_opt=paths=source_relative \
<path_to_your_proto_file>
```

Replace <path_to_your_proto_folder>, <path_to_your_output_folder>, and <path_to_your_proto_file> with the actual paths relevant to your project.

Ensure the protoc binary directory is included in your system's PATH environment variable. If not, navigate to the extracted bin directory and execute protoc from there.
This command will generate bws_grpc.pb.go and bws.pb.go files in the specified output directory, with paths relative to the source directory.
