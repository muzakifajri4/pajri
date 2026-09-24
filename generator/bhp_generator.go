package generator

import (
	"fmt"

	"pdf-to-excel/parser"

	"github.com/xuri/excelize/v2"
)

func GenerateBHP(path string, records []parser.BHPRecord, header parser.BHPHeaderInfo) error {
	f := excelize.NewFile()
	defer f.Close()

	customNumFmt := "#,##0"

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{WrapText: true},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	numberStyle, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &customNumFmt,
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "000000"},
			{Type: "right", Style: 1, Color: "000000"},
			{Type: "top", Style: 1, Color: "000000"},
			{Type: "bottom", Style: 1, Color: "000000"},
		},
	})

	// Set column widths
	f.SetColWidth("Sheet1", "A", "A", 14)
	f.SetColWidth("Sheet1", "B", "B", 12)
	f.SetColWidth("Sheet1", "C", "C", 18)
	f.SetColWidth("Sheet1", "D", "D", 10)
	f.SetColWidth("Sheet1", "E", "E", 18)
	f.SetColWidth("Sheet1", "F", "F", 45)
	f.SetColWidth("Sheet1", "G", "G", 14)
	f.SetColWidth("Sheet1", "H", "H", 14)
	f.SetColWidth("Sheet1", "I", "I", 14)

	// Title
	row := 1
	titleCell := fmt.Sprintf("A%d", row)
	titleEnd := fmt.Sprintf("I%d", row)
	f.MergeCell("Sheet1", titleCell, titleEnd)
	f.SetCellValue("Sheet1", titleCell, "REKAPITULASI REALISASI BELANJA DANA BOSP ( BARANG HABIS PAKAI )")
	f.SetCellStyle("Sheet1", titleCell, titleCell, titleStyle)

	row = 2
	monthCell := fmt.Sprintf("A%d", row)
	monthEnd := fmt.Sprintf("I%d", row)
	f.MergeCell("Sheet1", monthCell, monthEnd)
	f.SetCellValue("Sheet1", monthCell, fmt.Sprintf("BULAN : %s TAHUN : %s", header.Month, header.Year))
	f.SetCellStyle("Sheet1", monthCell, monthCell, titleStyle)

	row = 3
	npsnCell := fmt.Sprintf("A%d", row)
	npsnEnd := fmt.Sprintf("D%d", row)
	f.MergeCell("Sheet1", npsnCell, npsnEnd)
	f.SetCellValue("Sheet1", npsnCell, fmt.Sprintf("NPSN : %s", header.NPSN))
	f.SetCellStyle("Sheet1", npsnCell, npsnCell, infoStyle)

	row = 4
	sekolahCell := fmt.Sprintf("A%d", row)
	sekolahEnd := fmt.Sprintf("D%d", row)
	f.MergeCell("Sheet1", sekolahCell, sekolahEnd)
	f.SetCellValue("Sheet1", sekolahCell, fmt.Sprintf("Nama Sekolah : %s", header.NamaSekolah))
	f.SetCellStyle("Sheet1", sekolahCell, sekolahCell, infoStyle)

	// Header row
	row = 6
	headers := []string{"Tanggal", "Kode\nKegiatan", "Kode Rekening", "No.\nBukti", "ID Barang", "Uraian", "Jumlah\nBarang", "Harga\nSatuan", "Realisasi"}
	for i, h := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s%d", col, row)
		f.SetCellValue("Sheet1", cell, h)
		f.SetCellStyle("Sheet1", cell, cell, headerStyle)
	}

	// Data rows
	row = 7
	for _, record := range records {
		colA := fmt.Sprintf("A%d", row)
		colB := fmt.Sprintf("B%d", row)
		colC := fmt.Sprintf("C%d", row)
		colD := fmt.Sprintf("D%d", row)
		colE := fmt.Sprintf("E%d", row)
		colF := fmt.Sprintf("F%d", row)
		colG := fmt.Sprintf("G%d", row)
		colH := fmt.Sprintf("H%d", row)
		colI := fmt.Sprintf("I%d", row)

		f.SetCellValue("Sheet1", colA, record.Date)
		f.SetCellStyle("Sheet1", colA, colA, dataStyle)

		f.SetCellValue("Sheet1", colB, record.KodeKeg)
		f.SetCellStyle("Sheet1", colB, colB, dataStyle)

		f.SetCellValue("Sheet1", colC, record.KodeRek)
		f.SetCellStyle("Sheet1", colC, colC, dataStyle)

		f.SetCellValue("Sheet1", colD, record.NoBukti)
		f.SetCellStyle("Sheet1", colD, colD, dataStyle)

		f.SetCellValue("Sheet1", colE, record.IDBarang)
		f.SetCellStyle("Sheet1", colE, colE, dataStyle)

		f.SetCellValue("Sheet1", colF, record.Uraian)
		f.SetCellStyle("Sheet1", colF, colF, dataStyle)

		f.SetCellValue("Sheet1", colG, record.JumlahBarang)
		f.SetCellStyle("Sheet1", colG, colG, numberStyle)

		f.SetCellValue("Sheet1", colH, record.HargaSatuan)
		f.SetCellStyle("Sheet1", colH, colH, numberStyle)

		f.SetCellValue("Sheet1", colI, record.Realisasi)
		f.SetCellStyle("Sheet1", colI, colI, numberStyle)

		row++
	}

	// Total row
	totalCell := fmt.Sprintf("A%d", row)
	totalEnd := fmt.Sprintf("F%d", row)
	f.MergeCell("Sheet1", totalCell, totalEnd)
	f.SetCellValue("Sheet1", totalCell, "Jumlah")
	f.SetCellStyle("Sheet1", totalCell, totalCell, headerStyle)

	var totalRealisasi float64
	for _, record := range records {
		totalRealisasi += record.Realisasi
	}
	totalValCell := fmt.Sprintf("I%d", row)
	f.SetCellValue("Sheet1", totalValCell, totalRealisasi)
	f.SetCellStyle("Sheet1", totalValCell, totalValCell, numberStyle)

	// Signature section
	row += 3
	setujuCell := fmt.Sprintf("A%d", row)
	setujuEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", setujuCell, setujuEnd)
	f.SetCellValue("Sheet1", setujuCell, "Menyetujui,")
	f.SetCellStyle("Sheet1", setujuCell, setujuCell, infoStyle)

	tanggalCell := fmt.Sprintf("E%d", row)
	tanggalEnd := fmt.Sprintf("F%d", row)
	f.MergeCell("Sheet1", tanggalCell, tanggalEnd)
	f.SetCellValue("Sheet1", tanggalCell, fmt.Sprintf("Kec. Arjawinangun, %s %s", header.Month, header.Year))
	f.SetCellStyle("Sheet1", tanggalCell, tanggalCell, infoStyle)

	row++
	kepsekLabelCell := fmt.Sprintf("A%d", row)
	kepsekLabelEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekLabelCell, kepsekLabelEnd)
	f.SetCellValue("Sheet1", kepsekLabelCell, "Kepala Sekolah")
	f.SetCellStyle("Sheet1", kepsekLabelCell, kepsekLabelCell, infoStyle)

	bendaharaLabelCell := fmt.Sprintf("E%d", row)
	bendaharaLabelEnd := fmt.Sprintf("F%d", row)
	f.MergeCell("Sheet1", bendaharaLabelCell, bendaharaLabelEnd)
	f.SetCellValue("Sheet1", bendaharaLabelCell, "Bendahara,")
	f.SetCellStyle("Sheet1", bendaharaLabelCell, bendaharaLabelCell, infoStyle)

	row++
	kepsekCell := fmt.Sprintf("A%d", row)
	kepsekEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekCell, kepsekEnd)
	f.SetCellValue("Sheet1", kepsekCell, header.KepalaSekolah)
	f.SetCellStyle("Sheet1", kepsekCell, kepsekCell, infoStyle)

	bendaharaCell := fmt.Sprintf("E%d", row)
	bendaharaEnd := fmt.Sprintf("F%d", row)
	f.MergeCell("Sheet1", bendaharaCell, bendaharaEnd)
	f.SetCellValue("Sheet1", bendaharaCell, header.Bendahara)
	f.SetCellStyle("Sheet1", bendaharaCell, bendaharaCell, infoStyle)

	row++
	kepsekNipCell := fmt.Sprintf("A%d", row)
	kepsekNipEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekNipCell, kepsekNipEnd)
	f.SetCellValue("Sheet1", kepsekNipCell, header.NIPKepsek)
	f.SetCellStyle("Sheet1", kepsekNipCell, kepsekNipCell, infoStyle)

	bendaharaNipCell := fmt.Sprintf("E%d", row)
	bendaharaNipEnd := fmt.Sprintf("F%d", row)
	f.MergeCell("Sheet1", bendaharaNipCell, bendaharaNipEnd)
	f.SetCellValue("Sheet1", bendaharaNipCell, header.NIPBendahara)
	f.SetCellStyle("Sheet1", bendaharaNipCell, bendaharaNipCell, infoStyle)

	// Footer
	row += 2
	footerCell := fmt.Sprintf("A%d", row)
	footerEnd := fmt.Sprintf("I%d", row)
	footerText := fmt.Sprintf("BHP %s %s - NPSN : %s, Nama Sekolah : %s", header.Month, header.Year, header.NPSN, header.NamaSekolah)
	f.MergeCell("Sheet1", footerCell, footerEnd)
	f.SetCellValue("Sheet1", footerCell, footerText)
	f.SetCellStyle("Sheet1", footerCell, footerCell, infoStyle)

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("failed to save BHP Excel: %w", err)
	}

	return nil
}
