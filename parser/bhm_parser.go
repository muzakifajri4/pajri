package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

type BHMRecord struct {
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

type BHMHeaderInfo struct {
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

func ParseBHM(path string) ([]BHMRecord, BHMHeaderInfo, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, BHMHeaderInfo{}, fmt.Errorf("failed to open PDF: %w", err)
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
	headerInfo := extractBHMHeader(rawText)
	records := extractBHMRecords(rawText)

	return records, headerInfo, nil
}

func extractBHMHeader(text string) BHMHeaderInfo {
	h := BHMHeaderInfo{}

	// Inline formats (label and value on the same line)
	if m := regexp.MustCompile(`BULAN\s*:\s*(\w+)\s*TAHUN\s*:\s*(\d{4})`).FindStringSubmatch(text); len(m) > 2 {
		h.Month = m[1]
		h.Year = m[2]
	}
	if m := regexp.MustCompile(`NPSN\s*:\s*(\d+)`).FindStringSubmatch(text); len(m) > 1 {
		h.NPSN = m[1]
	}

	lines := strings.Split(text, "\n")

	// Stacked format: "Nama Sekolah" <nl> ": SD NEGERI 5 ..."
	// Labels appear first, then the ": value" lines in the same order.
	var pendingLabels []string
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" {
			continue
		}
		switch {
		case t == "Nama Sekolah":
			pendingLabels = append(pendingLabels, "nama")
		case t == "Desa/Kecamatan":
			pendingLabels = append(pendingLabels, "desa")
		case t == "Kabupaten / Kota" || t == "Kabupaten/Kota":
			pendingLabels = append(pendingLabels, "kab")
		case t == "Provinsi":
			pendingLabels = append(pendingLabels, "prov")
		case t == "NPSN":
			pendingLabels = append(pendingLabels, "npsn")
		case t == "Sumber Dana":
			pendingLabels = append(pendingLabels, "sumber")
		default:
			if strings.HasPrefix(t, ":") && len(pendingLabels) > 0 {
				val := strings.TrimSpace(strings.TrimPrefix(t, ":"))
				label := pendingLabels[0]
				switch label {
				case "nama":
					if h.NamaSekolah == "" {
						h.NamaSekolah = val
					}
				case "desa":
					if h.DesaKecamatan == "" {
						h.DesaKecamatan = val
					}
				case "kab":
					if h.KabupatenKota == "" {
						h.KabupatenKota = val
					}
				case "prov":
					if h.Provinsi == "" {
						h.Provinsi = val
					}
				case "npsn":
					if h.NPSN == "" {
						h.NPSN = val
					}
				case "sumber":
					if h.SumberDana == "" {
						h.SumberDana = val
					}
				}
				pendingLabels = pendingLabels[1:]
			}
		}
	}

	// Fallback: footer lines like "Nama Sekolah : SD NEGERI ..." on one line
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.Contains(t, "Nama Sekolah") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, "Nama Sekolah", 2); len(p) == 2 {
				rest := strings.TrimSpace(strings.TrimPrefix(p[1], ":"))
				if rest != "" && h.NamaSekolah == "" {
					h.NamaSekolah = rest
				}
			}
		}
		if strings.Contains(t, "Sumber Dana") && strings.Contains(t, ":") {
			if p := strings.SplitN(t, "Sumber Dana", 2); len(p) == 2 {
				rest := strings.TrimSpace(strings.TrimPrefix(p[1], ":"))
				if rest != "" && h.SumberDana == "" {
					h.SumberDana = rest
				}
			}
		}
	}

	// Signature block on the last page:
	//   Menyetujui,
	//   Kepala Sekolah
	//   RUSIYANTO, S.Pd.SD
	//   NIP. 196902192000031002
	//   Bendahara,
	//   SEVI SEPTRIYANI, S.Pd
	//   NIP. 199809052022212001
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "Kepala Sekolah" && h.KepalaSekolah == "" {
			if i+1 < len(lines) {
				name := strings.TrimSpace(lines[i+1])
				if name != "" && !strings.HasPrefix(name, "NIP.") && name != "Menyetujui," {
					h.KepalaSekolah = name
					if i+2 < len(lines) {
						nip := strings.TrimSpace(lines[i+2])
						if strings.HasPrefix(nip, "NIP.") {
							h.NIPKepsek = nip
						}
					}
				}
			}
		}
		if (t == "Bendahara," || t == "Bendahara") && h.Bendahara == "" {
			if i+1 < len(lines) {
				name := strings.TrimSpace(lines[i+1])
				if name != "" && !strings.HasPrefix(name, "NIP.") {
					h.Bendahara = name
					if i+2 < len(lines) {
						nip := strings.TrimSpace(lines[i+2])
						if strings.HasPrefix(nip, "NIP.") {
							h.NIPBendahara = nip
						}
					}
				}
			}
		}
	}

	return h
}
func extractBHMRecords(text string) []BHMRecord {
	lines := strings.Split(text, "\n")

	dateRe := regexp.MustCompile(`^(\d{2}-\d{2}-\d{4})$`)
	kodeKegRe := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\.$`)
	kodeRekRe := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+\.\d+\.\d+$`)
	noBuktiRe := regexp.MustCompile(`^(BNU|BPU)(\d+)$`)
	idBarangLine1Re := regexp.MustCompile(`^1\.3\.`)
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

		skipWords := []string{"Halaman", "dari", "NPSN", "Nama Sekolah", "BHM",
			"Tanggal", "Kode", "Kegiatan", "Rekening", "Bukti", "ID Barang",
			"Uraian", "Barang", "Harga", "Satuan", "Realisasi", "Sumber Dana"}
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

	var records []BHMRecord
	var cur BHMRecord
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
					cur.IDBarang = strings.Join(idLines, "\n")
				}
				records = append(records, cur)
			}

			// Start new record
			cur = BHMRecord{Date: line}
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
			cur.IDBarang = strings.Join(idLines, "\n")
		}
		records = append(records, cur)
	}

	return records
}