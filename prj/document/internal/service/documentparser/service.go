package documentparser

import (
	"github.com/zuzuka28/simreport/lib/tikaclient"
)

type Service struct {
	tika *tikaclient.Client
}

func NewService(tika *tikaclient.Client) *Service {
	return &Service{
		tika: tika,
	}
}
