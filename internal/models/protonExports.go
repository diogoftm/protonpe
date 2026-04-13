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
	ItemEmail            string   `json:"itemEmail"`
	Password             string   `json:"password"`
	Urls                 []string `json:"url"`
	TotpUri              string   `json:"totpUri"`
	ItemUsername         string   `json:"itemUsername"`
	CardholderName       string   `json:"cardholderName"`
	Number               string   `json:"number"`
	VerificationNumber   string   `json:"verificationNumber"`
	ExpirationDate       string   `json:"expirationDate"`
	Pin                  string   `json:"pin"`
	FullName             string   `json:"fullName"`
	Email                string   `json:"email"`
	PhoneNumber          string   `json:"phoneNumber"`
	FirstName            string   `json:"firstName"`
	MiddleName           string   `json:"middleName"`
	LastName             string   `json:"lastName"`
	Birthdate            string   `json:"birthdate"`
	Gender               string   `json:"gender"`
	Organization         string   `json:"organization"`
	StreetAddress        string   `json:"streetAddress"`
	ZipOrPostalCode      string   `json:"zipOrPostalCode"`
	City                 string   `json:"city"`
	StateOrProvince      string   `json:"stateOrProvince"`
	CountryOrRegion      string   `json:"countryOrRegion"`
	Floor                string   `json:"floor"`
	County               string   `json:"county"`
	SocialSecurityNumber string   `json:"socialSecurityNumber"`
	PassportNumber       string   `json:"passportNumber"`
	LicenseNumber        string   `json:"licenseNumber"`
	Website              string   `json:"website"`
	XHandle              string   `json:"xHandle"`
	SecondPhoneNumber    string   `json:"secondPhoneNumber"`
	Linkedin             string   `json:"linkedin"`
	Reddit               string   `json:"reddit"`
	Facebook             string   `json:"facebook"`
	Yahoo                string   `json:"yahoo"`
	Instagram            string   `json:"instagram"`
	Company              string   `json:"company"`
	JobTitle             string   `json:"jobTitle"`
	PersonalWebsite      string   `json:"personalWebsite"`
	WorkPhoneNumber      string   `json:"workPhoneNumber"`
	WorkEmail            string   `json:"workEmail"`
}
