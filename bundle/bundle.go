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

/***************** SelfserverReadSeeker - ReadSeeker implementation *****************/

type SelfservingReadSeeker struct {
	seek   int64
	curpos int64
	f      *zip.File
	fSize  int64
	rc     io.ReadCloser
}

func NewSelfservingReadSeeker(reader *zip.Reader, filespec string) (*SelfservingReadSeeker, error) {
	var f *zip.File
	for _, f = range reader.File {
		if f.Name == filespec {
			return &SelfservingReadSeeker{seek: 0, curpos: 0, f: f, fSize: int64(f.UncompressedSize64), rc: nil}, nil
		}
	}
	return nil, errors.New("file \"" + filespec + "\" not found in bundle")
}

func (ssrs *SelfservingReadSeeker) Close() {
	if ssrs.rc != nil {
		ssrs.rc.Close()
		ssrs.rc = nil
		ssrs.seek = 0
		ssrs.curpos = 0
	}
}

func (ssrs *SelfservingReadSeeker) Read(p []byte) (int, error) {
	var pLen int = len(p)

	if pLen == 0 {
		// asking for no data
		return 0, nil
	} else if ssrs.fSize <= ssrs.seek {
		// cannot read beyond EOF
		return 0, io.EOF
	} else {

		var err error

		if ssrs.rc == nil {
			if ssrs.rc, err = ssrs.f.Open(); err != nil {
				return 0, err
			}
		}

		var readCount int

		// if not currently at the read position, then skip to there
		if ssrs.curpos != ssrs.seek {
			var bigBuffer [100000]byte
			for ssrs.curpos < ssrs.seek {
				var wantToReadCount int64 = ssrs.seek - ssrs.curpos
				if int64(len(bigBuffer)) < wantToReadCount {
					wantToReadCount = int64(len(bigBuffer))
				}
				readCount, err = io.ReadFull(ssrs.rc, bigBuffer[0:wantToReadCount])
				if (err != nil) || (readCount != int(wantToReadCount)) {
					return 0, io.EOF
				}
				ssrs.curpos += wantToReadCount
			}
		}

		readCount, err = io.ReadFull(ssrs.rc, p)
		ssrs.curpos += int64(readCount)
		ssrs.seek = ssrs.curpos
		return readCount, err
	}
}

func (ssrs *SelfservingReadSeeker) Seek(offset int64, whence int) (int64, error) {
	var err error
	if whence == 1 {
		offset += ssrs.seek
	} else if whence == 2 {
		offset = ssrs.fSize + offset
	}
	/* else whence is 0 and sets file position absolutely */
	if offset < 0 {
		err = errors.New("invalid negative file offset")
	} else {
		if ssrs.rc != nil {
			if offset < ssrs.curpos {
				ssrs.Close()
			}
		}
		ssrs.seek = offset
	}
	return ssrs.seek, err
}

/***************** selfserverReaderAt - ReaderAt implementation *****************/

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
