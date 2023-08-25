package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AbstractGroupService handles Groups in Kong.
type AbstractGroupService interface {
	// Create creates a Group in Kong.
	Create(ctx context.Context, group *Group) (*Group, error)
	// Get fetches a Group in Kong.
	Get(ctx context.Context, emailOrID *string) (*Group, error)
	// GetByCustomID fetches a Group in Kong.
	GetByCustomID(ctx context.Context, customID *string) (*Group, error)
	// Update updates a Group in Kong
	Update(ctx context.Context, Group *Group) (*Group, error)
	// Delete deletes a Group in Kong
	Delete(ctx context.Context, emailOrID *string) error
	// List fetches a list of Groups in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Group, *ListOpt, error)
	// ListAll fetches all Groups in Kong.
	ListAll(ctx context.Context) ([]*Group, error)
}

// GroupService handles Groups in Kong.
type GroupService service

// Create creates a Group in Kong.
// If an ID is specified, it will be used to
// create a Group in Kong, otherwise an ID
// is auto-generated.
// This call does _not_ use a PUT when provided an ID.
// Although /Groups accepts PUTs, PUTs do not accept passwords and do not create
// the hidden consumer that backs the Group. Subsequent attempts to use such Groups
// result in fatal errors.
func (s *GroupService) Create(ctx context.Context,
	group *Group,
) (*Group, error) {
	queryPath := "/groups"
	method := "POST"
	req, err := s.client.NewRequest(method, queryPath, nil, group)
	if err != nil {
		return nil, err
	}

	createdGroup := Group{}
	_, err = s.client.Do(ctx, req, &createdGroup)
	if err != nil {
		return nil, err
	}
	return &createdGroup, nil
}

// Get fetches a Group in Kong.
func (s *GroupService) Get(ctx context.Context,
	emailOrID *string,
) (*Group, error) {
	if isEmptyString(emailOrID) {
		return nil, fmt.Errorf("emailOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/Groups/%v", *emailOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var group Group
	_, err = s.client.Do(ctx, req, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetByCustomID fetches a Group in Kong.
func (s *GroupService) GetByCustomID(ctx context.Context,
	customID *string,
) (*Group, error) {
	if isEmptyString(customID) {
		return nil, fmt.Errorf("customID cannot be nil for Get operation")
	}

	type QS struct {
		CustomID string `url:"custom_id,omitempty"`
	}

	req, err := s.client.NewRequest("GET", "/Groups",
		&QS{CustomID: *customID}, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data []Group
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, NewAPIError(http.StatusNotFound, "Not found")
	}

	return &resp.Data[0], nil
}

// Update updates a Group in Kong
func (s *GroupService) Update(ctx context.Context,
	group *Group,
) (*Group, error) {
	if isEmptyString(group.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/groups/%v", *group.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, group)
	if err != nil {
		return nil, err
	}
	type Response struct {
		group Group
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.group, nil
}

// Delete deletes a Group in Kong
func (s *GroupService) Delete(ctx context.Context,
	emailOrID *string,
) error {
	if isEmptyString(emailOrID) {
		return fmt.Errorf("emailOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/groups/%v", *emailOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Groups in Kong.
// opt can be used to control pagination.
func (s *GroupService) List(ctx context.Context,
	opt *ListOpt,
) ([]*Group, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/groups", opt)
	if err != nil {
		return nil, nil, err
	}
	var Groups []*Group

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var Group Group
		err = json.Unmarshal(b, &Group)
		if err != nil {
			return nil, nil, err
		}
		Groups = append(Groups, &Group)
	}

	return Groups, next, nil
}

// ListAll fetches all Groups in Kong.
// This method can take a while if there
// a lot of Groups present.
func (s *GroupService) ListAll(ctx context.Context) ([]*Group, error) {
	var Groups, data []*Group
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		Groups = append(Groups, data...)
	}
	return Groups, nil
}
