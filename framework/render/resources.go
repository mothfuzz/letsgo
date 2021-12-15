package render

import (
	"embed"
)

//to be set by user
var Resources embed.FS

//go:embed base_resources
var BaseResources embed.FS

//convenience function to check both Resources and BaseResources
func ReadResource(path string) ([]byte, error) {
	bytes, err := Resources.ReadFile("resources/" + path)
	if err != nil {
		bytes, err = BaseResources.ReadFile("base_resources/" + path)
	}
	return bytes, err
}
