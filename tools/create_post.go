package tools

import (
	"bufio"
	"bytes"
	"github.com/astaxie/beego"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)
type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

type ArticlePoster struct {
	PosterName string
	*Article
}

func NewArticlePoster(posterName string, article *Article) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool {
	if CheckNotExist(path+a.PosterName) == true {
		return false
	}

	return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

//在图片上写文字
type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	fontSource := "static/fonts/" + fontName
	beego.Info("font source is",fontSource)
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}

	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.White)

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}

//先不添加文字，之后添加文字因为需要用户名
func (a *ArticlePosterBg) Generate(path string,oldfileName string,fileName string,xsize int ,ysize int) (string, string, error) {
	//path :="static/img/"
	//oldfileName := "showqrcode_yougu.jpeg"
	//fileName :="new_showqrcode_yougu.jpeg"
	//检查新的二维码图片是否存在，存在重新生成
	if CheckNotExist(path+fileName) == false {
		err := Del(path+fileName)
		if err != nil {
			return "", "", nil
		}
	}

	// open "test.jpg"
	file, err := os.Open(path+oldfileName)
	if err != nil {
		return "", "", nil
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return "", "", nil
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio

	//kuan:=img.Bounds().Dx()//



	//m := resize.Resize(uint(kuan), 0, img, resize.Lanczos3)
	m := resize.Resize(uint(xsize), uint(ysize), img, resize.Lanczos3)
	//m := resize.Resize(100, 100, img, resize.Lanczos3)
	//m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create(path+fileName)
	if err != nil {
		return "", "", nil
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)



	//检查post123是否存在，存在删除，重新生成
	if a.CheckMergedImage(path) {
		err := Del(path+a.PosterName)
		if err != nil {
			return "", "", nil
		}
	}
	mergedF, err := a.OpenMergedImage(path)
	if err != nil {
		return "", "", err
	}
	defer mergedF.Close()

	bgF, err := MustOpen(a.Name, path)
	if err != nil {
		return "", "", err
	}
	defer bgF.Close()


	qrF, err := MustOpen(fileName, path)
	if err != nil {
		return "", "", err
	}
	defer qrF.Close()

	bgImage, err := jpeg.Decode(bgF)
	if err != nil {
		return "", "", err
	}

	qrImage, err := jpeg.Decode(qrF)
	if err != nil {
		return "", "", err
	}

	jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

	draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
	draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)
	//如果不添加文字直接把这个解注释掉，否则不生成文字
	jpeg.Encode(mergedF, jpg, nil)

	//生成有文字加二维码的jpg
	//err = a.DrawPoster(&DrawText{
	//	JPG:    jpg,
	//	Merged: mergedF,
	//
	//	Title: "陈君邀请你来一起参加",
	//	X0:    370,
	//	Y0:    370,
	//	Size0: 50,
	//	//
	//	//SubTitle: "---溢修瑜珈body&soul",
	//	//X1:       350,
	//	//Y1:       350,
	//	//Size1:    30,
	//}, "handan.ttf")
	//
	//if err != nil {
	//	return "", "", err
	//}

	return fileName, path, nil
}


//先添加qrcode 到海报上，生成bg.jpg
func GetQrPoster(qrcodename string,bgname string,postname string){
	article := &Article{}

	// 先取二维码图片，oldfileName
	//生成新的二维码图片 fileName
	//新的二维码图片变成多大图片xsize,ysize
	//path 是总体的目录，无论是二维码还是生成的poster 目录，都在path目录下
	path :="static/img/"

	//oldfileName := "qrstatic/showqrcode_yougu.jpeg"
	//fileName :="qrstatic/new_showqrcode_yougu.jpeg"

	oldfileName := "qrstatic/"+qrcodename
	fileName :="qrstatic/new_"+qrcodename
	xsize :=250
	ysize :=250


	//背景图片
	//backgroundimg :="qrstatic/bg.jpeg"
	////合并之后的图片
	//posterName := "openid/poster123.jpg"

	//背景图片
	backgroundimg :="qrstatic/"+bgname
	//合并之后的图片
	posterName := "qrstatic/"+postname

	articlePoster := NewArticlePoster(posterName, article)
	articlePosterBgService := NewArticlePosterBg(
		backgroundimg,//背景图片
		articlePoster,
		//背景图片大小
		&Rect{
			X0: 0,
			Y0: 0,
			X1: 1079,
			Y1: 1604,
		},
		//二维码放在什么位置
		&Pt{
			X: 420,
			Y: 730,
		},
	)

	_, filePath, err := articlePosterBgService.Generate(path,oldfileName,fileName,xsize,ysize)

	beego.Info("filePath is:",filePath)
	if err != nil {
		beego.Info("something wrong",err)
		return
	}
}
//再去添加用户头像
func GetFinalPoster(avatarname string,bgname string,postname string,title string){//static 下一层应该是appid
	article := &Article{}

	//path 是总体的目录，无论是二维码还是生成的poster 目录，都在path目录下
	path :="static/img/"
	//openid 原始头像
	fileName :="openid/"+avatarname

	//背景图片
	backgroundimg :="qrstatic/"+bgname
	//合并之后的图片
	posterName := "openid/"+postname

	articlePoster := NewArticlePoster(posterName, article)
	articlePosterBgService := NewArticlePosterBg(
		backgroundimg,//背景图片
		articlePoster,
		//背景图片大小
		&Rect{
			X0: 0,
			Y0: 0,
			X1: 1079,
			Y1: 1604,
		},
		//头像放在什么位置
		&Pt{
			X: 250,
			Y: 300,
		},
	)
	_, filePath, err := articlePosterBgService.GenerateFinal(path,fileName,title)

	beego.Info("filePath is:",filePath)
	if err != nil {
		beego.Info("something wrong",err)
		return
	}
}

//添加头像，并且添加文字
func (a *ArticlePosterBg) GenerateFinal(path string,fileName string,title string) (string, string, error) {
	//path :="static/img/"
	//fileName :="new_showqrcode_yougu.jpeg"



	//检查post123是否存在，存在删除，重新生成,之前版本如果存在就生成，现在存在的话，就直接返回，不删除
	//if a.CheckMergedImage(path) {
	//	err := Del(path+a.PosterName)
	//	if err != nil {
	//		return "", "", nil
	//	}
	//}
	if a.CheckMergedImage(path) {//存在则返回
		return path, fileName, nil
	}
	beego.Info("头像路径：",path+fileName)
	beego.Info("海报不存在，重新生成的路径是：",path+a.PosterName)
	mergedF, err := a.OpenMergedImage(path)
	if err != nil {
		return "", "", err
	}
	defer mergedF.Close()

	//背景图
	bgF, err := MustOpen(a.Name, path)
	if err != nil {
		return "", "", err
	}
	defer bgF.Close()

	//头像图片//检查新的头像图片是否存在，存在重新生成
	qrF, err := MustOpen(fileName, path)
	if err != nil {
		return "", "", err
	}
	defer qrF.Close()

	bgImage, err := jpeg.Decode(bgF)
	if err != nil {
		return "", "", err
	}

	qrImage, err := jpeg.Decode(qrF)
	if err != nil {
		return "", "", err
	}

	jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

	draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
	draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)
	//如果不添加文字直接把这个解注释掉
	jpeg.Encode(mergedF, jpg, nil)
	//不需要生成文字

	err = a.DrawPoster(&DrawText{
		JPG:    jpg,
		Merged: mergedF,

		Title: title,
		X0:    50,
		Y0:    50,
		Size0: 20,

		SubTitle: "---溢修瑜珈body&soul",
		X1:       80,
		Y1:       80,
		Size1:    15,
	}, "handan.ttf")

	if err != nil {
		return "", "", err
	}

	return fileName, path, nil
}

func WrongDownloadPic(imgUrl string,imgPath string,fileName string) error{
	//imgPath := "static/img/openid/"
	//imgUrl := userinfo.HeadImgURL
	//
	////fileName := path.Base(imgUrl)
	//fileName := openid+".jpeg"

	res, err := http.Get(imgUrl)
	if err != nil {
		beego.Info("图片信息获取错误")
		return err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32 * 1024)


	file, err := os.Create(imgPath + fileName)// 如果文件已存在，会将文件清空。
	if err != nil {
		beego.Info("图片创建失败",err)
		panic(err)
		return err
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	written, _ := io.Copy(writer, reader)
	beego.Info("Total length:", written)
	beego.Info("图片获取正确。。。。。")
	return nil
}

func DownloadPic(imagPath string,destPath string,width int,height int,option int) (err error) {
	resp, _ := http.Get(imagPath)

	var body []byte
	body, _ = ioutil.ReadAll(resp.Body)

	//var err error
	var data *bytes.Buffer
	if data, err = SetScaleImage(body, width, height, option); err != nil {
		beego.Info("err is:",err)
		return err
	}
	//beego.Info("data is:",data)

	if data != nil {
		if err := SaveImage(data.Bytes(), destPath); err != nil {
			beego.Info("second",err)
			return err
		}
	} else {
		beego.Info("data is nil")
		return nil
	}
	return nil
}
