package cachethq

import (
	"github.com/andygrunwald/cachet"
	"github.com/tixu/alch/errors"
)

type Instance interface {
	ChangeComponentStatus(name string, groupName string, stringStatus string) error
}

// instance CachetHQ instance
type instance struct {
	client *cachet.Client
}

// ChangeComponentStatus Change component status
func (ctx *instance) ChangeComponentStatus(name string, groupName string, stringStatus string) error {
	// Find component
	compo, err := ctx.findComponent(name, groupName)
	if err != nil {
		return err
	}
	status, err := getCachetHQComponentStatus(stringStatus)
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	// Store component ID
	id := compo.ID
	// Change component status
	compo.Status = status
	// Run update request
	_, _, err = ctx.client.Components.Update(id, compo)
	if err != nil {
		return errors.NewInternalServerError(err)
	}
	return nil
}

func (ctx *instance) findComponent(name string, groupName string) (*cachet.Component, error) {
	// Creating default group id
	grpID := 0
	// Find group if possible
	if groupName != "" {
		grp, err := ctx.findComponentGroup(groupName)
		// Check error
		if err != nil {
			return nil, err
		}

		grpID = grp.ID
	}

	// Create query params for name filter
	queryParams := &cachet.ComponentsQueryParams{
		Name:         name,
		GroupID:      grpID,
		QueryOptions: cachet.QueryOptions{PerPage: 10000},
	}
	c, _, err := ctx.client.Components.GetAll(queryParams)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}

	// Check if length is 0
	if len(c.Components) == 0 {
		return nil, errors.NewNotFoundError(ErrComponentNotFound)
	}

	// Filter components by groups
	// Client doesn't manage group id equal to 0 for no groups...
	for _, comp := range c.Components {
		if comp.GroupID == grpID {
			return &comp, nil
		}
	}

	// Default case
	return nil, errors.NewNotFoundError(ErrComponentNotFound)
}
func getCachetHQComponentStatus(statusString string) (int, error) {
	switch statusString {
	case config.ComponentMajorOutageStatus:
		return cachet.ComponentStatusMajorOutage, nil
	case config.ComponentOperationalStatus:
		return cachet.ComponentStatusOperational, nil
	case config.ComponentPartialOutageStatus:
		return cachet.ComponentStatusPartialOutage, nil
	case config.ComponentPerformanceIssuesStatus:
		return cachet.ComponentStatusPerformanceIssues, nil
	default:
		return 0, ErrStatusNotFound
	}
}
