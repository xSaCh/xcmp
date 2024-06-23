//Wrapper around rastream (not working properly for now)

package image

import (
	"bytes"
	"errors"
	"image"
	"os"

	"github.com/BourgeoisBear/rasterm"
)

type KittyImage struct {
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

func Display(filePath string, kittyOpts KittyOpts) (string, error) {

	fIn, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fIn.Close()

	imgConfig, fmtName, err := image.DecodeConfig(fIn)
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

	// ws := util.GetTermSize()
	// xs := ws.Xpixel / ws.Col
	// ys := ws.Ypixel / ws.Row
	// if kittyOpts.DisplayHeight != 0 {
	// 	opts.DstRows = kittyOpts.DisplayHeight / uint32(ys)

	// 	if kittyOpts.MaintainAspectRatio {
	// 		ratio := float32(imgConfig.Width) / float32(imgConfig.Height)
	// 		newWidth := float32(kittyOpts.DisplayHeight) * ratio

	// 		opts.DstCols = uint32(newWidth) / uint32(xs)
	// 	}
	// }
	// if kittyOpts.DisplayWidth != 0 {
	// 	opts.DstCols = kittyOpts.DisplayWidth / uint32(xs)

	// 	if kittyOpts.MaintainAspectRatio {
	// 		ratio := float32(imgConfig.Height) / float32(imgConfig.Width)
	// 		newHeight := float32(kittyOpts.DisplayWidth) * ratio

	// 		opts.DstRows = uint32(newHeight) / uint32(ys)
	// 	}
	// }
	_ = imgConfig
	// opts.CellOffsetX = kittyOpts.X
	// opts.CellOffsetY = kittyOpts.Y

	// fmt.Print(strings.Repeat())
	opts.DstCols = kittyOpts.DisplayWidth
	opts.DstRows = kittyOpts.DisplayHeight

	bo := bytes.NewBufferString("")

	// fmt.Fprint(bo, rasterm.ESC_ERASE_DISPLAY)

	if fmtName == "png" {
		eF := rasterm.KittyWritePNGLocal(bo, filePath, opts)
		eI := rasterm.KittyCopyPNGInline(bo, fIn, opts)
		err = errors.Join(eI, eF)

	} else {
		err = rasterm.KittyWriteImage(bo, iImg, opts)
	}

	return bo.String(), err

}
