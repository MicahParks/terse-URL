package storage

import (
	"context"
	"fmt"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terseurl/models"
)

type StoreManager struct {
	createCtx    ctxCreator
	group        ctxerrgroup.Group
	summaryStore SummaryStore
	terseStore   TerseStore
	visitsStore  VisitsStore
}

// NewStoreManager creates a new manager for the data stores.
func NewStoreManager(createCtx ctxCreator, group ctxerrgroup.Group, summaryStore SummaryStore, terseStore TerseStore, visitsStore VisitsStore) (manager StoreManager) {
	return StoreManager{
		createCtx:    createCtx,
		group:        group,
		summaryStore: summaryStore,
		terseStore:   terseStore,
		visitsStore:  visitsStore,
	}
}

// Close closes the ctxerrgroup and all the underlying data stores.
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

// DeleteShortened deletes the all data for the given shortened URLs. If shortenedURLs is nil, all shortened URL data
// are deleted. There should be no error if a shortened URL is not found.
func (s StoreManager) DeleteShortened(ctx context.Context, shortenedURLs []string) (err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	// Delete the Terse data for the shortened URL.
	if err = s.terseStore.Delete(ctx, shortenedURLs); err != nil {
		return err
	}

	// Delete the Visits data for the shortened URL.
	s.VisitsStore(func(store VisitsStore) {
		err = store.Delete(ctx, shortenedURLs)
	})
	if err != nil {
		return err
	}

	// Delete the Summary data for the shortened URL.
	s.SummaryStore(func(store SummaryStore) {
		err = store.Delete(ctx, shortenedURLs)
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteVisits deletes Visits data for the given shortened URLs. If shortenedURLs is nil, then all Visits data are
// deleted. No error should be given if a shortened URL is not found.
func (s StoreManager) DeleteVisits(ctx context.Context, shortenedURLs []string) (err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	// Delete the Visits data from the VisitsStore.
	s.VisitsStore(func(store VisitsStore) {
		err = store.Delete(ctx, shortenedURLs)
	})
	if err != nil {
		return err
	}

	return nil
}

// Export exports the Terse data and Visits data for the given shortened URLs. If shortenedURLs is nil, then all
// shortened URLs are exported.
func (s StoreManager) Export(ctx context.Context, shortenedURLs []string) (export map[string]*models.Export, err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	// Create the return map.
	export = make(map[string]*models.Export, len(shortenedURLs))

	// Get the Terse data for the export.
	var terse map[string]*models.Terse
	if terse, err = s.terseStore.Read(ctx, shortenedURLs); err != nil {
		return nil, err
	}

	// Get the Visits data for the export.
	visits := make(map[string][]models.Visit)
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
func (s StoreManager) Import(ctx context.Context, data map[string]*models.Export) (err error) {

	// Iterate through the given import data and put it in the proper format.
	terse := make(map[string]*models.Terse)
	visits := make(map[string][]models.Visit)
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

// Redirect is called when a visit to a shortened URL has occurred. It will keep track of the visit and return the
// required information for a redirect.
func (s StoreManager) Redirect(ctx context.Context, shortened string, visit models.Visit) (terse *models.Terse, err error) {

	// Handle the visit in another goroutine for a faster response.
	go s.handleVisit(shortened, visit)

	// Get the Terse data from the TerseStore.
	var terseData map[string]*models.Terse
	if terseData, err = s.Terse(ctx, []string{shortened}); err != nil {
		return nil, err
	}

	return terseData[shortened], nil
}

// Summary retrieves the Summary data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL
// summary data will be returned.
func (s StoreManager) Summary(ctx context.Context, shortenedURLs []string) (summaries map[string]*models.Summary, err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	// Create the return map.
	summaries = make(map[string]*models.Summary, len(shortenedURLs))

	// Retrieve the Summary data from the SummaryStore.
	s.SummaryStore(func(store SummaryStore) {
		summaries, err = store.Read(ctx, shortenedURLs)
	})
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

// SummaryStore accepts a function to do if the SummaryStore is not nil
func (s StoreManager) SummaryStore(doThis func(store SummaryStore)) {
	if s.summaryStore != nil {
		doThis(s.summaryStore)
	}
}

// Terse returns a map of shortened URLs to Terse data. If shortenedURLs is nil, all shortened URL Terse data are
// expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (s StoreManager) Terse(ctx context.Context, shortenedURLs []string) (terse map[string]*models.Terse, err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	return s.terseStore.Read(ctx, shortenedURLs)
}

// Visits exports the Visits data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL Visits
// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (s StoreManager) Visits(ctx context.Context, shortenedURLs []string) (visits map[string][]models.Visit, err error) {

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	// Create the return map.
	visits = make(map[string][]models.Visit, len(shortenedURLs))

	// Get the Visits data from the VisitsStore.
	s.VisitsStore(func(store VisitsStore) {
		visits, err = store.Read(ctx, shortenedURLs)
	})
	if err != nil {
		return nil, err
	}

	return visits, nil
}

// VisitsStore accepts a function to do if the VisitsStore is not nil.
func (s StoreManager) VisitsStore(doThis func(store VisitsStore)) {
	if s.visitsStore != nil {
		doThis(s.visitsStore)
	}
}

// WriteTerse Write writes the given Terse data according to the given operation. The error must be
// storage.ErrShortenedExists if an Insert operation cannot be performed due to the Terse data already existing. The
// error must be storage.ErrShortenedNotFound if an Update operation cannot be performed due to the Terse data not
// existing.
func (s StoreManager) WriteTerse(ctx context.Context, terse map[string]*models.Terse, operation WriteOperation) (err error) {
	return s.terseStore.Write(ctx, terse, operation)
}

// handleVisit happens asynchronously when a redirect occurs. It updates the appropriate data stores with the required
// information.
func (s StoreManager) handleVisit(shortened string, visit models.Visit) {

	// Add the Visits data to the VisitsStore.
	{
		ctx, cancel := s.createCtx()
		s.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
			visits := map[string][]models.Visit{shortened: {visit}}
			s.VisitsStore(func(store VisitsStore) {
				err = store.Insert(workCtx, visits)
			})

			return err
		})
	}

	// Update the count in the SummaryStore.
	{
		ctx, cancel := s.createCtx()
		s.group.AddWorkItem(ctx, cancel, func(workCtx context.Context) (err error) {
			s.SummaryStore(func(store SummaryStore) {
				err = store.IncrementVisitCount(workCtx, shortened)
			})

			return err
		})
	}
}

// makeStringSliceSet makes a slice of strings a set by removing duplicate elements.
func makeStringSliceSet(slice []string) (set []string) {

	// Create a map to serve as a set.
	m := make(map[string]struct{})

	// Turn the slice into a map.
	for _, str := range slice {
		m[str] = struct{}{}
	}

	// Preallocate the return slice memory for faster insertion.
	set = make([]string, len(m))

	// Turn the map back into a slice.
	i := 0
	for str := range m {
		set[i] = str
		i++
	}

	return set
}
