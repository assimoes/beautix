	// Relationships
	Staff                []Staff              `gorm:"foreignKey:UserID" json:"staff,omitempty"`

    Why do Users have relations with Staff? Staff should be connected with businesses, no?