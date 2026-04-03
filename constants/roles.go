package constants

// User roles used across the system for access control
const (
	RoleAdmin   = "admin"   // Full access (manage users, data, settings)
	RoleAnalyst = "analyst" // Can analyze and manage data
	RoleViewer  = "viewer"  // Read-only access
)

// Common status values for entities (users, records, etc.)
const (
	StatusActive   = "active"   // Currently active/usable
	StatusInactive = "inactive" // Disabled or not in use
)