package data

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"net/url"
	"strings"
)

func ValidateStructure(jsonData []byte) []error {
	var data Structure
	var errors []error

	// Unmarshal the JSON into the structure
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	// Use a map to hold the raw JSON
	var rawData map[string]interface{}
	err = json.Unmarshal(jsonData, &rawData)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	// Validate for extra keys in the root structure
	expectedFields := map[string]struct{}{
		"name":        {},
		"description": {},
		"websites":    {},
		"contracts":   {},
		"categories":  {},
		"email":       {},
		"social":      {},
	}

	for key := range rawData {
		if _, ok := expectedFields[key]; !ok {
			errors = append(errors, fmt.Errorf("extra key: %s", key))
		}
	}

	// Validate Websites for extra keys and field correctness
	if websitesRaw, ok := rawData["websites"].([]interface{}); ok {
		expectedWebsiteFields := map[string]struct{}{
			"url":         {},
			"description": {},
		}

		var usedUrls = make(map[string]bool)
		for i, website := range websitesRaw {
			if websiteMap, ok := website.(map[string]interface{}); ok {
				for key := range websiteMap {
					if _, ok := expectedWebsiteFields[key]; !ok {
						errors = append(errors, fmt.Errorf("extra key in websites[%d]: %s", i, key))
					}
				}

				// Validate URL
				if urlStr, ok := websiteMap["url"].(string); ok {
					if urlStr == "" {
						errors = append(errors, fmt.Errorf("websites[%d].url is required", i))
					}
					parsedUrl, err := url.Parse(urlStr)
					if err != nil || parsedUrl.Scheme != "https" {
						errors = append(errors, fmt.Errorf("websites[%d].url is not a valid URL", i))
					}
					if usedUrls[urlStr] {
						errors = append(errors, fmt.Errorf("websites[%d].url is a duplicate", i))
					}
					usedUrls[urlStr] = true
				} else {
					errors = append(errors, fmt.Errorf("websites[%d].url is required", i))
				}

				// Validate Description
				if desc, ok := websiteMap["description"].(string); !ok || desc == "" {
					errors = append(errors, fmt.Errorf("websites[%d].description is required", i))
				}
			}
		}
	}

	// Validate Contracts for extra keys and field correctness
	if contractsRaw, ok := rawData["contracts"].([]interface{}); ok {
		expectedContractFields := map[string]struct{}{
			"address":     {},
			"label":       {},
			"description": {},
		}

		usedContracts := make(map[string]bool)
		for i, contract := range contractsRaw {
			if contractMap, ok := contract.(map[string]interface{}); ok {
				for key := range contractMap {
					if _, ok := expectedContractFields[key]; !ok {
						errors = append(errors, fmt.Errorf("extra key in contracts[%d]: %s", i, key))
					}
				}

				// Validate Address
				if address, ok := contractMap["address"].(string); ok {
					if address == "" {
						errors = append(errors, fmt.Errorf("contracts[%d].address is required", i))
					}
					if !strings.HasPrefix(address, "0x") || !common.IsHexAddress(address) {
						errors = append(errors, fmt.Errorf("contracts[%d].address is not a valid Ronin address", i))
					}
					if usedContracts[address] {
						errors = append(errors, fmt.Errorf("contracts[%d].address is a duplicate", i))
					}
					usedContracts[address] = true
				} else {
					errors = append(errors, fmt.Errorf("contracts[%d].address is required", i))
				}

				// Validate Label
				if label, ok := contractMap["label"].(string); !ok || label == "" {
					errors = append(errors, fmt.Errorf("contracts[%d].label is required", i))
				}

				// Validate Description
				if desc, ok := contractMap["description"].(string); !ok || desc == "" {
					errors = append(errors, fmt.Errorf("contracts[%d].description is required", i))
				}
			}
		}
	}

	// Validate Categories
	for i, category := range data.Categories {
		if !isValidCategory(category) {
			errors = append(errors, fmt.Errorf("categories[%d] is not a valid category", i))
		}
	}

	// Validate Social
	if data.Social != nil {

		keys := make([]Platform, 0, len(*data.Social))

		for key := range *data.Social {
			keys = append(keys, key)
		}

		for i, platform := range keys {
			if !isValidPlatform(platform) {
				errors = append(errors, fmt.Errorf("invalid social platform: %s", platform))
			}
			socialUrl := (*data.Social)[platform]
			if socialUrl == "" {
				errors = append(errors, fmt.Errorf("social.%d can't be empty", i))
			}
			parsedUrl, err := url.Parse(socialUrl)
			if err != nil || parsedUrl == nil || parsedUrl.Scheme != "https" {
				errors = append(errors, fmt.Errorf("social.%d is not a valid URL", i))
			}
			if !isValidPlatformHost(platform, parsedUrl) {
				errors = append(errors, fmt.Errorf("social.%d is not a valid Host for platform %s", i, platform))
			}
		}
	}

	// Validate the structure with go-playground/validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(data)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Errorf("validation error: %s", err))
		}
	}

	return errors
}

func isValidPlatform(platform Platform) bool {
	switch platform {
	case PlatformFacebook, PlatformInstagram, PlatformTwitter, PlatformX, PlatformLinkedIn, PlatformThreads, PlatformMastodon, PlatformTelegram, PlatformDiscord:
		return true
	default:
		return false
	}
}

func isValidCategory(category Category) bool {
	switch category {
	case CategoryGame, CategoryNFT, CategoryFinance, CategoryDAO, CategoryTool, CategoryOther:
		return true
	default:
		return false
	}
}

func isValidPlatformHost(platform Platform, url *url.URL) bool {
	switch platform {
	case PlatformFacebook:
		if url.Host != "facebook.com" && url.Host != "www.facebook.com" && url.Host != "fb.me" {
			return false
		}
	case PlatformInstagram:
		if url.Host != "instagram.com" && url.Host != "www.instagram.com" && url.Host != "instagr.am" {
			return false
		}
	case PlatformTwitter, PlatformX:
		if url.Host != "twitter.com" && url.Host != "www.twitter.com" && url.Host != "x.com" && url.Host != "www.x.com" && url.Host != "t.co" {
			return false
		}
	case PlatformLinkedIn:
		if url.Host != "linkedin.com" && url.Host != "www.linkedin.com" && url.Host != "lnkd.in" {
			return false
		}
	case PlatformThreads:
		if url.Host != "threads.net" && url.Host != "www.threads.net" {
			return false
		}
	case PlatformMastodon:
		if url.Host != "mastodon.social" && url.Host != "www.mastodon.social" {
			return false
		}
	case PlatformTelegram:
		if url.Host != "t.me" && url.Host != "www.t.me" {
			return false
		}
	case PlatformDiscord:
		if url.Host != "discord.com" && url.Host != "www.discord.com" && url.Host != "discord.gg" && url.Host != "www.discord.gg" {
			return false
		}
	default:
		return false
	}

	return true
}
