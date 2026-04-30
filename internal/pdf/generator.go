package pdf

import (
	"bam-bo/internal/models"
	"bytes"
	"fmt"
	"strings"
)

func Build(projectName string, cabinets []models.Cabinet, pieces []models.Piece, materials []models.PlyMaterial) ([]byte, error) {
	title := "Bam-Bo Cutting List"
	if strings.TrimSpace(projectName) != "" {
		title = "Bam-Bo Cutting List - " + projectName
	}

	lines := []string{
		title,
		fmt.Sprintf("Cabinets: %d", len(cabinets)),
		fmt.Sprintf("Pieces: %d", TotalPieceCount(pieces)),
		"",
	}

	for _, group := range groupPiecesByMaterial(pieces, materials) {
		lines = append(lines, group.name)
		for _, piece := range group.pieces {
			line := fmt.Sprintf(
				"%s | %s | qty %d | %.2f x %.2f x %.2f",
				piece.CabinetName,
				piece.Name,
				piece.Count,
				piece.Length,
				piece.Width,
				piece.Thickness,
			)
			lines = append(lines, line)
		}
		lines = append(lines, "")
	}

	return renderPDF(buildPDFContent(lines)), nil
}

type materialGroup struct {
	id     int
	name   string
	pieces []models.Piece
}

func groupPiecesByMaterial(pieces []models.Piece, materials []models.PlyMaterial) []materialGroup {
	indexes := map[int]int{}
	groups := make([]materialGroup, 0)
	materials = models.NormalizeMaterials(materials)

	for _, piece := range pieces {
		groupIndex, ok := indexes[piece.MaterialID]
		if !ok {
			materialName := fmt.Sprintf("Material %d", piece.MaterialID)
			if material, found := models.MaterialByID(materials, piece.MaterialID); found {
				materialName = material.Name
			}

			groupIndex = len(groups)
			indexes[piece.MaterialID] = groupIndex
			groups = append(groups, materialGroup{id: piece.MaterialID, name: materialName})
		}

		groups[groupIndex].pieces = append(groups[groupIndex].pieces, piece)
	}

	return groups
}

func buildPDFContent(lines []string) string {
	var b strings.Builder
	y := 810

	for index, line := range lines {
		if y < 50 {
			break
		}

		fontSize := 10
		if index == 0 {
			fontSize = 16
		}

		b.WriteString("BT\n")
		b.WriteString(fmt.Sprintf("/F1 %d Tf\n", fontSize))
		b.WriteString(fmt.Sprintf("1 0 0 1 50 %d Tm\n", y))
		b.WriteString(pdfText(line))
		b.WriteString(" Tj\nET\n")
		y -= 18
	}

	return b.String()
}

func renderPDF(content string) []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", len(content), content),
	}

	var buf bytes.Buffer
	offsets := make([]int, len(objects)+1)

	buf.WriteString("%PDF-1.4\n")
	for i, obj := range objects {
		offsets[i+1] = buf.Len()
		buf.WriteString(fmt.Sprintf("%d 0 obj\n%s\nendobj\n", i+1, obj))
	}

	xrefOffset := buf.Len()
	buf.WriteString(fmt.Sprintf("xref\n0 %d\n", len(objects)+1))
	buf.WriteString("0000000000 65535 f \n")
	for i := 1; i <= len(objects); i++ {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}

	buf.WriteString(fmt.Sprintf(
		"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF",
		len(objects)+1,
		xrefOffset,
	))

	return buf.Bytes()
}

func pdfText(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "(", `\(`)
	value = strings.ReplaceAll(value, ")", `\)`)
	return "(" + value + ")"
}
