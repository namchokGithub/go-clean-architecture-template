package enums

// Enum is the base type for typed enumerations with bilingual labels.
type Enum struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`    // Thai
	NameEn string `json:"name_en"` // English
}

// UserType enumerates valid user type values.
var UserType = struct {
	Admin  Enum
	Driver Enum
}{
	Admin:  Enum{ID: 1, Name: "แอดมิน", NameEn: "Admin"},
	Driver: Enum{ID: 2, Name: "คนขับ", NameEn: "Driver"},
}
