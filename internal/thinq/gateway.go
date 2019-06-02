package thinq

import (
	"context"
)

type GatewayService service

func (s *GatewayService) Discover(ctx context.Context) error {
	return nil
}