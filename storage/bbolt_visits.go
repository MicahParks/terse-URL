package storage

import (
	"context"

	"github.com/MicahParks/terse-URL/models"
)

type BboltVisits struct {
}

func (b *BboltVisits) AddVisit(ctx context.Context, shortened string, visit *models.Visit) (err error) {
	panic("implement me")
}

func (b *BboltVisits) Close(ctx context.Context) (err error) {
	panic("implement me")
}

func (b *BboltVisits) DeleteVisits(ctx context.Context, shortened string) (err error) {
	panic("implement me")
}

func (b *BboltVisits) ReadVisits(ctx context.Context, shortened string) (visits []*models.Visit, err error) {
	panic("implement me")
}
