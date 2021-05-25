package cachethq

import (
	"github.com/andygrunwald/cachet"
	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/config"
	"github.com/tixu/alch/errors"
)

type Instance interface {
	ChangeComponentStatus(name string, groupName string, status int) error
}

// instance CachetHQ instance
type instance struct {
	client *cachet.Client
	logger *logrus.Logger
}

// ChangeComponentStatus Change component status
func (ctx *instance) ChangeComponentStatus(name string, groupName string, status int) error {
	ctx.logger.Infof("looking for component %s in group %s", name, groupName)
	// Find component
	compo, err := ctx.findComponent(name, groupName)

	if err != nil {
		return err
	}
	//	status, err := getCachetHQComponentStatus(stringStatus)
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
	ctx.logger.Infof("Finding component %s in %s", name, groupName)
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
	ctx.logger.Infof("Finding component %s in %s found", name, groupName)
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

func (ctx *instance) findComponentGroup(name string) (*cachet.ComponentGroup, error) {
	// Create query params for name filter
	ctx.logger.Infof("looking for group %s", name)
	queryParams := &cachet.ComponentGroupsQueryParams{
		Name: name,
	}
	c, _, err := ctx.client.ComponentGroups.GetAll(queryParams)
	if err != nil {
		ctx.logger.Infof("looking for group %s error while lookup", name)
		return nil, err
	}

	// Check if length is 0
	if len(c.ComponentGroups) == 0 {
		ctx.logger.Infof("group %s not found", name)
		return nil, errors.NewNotFoundError(ErrComponentGroupNotFound)
	}

	// Get component group
	grp := c.ComponentGroups[0]
	ctx.logger.Infof("group found %v\n")
	return &grp, nil
}
func NewInstance(url string, key string, logger *logrus.Logger) (Instance, error) {
	client, err := cachet.NewClient(url, nil)
	if err != nil {
		return nil, err
	}
	client.Authentication.SetTokenAuth(key)
	return &instance{client: client, logger: logger}, nil
}

func NewInstanceConf(conf *config.CachetConfig, logger *logrus.Logger) (Instance, error) {
	return NewInstance(conf.URL, conf.APIKey, logger)

}
