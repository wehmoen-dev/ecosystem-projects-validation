package data

type Structure struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Websites    []Website  `json:"websites" validate:"required"`
	Contracts   []Contract `json:"contracts" validate:"required"`
	Categories  []Category `json:"categories" validate:"required"`
	Email       *string    `json:"email,omitempty" validate:"email,omitempty"`
	Social      *Social    `json:"social,omitempty"`
}

type Website struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Contract struct {
	Address     string `json:"address"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type Platform string

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformTwitter   Platform = "twitter"
	PlatformX         Platform = "x"
	PlatformLinkedIn  Platform = "linkedin"
	PlatformThreads   Platform = "threads"
	PlatformMastodon  Platform = "mastodon"
	PlatformTelegram  Platform = "telegram"
	PlatformDiscord   Platform = "discord"
)

type Category string

const (
	CategoryGame    Category = "game"
	CategoryNFT     Category = "nft"
	CategoryFinance Category = "finance"
	CategoryDAO     Category = "dao"
	CategoryTool    Category = "tool"
	CategoryOther   Category = "other"
)

type Social map[Platform]string
