// Package state provides persistence for port scan snapshots.
//
// A Store serialises and deserialises Snapshot values to a JSON file so that
// successive portwatch runs can detect changes between scans.
//
// A Manager wraps a Store and ties together the load-diff-save cycle used by
// the main monitoring loop:
//
//	store := state.NewStore("/var/lib/portwatch/state.json")
//	mgr   := state.NewManager(store)
//
//	changes, err := mgr.Update(currentPorts)
//
On the very first run, when no state file exists, Manager.IsFirstRun returns
// true and Update stores the baseline without reporting any changes.
package state
