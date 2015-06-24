package main
//TODO: change this to a lib to use with s3 uploader go program
import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"log"
	"fmt"
	"io"
	"strings"
	"os"
)

func checkerror(err error) {

	if err != nil {
		log.Fatal(err)
	}
}


// this function does what it reads 
func DeepFileGet(dirname string, tw *tar.Writer )  {
	dir, err := os.Open(dirname)
	checkerror(err)
	defer dir.Close()
	fi, err := dir.Readdir(0)
	checkerror(err)
	for _, file := range fi {
		curPath := dirname + "/" + file.Name()
		if file.IsDir() {
			DeepFileGet( curPath, tw)
		} else {
			fmt.Printf( "adding ... %s\n", curPath)
			TarGzWrite( curPath, tw, file )
		}
	}
}


func TarGz( destinationFile string, sourceDir string) {
	// write file
	fw, err := os.Create( destinationFile )
	checkerror(err)
	defer fw.Close()


	// gzip the ish up

	gw := gzip.NewWriter( fw )
	defer gw.Close()

	// write that tar ish

	tw := tar.NewWriter( gw)
	defer tw.Close()

	DeepFileGet( sourceDir, tw )
	fmt.Println( "tar ok")
}



// i lika leave comments 
func parseFlags() ( string,  string) {
	flag.Parse() // get the arguments from command line
	destinationFile := flag.Arg(0)
	if destinationFile == "" {
		fmt.Println("Usage : compressr destinationfile.tar.gz source")
		os.Exit(1)
	}
	sourceDir := flag.Arg(1)
	if sourceDir == "" {
		fmt.Println("Usage : compressr destinationfile.tar.gz source-directory")
		os.Exit(1)

	}
	return  destinationFile, sourceDir
}

func TarGzWrite( _path string, tw *tar.Writer, fi os.FileInfo) {

	fr, err := os.Open( _path )
	checkerror(err)

	h := new( tar.Header )
	h.Name = _path
	h.Size = fi.Size()
	h.Mode = int64( fi.Mode() )
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	checkerror(err)

	_, err = io.Copy( tw, fr )
	checkerror(err)

}




func main() {
// TODO: Get soft error when path to directory is too long
// need to fix / truncate it to just the target so 
// the tar archive header isn't too long. This doesn't cause it to fail
// just throws soft error

// parse cmd line options and return targets and source
	destinationFile, sourceDir := parseFlags()
	TarGz( destinationFile, strings.TrimRight( sourceDir, "/") )
}
