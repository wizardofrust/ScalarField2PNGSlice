package main

import (
    "image"
    "image/color"
    "image/png"
    "encoding/binary"
    "archive/zip"
    "math"
    "flag"
    "os"
    "fmt"
    "errors"
    "strings"
    "strconv"
)

func readFileToMemory(inputPath string, invert bool, numLogs int) ([]float32, error) {
    fmt.Println("Loading data into memory")
    
    // Read the entire file into a byte slice
    rawData, err := os.ReadFile(inputPath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return nil, err
    }
    
    var packedData []float32
    
    // Each four bytes is a float, stored in little endian format
    for i := 0; i < len(rawData); i += 4 {
        rawBits := binary.LittleEndian.Uint32(rawData[i:i+4])
        amplitude := math.Float32frombits(rawBits)
        if invert {
            amplitude = -amplitude
        }
        for log := 0; log < numLogs; log++ {
            amplitude = float32(math.Log(float64(amplitude+1)))
        }
        packedData = append(packedData, amplitude)
    }
    
    return packedData, nil
}

func getMinAndMax(packedData []float32) (min float32, max float32) {
    min = packedData[0]
    max = packedData[0]
    for _, element := range packedData {
        if element < min {
            min = element
        }
        if element > max {
            max = element
        }
    }
    return min, max
}

func saveImagesFromArray(packedData []float32, width int, height int, depth int, permutation int, outputPath string) error {
    fmt.Println("Processing")
    // Create a new ZIP file
    zipFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer zipFile.Close()
    
    zw := zip.NewWriter(zipFile)
    defer zw.Close()
    
    min, max := getMinAndMax(packedData)
    fmt.Println("Min:", min)
    fmt.Println("Max:", max)
    
    var length1 int
    var length2 int
    var length3 int
    switch permutation {
        case 123:
            length1 = width
            length2 = height
            length3 = depth
        case 132:
            length1 = width
            length2 = depth
            length3 = height
        case 213:
            length1 = height
            length2 = width
            length3 = depth
        case 231:
            length1 = height
            length2 = depth
            length3 = width
        case 312:
            length1 = depth
            length2 = width
            length3 = height
        case 321:
            length1 = depth
            length2 = height
            length3 = width
        default:
            return errors.New("Error using permutation. Must be a three digit number containing 1, 2, and 3")
    }
    
    for d := 0; d < length3; d++ {
        img := image.NewGray(image.Rect(0, 0, length1, length2))
        for y := 0; y < length2; y++ {
            for x := 0; x < length1; x++ {
                var element float32
                switch permutation {
                    case 123:
                        element = packedData[d*width*height + y*width + x]
                    case 132:
                        element = packedData[y*width*height + d*width + x]
                    case 213:
                        element = packedData[d*width*height + x*width + y]
                    case 231:
                        //312 inverse
                        element = packedData[y*width*height + x*width + d]
                    case 312:
                        //231 inverse
                        element = packedData[x*width*height + d*width + y]
                    case 321:
                        element = packedData[x*width*height + y*width + d]
                    default:
                        return errors.New("Error using permutation. Must be a three digit number containing 1, 2, and 3")
                }
                element = (element-min)/(max-min) * math.MaxUint8
                img.SetGray(x, y, color.Gray{Y: uint8(element)})
            }
        }
        // Create a ZIP file writer for the image
        w, err := zw.Create(fmt.Sprintf("image.%0*d.png", 5, d))
        if err != nil {
            return err
        }
        
        // Encode the image into the ZIP file writer
        png.Encode(w, img)
    }
    
    return nil
}

func main() {
    if len(os.Args) == 1 {
        fmt.Println("Please enter command line arguments. To list the options, add \" --help\" to your command.")
        os.Exit(0)
    }
    
    // Set up input arguments
    inputPathPtr := flag.String("inp", "", "Path to the input file")
    inputDimensionsPtr := flag.String("dim", "", "Stored dimensions of the input file. Enter in the format x,y,z. Example: 192,2048,192")
    outputPathPtr := flag.String("out", "", "Path or name of the output file")
    invertPtr := flag.Bool("inv", false, "Multiply the data by -1. Just include the flag. Do not supply a value to it")
    numLogsPtr := flag.Int("log", 0, "Apply the target number of log scales to the data")
    permutationsPtr := flag.Int("perm", 123, "Permute the output dimensions. E.g. 312 changes output from xyz to zxy")
    flag.Parse()
    
    // Check arguments
    if *inputPathPtr == "" {
        fmt.Println("Error: input path is required")
        os.Exit(1)
    }
    if *outputPathPtr == "" {
        fmt.Println("Error: output path is required")
        os.Exit(1)
    }
    inputDimensionsSplit := strings.Split(*inputDimensionsPtr, ",")
    if len(inputDimensionsSplit) != 3 {
        fmt.Println("Error: dim argument must have exactly 3 numbers. Do you have 2 commas and no spaces?")
        os.Exit(1)
    }
    inputDimensionX, err := strconv.Atoi(inputDimensionsSplit[0])
    if err != nil {
        fmt.Println("Error parsing dimension 1")
        os.Exit(1)
    }
    inputDimensionY, err := strconv.Atoi(inputDimensionsSplit[1])
    if err != nil {
        fmt.Println("Error parsing dimension 2")
        os.Exit(1)
    }
    inputDimensionZ, err := strconv.Atoi(inputDimensionsSplit[2])
    if err != nil {
        fmt.Println("Error parsing dimension 3")
        os.Exit(1)
    }
    switch *permutationsPtr {
        case 123:
        case 132:
        case 213:
        case 231:
        case 312:
        case 321:
        default:
            fmt.Println("Error parsing dimension permutation. Must be a three digit number containing 1, 2, and 3")
            os.Exit(1)
    }
    
    //Begin processing
    packedData, err := readFileToMemory(*inputPathPtr, *invertPtr, *numLogsPtr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    err = saveImagesFromArray(packedData, inputDimensionX, inputDimensionY, inputDimensionZ, *permutationsPtr, *outputPathPtr)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}
