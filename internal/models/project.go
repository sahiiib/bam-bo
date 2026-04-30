package models

import "time"

type Texture struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Image string `json:"image"`
}

type PlyMaterial struct {
	Name      string  `json:"name"`
	ID        int     `json:"id"`
	Thickness float64 `json:"thickness"`
	ColorID   int     `json:"colorID"`
	TextureID int     `json:"textureID"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
}

func DefaultPlyMaterials() []PlyMaterial {
	return []PlyMaterial{
		{ID: 1, Name: "Cabinet Ply 18mm", Thickness: 1.8, TextureID: 1, Width: 122, Height: 244},
		{ID: 2, Name: "Back Panel 3mm", Thickness: 0.3, TextureID: 2, Width: 122, Height: 244},
		{ID: 3, Name: "Face Ply 18mm", Thickness: 1.8, TextureID: 1, Width: 122, Height: 244},
	}
}

func DefaultTextures() []Texture {
	return []Texture{
		{ID: 1, Name: "Natural Wood", Image: ""},
		{ID: 2, Name: "Back Panel Plain", Image: ""},
	}
}

func NormalizeMaterials(materials []PlyMaterial) []PlyMaterial {
	if len(materials) == 0 {
		return DefaultPlyMaterials()
	}

	normalized := make([]PlyMaterial, 0, len(materials))
	nextID := 1
	for _, material := range materials {
		if material.ID >= nextID {
			nextID = material.ID + 1
		}
	}

	for _, material := range materials {
		if material.ID == 0 {
			material.ID = nextID
			nextID++
		}
		if material.Name == "" {
			material.Name = "Material"
		}
		if material.Thickness <= 0 {
			material.Thickness = 1.8
		}
		if material.Width <= 0 {
			material.Width = 122
		}
		if material.Height <= 0 {
			material.Height = 244
		}
		normalized = append(normalized, material)
	}

	return normalized
}

func NormalizeTextures(textures []Texture) []Texture {
	if len(textures) == 0 {
		return DefaultTextures()
	}

	normalized := make([]Texture, 0, len(textures))
	nextID := 1
	for _, texture := range textures {
		if texture.ID >= nextID {
			nextID = texture.ID + 1
		}
	}

	for _, texture := range textures {
		if texture.ID == 0 {
			texture.ID = nextID
			nextID++
		}
		if texture.Name == "" {
			texture.Name = "Texture"
		}
		normalized = append(normalized, texture)
	}

	return normalized
}

func MaterialByID(materials []PlyMaterial, id int) (PlyMaterial, bool) {
	for _, material := range NormalizeMaterials(materials) {
		if material.ID == id {
			return material, true
		}
	}

	return PlyMaterial{}, false
}

type Piece struct {
	CabinetName string  `json:"cabinetName"`
	Name        string  `json:"name"`
	ID          int     `json:"id"`
	Count       int     `json:"count"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Thickness   float64 `json:"thickness"`
	MaterialID  int     `json:"materialID"`
	PVCSides    string  `json:"pvcSides"`
}

type Cabinet struct {
	Name              string  `json:"name"`
	ID                string  `json:"id,omitempty"`
	Width             float64 `json:"width"`
	Depth             float64 `json:"depth"`
	Height            float64 `json:"height"`
	Mode              int     `json:"mode"`
	CabinetMaterialID int     `json:"cabinetMaterialID"`
	BackMaterialID    int     `json:"backMaterialID"`
	FaceMaterialID    int     `json:"faceMaterialID"`
	Thickness         float64 `json:"thickness"`
	Back              float64 `json:"back"`
}

type Project struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Cabinets  []Cabinet     `json:"cabinets"`
	Materials []PlyMaterial `json:"materials"`
	Textures  []Texture     `json:"textures"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type ProjectSummary struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	CabinetCount int       `json:"cabinetCount"`
	PieceCount   int       `json:"pieceCount"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ProjectPayload struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Cabinets  []Cabinet     `json:"cabinets"`
	Materials []PlyMaterial `json:"materials"`
	Textures  []Texture     `json:"textures"`
}

type ExportRequest struct {
	ProjectName string        `json:"projectName"`
	Cabinets    []Cabinet     `json:"cabinets"`
	Materials   []PlyMaterial `json:"materials"`
	Textures    []Texture     `json:"textures"`
}

type ExportResponse struct {
	ProjectName string        `json:"projectName"`
	Cabinets    []Cabinet     `json:"cabinets"`
	Pieces      []Piece       `json:"pieces"`
	Materials   []PlyMaterial `json:"materials"`
	Textures    []Texture     `json:"textures"`
}
