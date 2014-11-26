package main

import (
	"archive/zip"
	"bitbucket.org/kardianos/osext"
	"errors"
	"flag"
	"fmt"
	"github.com/BrentNoorda/slfsrv/browser"
	"github.com/BrentNoorda/slfsrv/bundle"
	"github.com/BrentNoorda/slfsrv/compiler"
	"github.com/BrentNoorda/slfsrv/config"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"github.com/BrentNoorda/slfsrv/tempdir"
	"github.com/BrentNoorda/slfsrv/webserver"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func parse_command_line() (helpWanted bool,
	bundleOut string,
	compile string, compileFrom string, compileReplaceOK bool,
	port int, verbose bool, configFile string, args []string, err error) {
	err = nil
	flag.BoolVar(&helpWanted, "help", false, "show help screen")
	flag.StringVar(&bundleOut, "bundle", "", "bundle directory into single file")
	flag.StringVar(&compile, "compile", "", "compile into standalone executable")
	flag.StringVar(&compileFrom, "compile-from", "", "executable source for standalone")
	flag.BoolVar(&compileReplaceOK, "-replace", false, "OK to replace bundle or compile output")
	flag.IntVar(&port, "port", 0, "port to server localhost app from")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.StringVar(&configFile, "config-file", "", "filename of configuration json file")

	flag.Parse()
	args = flag.Args()

	if 2 < len(args) {
		err = errors.New("Huh? too many paths specified")
	} else if (bundleOut != "") && (compile != "") {
		err = errors.New("Cannot use both -bundle and -compile flags")
	} else if compileReplaceOK && (compile == "") && (bundleOut == "") {
		err = errors.New("Huh? -replace has no meaning without -compile or -bundle")
	} else if (compileFrom != "") && (compile == "") {
		err = errors.New("Huh? compile-from has no meaning without -compile")
	} else {
		if bundleOut != "" {
			lowerBundle := strings.ToLower(bundleOut)
			if !strings.HasSuffix(lowerBundle, ".slfsrv") {
				bundleOut = bundleOut + ".slfsrv"
			}
		}
	}
	return
}

func parse_compiled_command_line() (port int, verbose bool, err error) {
	err = nil
	flag.IntVar(&port, "port", 0, "port to server localhost app from")
	flag.BoolVar(&verbose, "verbose", false, "verbose")

	flag.Parse()

	return
}

func generate_secret_key() string {
	s := ""
	for i := 0; i < 5; i++ {
		var seg int = 100000 + rand.Intn(100000) // 100000 -> 199999
		s += strconv.Itoa(seg)
	}
	return s
}

func usage() {
	for _, t := range [...]string{
		"",
		"slfsrv - Launch web browser and serve local HTML/CSS/JS/ASSET files to",
		"         it; or compile everything into a standalone executable, for",
		"         multiple platforms, as distributable localhost apps;",
		"         or wrap a GUI around a command-line program.",
		"         For more info, see https://github.com/BrentNoorda/slfsrv/",
		"",
		"USAGE: slfsrv [flags...] [filePath [initialUrlPath]",
		"WHERE: filePath - path to initial file or directory to serve; if file (e.g.",
		"           index.html) then start with that file; if directory then",
		"           look for index.html (or index.htm) in that directory; if not",
		"           specified then look for index.html (or index.htm) in the",
		"           current directory",
		"           path may also be a *.slfsrv bundle created with -bundle",
		"       initialUrlPath - when filePath is simply a directory, then initialUrlPath",
		"           is the path to the first page to load within that directory",
		"flags:",
		"  -help - this help text",
		"  -bundle <out.slfsrv> - zip all contents of path into self-contained file",
		"                         to be executed later with slfsrv - the .slfsrv",
		"                         extension is required",
		"  -compile <output-name> - compile slfsrv and all contents of directory",
		"                           tree into a standalone executable",
		"  -compile-from <ss-source> - used with the '-compile' flag, this specifies a",
		"                              slfsrv executable to use (else will use",
		"                              the current slfsrv executable); this flag is",
		"                              useful for making cross-platform distributals (e.g.",
		"                              creating a windows executable from OSX)",
		"  -replace - used with the '-compile' or '-bundle' flags, this specifies",
		"             whether it is OK to replace the <output-name> file",
		"  -port <port#> - specify port to run server on, else will look for a random port",
		"                  that is not in use",
		"  -verbose - write out lots of message about what's going on (else silent)",
		"EXAMPLES:",
		"  slfsrv",
		"  slfsrv /myapp",
		"  slfsrv ../htmlgame/index.html",
		"  slfsrv -compile /myapp",
		"  slfsrv -compile /myapp -compile-from /tools/ss/win/slfsrv.exe",
		"  slfsrv c:\\user\\me photo/tools/photoview.html",
		"VERSION: " + bundle.Version,
		"",
	} {
		fmt.Println(filepath.FromSlash(t))
	}
}

func main() {

	var precompiled int = compiler.OffsetOfAppendedData()
	var err error = nil
	var bundleOut, compile, compileFrom, patharg, fullpatharg, initPathUrl, initQuery, storeFilespec, configFile string
	var helpWanted, compileReplaceOK, verbose, isdir bool
	var port int
	var cwd string
	var args []string
	var zipReader *zip.Reader

	rand.Seed(time.Now().Unix())

	var myselfExecutable string
	myselfExecutable, err = osext.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine location of self-server executable\n")
		return
	}

	if precompiled != 0 {
		helpWanted = false
		bundleOut = ""
		compile = ""
		compileFrom = ""
		configFile = ""
		fullpatharg = ""
		patharg = ""
		compileReplaceOK = false
		isdir = false
		args = []string{}

		port, verbose, err = parse_compiled_command_line()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		zipReader, initPathUrl, err = bundle.Reader(myselfExecutable, precompiled)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		storeFilespec = myselfExecutable + ".slfsrv-data"

	} else {

		// parse command-line inputs
		helpWanted, bundleOut, compile, compileFrom, compileReplaceOK, port, verbose, configFile, args, err = parse_command_line()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			helpWanted = true
		}
		if helpWanted {
			usage()
			return
		}

		// figure out the initial path and root file to serve from
		if cwd, err = os.Getwd(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		cwd = filepath.ToSlash(cwd)

		config.ParseConfigFile(configFile)
		settings := config.GetSettings()
		fmt.Printf("config: %+v\n", *settings)

		if settings.Port != 0 {
			port = settings.Port
		}

		if len(args) < 2 {
			initPathUrl = ""
		} else {
			initPathUrl = args[1]
		}
		if len(args) == 0 {
			patharg = cwd
		} else {
			patharg = filepath.ToSlash(path.Clean(args[0]))
		}

		var queryIdx int = strings.Index(initPathUrl, "?")
		if queryIdx == -1 {
			initQuery = ""
		} else {
			initQuery = initPathUrl[queryIdx+1:]
			initPathUrl = initPathUrl[:queryIdx]
		}

		if initPathUrl == "" {
			// figure out if path arg is just a path (in which case add index.htm or index.html) or a full file
			if isdir, err = ssutil.IsDir(patharg); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			if !isdir {
				patharg, initPathUrl = path.Split(patharg)
			}
		}
		fullpatharg = path.Clean(patharg)

		// if patharg is a subfolder of cwd then skip the cwd part so it is prettier
		patharg = ssutil.CleanPathRelativeToCwd(fullpatharg, cwd)

		// if initPathUrl not specified then guess at "index.htm" or "index.html"
		if initPathUrl == "" {
			var exists bool
			initPathUrl = "index.htm"
			if exists, _ = ssutil.FileOrDirExists(filepath.FromSlash(fullpatharg + "/" + initPathUrl)); !exists {
				initPathUrl = "index.html"
			}
		}
		storeFilespec = filepath.FromSlash(path.Join(fullpatharg, initPathUrl+".slfsrv-data"))

	}

	if bundleOut != "" {
		// BUNDLE into a standalone zip file
		if verbose {
			fmt.Printf("create compiled version from %s\n", filepath.FromSlash(patharg))
		}
		bundleOut = ssutil.CleanPathRelativeToCwd(bundleOut, cwd)
		err = bundle.Bundle(bundleOut, compileReplaceOK, patharg, initPathUrl, verbose)
		if err == nil {
			if verbose {
				fmt.Println("bundle SUCCESS! Created:", bundleOut)
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}

	} else if compile != "" {
		// COMPILE
		if compileFrom == "" {
			compileFrom = myselfExecutable
		}
		compileFrom = ssutil.CleanPathRelativeToCwd(compileFrom, cwd)
		if verbose {
			fmt.Printf("create compiled version from %s executable\n", compileFrom)
			fmt.Printf("with source files from %s/%s\n", filepath.FromSlash(patharg), initPathUrl)
		}

		compile = ssutil.CleanPathRelativeToCwd(compile, cwd)
		err = compiler.Compile(compile, compileReplaceOK, compileFrom, patharg, initPathUrl, verbose)
		if err == nil {
			if verbose {
				fmt.Println("compilation SUCCESS! Created:", compile)
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
	} else {
		// SERVE
		exitChan := make(chan int)

		// we may be reading from a pre-bundled file, determine now if that is the case
		if (precompiled == 0) && !isdir {
			var zip_initPathUrl string
			zipReader, zip_initPathUrl, err = bundle.Reader(path.Join(patharg, initPathUrl), 0)

			if err == nil {
				storeFilespec = filepath.FromSlash(path.Join(fullpatharg, initPathUrl) + ".slfsrv-data")
				initPathUrl = zip_initPathUrl
			} else {
				zipReader = nil
			}
		}

		var secretKey string
		var keepAliveSeconds int64
		settings := config.GetSettings()

		if settings.SecretKey != "" {
			secretKey = settings.SecretKey
		} else {
			secretKey = generate_secret_key()
		}

		if settings.KeepAliveSeconds != 0 {
			keepAliveSeconds = settings.KeepAliveSeconds
		} else {
			settings.KeepAliveSeconds = 2
		}

		port = webserver.ListenAndServe(port, secretKey, keepAliveSeconds, zipReader, patharg, fullpatharg, initPathUrl,
			verbose, exitChan, myselfExecutable, storeFilespec)

		browser.LaunchDefaultBrowser(port, secretKey, initPathUrl, initQuery, verbose)

		<-exitChan

		tempdir.Cleanup(verbose)
	}
}
