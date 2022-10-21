package model

import "google.golang.org/protobuf/types/known/timestamppb"

type Parcel struct {
	ID      uint64
	Name    string
	Created *timestamppb.Timestamp
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type ParcelEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Parcel
}
