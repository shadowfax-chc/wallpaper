package wallpaper

// Repository is a collection of wallpapers that can be used for a background.
type Repository interface {
	// Load initializes the Repository.
	Load() error
	// Reload refresh the Repository.
	Reload() error
	// Next returns the next Wallpaper to be used.
	Next() Image
	// SetShuffle sets the Repository to randomize images or not.
	SetShuffle(bool)
	// SetLocation sets a new location string for the Repository source.
	SetLocation(string)
}
