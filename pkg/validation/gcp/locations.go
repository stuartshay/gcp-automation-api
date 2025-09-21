package gcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// LocationValidator provides dynamic validation against live GCP APIs
type LocationValidator struct {
	projectID   string
	regions     map[string]bool
	zones       map[string]bool
	lastUpdated time.Time
	cacheTTL    time.Duration
	mu          sync.RWMutex
	clientOpts  []option.ClientOption
}

// NewLocationValidator creates a new location validator that uses live GCP APIs
func NewLocationValidator(projectID string, opts ...option.ClientOption) *LocationValidator {
	return &LocationValidator{
		projectID:  projectID,
		regions:    make(map[string]bool),
		zones:      make(map[string]bool),
		cacheTTL:   time.Hour, // Cache for 1 hour
		clientOpts: opts,
	}
}

// ValidateLocationDynamic validates a location against live GCP APIs
func (lv *LocationValidator) ValidateLocationDynamic(ctx context.Context, location string) error {
	if location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	// Update cache if needed
	if err := lv.updateCacheIfNeeded(ctx); err != nil {
		return fmt.Errorf("failed to fetch location data: %w", err)
	}

	lv.mu.RLock()
	defer lv.mu.RUnlock()

	// Check if it's a valid region
	if lv.regions[location] {
		return nil
	}

	// Check if it's a valid zone
	if lv.zones[location] {
		return nil
	}

	return fmt.Errorf("invalid GCP location: %s. Location not found in available regions or zones", location)
}

// GetAvailableLocations returns all available regions and zones
func (lv *LocationValidator) GetAvailableLocations(ctx context.Context) (regions, zones []string, err error) {
	if err := lv.updateCacheIfNeeded(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch location data: %w", err)
	}

	lv.mu.RLock()
	defer lv.mu.RUnlock()

	regions = make([]string, 0, len(lv.regions))
	for region := range lv.regions {
		regions = append(regions, region)
	}

	zones = make([]string, 0, len(lv.zones))
	for zone := range lv.zones {
		zones = append(zones, zone)
	}

	return regions, zones, nil
}

// updateCacheIfNeeded updates the location cache if it's stale
func (lv *LocationValidator) updateCacheIfNeeded(ctx context.Context) error {
	lv.mu.RLock()
	needsUpdate := time.Since(lv.lastUpdated) > lv.cacheTTL
	lv.mu.RUnlock()

	if !needsUpdate {
		return nil
	}

	return lv.updateCache(ctx)
}

// updateCache fetches fresh location data from GCP APIs
func (lv *LocationValidator) updateCache(ctx context.Context) error {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	// Double-check pattern: another goroutine might have updated while we waited
	if time.Since(lv.lastUpdated) <= lv.cacheTTL {
		return nil
	}

	// Create context with timeout for API calls
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	newRegions := make(map[string]bool)
	newZones := make(map[string]bool)

	// Fetch regions
	if err := lv.fetchRegions(timeoutCtx, newRegions); err != nil {
		return fmt.Errorf("failed to fetch regions: %w", err)
	}

	// Fetch zones
	if err := lv.fetchZones(timeoutCtx, newZones); err != nil {
		return fmt.Errorf("failed to fetch zones: %w", err)
	}

	// Update cache atomically
	lv.regions = newRegions
	lv.zones = newZones
	lv.lastUpdated = time.Now()

	return nil
}

// fetchRegions retrieves all available regions from Compute Engine API
func (lv *LocationValidator) fetchRegions(ctx context.Context, regions map[string]bool) error {
	client, err := compute.NewRegionsRESTClient(ctx, lv.clientOpts...)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.ListRegionsRequest{
		Project: lv.projectID,
	}

	it := client.List(ctx, req)
	for {
		region, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating regions: %w", err)
		}

		if region.GetName() != "" {
			regions[region.GetName()] = true
		}
	}

	return nil
}

// fetchZones retrieves all available zones from Compute Engine API
func (lv *LocationValidator) fetchZones(ctx context.Context, zones map[string]bool) error {
	client, err := compute.NewZonesRESTClient(ctx, lv.clientOpts...)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.ListZonesRequest{
		Project: lv.projectID,
	}

	it := client.List(ctx, req)
	for {
		zone, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error iterating zones: %w", err)
		}

		if zone.GetName() != "" {
			zones[zone.GetName()] = true
		}
	}

	return nil
}

// ValidateLocationWithFallback uses dynamic validation with static fallback
func ValidateLocationWithFallback(ctx context.Context, projectID, location string) error {
	// Try static validation first (fast)
	if err := ValidateLocation(location); err == nil {
		return nil
	}

	// If static validation fails, try dynamic validation
	validator := NewLocationValidator(projectID)
	return validator.ValidateLocationDynamic(ctx, location)
}
