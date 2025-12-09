package resp

// RunePage 符文页
type RunePage struct {
	Id                     int    `json:"id"`
	Name                   string `json:"name"`
	Current                bool   `json:"current"`
	Order                  int    `json:"order"`
	PrimaryStyleId         int    `json:"primaryStyleId"`
	SubStyleId             int    `json:"subStyleId"`
	SelectedPerkIds        []int  `json:"selectedPerkIds"`
	IsDeletable            bool   `json:"isDeletable"`
	IsEditable             bool   `json:"isEditable"`
	IsActive               bool   `json:"isActive"`
	IsValid                bool   `json:"isValid"`
	IsTemporary            bool   `json:"isTemporary"`
	AutoModifiedSelections []int  `json:"autoModifiedSelections"`
}

// RuneStyle 符文系统样式
type RuneStyle struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Icon  string     `json:"icon"`
	Slots []RuneSlot `json:"slots"`
}

// RuneSlot 符文槽位
type RuneSlot struct {
	Runes []Rune `json:"runes"`
}

// Rune 单个符文
type Rune struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}
