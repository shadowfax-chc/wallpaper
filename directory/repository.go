package directory

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/shadowfax-chc/wallpaper/wallpaper"
)

// Config is used to configure a repository.
type Config struct {
	Root    string
	Shuffle bool
}

// Repository of wallpapers in a file system's directory/folder.
type Repository struct {
	Root string // The path to the directory

	shuffle bool // If the repository should randomize the order of the images.

	images    []wallpaper.Image
	lastIndex int
}

// NewRepository creates a repository.
func NewRepository(c *Config) (wallpaper.Repository, error) {
	info, err := os.Lstat(c.Root)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		realpath, _ := filepath.EvalSymlinks(c.Root)
		log.Printf("[DEBUG] root %s is symlink to %s", c.Root, realpath)
		c.Root = realpath
	}
	return &Repository{
		Root:    c.Root,
		shuffle: c.Shuffle,
	}, nil
}

func (r *Repository) scandir() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		log.Printf("[DEBUG] checking if file %s is an image", path)
		if wallpaper.IsImage(path) {
			r.images = append(r.images, wallpaper.Image(path))
		}
		return nil
	}
}

func (r *Repository) scan(dir string) error {
	log.Printf("[DEBUG] scanning directory: %s", dir)
	return filepath.Walk(dir, r.scandir())
}

func (r *Repository) reshuffle() {
	shuffled := make([]wallpaper.Image, len(r.images))
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := random.Perm(len(r.images))
	for i, randIndex := range perm {
		shuffled[i] = r.images[randIndex]
	}
	r.images = shuffled
}

// Load initializes the Repository, scanning the directory for images.
func (r *Repository) Load() error {
	if err := r.scan(r.Root); err != nil {
		return err
	}
	if r.shuffle {
		log.Printf("[DEBUG] shuffling images")
		r.reshuffle()
	}
	return nil
}

// Reload refreshes the Repository. It will rescan the directory for images.
func (r *Repository) Reload() error {
	err := r.Load()
	if err != nil {
		return err
	}
	r.lastIndex = 0
	return nil
}

// Next returns the next repository in the list.
func (r *Repository) Next() wallpaper.Image {
	r.lastIndex++
	if r.lastIndex >= len(r.images) {
		r.lastIndex = 0
	}
	return r.images[r.lastIndex]
}

// SetLocation updates the root point of the repository
func (r *Repository) SetLocation(root string) {
	r.Root = root
}

// SetShuffle sets the repository to shuffle
func (r *Repository) SetShuffle(shuffle bool) {
	r.shuffle = shuffle
}
