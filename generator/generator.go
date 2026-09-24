package generator

import (
	"fmt"
	"strings"

	"pdf-to-excel/parser"

	"github.com/xuri/excelize/v2"
)

func GenerateBKU(path string, records []parser.Record, header parser.HeaderInfo) error {
	f := excelize.NewFile()
	defer f.Close()

	// Create styles
	customNumFmt := "#,##0"

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{WrapText: true},
	})

	headerRowStyle, _ := f.NewStyle(&excelize.Style{
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
	f.SetColWidth("Sheet1", "E", "E", 50)
	f.SetColWidth("Sheet1", "F", "F", 16)
	f.SetColWidth("Sheet1", "G", "G", 16)
	f.SetColWidth("Sheet1", "H", "H", 16)

	// Title
	row := 1
	f.MergeCell("Sheet1", "A1", "H1")
	f.SetCellValue("Sheet1", "A1", "B U K U  K A S  U M U M")
	f.SetCellStyle("Sheet1", "A1", "A1", titleStyle)

	row = 2
	f.MergeCell("Sheet1", "A2", "H2")
	monthYear := fmt.Sprintf("BULAN : %s TAHUN : %s", strings.ToUpper(header.Month), header.Year)
	f.SetCellValue("Sheet1", "A2", monthYear)
	f.SetCellStyle("Sheet1", "A2", "A2", headerStyle)

	row = 3
	f.MergeCell("Sheet1", "A3", "H3")
	f.SetCellValue("Sheet1", "A3", fmt.Sprintf("NPSN : %s", header.NPSN))
	f.SetCellStyle("Sheet1", "A3", "A3", infoStyle)

	row = 4
	f.MergeCell("Sheet1", "A4", "H4")
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf("Nama Sekolah : %s", header.NamaSekolah))
	f.SetCellStyle("Sheet1", "A4", "A4", infoStyle)

	row = 5
	f.MergeCell("Sheet1", "A5", "H5")
	f.SetCellValue("Sheet1", "A5", fmt.Sprintf("Desa/Kecamatan : %s", header.DesaKecamatan))
	f.SetCellStyle("Sheet1", "A5", "A5", infoStyle)

	row = 6
	f.MergeCell("Sheet1", "A6", "H6")
	f.SetCellValue("Sheet1", "A6", fmt.Sprintf("Kabupaten / Kota : %s", header.KabupatenKota))
	f.SetCellStyle("Sheet1", "A6", "A6", infoStyle)

	row = 7
	f.MergeCell("Sheet1", "A7", "H7")
	f.SetCellValue("Sheet1", "A7", fmt.Sprintf("Provinsi : %s", header.Provinsi))
	f.SetCellStyle("Sheet1", "A7", "A7", infoStyle)

	row = 8
	f.MergeCell("Sheet1", "A8", "H8")
	f.SetCellValue("Sheet1", "A8", fmt.Sprintf("Sumber Dana : %s", header.SumberDana))
	f.SetCellStyle("Sheet1", "A8", "A8", infoStyle)

	// Header row
	row = 9
	headers := []string{"TANGGAL", "KODE\nKEGIATAN", "KODE\nREKENING", "NO. BUKTI", "URAIAN", "PENERIMAAN", "PENGELUARAN", "SALDO"}
	for i, h := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s%d", col, row)
		f.SetCellValue("Sheet1", cell, h)
		f.SetCellStyle("Sheet1", cell, cell, headerRowStyle)
	}

	// Sub header row
	row = 10
	subHeaders := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for i, h := range subHeaders {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s%d", col, row)
		f.SetCellValue("Sheet1", cell, h)
		f.SetCellStyle("Sheet1", cell, cell, headerRowStyle)
	}

	// Data rows
	row = 11
	for _, record := range records {
		colA := fmt.Sprintf("A%d", row)
		colB := fmt.Sprintf("B%d", row)
		colC := fmt.Sprintf("C%d", row)
		colD := fmt.Sprintf("D%d", row)
		colE := fmt.Sprintf("E%d", row)
		colF := fmt.Sprintf("F%d", row)
		colG := fmt.Sprintf("G%d", row)
		colH := fmt.Sprintf("H%d", row)

		f.SetCellValue("Sheet1", colA, record.Date)
		f.SetCellStyle("Sheet1", colA, colA, dataStyle)

		f.SetCellValue("Sheet1", colB, record.KodeKeg)
		f.SetCellStyle("Sheet1", colB, colB, dataStyle)

		f.SetCellValue("Sheet1", colC, record.KodeRek)
		f.SetCellStyle("Sheet1", colC, colC, dataStyle)

		f.SetCellValue("Sheet1", colD, record.NoBukti)
		f.SetCellStyle("Sheet1", colD, colD, dataStyle)

		f.SetCellValue("Sheet1", colE, record.Uraian)
		f.SetCellStyle("Sheet1", colE, colE, dataStyle)

		f.SetCellValue("Sheet1", colF, record.Penerimaan)
		f.SetCellStyle("Sheet1", colF, colF, numberStyle)

		f.SetCellValue("Sheet1", colG, record.Pengeluaran)
		f.SetCellStyle("Sheet1", colG, colG, numberStyle)

		f.SetCellValue("Sheet1", colH, record.Saldo)
		f.SetCellStyle("Sheet1", colH, colH, numberStyle)

		row++
	}

	// Summary row
	colA := fmt.Sprintf("A%d", row)
	colE := fmt.Sprintf("E%d", row)
	colF := fmt.Sprintf("F%d", row)
	colG := fmt.Sprintf("G%d", row)
	colH := fmt.Sprintf("H%d", row)

	f.SetCellValue("Sheet1", colA, "Jumlah")
	f.SetCellStyle("Sheet1", colA, colA, headerRowStyle)
	f.MergeCell("Sheet1", colA, colE)

	// Calculate totals
	var totalPenerimaan, totalPengeluaran float64
	for _, record := range records {
		totalPenerimaan += record.Penerimaan
		totalPengeluaran += record.Pengeluaran
	}

	f.SetCellValue("Sheet1", colF, totalPenerimaan)
	f.SetCellStyle("Sheet1", colF, colF, numberStyle)
	f.SetCellValue("Sheet1", colG, totalPengeluaran)
	f.SetCellStyle("Sheet1", colG, colG, numberStyle)
	f.SetCellValue("Sheet1", colH, 0)
	f.SetCellStyle("Sheet1", colH, colH, numberStyle)

	// Add closing section
	row += 2
	closingText := fmt.Sprintf("Pada hari ini %s Buku Kas Umum Ditutup dengan keadaan/posisi buku sebagai berikut :", header.TanggalTutup)
	closingCell := fmt.Sprintf("A%d", row)
	closingEnd := fmt.Sprintf("H%d", row)
	f.MergeCell("Sheet1", closingCell, closingEnd)
	f.SetCellValue("Sheet1", closingCell, closingText)
	f.SetCellStyle("Sheet1", closingCell, closingCell, infoStyle)

	row++
	saldoCell := fmt.Sprintf("A%d", row)
	saldoEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", saldoCell, saldoEnd)
	f.SetCellValue("Sheet1", saldoCell, "Saldo Buku Kas Umum : Rp. 0")
	f.SetCellStyle("Sheet1", saldoCell, saldoCell, infoStyle)

	row++
	terdiriCell := fmt.Sprintf("A%d", row)
	terdiriEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", terdiriCell, terdiriEnd)
	f.SetCellValue("Sheet1", terdiriCell, "Terdiri Dari :")
	f.SetCellStyle("Sheet1", terdiriCell, terdiriCell, infoStyle)

	row++
	bankCell := fmt.Sprintf("A%d", row)
	bankEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", bankCell, bankEnd)
	f.SetCellValue("Sheet1", bankCell, "- Saldo Bank : Rp. 0")
	f.SetCellStyle("Sheet1", bankCell, bankCell, infoStyle)

	row++
	tunaiCell := fmt.Sprintf("A%d", row)
	tunaiEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", tunaiCell, tunaiEnd)
	f.SetCellValue("Sheet1", tunaiCell, "- Saldo Kas Tunai : Rp. 0")
	f.SetCellStyle("Sheet1", tunaiCell, tunaiCell, infoStyle)

	row++
	jumlahCell := fmt.Sprintf("A%d", row)
	jumlahEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", jumlahCell, jumlahEnd)
	f.SetCellValue("Sheet1", jumlahCell, "Jumlah : Rp. 0")
	f.SetCellStyle("Sheet1", jumlahCell, jumlahCell, infoStyle)

	// Signature section
	row += 2
	setujuCell := fmt.Sprintf("A%d", row)
	setujuEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", setujuCell, setujuEnd)
	f.SetCellValue("Sheet1", setujuCell, "Menyetujui,")
	f.SetCellStyle("Sheet1", setujuCell, setujuCell, infoStyle)

	tanggalCell := fmt.Sprintf("D%d", row)
	tanggalEnd := fmt.Sprintf("E%d", row)
	f.MergeCell("Sheet1", tanggalCell, tanggalEnd)
	f.SetCellValue("Sheet1", tanggalCell, fmt.Sprintf("Kec. Arjawinangun, %s", header.TanggalTutup))
	f.SetCellStyle("Sheet1", tanggalCell, tanggalCell, infoStyle)

	row++
	kepsekCell := fmt.Sprintf("A%d", row)
	kepsekEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekCell, kepsekEnd)
	f.SetCellValue("Sheet1", kepsekCell, "Kepala Sekolah")
	f.SetCellStyle("Sheet1", kepsekCell, kepsekCell, infoStyle)

	bendaharaLabelCell := fmt.Sprintf("D%d", row)
	bendaharaLabelEnd := fmt.Sprintf("E%d", row)
	f.MergeCell("Sheet1", bendaharaLabelCell, bendaharaLabelEnd)
	f.SetCellValue("Sheet1", bendaharaLabelCell, "Bendahara,")
	f.SetCellStyle("Sheet1", bendaharaLabelCell, bendaharaLabelCell, infoStyle)

	row++
	kepsekNameCell := fmt.Sprintf("A%d", row)
	kepsekNameEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekNameCell, kepsekNameEnd)
	f.SetCellValue("Sheet1", kepsekNameCell, header.KepalaSekolah)
	f.SetCellStyle("Sheet1", kepsekNameCell, kepsekNameCell, infoStyle)

	bendaharaNameCell := fmt.Sprintf("D%d", row)
	bendaharaNameEnd := fmt.Sprintf("E%d", row)
	f.MergeCell("Sheet1", bendaharaNameCell, bendaharaNameEnd)
	f.SetCellValue("Sheet1", bendaharaNameCell, header.Bendahara)
	f.SetCellStyle("Sheet1", bendaharaNameCell, bendaharaNameCell, infoStyle)

	row++
	kepsekNipCell := fmt.Sprintf("A%d", row)
	kepsekNipEnd := fmt.Sprintf("B%d", row)
	f.MergeCell("Sheet1", kepsekNipCell, kepsekNipEnd)
	f.SetCellValue("Sheet1", kepsekNipCell, header.NIPKepsek)
	f.SetCellStyle("Sheet1", kepsekNipCell, kepsekNipCell, infoStyle)

	bendaharaNipCell := fmt.Sprintf("D%d", row)
	bendaharaNipEnd := fmt.Sprintf("E%d", row)
	f.MergeCell("Sheet1", bendaharaNipCell, bendaharaNipEnd)
	f.SetCellValue("Sheet1", bendaharaNipCell, header.NIPBendahara)
	f.SetCellStyle("Sheet1", bendaharaNipCell, bendaharaNipCell, infoStyle)

	// Footer
	row += 2
	footerCell := fmt.Sprintf("A%d", row)
	footerEnd := fmt.Sprintf("H%d", row)
	footerText := fmt.Sprintf("BKU %s %s - NPSN : %s, Nama Sekolah : %s", header.Month, header.Year, header.NPSN, header.NamaSekolah)
	f.MergeCell("Sheet1", footerCell, footerEnd)
	f.SetCellValue("Sheet1", footerCell, footerText)
	f.SetCellStyle("Sheet1", footerCell, footerCell, infoStyle)

	// Save
	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("failed to save BKU Excel: %w", err)
	}

	return nil
}

func GenerateTemplateImport(path string, records []parser.Record, header parser.HeaderInfo) error {
	f := excelize.NewFile()
	defer f.Close()

	// Rename default sheet
	f.SetSheetName("Sheet1", "Sheet2")

	// Create styles
	customNumFmt := "#,##0"

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
	f.SetColWidth("Sheet2", "A", "A", 14)
	f.SetColWidth("Sheet2", "B", "B", 12)
	f.SetColWidth("Sheet2", "C", "C", 18)
	f.SetColWidth("Sheet2", "D", "D", 10)
	f.SetColWidth("Sheet2", "E", "E", 50)
	f.SetColWidth("Sheet2", "F", "F", 16)

	// Header row
	row := 1
	headers := []string{"TANGGAL", "KODE KEGIATAN", "KODE REKENING", "NO BUKTI", "URAIAN", "PENGELUARAN"}
	for i, h := range headers {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s%d", col, row)
		f.SetCellValue("Sheet2", cell, h)
		f.SetCellStyle("Sheet2", cell, cell, headerStyle)
	}

	// Data rows
	row = 2
	for _, record := range records {
		colA := fmt.Sprintf("A%d", row)
		colB := fmt.Sprintf("B%d", row)
		colC := fmt.Sprintf("C%d", row)
		colD := fmt.Sprintf("D%d", row)
		colE := fmt.Sprintf("E%d", row)
		colF := fmt.Sprintf("F%d", row)

		f.SetCellValue("Sheet2", colA, record.Date)
		f.SetCellStyle("Sheet2", colA, colA, dataStyle)

		f.SetCellValue("Sheet2", colB, record.KodeKeg)
		f.SetCellStyle("Sheet2", colB, colB, dataStyle)

		f.SetCellValue("Sheet2", colC, record.KodeRek)
		f.SetCellStyle("Sheet2", colC, colC, dataStyle)

		f.SetCellValue("Sheet2", colD, record.NoBukti)
		f.SetCellStyle("Sheet2", colD, colD, dataStyle)

		f.SetCellValue("Sheet2", colE, record.Uraian)
		f.SetCellStyle("Sheet2", colE, colE, dataStyle)

		f.SetCellValue("Sheet2", colF, record.Pengeluaran)
		f.SetCellStyle("Sheet2", colF, colF, numberStyle)

		row++
	}

	// Save
	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("failed to save Template Import Excel: %w", err)
	}

	return nil
}
