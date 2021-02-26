package terseurl

import (
	"embed"
	"io/fs"
	"os"
)

const (

	// frontendDirName is the name of the directory to serve frontend assets.
	frontendDirName = "frontend"
)

// frontend is the embedded file system for frontend assets.
//go:embed frontend
var frontend embed.FS

// TODO Group var declaration?

// RedirectTemplate is the embedded HTML template for issuing redirects.
//go:embed redirect.gohtml
var RedirectTemplate string

// FrontendFS is the file system for the frontend assets: HTML, JS, etc. It can use either the embedded assets or a
// directory on the host OS.
func FrontendFS(dirName string) (fileSystem fs.FS, err error) {

	// Check to see if the live file system's directory should be used.
	if dirName != "" {

		// Stat the directory in the OS to confirm it's good to go.
		if _, err = os.Stat(dirName); err != nil {
			return nil, err
		}
		fileSystem = os.DirFS(dirName)
	} else {

		// Turn the embedded file system into the appropriate Go type.
		if fileSystem, err = fs.Sub(frontend, frontendDirName); err != nil {
			return nil, err
		}
	}

	return fileSystem, nil
}
