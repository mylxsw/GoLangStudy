package main

import "github.com/Collinux/watermark"

func main() {

	logo := watermark.Watermark{
		Source:   "/Users/mylxsw/codes/work/yunsom-scm-site/public/assets/yunsom-scm/images/ys-seal.png",
		Position: watermark.BOTTOM_RIGHT,
	}

	logo.Apply("/Users/mylxsw/Downloads/FullSizeRender.jpg")

}
