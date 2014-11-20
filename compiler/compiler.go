package compiler

import (
	"bytes"
	"errors"
	"github.com/BrentNoorda/slfsrv/bundle"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const secret_offset_of_appended_data = "ApPenDeDdAtAOffSeT_0___________"

func secret_offset_prefix() (prefix string) {
	prefix = secret_offset_of_appended_data[:19]
	return
}

func OffsetOfAppendedData() int {
	var length_string string
	var offset int64
	length_string = secret_offset_of_appended_data[len(secret_offset_prefix()):]
	length_string = length_string[:strings.Index(length_string, "_")]
	offset, _ = strconv.ParseInt(length_string, 10, 32)
	return int(offset)
}

func mustBeARawSelfServingSourceFile(filespec string) (exeFileData []byte, err error) {

	var fileOrDirExists, isDir bool
	if fileOrDirExists, err = ssutil.FileOrDirExists(filespec); err != nil {
		return exeFileData, err
	}
	if isDir, err = ssutil.IsDir(filespec); err != nil {
		return exeFileData, err
	}
	if !fileOrDirExists || isDir {
		return exeFileData, errors.New("compile-from file \"" + filespec + "\" does not look valid")
	}

	// read binary file - if amIAlreadyCompiled is not in there this this isn't right
	if exeFileData, err = ioutil.ReadFile(filespec); err != nil {
		return exeFileData, err
	}
	if !bytes.Contains(exeFileData, []byte(secret_offset_of_appended_data)) {
		return exeFileData, errors.New("compile-from file \"" + filespec + "\" is not valid slfsrv executable")
	}

	return
}

func hardcode_startup_data_into_output_file(data []byte) []byte {

	// get the offset in data where the secret_offset is stored
	var prefix_bytes = []byte(secret_offset_prefix())
	var offset_offset int = bytes.Index(data, prefix_bytes) + len(prefix_bytes)

	// replace the next bytes with the length of the exe file itself (which is the length of data)
	var length_as_bytes []byte = []byte(strconv.Itoa(len(data)))
	copy(data[offset_offset:], length_as_bytes)

	return data
}

func saveAlteredOutputFile(filespec string, data []byte) error {
	return ioutil.WriteFile(filespec, data, 0)
}

func Compile(exeOutPath string, replaceOK bool, exeInPath string, rootPath string, initialUrl string, verbose bool) (err error) {

	var exeFileData []byte

	// verify that exeInPath exists and is a file
	if _, err = mustBeARawSelfServingSourceFile(exeInPath); err != nil {
		return err
	}

	// output must be file name, not a directory
	var isDir bool
	if isDir, err = ssutil.IsDir(exeOutPath); err != nil {
		return err
	}
	if isDir {
		return errors.New("compile file \"" + exeOutPath + "\" must not be a directory")
	}

	// verify that there is not already an output file, or it's OK to create that file
	if !replaceOK {
		var fileOrDirExists bool
		if fileOrDirExists, err = ssutil.FileOrDirExists(exeOutPath); err != nil {
			return err
		}
		if fileOrDirExists {
			return errors.New("compile file \"" + exeOutPath + "\" already exists; cannot overwrite without -replace flag")
		}
	}

	// copy source file to destination (to make sure we get all attribute)
	switch runtime.GOOS {
	case "darwin", "linux":
		var cmd *exec.Cmd = exec.Command("cp", exeInPath, exeOutPath)
		if err = cmd.Run(); err != nil {
			return err
		}
	case "windows":
		var cmd *exec.Cmd = exec.Command("cmd.exe", "/C", "copy", exeInPath, exeOutPath)
		if err = cmd.Run(); err != nil {
			return err
		}
	default:
		return errors.New("Sorry, don't know how to copy file for OS " + runtime.GOOS)
	}
	// verify that exeInPath exists and is a file
	if exeFileData, err = mustBeARawSelfServingSourceFile(exeInPath); err != nil {
		return errors.New("error copying file \"" + exeInPath + "\" to " + "\"" + exeOutPath + "\"")
	}

	exeFileData = hardcode_startup_data_into_output_file(exeFileData)

	// zip all file data
	var zipContents *bytes.Buffer
	if zipContents, err = bundle.BundleDirectory(rootPath, initialUrl); err != nil {
		return err
	}

	// append zipContents to buffer
	var exeFileBuf *bytes.Buffer = bytes.NewBuffer(exeFileData)
	if _, err = exeFileBuf.Write(zipContents.Bytes()); err != nil {
		return err
	}

	err = saveAlteredOutputFile(exeOutPath, exeFileBuf.Bytes())

	return
}
