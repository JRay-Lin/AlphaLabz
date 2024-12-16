package pdf

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// CompressionLevel represents the compression quality settings
type CompressionLevel string

const (
	CompressMax    CompressionLevel = "max"    // Maximum compression, lowest quality
	CompressHigh   CompressionLevel = "high"   // High compression
	CompressMedium CompressionLevel = "medium" // Balanced compression
	CompressLow    CompressionLevel = "low"    // Light compression, better quality
)

// CompressPDF compresses a PDF file using pdfcpu
func CompressPDF(inputPath, outputPath string, level CompressionLevel) error {
	// Open input file
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inFile.Close()

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Create default configuration
	conf := model.NewDefaultConfiguration()

	// Adjust configuration based on compression level
	switch level {
	case CompressMax:
		conf.Cmd = model.OPTIMIZE
		conf.ValidationMode = model.ValidationStrict
	case CompressHigh:
		conf.Cmd = model.OPTIMIZE
		conf.ValidationMode = model.ValidationRelaxed
	case CompressMedium:
		conf.Cmd = model.OPTIMIZE
		conf.ValidationMode = model.ValidationRelaxed
	case CompressLow:
		conf.Cmd = model.OPTIMIZE
		conf.ValidationMode = model.ValidationRelaxed
	default:
		return fmt.Errorf("invalid compression level. Use: max, high, medium, or low")
	}

	// Perform the compression
	if err := api.Optimize(inFile, outFile, conf); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	return nil
}

// GetCompressionStats returns statistics about the compression
func GetCompressionStats(originalPath, compressedPath string) (*CompressionStats, error) {
	originalSize, err := getFileSize(originalPath)
	if err != nil {
		return nil, fmt.Errorf("error getting original file size: %v", err)
	}

	compressedSize, err := getFileSize(compressedPath)
	if err != nil {
		return nil, fmt.Errorf("error getting compressed file size: %v", err)
	}

	savings := float64(originalSize-compressedSize) / float64(originalSize) * 100

	return &CompressionStats{
		OriginalSize:   originalSize,
		CompressedSize: compressedSize,
		SavingsPercent: savings,
	}, nil
}

// CompressionStats holds information about the compression results
type CompressionStats struct {
	OriginalSize   int64
	CompressedSize int64
	SavingsPercent float64
}

// getFileSize returns the size of a file in bytes
func getFileSize(filepath string) (int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	return size, nil
}
