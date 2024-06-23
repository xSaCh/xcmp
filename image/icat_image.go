package image

import (
	"fmt"
	"image"
	"os"
	"os/exec"

	"github.com/xSaCh/xcmp/util"
)

type IcatOpts struct {
	Dx                  uint32
	Dy                  uint32
	W                   uint32
	H                   uint32
	MaintainAspectRatio bool
}

func DisplayA(filePath string, icatOpts IcatOpts) error {

	fIn, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fIn.Close()

	imgConfig, _, err := image.DecodeConfig(fIn)
	if err != nil {
		return err
	}

	if icatOpts.H == 0 && icatOpts.W == 0 {
		icatOpts.H = uint32(imgConfig.Height)
		icatOpts.W = uint32(imgConfig.Width)
	} else if icatOpts.MaintainAspectRatio {
		ratio := float32(imgConfig.Width) / float32(imgConfig.Height)
		if icatOpts.W != 0 {
			icatOpts.H = uint32(float32(icatOpts.W) * ratio)
		} else {
			icatOpts.W = uint32(float32(icatOpts.H) / ratio)
		}
	}

	// xp, yp, _ := util.GetCousorPos()

	// executeIcat(filePath, icatOpts.W, icatOpts.H, xp+icatOpts.Dx, yp+icatOpts.Dy)
	executeIcat(filePath, icatOpts.W, icatOpts.H, icatOpts.Dx, icatOpts.Dy)

	return err

}
func DisplayB(filePath string, icatOpts IcatOpts) error {

	fIn, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fIn.Close()

	imgConfig, _, err := image.DecodeConfig(fIn)
	if err != nil {
		return err
	}

	t := util.GetTermSize()
	xM := float32(t.Xpixel) / float32(t.Col)
	yM := float32(t.Ypixel) / float32(t.Row)
	_ = imgConfig

	icatOpts.W *= uint32(xM)
	icatOpts.H *= uint32(yM)

	// if icatOpts.H == 0 && icatOpts.W == 0 {
	// 	icatOpts.H = uint32(imgConfig.Height)
	// 	icatOpts.W = uint32(imgConfig.Width)
	// } else if icatOpts.MaintainAspectRatio {
	// 	ratio := float32(imgConfig.Width) / float32(imgConfig.Height)
	// 	if icatOpts.W != 0 {
	// 		icatOpts.H = uint32(float32(icatOpts.W) * ratio)
	// 	} else {
	// 		icatOpts.W = uint32(float32(icatOpts.H) / ratio)
	// 	}
	// }

	// xp, yp, _ := util.GetCousorPos()

	// executeIcat(filePath, icatOpts.W, icatOpts.H, xp+icatOpts.Dx, yp+icatOpts.Dy)
	executeIcat(filePath, icatOpts.W, icatOpts.H, icatOpts.Dx, icatOpts.Dy)

	return err

}

var isEx = false

func executeIcat(filePath string, w, h, x, y uint32) {

	// if isEx {
	// 	return
	// }
	pos := fmt.Sprintf("%dx%d@%dx%d", w, h, x, y)
	fmt.Printf("pos: %v\n", pos)
	command := []string{"kitty", "+kitten", "icat", "--scale-up", "--align", "left", "--place", pos, filePath}
	err := exec.Command(command[0], command[1:]...).Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	isEx = true

}
