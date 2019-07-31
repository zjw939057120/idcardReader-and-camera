package sdk

import (
	"encoding/base64"
	"gocv.io/x/gocv"
	"image"
)

var (
	webcam    *gocv.VideoCapture
	webcamErr error
	img       gocv.Mat
	Image     string
	Camera    int
)

//调用摄像头
func OpenCamera() {
	if Camera != 1 {
		// open webcam
		webcam, webcamErr = gocv.OpenVideoCapture(0)
		if webcamErr != nil {
			Camera = -1
			return
		}
		defer webcam.Close()
		img = gocv.NewMat()
		defer img.Close()
		Camera = 1

		for {
			if Camera == -1 {
				return
			}
			webcam.Read(&img)
			if img.Empty() {
				continue
			}
			gocv.Resize(img, &img, image.Point{480, 640}, 0, 0, gocv.InterpolationDefault)
			buf, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				continue
			}
			Image = base64.StdEncoding.EncodeToString(buf)
		}
	}

}

//关闭摄像头
func CloseCamera() {
	if Camera != -1 {
		Camera = -1
	}
	//fmt.Println(Camera)
	return
}
