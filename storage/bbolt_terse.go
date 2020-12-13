package storage

import (
	"context"

	"github.com/MicahParks/terse-URL/models"
)

type BboltTerse struct {
}

func (b *BboltTerse) InsertTerse(ctx context.Context, terse *models.Terse) (err error) {
	panic("implement me")
}

func (b *BboltTerse) Close(ctx context.Context) (err error) {
	panic("implement me")
}

func (b *BboltTerse) DeleteTerse(ctx context.Context, shortened string) (err error) {
	panic("implement me")
}

func (b *BboltTerse) Export(ctx context.Context, shortened string) (export models.Export, err error) {
	panic("implement me")
}

func (b *BboltTerse) ExportAll(ctx context.Context) (export map[string]models.Export, err error) {
	panic("implement me")
}

func (b *BboltTerse) GetTerse(ctx context.Context, shortened string, visit *models.Visit) (terse *models.Terse, err error) {
	panic("implement me")
}

func (b *BboltTerse) UpdateTerse(ctx context.Context, terse *models.Terse) (err error) {
	panic("implement me")
}

func (b *BboltTerse) UpsertTerse(ctx context.Context, terse *models.Terse) (err error) {
	panic("implement me")
}

func (b *BboltTerse) VisitsStore() VisitsStore {
	panic("implement me")
}
