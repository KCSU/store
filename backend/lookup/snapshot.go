package lookup

import (
	"fmt"

	"github.com/kcsu/store/db"
	"github.com/kcsu/store/logging"
	"github.com/kcsu/store/model"
	lookupv1 "github.com/kcsu/store/proto/kcsu/lookup/v1"
)

func ApplySnapshot(store db.GroupStore, s *lookupv1.GroupSnapshot) error {
	users := make([]model.GroupUser, 0, len(s.Members))
	for _, member := range s.Members {
		if member.Crsid != "" {
			users = append(users, model.GroupUser{UserEmail: member.Crsid + "@cam.ac.uk"})
		}
	}
	// Sanity check: lookup groups shouldn't be empty, as this would remove everyone
	if len(users) == 0 {
		return fmt.Errorf("snapshot of %s has no members", s.Name)
	}
	groups, err := store.Get()
	if err != nil {
		return err
	}
	applied := 0
	for _, group := range groups {
		if group.Type != "group" || group.Lookup != s.Name {
			continue
		}
		if err := store.ReplaceLookupUsers(&group, users); err != nil {
			return err
		}
		applied++
	}
	logging.Infof("applied Lookup snapshot %s (%d members) to %d groups", s.Name, len(users), applied)
	return nil
}
