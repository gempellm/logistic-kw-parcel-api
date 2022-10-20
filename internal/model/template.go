package model

// parcel - parcel entity.
type parcel struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
}
