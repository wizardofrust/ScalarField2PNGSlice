# ScalarField2PNGSlice

This is a specialized program for taking a binary blob containing scalar field data (voxels) and converting it to a stack of PNG slices. It was designed for one application, and may need to be modified for any others.

## Building

Install Go and navigate a terminal into the folder containing main.go. Run this command:

```
go build .
```

Cross-compilation, or compiling for a platform other than the one you are on, is easy. Follow the normal directions for cross-compilation in Go. Here is an example using Powershell on Windows:

```
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build .
```

## Usage

Run the program through the command line or terminal. Ensure that your input data is a raw data blob of 32-bit floats (4 bytes each). For a list of options, run the program with the --help flag:

```
ScalarField2PNGSlice.exe --help
```

Here is an example command for running the program:

```
ScalarField2PNGSlice.exe -inp raw_data_out.bin -dim 192,2048,192 -out slices.zip -inv -log 2 -perm 213
```

The -inp flag is required, and specifies the path to the input data.  
The -dim flag is required, and specifies the dimensions of the input data, from least to most significant. If your data were a book, the dimensions would be: number of letters per line, number of lines per page, and number of pages.  
The -out flag is required, and specifies the path to store the output. It is structured as a zip file, so ensure you add a .zip extension.  
The -inv flag inverts the data before processing.  
The -log flag specifies how many times to apply a log function. This is useful for lowering the contrast.  
The -perm flag swaps the dimensions before saving. As an example, a value of 312 takes xyz input and changes it to zxy output. Swapping dimensions is equivalent to rotating and/or mirroring. You may choose to leave this at the default of 123, and simply rotate and/or mirror the object in your final program.
