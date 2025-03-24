package osm

// UserID is the primary key for a user.
// This is unique as the display name may change.
type UserID int64

// ObjectID is a helper returning the object id for this user id.
func (id UserID) ObjectID() ObjectID {
	return ObjectID(userMask | (id << versionBits))
}
