package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const Version string = "0.0.1"

type selfservingReaderAt struct {
	offset int64
	f      *os.File
}

func (s selfservingReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	return s.f.ReadAt(p, off+s.offset)
}

func ReadFile(reader *zip.Reader, filespec string) ([]byte, error) {

	var f *zip.File
	for _, f = range reader.File {
		if f.Name == filespec {
			var readcount int
			var err error
			var filesize uint64 = f.UncompressedSize64
			var data []byte
			var rc io.ReadCloser

			if rc, err = f.Open(); err != nil {
				return data, err
			}
			if data, err = ioutil.ReadAll(rc); err != nil {
				rc.Close()
			} else {
				err = rc.Close()
			}
			if err != nil {
				return data, err
			}

			readcount = len(data)
			if uint64(readcount) != filesize {
				return data, errors.New("bytes missing from file")
			}

			return data, nil

		}
	}
	return nil, errors.New("file \"" + filespec + "\" not found in bundle")
}

func Reader(zipFilespec string, offsetByte0 int) (reader *zip.Reader, initialUrl string, err error) {

	var fileSize int64
	var f *os.File

	if f, err = os.Open(zipFilespec); err != nil {
		return
	}

	// determine total size of file zipFilespec
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return
	}
	fileSize = fi.Size()

	var readerAt selfservingReaderAt = selfservingReaderAt{offset: int64(offsetByte0), f: f}

	if reader, err = zip.NewReader(readerAt, fileSize-readerAt.offset); err != nil {
		return
	}

	// check that the version string is correct, and pull out initialUrl
	var data []byte
	if data, err = ReadFile(reader, "slfsrv"); err != nil {
		return
	}
	/* not checking the slfsrv version yet */

	if data, err = ReadFile(reader, "initialUrl"); err != nil {
		return
	}
	initialUrl = string(data)

	return
}

func BundleDirectory(rootPath string, initialUrl string) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	var w *zip.Writer = zip.NewWriter(buf)

	// the first entry is always version, then all of the files
	var f io.Writer
	f, _ = w.Create("slfsrv")
	f.Write([]byte(Version))

	// the second entry is always initialUrl, then all of the files
	f, _ = w.Create("initialUrl")
	f.Write([]byte(initialUrl))

	// rootpath must be a directory
	var isDir bool
	if isDir, err = ssutil.IsDir(rootPath); err != nil {
		return
	}
	if !isDir {
		err = errors.New("path \"" + rootPath + "\" must be a directory")
		return
	}

	// add ALL of the files in rootPath
	var skipLen int = len(rootPath) + 1
	var fileData []byte
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				var relativePath string = path[skipLen:]

				if f, err = w.Create(filepath.ToSlash(relativePath)); err != nil {
					return err
				}

				if fileData, err = ioutil.ReadFile(path); err != nil {
					return err
				}

				if _, err = f.Write(fileData); err != nil {
					return err
				}

			}
		}
		return err
	}
	if err = filepath.Walk(rootPath, walkFunc); err != nil {
		return
	}

	err = w.Close()

	return
}

func saveOutputFile(filespec string, data []byte) error {
	return ioutil.WriteFile(filespec, data, 0644)
}

func Bundle(outFilespec string, replaceOK bool, rootPath string, initialUrl string, verbose bool) (err error) {

	// output must be file name, not a directory
	var isDir bool
	if isDir, err = ssutil.IsDir(outFilespec); err != nil {
		return err
	}
	if isDir {
		return errors.New("bundle file \"" + outFilespec + "\" must not be a directory")
	}

	// verify that there is not already an output file, or it's OK to create that file
	if !replaceOK {
		var fileOrDirExists bool
		if fileOrDirExists, err = ssutil.FileOrDirExists(outFilespec); err != nil {
			return err
		}
		if fileOrDirExists {
			return errors.New("bundle file \"" + outFilespec + "\" already exists; cannot overwrite without -replace flag")
		}
	}

	var zipContents *bytes.Buffer
	if zipContents, err = BundleDirectory(rootPath, initialUrl); err != nil {
		return err
	}

	err = saveOutputFile(outFilespec, zipContents.Bytes())

	return
}
