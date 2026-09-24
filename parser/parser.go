package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Record struct {
	Date        string
	KodeKeg     string
	KodeRek     string
	NoBukti     string
	Uraian      string
	Penerimaan  float64
	Pengeluaran float64
	Saldo       float64
}

type HeaderInfo struct {
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
	TanggalTutup  string
}

func ParsePDF(path string) ([]Record, HeaderInfo, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, HeaderInfo{}, fmt.Errorf("failed to open PDF: %w", err)
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

	headerInfo := extractHeader(rawText)
	records := extractRecords(rawText)

	return records, headerInfo, nil
}

func extractHeader(text string) HeaderInfo {
	h := HeaderInfo{}

	reBulan := regexp.MustCompile(`BULAN\s*:\s*(\w+)\s*TAHUN\s*:\s*(\d{4})`)
	if match := reBulan.FindStringSubmatch(text); len(match) > 2 {
		h.Month = match[1]
		h.Year = match[2]
	}

	reNPSN := regexp.MustCompile(`NPSN\s*:\s*(\d+)`)
	if match := reNPSN.FindStringSubmatch(text); len(match) > 1 {
		h.NPSN = match[1]
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "Nama Sekolah") && strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				if strings.Contains(val, "Nama Sekolah") {
					idx := strings.LastIndex(val, "Nama Sekolah")
					val = strings.TrimSpace(val[idx+13:])
				}
				if val != "" {
					h.NamaSekolah = val
				}
			}
		}
		if strings.Contains(trimmed, "Desa") || strings.Contains(trimmed, "Kecamatan") {
			if strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					if val != "" && h.DesaKecamatan == "" {
						h.DesaKecamatan = val
					}
				}
			}
		}
		if strings.Contains(trimmed, "Kabupaten") || strings.Contains(trimmed, "Kota") {
			if strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					if val != "" && h.KabupatenKota == "" {
						h.KabupatenKota = val
					}
				}
			}
		}
		if strings.Contains(trimmed, "Provinsi") || strings.Contains(trimmed, "Prov.") {
			if strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					h.Provinsi = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	reSumber := regexp.MustCompile(`Sumber Dana\s*:\s*(.+)`)
	if match := reSumber.FindStringSubmatch(text); len(match) > 1 {
		h.SumberDana = strings.TrimSpace(match[1])
	}

	reTanggal := regexp.MustCompile(`(\d{2}\s+\w+\s+\d{4})\s*Buku Kas Umum Ditutup`)
	if match := reTanggal.FindStringSubmatch(text); len(match) > 1 {
		h.TanggalTutup = match[1]
	}

	// Extract Kepala Sekolah and Bendahara from signature section at end of document
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "Kepala Sekolah") {
			// Next non-empty lines contain name and NIP
			for j := i + 1; j < len(lines) && j < i+10; j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}
				if h.KepalaSekolah == "" && !strings.Contains(nextLine, "NIP.") {
					h.KepalaSekolah = nextLine
				}
				if strings.Contains(nextLine, "NIP.") || strings.HasPrefix(nextLine, "NIP") {
					h.NIPKepsek = nextLine
					break
				}
			}
		}
		if strings.EqualFold(trimmed, "Bendahara,") || strings.EqualFold(trimmed, "Bendahara") {
			for j := i + 1; j < len(lines) && j < i+10; j++ {
				nextLine := strings.TrimSpace(lines[j])
				if nextLine == "" {
					continue
				}
				if h.Bendahara == "" && !strings.Contains(nextLine, "NIP.") {
					h.Bendahara = nextLine
				}
				if strings.Contains(nextLine, "NIP.") || strings.HasPrefix(nextLine, "NIP") {
					h.NIPBendahara = nextLine
					break
				}
			}
		}
	}

	return h
}

func extractRecords(text string) []Record {
	dateRe := regexp.MustCompile(`^(\d{2}-\d{2}-\d{4})$`)
	numRe := regexp.MustCompile(`^(0|[\d]{1,3}(\.\d{3})+)$`)
	kodeKegRe := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\.$`)
	kodeRekRe := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+\.\d+\.\d+$`)
	kodeRekContRe := regexp.MustCompile(`^\d{2,4}$`)
	noBuktiRe := regexp.MustCompile(`^(BNU|BPU)\d+$`)
	jumlahRe := regexp.MustCompile(`(?i)^jumlah`)

	lines := strings.Split(text, "\n")

	// First, collect all non-header lines with their types
	type Line struct {
		text string
		t    string // "date", "desc", "num", "kodekeg", "koderek", "nobukti"
	}

	var dataLines []Line
	inData := false
	headerNumbersSeen := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Track header numbers 1-8 to detect when data starts
		if !inData {
			if trimmed >= "1" && trimmed <= "8" && len(trimmed) == 1 {
				headerNumbersSeen[trimmed] = true
				if len(headerNumbersSeen) >= 6 { // we've seen most number headers
					inData = true
				}
				continue
			}
			// Skip all header lines until first date
			if dateRe.MatchString(trimmed) {
				inData = true
				// fall through to process this date line
			} else {
				continue
			}
		}

		// Skip headers/footers - be careful not to filter actual data
		if strings.Contains(trimmed, "Halaman") || strings.Contains(trimmed, "dari") ||
			strings.Contains(trimmed, "NPSN") || strings.Contains(trimmed, "Nama Sekolah") ||
			strings.Contains(trimmed, "Sumber Dana") || strings.Contains(trimmed, "BULAN") ||
			strings.Contains(trimmed, "TANGGAL") || strings.Contains(trimmed, "KEGIATAN") ||
			strings.Contains(trimmed, "REKENING") || strings.Contains(trimmed, "BUKTI") ||
			strings.Contains(trimmed, "URAIAN") || strings.Contains(trimmed, "PENERIMAAN") ||
			strings.Contains(trimmed, "PENGELUARAN") || strings.Contains(trimmed, "Menyetujui") ||
			strings.Contains(trimmed, "Kepala Sekolah") || strings.Contains(trimmed, "Bendahara") ||
			strings.Contains(trimmed, "NIP.") || strings.Contains(trimmed, "Kec. Arjawinangun") ||
			strings.Contains(trimmed, "BKU") || strings.Contains(trimmed, "SD NEGERI") ||
			strings.Contains(trimmed, "JUNGJANG") || strings.Contains(trimmed, "RUSIYANTO") ||
			strings.Contains(trimmed, "SEVI") {
			continue
		}

		// Skip closing section lines
		if strings.HasPrefix(trimmed, "Saldo Buku Kas") || strings.HasPrefix(trimmed, "Saldo  Buku Kas") ||
			strings.HasPrefix(trimmed, "Terdiri Dari") || strings.HasPrefix(trimmed, "- Saldo Bank") ||
			strings.HasPrefix(trimmed, "- Saldo Kas") || strings.HasPrefix(trimmed, "Jumlah :") ||
			strings.HasPrefix(trimmed, "Pada hari ini") || strings.HasPrefix(trimmed, "Ditutup") ||
			strings.HasPrefix(trimmed, ": Rp.") {
			continue
		}

		if jumlahRe.MatchString(trimmed) {
			break
		}

		// Detect line type - order matters!
		// Check for kode rekening continuation FIRST (2-4 digit numbers that follow a koderek)
		isKodeRekCont := kodeRekContRe.MatchString(trimmed) && len(dataLines) > 0 && dataLines[len(dataLines)-1].t == "koderek"

		switch {
		case dateRe.MatchString(trimmed):
			dataLines = append(dataLines, Line{trimmed, "date"})
		case kodeKegRe.MatchString(trimmed):
			dataLines = append(dataLines, Line{trimmed, "kodekeg"})
		case kodeRekRe.MatchString(trimmed):
			dataLines = append(dataLines, Line{trimmed, "koderek"})
		case isKodeRekCont:
			dataLines = append(dataLines, Line{trimmed, "koderek"})
		case noBuktiRe.MatchString(trimmed):
			dataLines = append(dataLines, Line{trimmed, "nobukti"})
		case numRe.MatchString(trimmed):
			dataLines = append(dataLines, Line{trimmed, "num"})
		default:
			dataLines = append(dataLines, Line{trimmed, "desc"})
		}
	}

	// Now reconstruct records from the typed lines
	var records []Record
	var currentRecord *Record
	var pendingKodeKeg, pendingKodeRek, pendingNoBukti string
	numCount := 0
	descParts := []string{}

	for _, line := range dataLines {
		switch line.t {
		case "date":
			// Save previous record
			if currentRecord != nil {
				// Add any pending metadata
				currentRecord.KodeKeg = pendingKodeKeg
				currentRecord.KodeRek = pendingKodeRek
				currentRecord.NoBukti = pendingNoBukti
				records = append(records, *currentRecord)
			}

			// Clear pending for new record
			pendingKodeKeg = ""
			pendingKodeRek = ""
			pendingNoBukti = ""

			// Start new record
			currentRecord = &Record{Date: line.text}
			descParts = []string{}
			numCount = 0

		case "desc":
			if currentRecord != nil {
				descParts = append(descParts, line.text)
				numCount = 0
			}

		case "num":
			if currentRecord != nil {
				numCount++
				cleaned := strings.ReplaceAll(line.text, ".", "")
				val, _ := strconv.ParseFloat(cleaned, 64)

				switch numCount {
				case 1:
					currentRecord.Penerimaan = val
				case 2:
					currentRecord.Pengeluaran = val
				case 3:
					currentRecord.Saldo = val
					// Record is complete, save description
					currentRecord.Uraian = strings.Join(descParts, " ")
				}
			}

		case "kodekeg":
			pendingKodeKeg = line.text
		case "koderek":
			// Check if this is a continuation of a previous koderek
			if kodeRekContRe.MatchString(line.text) && pendingKodeRek != "" {
				pendingKodeRek += line.text
			} else {
				pendingKodeRek = line.text
			}
		case "nobukti":
			pendingNoBukti = line.text
		}
	}

	// Don't forget the last record
	if currentRecord != nil {
		currentRecord.KodeKeg = pendingKodeKeg
		currentRecord.KodeRek = pendingKodeRek
		currentRecord.NoBukti = pendingNoBukti
		records = append(records, *currentRecord)
	}

	return records
}
