package request

type ItemGetAllRequest struct {
	ExclusionZeroQuantity bool `json:"exclusion_zero_quantity" default:"false"`
}
