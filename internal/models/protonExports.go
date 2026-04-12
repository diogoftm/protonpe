package models

type VaultFile struct {
	Version string           `json:"version"`
	Vaults  map[string]Vault `json:"vaults"`
	UserId  string           `json:"userId"`
}

type Vault struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Display     DisplayOptions `json:"display"`
	Items       []Item         `json:"items"`
}

type DisplayOptions struct {
	Color int `json:"color"`
	Icon  int `json:"icon"`
}

type Item struct {
	ItemId               string   `json:"itemId"`
	ShareId              string   `json:"shareId"`
	AliasEmail           string   `json:"aliasEmail"`
	ContentFormatVersion int      `json:"contentFormatVersion"`
	CreateTime           int      `json:"createTime"`
	ModifyTime           int      `json:"modifyTime"`
	Pinned               bool     `json:"pinned"`
	ShareCount           int      `json:"shareCount"`
	State                int      `json:"state"`
	Data                 ItemData `json:"data"`
}

type ItemData struct {
	Metadata ItemMetadata `json:"metadata"`
	Type     string       `json:"type"`
	Content  ItemContent  `json:"content"`
}

// files ignored
type ItemMetadata struct {
	Name     string `json:"name"`
	Note     string `json:"note"`
	ItemUuid string `json:"itemUuid"`
}

// passkeys ignored
type ItemContent struct {
	ItemEmail          string   `json:"itemEmail"`
	Password           string   `json:"password"`
	Urls               []string `json:"url"`
	TotpUri            string   `json:"totpUri"`
	ItemUsername       string   `json:"itemUsername"`
	CardholderName     string   `json:"cardholderName"`
	Number             string   `json:"number"`
	VerificationNumber string   `json:"verificationNumber"`
	ExpirationDate     string   `json:"expirationDate"`
	Pin                string   `json:"pin"`
}
