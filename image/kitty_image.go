//Wrapper around rastream (not working properly for now)

package image

import (
	"bytes"
	"errors"
	"image"
	"os"

	"github.com/BourgeoisBear/rasterm"
)

const KITTY_CLEAR_IMAGE = "\x1b_Ga=d\x1b\\"

// TODO: ADD Support for Image Id
type KittyImage struct {
	Id        string
	SrcWidth  uint64
	SrcHeight uint64
	Filepath  string
}

type KittyOpts struct {
	X                   uint32
	Y                   uint32
	DisplayWidth        uint32
	DisplayHeight       uint32
	MaintainAspectRatio bool
}

// TODO: Delete img with ID
func ClearImage() string {
	return KITTY_CLEAR_IMAGE
}

func Display(filePath string, kittyOpts KittyOpts) (string, error) {

	fIn, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fIn.Close()

	_, fmtName, err := image.DecodeConfig(fIn)
	if err != nil {
		return "", err
	}
	_, err = fIn.Seek(0, 0)
	if err != nil {
		return "", err
	}

	iImg, _, err := image.Decode(fIn)
	if err != nil {
		return "", err
	}

	_, err = fIn.Seek(0, 0)
	if err != nil {
		return "", err
	}

	opts := rasterm.KittyImgOpts{}

	opts.DstCols = kittyOpts.DisplayWidth
	opts.DstRows = kittyOpts.DisplayHeight

	bo := bytes.NewBufferString("")

	if fmtName == "png" {
		eF := rasterm.KittyWritePNGLocal(bo, filePath, opts)
		eI := rasterm.KittyCopyPNGInline(bo, fIn, opts)
		err = errors.Join(eI, eF)

	} else {
		err = rasterm.KittyWriteImage(bo, iImg, opts)
	}

	return bo.String(), err

}
