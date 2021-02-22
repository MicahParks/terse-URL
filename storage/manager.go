package storage

import (
	"context"
	"fmt"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terseurl/models"
)

// TODO Remove VisitsStore and SummaryStore from TerseStore implementations.

type StoreManager struct { // TODO Rename.
	group        ctxerrgroup.Group
	summaryStore SummaryStore
	terseStore   TerseStore
	visitsStore  VisitsStore
}

// TODO Turn shortenedURLs argument into a map, then back into a slice before giving to stores. This will ensure it is a
// set and helps prevent DOS.

// Close TODO
func (s StoreManager) Close(ctx context.Context) (err error) {

	// Kill the worker pool.
	s.group.Kill()

	// TODO See if you can stick more than one error in with %w somehow...

	// Close the SummaryStore.
	var closeErr error
	s.SummaryStore(func(store SummaryStore) {
		closeErr = store.Close(ctx)
	})
	if closeErr != nil {
		err = fmt.Errorf("%v SummaryStore: %v", err, closeErr)
	}

	// Close the TerseStore.
	if closeErr = s.terseStore.Close(ctx); closeErr != nil {
		err = fmt.Errorf("%v TerseStore: %v", err, closeErr)
	}

	// Close the VisitsStore.
	s.VisitsStore(func(store VisitsStore) {
		closeErr = store.Close(ctx)
	})
	if closeErr != nil {
		err = fmt.Errorf("%v VisitsStore: %v", err, closeErr)
	}

	return err
}

// Export exports the Terse data and Visits data for the given shortened URLs. If shortenedURLs is nil, then all
// shortened URLs are exported.
func (s StoreManager) Export(ctx context.Context, shortenedURLs []string) (export map[string]*models.Export, err error) {

	// Create the return map.
	export = make(map[string]*models.Export, len(shortenedURLs))

	// Get the Terse data for the export.
	var terse map[string]*models.Terse
	if terse, err = s.terseStore.Read(ctx, shortenedURLs); err != nil {
		return nil, err
	}

	// Get the Visits data for the export.
	visits := make(map[string][]*models.Visit)
	s.VisitsStore(func(store VisitsStore) {
		visits, err = store.Read(ctx, shortenedURLs)
	})
	if err != nil {
		return nil, err
	}

	// Combine the Terse data and Visits data for the export.
	for shortened, t := range terse {
		export[shortened] = &models.Export{
			Terse:  t,
			Visits: visits[shortened],
		}
	}

	return export, nil
}

// Import imports the given Terse data and Visits data to the TerseStore and VisitsStore respectively. Terse data will
// be overwritten, Visits data will be appended.
func (s StoreManager) Import(ctx context.Context, data map[string]models.Export) (err error) {

	// Iterate through the given import data and put it in the proper format.
	terse := make(map[string]*models.Terse)
	visits := make(map[string][]*models.Visit)
	for shortened, export := range data {
		terse[shortened] = export.Terse
		visits[shortened] = export.Visits
	}

	// Write the Terse data to the TerseStore.
	if err = s.terseStore.Write(ctx, terse, Upsert); err != nil {
		return err
	}

	// Write the Visits data to the VisitsStore.
	s.VisitsStore(func(store VisitsStore) {
		err = store.Insert(ctx, visits)
	})
	if err != nil {
		return err
	}

	return nil
}

// InitializeSummaryStore initializes the SummaryStore with SummaryData gathered from the TerseStore and VisitsStore.
func (s StoreManager) InitializeSummaryStore(ctx context.Context) (err error) {

	// Get the Visits Summary data.
	var visitsSummary map[string]*models.VisitsSummary
	s.VisitsStore(func(store VisitsStore) {
		visitsSummary, err = store.Summary(ctx, nil)
	})
	if err != nil {
		return err
	}
	haveVisits := visitsSummary != nil

	// Get the Summary data.
	var terseSummary map[string]*models.TerseSummary
	if terseSummary, err = s.terseStore.Summary(ctx, nil); err != nil {
		return err
	}

	// Populate the SummaryStore with the Summary data.
	summaryData := make(map[string]*models.Summary)
	for shortened, terse := range terseSummary {

		// Check if there is Visits data.
		var visits *models.VisitsSummary
		if haveVisits {
			visits = visitsSummary[shortened]
		}

		// Assign the shortened URL's summary data to the return map.
		summaryData[shortened] = &models.Summary{
			Terse:  terse,
			Visits: visits,
		}
	}

	// If Visits data is allowed to be present when Terse data is not present for a shortened URL, then it would need to
	// be looped through.

	// Delete all existing Summary data and import the most recent summary data.
	s.SummaryStore(func(store SummaryStore) {
		if err = store.Delete(ctx, nil); err != nil { // Not necessary if only used on startup.
			return
		}
		err = store.Upsert(ctx, summaryData)
	})

	return err
}

// SummaryStore accepts a function to do if the SummaryStore is not nil
func (s StoreManager) SummaryStore(doThis func(store SummaryStore)) {
	if s.summaryStore != nil {
		doThis(s.summaryStore)
	}
}

// VisitsStore accepts a function to do if the VisitsStore is not nil.
func (s StoreManager) VisitsStore(doThis func(store VisitsStore)) {
	if s.visitsStore != nil {
		doThis(s.visitsStore)
	}
}
