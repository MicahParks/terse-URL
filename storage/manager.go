package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/MicahParks/ctxerrgroup"

	"github.com/MicahParks/terseurl/models"
)

// StoreManager holds all data stores and coordinates operations on them.
type StoreManager struct {
	authStore    AuthorizationStore
	createCtx    CtxCreator
	group        ctxerrgroup.Group
	summaryStore SummaryStore
	terseStore   TerseStore
	visitsStore  VisitsStore
}

// NewStoreManager creates a new manager for the data stores.
func NewStoreManager(createCtx CtxCreator, group ctxerrgroup.Group, summaryStore SummaryStore, terseStore TerseStore, visitsStore VisitsStore) (manager StoreManager) {
	return StoreManager{
		createCtx:    createCtx,
		group:        group,
		summaryStore: summaryStore,
		terseStore:   terseStore,
		visitsStore:  visitsStore,
	}
}

// AuthStore accepts a function to do if the AuthorizationStore is not nil.
func (s StoreManager) AuthStore(doThis func(store AuthorizationStore)) {
	if s.authStore != nil {
		doThis(s.authStore)
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
func (s StoreManager) DeleteShortened(ctx context.Context, principal *models.Principal, shortenedURLs []string) (err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		Owner: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return err
	}

	// Check if the request is authorized.
	if !authorized {
		return ErrUnauthorized
	}

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

	// TODO Delete from auth store in both places.

	return nil
}

// DeleteVisits deletes Visits data for the given shortened URLs. If shortenedURLs is nil, then all Visits data are
// deleted. No error should be given if a shortened URL is not found.
func (s StoreManager) DeleteVisits(ctx context.Context, principal *models.Principal, shortenedURLs []string) (err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		Owner: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return err
	}

	// Check if the request is authorized.
	if !authorized {
		return ErrUnauthorized
	}

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
func (s StoreManager) Export(ctx context.Context, principal *models.Principal, shortenedURLs []string) (export map[string]*models.Export, err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		ReadVisits: true,
		ReadTerse:  true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return nil, err
	}

	// Check if the request is authorized.
	if !authorized {
		return nil, ErrUnauthorized
	}

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
		v := visits[shortened]
		if v == nil {
			v = make([]models.Visit, 0)
		}
		export[shortened] = &models.Export{
			Terse:  t,
			Visits: v,
		}
	}

	return export, nil
}

// Import imports the given Terse data and Visits data to the TerseStore and VisitsStore respectively. Terse data will
// be overwritten, Visits data will be appended.
func (s StoreManager) Import(ctx context.Context, data map[string]*models.Export, principal *models.Principal) (err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		Owner: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, nil, neededAuth); err != nil {
		return err
	}

	// Check if the request is authorized.
	if !authorized {
		return ErrUnauthorized
	}

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
		visits := &models.VisitsSummary{}
		if haveVisits {
			visits = visitsSummary[shortened]
		}

		// Assign the shortened URL's summary data to the return map.
		summaryData[shortened] = &models.Summary{
			Terse:  terse,
			Visits: visits,
		}
	}

	// If Visits data are allowed to be present when Terse data are not present for a shortened URL, then it would need to
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

	// Get the Terse data from the TerseStore.
	var terseData map[string]*models.Terse
	if terseData, err = s.Terse(ctx, nil, []string{shortened}); err != nil {
		return nil, err
	}

	// Handle the visit in another goroutine for a faster response.
	go s.handleVisit(shortened, visit)

	return terseData[shortened], nil
}

// Summary retrieves the Summary data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL
// summary data will be returned.
func (s StoreManager) Summary(ctx context.Context, principal *models.Principal, shortenedURLs []string) (summaries map[string]*models.Summary, err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		ReadSummary: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return nil, err
	}

	// Check if the request is authorized.
	if !authorized {
		return nil, ErrUnauthorized
	}

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
func (s StoreManager) Terse(ctx context.Context, principal *models.Principal, shortenedURLs []string) (terse map[string]*models.Terse, err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		ReadTerse: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return nil, err
	}

	// Check if the request is authorized.
	if !authorized {
		return nil, ErrUnauthorized
	}

	// Turn the input slice into a set.
	shortenedURLs = makeStringSliceSet(shortenedURLs)

	return s.terseStore.Read(ctx, shortenedURLs)
}

// Visits exports the Visits data for the given shortened URLs. If shortenedURLs is nil, then all shortened URL Visits
// data are expected. The error must be storage.ErrShortenedNotFound if a shortened URL is not found.
func (s StoreManager) Visits(ctx context.Context, principal *models.Principal, shortenedURLs []string) (visits map[string][]models.Visit, err error) {

	// Create the needed authorization.
	neededAuth := Authorization{
		ReadVisits: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return nil, err
	}

	// Check if the request is authorized.
	if !authorized {
		return nil, ErrUnauthorized
	}

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
func (s StoreManager) WriteTerse(ctx context.Context, principal *models.Principal, terse map[string]*models.Terse, operation WriteOperation) (err error) {

	// Create a slice of affected shortened URLs.
	shortenedURLs := make([]string, len(terse))
	index := 0
	for shortened := range terse {
		shortenedURLs[index] = shortened
		index++
	}

	// Create the needed authorization.
	neededAuth := Authorization{
		WriteTerse: true,
	}

	// Make sure the principal is authorized.
	var authorized bool
	if authorized, err = s.authorize(ctx, principal, shortenedURLs, neededAuth); err != nil {
		return err
	}

	// Check if the request is authorized.
	if !authorized {
		return ErrUnauthorized
	}

	// Write the Terse data.
	if err = s.terseStore.Write(ctx, terse, operation); err != nil {
		return err
	}

	// Add the shortened URLs to the SummaryStore, if required.
	s.SummaryStore(func(store SummaryStore) {
		summaries := make(map[string]*models.Summary)

		// Iterate through the given shortened URLs.
		for shortened, terseData := range terse {

			// Create entries in the SummaryStore for any new shortened URLs.
			var summary map[string]*models.Summary
			var visitsData *models.VisitsSummary
			summary, err = store.Read(ctx, []string{shortened}) // TODO Loop efficiency could be improved if used zero values instead of ErrShortenedNotFound.
			if err != nil {
				if errors.Is(err, ErrShortenedNotFound) {

					// Use the empty visitsData.
					visitsData = &models.VisitsSummary{VisitCount: 0}
				} else {
					return
				}
			} else {

				// Use the existing visitsData.
				visitsData = summary[shortened].Visits
			}

			// Assign the Summary data for upsertion later.
			summaries[shortened] = &models.Summary{
				Terse:  summarizeTerse(*terseData),
				Visits: visitsData,
			}
		}

		// Upsert the new Summary data into the SummaryStore.
		if err = store.Upsert(ctx, summaries); err != nil {
			return
		}
	})
	if err != nil {
		return err
	}

	// Add the shortened URLs to the VisitsStore, if required.
	s.VisitsStore(func(store VisitsStore) {
		visitsData := make(map[string][]models.Visit)
		switch operation {
		case Insert:

			// Iterate through the given shortened URLs, which are new. Create Visits data for them.
			for shortened := range terse {
				visitsData[shortened] = make([]models.Visit, 0)
			}
		case Upsert:

			// Iterate through the given shortened URLs.
			for shortened := range terse {

				// Create entries in the VisitsStore for any new shortened URLs.
				_, err = store.Read(ctx, []string{shortened})
				if err != nil {
					if errors.Is(err, ErrShortenedNotFound) {

						// The shortened URL was not already in the VisitsStore. Add it to the map of new visits.
						visitsData[shortened] = make([]models.Visit, 0)
					} else {
						return
					}
				}
			}
		}

		// Upsert the new Visits data into the VisitsStore.
		if err = store.Insert(ctx, visitsData); err != nil {
			return
		}
	})
	if err != nil {
		return err
	}

	// TODO Add to the AuthStore.

	return nil
}

// authorize authorizes actions on the given shortened URLs based on the needed authorizations and principal. If the
// principal is nil, the request is authorized.
func (s StoreManager) authorize(ctx context.Context, principal *models.Principal, shortenedURLs []string, neededAuth Authorization) (authorized bool, err error) {
	if principal != nil {
		s.AuthStore(func(store AuthorizationStore) {

			// The user is the subject of the JWT.
			user := principal.Sub

			// Get the mapping of users to Authorization data.
			var usersShortened map[string]UserAuth
			if usersShortened, err = store.ReadUsers(ctx, []string{user}); err != nil {
				return
			}

			// If all shortened URLs are supposed to be used, find out which ones.
			if len(shortenedURLs) == 0 {

				// Turn the map of shortened URLs to Authorization data into a slice of shortened URLs.
				for shortened := range usersShortened[user] {
					shortenedURLs = append(shortenedURLs, shortened)
				}
			}

			// Confirm the user has the appropriate permissions.
			for _, shortened := range shortenedURLs {
				authData, ok := usersShortened[user][shortened]
				if !ok {
					// TODO Log.
				}

				// If the user is the owner of the shortenedURL, they are allowed to perform any action.
				if authData.Owner {
					continue
				}

				// Confirm the user has the required permissions.
				//
				// TODO Might be a better way to check this.
				if neededAuth.ReadSummary {
					if !authData.ReadSummary {
						authorized = false
						return
					}
				}
				if neededAuth.ReadTerse {
					if !authData.ReadTerse {
						authorized = false
						return
					}
				}
				if neededAuth.ReadVisits {
					if !authData.ReadVisits {
						authorized = false
						return
					}
				}
				if neededAuth.WriteTerse {
					if !authData.WriteTerse {
						authorized = false
						return
					}
				}
			}
		})
		if err != nil {
			return false, err
		}
	} else {
		authorized = true
	}

	return authorized, nil
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
	index := 0
	for str := range m {
		set[index] = str
		index++
	}

	return set
}
