//
// Copyright (c) 2018 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

import (
	"fmt"
)

// The pendingop.go file defines the basic structures needed to track
// life-cycle of database entries w/in Heketi. There are generally two
// levels of objects which we track: the pending operation a higher-level
// operation that should roughly correspond to an action in the Heketi API,
// and lower level change tracking where individual object creates, deletes,
// and modifications are logged. Note that objects such as bricks and
// volumes are still managed in their own buckets. The pending operation
// metadata tracks the IDs of these objects and these objects have references
// back to their associated pending operations (via IDs).

// PendingOperationType identifies what kind of high-level operation a
// PendingOperation will be.
type PendingOperationType int

const (
	OperationUnknown PendingOperationType = iota
	OperationCreateVolume
	OperationDeleteVolume
	OperationExpandVolume
	OperationCreateBlockVolume
	OperationDeleteBlockVolume
	OperationExpandBlockVolume
	OperationRemoveDevice
	OperationCloneVolume
	OperationBrickEvict
)

// PendingChangeType identifies what kind of lower-level new item or change
// is being made to the system as part of a higher-level pending operation.
type PendingChangeType int

const (
	OpUnknown PendingChangeType = iota
	OpAddBrick
	OpAddVolume
	OpDeleteBrick
	OpDeleteVolume
	OpExpandVolume
	OpAddBlockVolume
	OpDeleteBlockVolume
	OpExpandBlockVolume
	OpRemoveDevice
	OpCloneVolume
	OpSnapshotVolume
	OpAddVolumeClone
	OpChildOperation
	OpParentOperation
)

// PendingOperationAction tracks individual changes to entries within the
// heketi db. It consists of a required change type and (heketi uuid) id,
// as well as an optional delta object for extra metadata.
type PendingOperationAction struct {
	Change PendingChangeType
	Id     string
	Delta  interface{}
}

// PendingItem encapsulates the common pending item ID field.
type PendingItem struct {
	Id string
}

// PendingOperation tracks higher-level changes to the heketi system, such
// as volume creation or deletion.
type PendingOperation struct {
	PendingItem
	Timestamp int64
	Type      PendingOperationType
	Actions   []PendingOperationAction
}

// ExpandSize extracts an int value for a pending size expansion from the
// PendingOperationAction if the change type is correct. If the type is
// not correct error will be non-nil.
func (a PendingOperationAction) ExpandSize() (int, error) {
	if a.Change == OpExpandVolume {
		if v, ok := a.Delta.(int); ok {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Action delta for ExpandSize is missing/invalid")
}

// Name returns the pending operation type as a brief string.
// NOTE: Stringer was considered but not used as the literal
// names of the variables were not desired. Thus to avoid
// conflict with possible future use of stringer for other
// purposes, this was called Name.
func (v PendingOperationType) Name() string {
	switch v {
	case OperationCreateVolume:
		return "create-volume"
	case OperationDeleteVolume:
		return "delete-volume"
	case OperationExpandVolume:
		return "expand-volume"
	case OperationCreateBlockVolume:
		return "create-block-volume"
	case OperationDeleteBlockVolume:
		return "delete-block-volume"
	case OperationExpandBlockVolume:
		return "expand-block-volume"
	case OperationRemoveDevice:
		return "remove-device"
	case OperationCloneVolume:
		return "clone-volume"
	case OperationBrickEvict:
		return "evict-brick"
	}
	return "unknown"
}

// Name returns a short description of a change action.
func (c PendingChangeType) Name() string {
	switch c {
	case OpAddBrick:
		return "Add brick"
	case OpAddVolume:
		return "Add volume"
	case OpDeleteBrick:
		return "Delete brick"
	case OpDeleteVolume:
		return "Delete volume"
	case OpExpandVolume:
		return "Expand volume"
	case OpAddBlockVolume:
		return "Add block volume"
	case OpDeleteBlockVolume:
		return "Delete block volume"
	case OpExpandBlockVolume:
		return "Expand block volume"
	case OpRemoveDevice:
		return "Remove device"
	case OpCloneVolume:
		return "Clone volume from"
	case OpSnapshotVolume:
		return "Snapshot volume"
	case OpAddVolumeClone:
		return "Expand volume to"
	case OpChildOperation:
		return "Performing child operation"
	case OpParentOperation:
		return "Belongs to parent operation"
	}
	return "Unknown"
}
