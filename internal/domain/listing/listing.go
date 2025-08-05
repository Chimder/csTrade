package listing

type ListingStatus string

const (
	ListingStatusPending    ListingStatus = "pending"
	ListingStatusInProgress ListingStatus = "in_progress"
	ListingStatusCompleted  ListingStatus = "completed"
	ListingStatusCanceled   ListingStatus = "canceled"
)
type Listing struct {
    ID             string
    SellerID       string
    InventorySkinID string
    Price          float64
    Status         ListingStatus
}
