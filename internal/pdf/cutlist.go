package pdf

import (
	"bam-bo/internal/models"
	"fmt"
	"strings"
)

type cabinetMaterials struct {
	cabinet models.PlyMaterial
	back    models.PlyMaterial
	face    models.PlyMaterial
}

func calculateMDFPieces(cab models.Cabinet, materials cabinetMaterials) []models.Piece {
	var pieces []models.Piece
	cabinetThickness := materials.cabinet.Thickness
	backThickness := materials.back.Thickness
	faceThickness := materials.face.Thickness

	switch cab.Mode {
	case 1:
		pieces = []models.Piece{
			{CabinetName: cab.Name, Name: "Bottom", ID: 1, Length: cab.Depth, Width: cab.Width, Thickness: cabinetThickness, Count: 1, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Sides", ID: 2, Length: cab.Depth, Width: cab.Height - cabinetThickness, Thickness: cabinetThickness, Count: 2, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Shelf", ID: 3, Length: cab.Width - 2*cabinetThickness, Width: cab.Depth - cabinetThickness - backThickness, Thickness: cabinetThickness, Count: 1, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Back Side", ID: 4, Length: cab.Width - cabinetThickness - 0.2, Width: cab.Height - cabinetThickness/2 - 0.1, Thickness: backThickness, Count: 1, MaterialID: materials.back.ID},
			{CabinetName: cab.Name, Name: "Back Tie", ID: 5, Length: cab.Width - 2*cabinetThickness, Width: 8, Thickness: cabinetThickness, Count: 2, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Top Tie", ID: 6, Length: cab.Width - 2*cabinetThickness, Width: 8, Thickness: cabinetThickness, Count: 2, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Face", ID: 7, Length: cab.Height, Width: cab.Width, Thickness: faceThickness, Count: 1, MaterialID: materials.face.ID},
		}
	case 2:
		pieces = []models.Piece{
			{CabinetName: cab.Name, Name: "Bottom", ID: 1, Length: cab.Depth, Width: cab.Width - 2*cabinetThickness, Thickness: cabinetThickness, Count: 1, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Sides", ID: 2, Length: cab.Depth, Width: cab.Height, Thickness: cabinetThickness, Count: 2, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Top/Shelf", ID: 3, Length: cab.Width - 2*cabinetThickness, Width: cab.Depth - cabinetThickness - backThickness, Thickness: cabinetThickness, Count: 1, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Back/Tie", ID: 4, Length: cab.Width - 2*cabinetThickness, Width: 8, Thickness: cabinetThickness, Count: 2, MaterialID: materials.cabinet.ID},
			{CabinetName: cab.Name, Name: "Back Side", ID: 5, Length: cab.Width - cabinetThickness - 0.2, Width: cab.Height - cabinetThickness/2 - 0.1, Thickness: backThickness, Count: 1, MaterialID: materials.back.ID},
			{CabinetName: cab.Name, Name: "Face", ID: 6, Length: cab.Height, Width: cab.Width, Thickness: faceThickness, Count: 1, MaterialID: materials.face.ID},
		}
	}

	return pieces
}

func GatherPieces(cabinets []models.Cabinet, materials []models.PlyMaterial) ([]models.Piece, error) {
	pieces := make([]models.Piece, 0)
	materials = models.NormalizeMaterials(materials)

	for i, cab := range cabinets {
		if strings.TrimSpace(cab.Name) == "" {
			cab.Name = fmt.Sprintf("Cabinet %d", i+1)
		}

		if cab.Width <= 0 || cab.Depth <= 0 || cab.Height <= 0 {
			return nil, fmt.Errorf("cabinet %q has invalid dimensions", cab.Name)
		}

		cabinetMaterials, err := resolveCabinetMaterials(cab, materials)
		if err != nil {
			return nil, fmt.Errorf("cabinet %q: %w", cab.Name, err)
		}

		cabPieces := calculateMDFPieces(cab, cabinetMaterials)
		if len(cabPieces) == 0 {
			return nil, fmt.Errorf("cabinet %q has unsupported mode; use 1 or 2", cab.Name)
		}

		pieces = append(pieces, cabPieces...)
	}

	return pieces, nil
}

func resolveCabinetMaterials(cab models.Cabinet, materials []models.PlyMaterial) (cabinetMaterials, error) {
	cabinet, err := materialOrDefault(materials, cab.CabinetMaterialID, 1, cab.Thickness)
	if err != nil {
		return cabinetMaterials{}, fmt.Errorf("invalid cabinet material: %w", err)
	}

	back, err := materialOrDefault(materials, cab.BackMaterialID, 2, cab.Back)
	if err != nil {
		return cabinetMaterials{}, fmt.Errorf("invalid back material: %w", err)
	}

	face, err := materialOrDefault(materials, cab.FaceMaterialID, 3, 0)
	if err != nil {
		return cabinetMaterials{}, fmt.Errorf("invalid face material: %w", err)
	}

	return cabinetMaterials{cabinet: cabinet, back: back, face: face}, nil
}

func materialOrDefault(materials []models.PlyMaterial, id, defaultID int, legacyThickness float64) (models.PlyMaterial, error) {
	usesLegacyID := id == 0
	if id == 0 {
		id = defaultID
	}

	material, ok := models.MaterialByID(materials, id)
	if !ok {
		return models.PlyMaterial{}, fmt.Errorf("unknown material id %d", id)
	}

	if usesLegacyID && legacyThickness > 0 {
		material.Thickness = legacyThickness
	}

	if material.Thickness <= 0 {
		return models.PlyMaterial{}, fmt.Errorf("material %q has invalid thickness", material.Name)
	}

	return material, nil
}

func TotalPieceCount(pieces []models.Piece) int {
	total := 0
	for _, piece := range pieces {
		total += piece.Count
	}
	return total
}
