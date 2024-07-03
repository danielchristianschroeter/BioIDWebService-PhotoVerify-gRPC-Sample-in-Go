package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "BioIDWebService-PhotoVerify-gRPC-Sample-In-Go/proto"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// Used for build version information
var version = "development"

// Command line flag variables
var (
	BWSHost      = "grpc.bws-eu.bioid.com"
	BWSCAFile    = "ISRG_Root_X1.pem"
	BWSClientID  string
	BWSKey       string
	photoPath    string
	image1Path   string
	image2Path   string
)

// Generates a JWT token for authentication
func generateToken(BWSClientID, BWSKey string, expireMinutes int) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(BWSKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %w", err)
	}

	signingKey := keyBytes
	claims := jwt.MapClaims{
		"iss": BWSClientID,
		"sub": BWSClientID,
		"aud": "BWS",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(signingKey)
}

// Creates an authenticated gRPC client
func createAuthenticatedClient() (*grpc.ClientConn, pb.BioIDWebServiceClient, error) {
	flag.Parse()
	var opts []grpc.DialOption

	creds, err := credentials.NewClientTLSFromFile(BWSCAFile, BWSHost)
	if err != nil {
		log.Fatalf("Failed to create TLS credentials: %v", err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.NewClient(BWSHost + ":443", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewBioIDWebServiceClient(conn)
	return conn, client, nil
}

// Reads file content and stores it in the provided byte slice
func readFileContent(filePath string, wg *sync.WaitGroup, content *[]byte, err *error) {
	defer wg.Done()
	*content, *err = os.ReadFile(filePath)
}

// Initialize command line flags
func init() {
	flag.StringVar(&BWSClientID, "BWSClientID", "", "BioIDWebService ClientID")
	flag.StringVar(&BWSKey, "BWSKey", "", "BioIDWebService Key")
	flag.StringVar(&photoPath, "photo", "", "reference photo image")
	flag.StringVar(&image1Path, "image1", "", "1st live image")
	flag.StringVar(&image2Path, "image2", "", "2nd live image (optional)")
}

func main() {
	startTime := time.Now() // Record the start time of the entire process
	log.SetFlags(0)

	flag.Usage = func() {
		log.Println("BioIDWebService PhotoVerify gRPC Sample in Go. Version: " + version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(BWSClientID) == 0 || len(BWSKey) == 0 || len(photoPath) == 0 || len(image1Path) == 0 {
		log.Fatal("Usage: -BWSClientID <BWSClientID> -BWSKey <BWSKey> -photo <photo> -image1 <image1> [-image2 <image2>]")
	}

	tokenStart := time.Now() // Measure time to generate the token
	token, err := generateToken(BWSClientID, BWSKey, 5)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	tokenDuration := time.Since(tokenStart)
	fmt.Printf("Token generation took: %v\n", tokenDuration)
	
	clientCreationStart := time.Now() // Measure time to create the authenticated client
	conn, client, err := createAuthenticatedClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()
	clientCreationDuration := time.Since(clientCreationStart)
	fmt.Printf("Client creation took: %v\n", clientCreationDuration)
	
	fileReadStart := time.Now() // Measure time to read all files
	var wg sync.WaitGroup
	wg.Add(2)
	if image2Path != "" {
		wg.Add(1)
	}

	var liveImage1, liveImage2, photo []byte
	var err1, err2, err3 error

	go readFileContent(image1Path, &wg, &liveImage1, &err1)
	go readFileContent(photoPath, &wg, &photo, &err3)
	if image2Path != "" {
		go readFileContent(image2Path, &wg, &liveImage2, &err2)
	}

	wg.Wait()
	fmt.Printf("File reading took: %v\n", time.Since(fileReadStart))

	if err1 != nil {
		log.Fatalf("Failed to read live image 1: %v", err1)
	}
	if err3 != nil {
		log.Fatalf("Failed to read photo: %v", err3)
	}

	// Create request with the images
	request := &pb.PhotoVerifyRequest{
		LiveImages: []*pb.ImageData{
			{Image: liveImage1},
		},
		Photo: photo,
	}
	if image2Path != "" {
		if err2 != nil {
			log.Fatalf("Failed to read live image 2: %v", err2)
		}
		request.LiveImages = append(request.LiveImages, &pb.ImageData{Image: liveImage2})
	}

	photoVerifyStart := time.Now() // Measure time to verify photo
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	response, err := client.PhotoVerify(ctx, request)
	if err != nil {
		log.Fatalf("Failed to verify photo: %v", err)
	}
	photoVerifyDuration := time.Since(photoVerifyStart)
	fmt.Printf("Photo verification took: %v\n", photoVerifyDuration)

	// Print the results
	fmt.Printf("Verification Status: %v\n", response.Status)
	fmt.Printf("Verification Errors: %v\n", response.Errors)
	fmt.Printf("Verification ImageProperties: %v\n", response.ImageProperties)
	fmt.Printf("Verification PhotoProperties: %v\n", response.PhotoProperties)
	fmt.Printf("Verification Level: %v\n", response.VerificationLevel)
	fmt.Printf("Verification Score: %v\n", response.VerificationScore)
	fmt.Printf("Verification Live: %v\n", response.Live)
	fmt.Printf("Verification LivenessScore: %v\n", response.LivenessScore)
	
	fmt.Printf("Total execution time: %v\n", time.Since(startTime)) // Print the total execution time
}