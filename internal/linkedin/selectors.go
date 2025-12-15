package linkedin

// LinkedIn UI selectors
// All selectors are centralized here for easy maintenance

// Login selectors
const (
	// Login page
	LoginEmailInput    = "input[name='session_key'], input#username"
	LoginPasswordInput = "input[name='session_password'], input#password"
	LoginSubmitButton  = "button[type='submit'], button[data-litms-control-urn*='login']"

	// Login state indicators
	GlobalNav          = "[data-testid='global-nav'], .global-nav, #global-nav"
	FeedIdentityModule = ".feed-identity-module"
)

// Search selectors
const (
	// Search input
	SearchBox    = "input[aria-label='Search']"
	SearchButton = "button[data-test-search-submit]"

	// Search results
	SearchResultsList = ".search-results-container"
	SearchResultItem  = ".reusable-search__result-container"
	SearchResultLink  = "a.app-aware-link"

	// Pagination
	NextPageButton   = "button[aria-label='Next']"
	PaginationButton = ".artdeco-pagination__indicator button"

	// Filters
	FilterButton = "button[aria-label*='filter']"
	FilterPanel  = ".search-reusables__filter-panel"
)

// Profile selectors
const (
	// Profile elements
	ProfileName       = ".text-heading-xlarge"
	ProfileHeadline   = ".text-body-medium"
	ProfileLocation   = ".text-body-small"
	ProfilePictureImg = "img[class*='pv-top-card-profile-picture']"

	// Connection button
	ConnectButton    = "button[aria-label*='Invite'][aria-label*='connect'], button.pvs-profile-actions__action:has-text('Connect')"
	ConnectButtonAlt = "button:has-text('Connect')"

	// More actions
	MoreButton = "button[aria-label='More actions']"

	// Profile sections
	ExperienceSection = "#experience"
	EducationSection  = "#education"
	SkillsSection     = "#skills"

	// Connection status
	PendingButton = "button[aria-label*='Pending']"
	MessageButton = "button[aria-label*='Message'], a[href*='/messaging/']"
)

// Connection request selectors
const (
	// Connection modal
	ConnectModal     = "div[role='dialog'][aria-labelledby*='send-invite']"
	AddNoteButton    = "button[aria-label='Add a note']"
	NoteTextarea     = "textarea[name='message'], textarea[id='custom-message']"
	SendInviteButton = "button[aria-label='Send'], button[aria-label='Send invitation']"
	SendButton       = "button[aria-label*='Send']"

	// Without note
	SendWithoutNoteButton = "button[aria-label='Send without a note']"

	// Close/Cancel
	CloseModalButton = "button[aria-label='Dismiss']"
	CancelButton     = "button[aria-label='Cancel']"
)

// Messaging selectors
const (
	// Message list
	MessagesList     = ".msg-conversations-container__conversations-list"
	ConversationItem = ".msg-conversation-listitem"

	// Message compose
	MessageComposer   = "div[role='textbox'][contenteditable='true']"
	MessageInput      = ".msg-form__contenteditable"
	SendMessageButton = "button[type='submit'].msg-form__send-button"

	// Conversation
	ConversationCard = ".msg-overlay-bubble-header"
	MessageThread    = ".msg-s-message-list__event"
)

// Navigation selectors
const (
	// Main navigation
	HomeNavButton      = "a[href*='/feed/']"
	NetworkNavButton   = "a[href*='/mynetwork/']"
	JobsNavButton      = "a[href*='/jobs/']"
	MessagingNavButton = "a[href*='/messaging/']"

	// Profile menu
	ProfileMenuButton = "button[id='global-nav-icon']"
	SignOutButton     = "a[href*='/logout']"
)

// Common UI elements
const (
	// Buttons
	PrimaryButton   = "button.artdeco-button--primary"
	SecondaryButton = "button.artdeco-button--secondary"

	// Loaders
	LoadingSpinner = ".artdeco-loader"
	PageLoader     = "[data-test-loading-indicator]"

	// Alerts/Toasts
	Toast        = ".artdeco-toast"
	ToastMessage = ".artdeco-toast-item__message"

	// Modals
	Modal            = "div[role='dialog']"
	ModalCloseButton = "button[data-test-modal-close-btn]"
)
