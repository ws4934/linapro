// This file defines descriptor wrappers published to source-plugin governance
// callbacks for menu and permission filtering.

package pluginhost

import "lina-core/pkg/plugin/pluginhost/internal/descriptor"

// MenuDescriptor exposes one published menu descriptor for plugin menu filtering.
type MenuDescriptor interface {
	// ID returns the menu id.
	ID() int
	// ParentID returns the parent menu id.
	ParentID() int
	// Name returns the menu display name.
	Name() string
	// Path returns the menu path.
	Path() string
	// Component returns the routed component name.
	Component() string
	// Permissions returns the permission key bound to the menu.
	Permissions() string
	// MenuKey returns the stable menu business key.
	MenuKey() string
	// Type returns the menu type.
	Type() string
	// Visible returns the visible status.
	Visible() int
	// Status returns the enabled status.
	Status() int
}

// PermissionDescriptor exposes one published permission descriptor for plugin permission filtering.
type PermissionDescriptor interface {
	// MenuKey returns the stable menu business key.
	MenuKey() string
	// MenuName returns the display name of the menu that owns the permission.
	MenuName() string
	// Permission returns the permission string.
	Permission() string
}

// NewMenuDescriptor creates one published menu descriptor wrapper for plugins.
func NewMenuDescriptor(
	id int,
	parentID int,
	name string,
	path string,
	component string,
	permission string,
	menuKey string,
	menuType string,
	visible int,
	status int,
) MenuDescriptor {
	return descriptor.NewMenu(descriptor.MenuInput{
		ID:         id,
		ParentID:   parentID,
		Name:       name,
		Path:       path,
		Component:  component,
		Permission: permission,
		MenuKey:    menuKey,
		MenuType:   menuType,
		Visible:    visible,
		Status:     status,
	})
}

// NewPermissionDescriptor creates one published permission descriptor wrapper for plugins.
func NewPermissionDescriptor(menuKey string, menuName string, permission string) PermissionDescriptor {
	return descriptor.NewPermission(menuKey, menuName, permission)
}
