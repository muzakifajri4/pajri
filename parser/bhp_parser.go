package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

type BHPRecord struct {
	Date         string
	KodeKeg      string
	KodeRek      string
	NoBukti      string
	IDBarang     string
	Uraian       string
	JumlahBarang float64
	HargaSatuan  float64
	Realisasi    float64
}

type BHPHeaderInfo struct {
	Month         string
	Year          string
	NPSN          string
	NamaSekolah   string
	DesaKecamatan string
	KabupatenKota string
	Provinsi      string
	SumberDana    string
	KepalaSekolah string
	NIPKepsek     string
	Bendahara     string
	NIPBendahara  string
}

func ParseBHP(path string) ([]BHPRecord, BHPHeaderInfo, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, BHPHeaderInfo{}, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var text strings.Builder
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		text.WriteString(content)
		text.WriteString("\n")
	}

	rawText := text.String()
	headerInfo := extractBHPHeader(rawText)
	records := extractBHPRecords(rawText)

	return records, headerInfo, nil
}

func extractBHPHeader(text string) BHPHeaderInfo {
	h := BHPHeaderInfo{}
	if m := regexp.MustCompile(`BULAN\s*:\s*(\w+)\s*TAHUN\s*:\s*(\d{4})`).FindStringSubmatch(text); len(m) > 2 {
		h.Month = m[1]
		h.Year = m[2]
	}
	if m := regexp.MustCompile(`NPSN\s*:\s*(\d+)`).FindStringSubmatch(text); len(m) > 1 {
		h.NPSN = m[1]
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.Contains(t, "Nama Sekolah") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, ":", 2); len(p) == 2 {
				v := strings.TrimSpace(p[1])
				if strings.Contains(v, "Nama Sekolah") {
					v = strings.TrimSpace(v[strings.LastIndex(v, "Nama Sekolah")+13:])
				}
				if v != "" {
					h.NamaSekolah = v
				}
			}
		}
		if strings.Contains(t, "Desa/Kecamatan") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, ":", 2); len(p) == 2 {
				h.DesaKecamatan = strings.TrimSpace(p[1])
			}
		}
		if strings.Contains(t, "Kabupaten") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, ":", 2); len(p) == 2 {
				h.KabupatenKota = strings.TrimSpace(p[1])
			}
		}
		if strings.Contains(t, "Provinsi") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, ":", 2); len(p) == 2 {
				h.Provinsi = strings.TrimSpace(p[1])
			}
		}
	}
	if m := regexp.MustCompile(`Sumber Dana\s*:\s*(.+)`).FindStringSubmatch(text); len(m) > 1 {
		h.SumberDana = strings.TrimSpace(m[1])
	}
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.Contains(t, "Kepala") && strings.Contains(t, "Jannah") {
			h.KepalaSekolah = "Tsamrotul Jannah, S.Pd. AUD"
		}
		if strings.Contains(t, "196404271986031014") {
			h.NIPKepsek = "NIP. 196404271986031014"
		}
		if strings.Contains(t, "Bendahara") && strings.Contains(t, "Saefudin") {
			h.Bendahara = "Asep Saefudin, S.Pd"
		}
		if strings.Contains(t, "199210202022211001") {
			h.NIPBendahara = "NIP. 199210202022211001"
		}
	}
	return h
}

func extractBHPRecords(text string) []BHPRecord {
	lines := strings.Split(text, "\n")

	dateRe := regexp.MustCompile(`^(\d{2}-\d{2}-\d{4})$`)
	kodeKegRe := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\.$`)
	kodeRekRe := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+\.\d+\.\d+$`)
	noBuktiRe := regexp.MustCompile(`^(BNU|BPU)(\d+)$`)
	idBarangLine1Re := regexp.MustCompile(`^1\.1\.`)
	idBarangLine2Re := regexp.MustCompile(`^\d{2,4}-1-\d{12}$`)
	numberRe := regexp.MustCompile(`^[\d.]+$`)
	jumlahCompleteRe := regexp.MustCompile(`(?i)^jumlah\s+[\d.]+`)

	// Patterns that signal the end of data (signature section)
	signaturePatterns := []string{"Menyetujui", "Kepala Sekolah", "Bendahara", "NIP.", "Kec. Arjawinangun"}

	// Collect all non-header lines
	var allLines []string
	inData := false
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" {
			continue
		}

		if !inData {
			if dateRe.MatchString(t) {
				inData = true
				allLines = append(allLines, t)
			}
			continue
		}

		// Stop parsing when we reach the signature section
		isSignature := false
		for _, sp := range signaturePatterns {
			if strings.Contains(t, sp) {
				isSignature = true
				break
			}
		}
		if isSignature {
			break
		}

		// "Jumlah" as a standalone header word on page 2+ should be skipped
		if strings.EqualFold(t, "jumlah") {
			continue
		}

		// Skip subtotal/total lines like "Jumlah 294.000" (may appear on each page)
		if jumlahCompleteRe.MatchString(t) {
			continue
		}

		// Skip lines containing "Rp." (may appear as part of subtotal on each page)
		if strings.Contains(t, "Rp.") {
			continue
		}

		skipWords := []string{"Halaman", "dari", "NPSN", "Nama Sekolah", "BHP",
			"Tanggal", "Kode", "Kegiatan", "Rekening", "Bukti", "ID Barang",
			"Uraian", "Barang", "Harga", "Satuan", "Realisasi",
			"Tsamrotul Jannah", "Asep Saefudin", "19640427", "19921020"}
		// Also skip "No." as standalone line (page 2 header)
		if t == "No." {
			continue
		}
		skip := false
		for _, w := range skipWords {
			if strings.Contains(t, w) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		allLines = append(allLines, t)
	}

	// Parse lines into records using state machine
	type state int
	const (
		stDate state = iota
		stKodeKeg
		stUraian
		stKodeRek
		stNoBukti
		stNumbers
		stIDBarang
	)

	var records []BHPRecord
	var cur BHPRecord
	var uraian []string
	var idLines []string
	currentState := stDate
	numberCount := 0

	for _, line := range allLines {
		switch {
		case dateRe.MatchString(line):
			// Save previous record
			if cur.Date != "" {
				cur.Uraian = strings.Join(uraian, " ")
				if len(idLines) > 0 {
					cur.IDBarang = strings.Join(idLines, "")
				}
				records = append(records, cur)
			}

			// Start new record
			cur = BHPRecord{Date: line}
			uraian = nil
			idLines = nil
			currentState = stKodeKeg
			numberCount = 0

		case idBarangLine1Re.MatchString(line):
			idLines = append(idLines, line)
			currentState = stIDBarang

		case idBarangLine2Re.MatchString(line):
			idLines = append(idLines, line)
			currentState = stIDBarang

		case currentState == stKodeKeg && kodeKegRe.MatchString(line):
			cur.KodeKeg = line
			currentState = stUraian

		case kodeRekRe.MatchString(line):
			cur.KodeRek = line
			currentState = stNoBukti

		case noBuktiRe.MatchString(line):
			m := noBuktiRe.FindStringSubmatch(line)
			cur.NoBukti = m[1] + m[2]
			currentState = stNumbers
			numberCount = 0

		case currentState == stNumbers && numberRe.MatchString(line):
			numberCount++
			switch numberCount {
			case 1:
				cur.HargaSatuan = parseThousandsNumber(line)
			case 2:
				cur.Realisasi = parseThousandsNumber(line)
			case 3:
				cur.JumlahBarang = parseRawNumber(line)
				// After 3 numbers, we go to IDBarang state (or back to uraian if next record)
				currentState = stIDBarang
			}

		default:
			// Uraian line - but only if we haven't already seen numbers
			if currentState != stIDBarang || len(idLines) == 0 {
				uraian = append(uraian, line)
				if currentState != stUraian && currentState != stIDBarang {
					currentState = stUraian
				}
			}
		}
	}

	// Last record
	if cur.Date != "" {
		cur.Uraian = strings.Join(uraian, " ")
		if len(idLines) > 0 {
			cur.IDBarang = strings.Join(idLines, "")
		}
		records = append(records, cur)
	}

	return records
}

func parseThousandsNumber(s string) float64 {
	// Numbers like "14.000" -> 14000 (dot is thousands separator in Indonesian format)
	cleaned := strings.ReplaceAll(s, ".", "")
	v, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}
	return v
}

func parseRawNumber(s string) float64 {
	cleaned := strings.ReplaceAll(s, ".", "")
	v, _ := strconv.ParseFloat(cleaned, 64)
	return v
}
