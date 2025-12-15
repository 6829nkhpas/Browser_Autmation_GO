package linkedin

// LinkedIn URL constants
const (
	// Base URLs
	BaseURL  = "https://www.linkedin.com"
	LoginURL = "https://www.linkedin.com/login"
	FeedURL  = "https://www.linkedin.com/feed/"

	// Search URLs
	SearchURL = "https://www.linkedin.com/search/results/people/"

	// Messaging
	MessagingURL = "https://www.linkedin.com/messaging/"

	// Network
	MyNetworkURL   = "https://www.linkedin.com/mynetwork/"
	ConnectionsURL = "https://www.linkedin.com/mynetwork/invite-connect/connections/"
)

// BuildSearchURL builds a LinkedIn people search URL with filters
func BuildSearchURL(keywords string, filters map[string]string) string {
	url := SearchURL + "?keywords=" + keywords

	// Add additional filters
	for key, value := range filters {
		url += "&" + key + "=" + value
	}

	return url
}

// BuildProfileURL builds a LinkedIn profile URL
func BuildProfileURL(profileID string) string {
	return BaseURL + "/in/" + profileID + "/"
}
