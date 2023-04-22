package code

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
	"tools/library/snowflake"
)

// CreateQrCodeImage ，生成二维码链接，
// 参数：
//     param 需要缓存的数据
//     key   作为缓存的 key 值之一
//     host  绑定微信需要跳转的地址
// 返回值：
//      image_url 二维码图片保存地址
//      report_code 上报码，用于在扫码跳转后
//      err         返回的错误
func CreateQrCodeImage(param, key interface{}, host string) (image_url string, err error) {
	// report_code 雪花算法ID
	//report_code = generateReportCode()
	//url := fmt.Sprintf("%s?report_code=%s", host, report_code)
	report_code := gconv.String(key)
	report_code_key := fmt.Sprintf("%s%s", REPORT_CODE_KEY, report_code)
	report_code_value := gconv.String(param)
	gcache.Set(report_code_key, report_code_value, time.Duration(180)*time.Second) //暂时有效时间为120秒方便测试
	glog.Info("report_code_key", report_code_key)
	// 测试一下缓存的
	v, _ := gcache.Get(report_code_key)
	glog.Info("get key v:", v)

	//然后生成一张二维码的图片
	fn, _ := gmd5.Encrypt(report_code)
	//image_url = fmt.Sprintf("public\\qrcode\\%s\\%s", gtime.Now().Format("Ymd"), fn)
	//image_url = fmt.Sprintf("public/qrcode/%s/%s", gtime.Now().Format("Ymd"), fn)
	image_url = fmt.Sprintf("/qrcode/%s/%s", gtime.Now().Format("Ymd"), fn+".png")
	save_image_url := fmt.Sprintf("%s%s", "public", image_url)

	//s.createQrCode(url, "", image_url)
	createQrCode(host, "", save_image_url)

	return
}

// generateReportCode 生成一个唯一的 report 代号
func generateReportCode() string {
	return gconv.String(snowflake.GetInstance().GetId())
}

func createQrCode(text string, logoPath string, outQrPath string) {
	glog.Info(fmt.Sprintf("生成二维码:%s,", outQrPath))
	size := 400
	percent := 15

	code, err := qrcode.New(text, qrcode.Highest)
	if err != nil {
		glog.Error(err)
	}

	srcImage := code.Image(size)
	if logoPath != "" {
		logoSize := float64(size) * float64(percent) / 100
		srcImage, err = addLogo(srcImage, logoPath, int(logoSize))
	}

	// outAbs 二维码生成文件的路径
	outAbs, err := filepath.Abs(outQrPath)

	os.MkdirAll(filepath.Dir(outAbs), 0777)
	outFile, err := os.Create(outAbs)
	defer outFile.Close()

	jpeg.Encode(outFile, srcImage, &jpeg.Options{Quality: 100})

	glog.Infof("二维码生成成功，文件路径：%s", outAbs)
}

func resizeLogo(logo string, size uint) (image.Image, error) {
	file, err := os.Open(logo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	img = resize.Resize(size, size, img, resize.Lanczos3)
	return img, nil
}

func addLogo(srcImage image.Image, logo string, size int) (image.Image, error) {
	logoImage, err := resizeLogo(logo, uint(size))
	if err != nil {
		return nil, err
	}

	offset := image.Pt((srcImage.Bounds().Dx()-logoImage.Bounds().Dx())/2, (srcImage.Bounds().Dy()-logoImage.Bounds().Dy())/2)
	b := srcImage.Bounds()
	m := image.NewNRGBA(b)
	draw.Draw(m, b, srcImage, image.ZP, draw.Src)
	draw.Draw(m, logoImage.Bounds().Add(offset), logoImage, image.ZP, draw.Over)

	return m, nil
}
