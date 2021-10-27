package imprimir_qr

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo"
	"github.com/jung-kurt/gofpdf/contrib/httpimg"
)

type Respuesta struct {
	Status  string `json:"status"`
	Data    Data   `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Data struct {
}

func QrImprimir(c echo.Context) error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	id := c.Param("id")

	fmt.Println("'id:'", id)

	x := 0.00
	y := 0.00

	pdf.SetHeaderFunc(func() {
	})

	pdf.AddPage()

	// Plantilla
	pdf.SetXY(x-20, y-20)
	var opt gofpdf.ImageOptions
	opt.ImageType = "jpg"
	pdf.ImageOptions("image/plantilla_pdf_qr.jpg", 0, 0, 190, 280, false, opt, 0, "") // x, y, ancho, alto

	// QR
	url := "https://www.vivircarlospaz.com/img/lugares/qr/" + id + ".png"
	httpimg.Register(pdf, url, "")
	pdf.Image(url, 74, 173, 50, 0, false, "", 0, "") // x, y, size

	pdf.SetXY(x-20, y-20)
	path := fmt.Sprintf("%v", "public/pdfqr/qr-"+c.Param("id")+".pdf")
	pdf.OutputFileAndClose(path)

	r := c.File(path)
	_ = os.Remove(path)

	return r
}
