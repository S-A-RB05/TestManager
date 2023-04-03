package converter

import (
	"encoding/base64"
	"io"
	"os"
)

func DecodeBase64(file string) {

	dec, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("myfilename.ex5")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	// go to begginng of file
	f.Seek(0, 0)

	// output file contents
	io.Copy(os.Stdout, f)
}
