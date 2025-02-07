package attribute

const (
	metricsError = "error"
)

//nolint:revive
type map_ map[string]any

type termsAggsBucket struct {
	Key string `json:"key"`
}

type termsAggs struct {
	Buckets []termsAggsBucket `json:"buckets"`
}

type attributeAggs struct {
	Attr *termsAggs `json:"attr"`
}
