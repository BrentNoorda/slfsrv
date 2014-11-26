slfsrv
============

### Create simple, cross-platform GUI applications, or wrap GUIs around command-line applications, using HTML/JS/CSS and your own browser.

__slfsrv__ is an executable that lets you create GUI applications (using simple HTML, CSS, and
Javascript) that launch the default web browser as both the logic system and graphical front end. A Javascript API is
available to those application to provide enough access to the local computer system, files, and
executable, to create basic computer applications or to wrap less-friendly applications within a browser front end.

Jump To:

* [primary features](#primary-features)
* [why was this made and what situations is it good for?](#who-for)
* [usage, helpscreen version](#usage-helpscreen)
* [javascript api](#jsapi)
* [config file](#configfile)
* [getting the source/executable files](#getting-files)
* [getting started with demos](#getting-started)
* [who to contact](#todo)
* [history](#history)
* [project homepage](#home)

------------------------------------------------------------------------------
<a name="primary-features"></a>
### Primary features:

* Run as server to deliver local HTML/CSS/JS in a local browser, as a localhost app.

* Provide JS API for access to the local computer system, including: data storage, file
and directory read/write, unpacking of bundled data (including other executables), and
ability to execute other programs.

* Run the same applications from within Windows, Mac/OSX, and Linux.

* Bundle applications into either single-file "slfsrv" files, or compile into
standalone executables. (Note, a Mac user can create a standalone Windows or Linux executable,
and vice versa).

* Blur the lines between what runs on the web, and what runs on the local machine. (i.e.
no longer must I launch a browser to run certain kinds of tasks, and some local start shell
to run local tasks, when now they can all launch from the browser).

------------------------------------------------------------------------------
<a name="who-for"></a>
### Why was this made and what situations is it good for?

Putting a GUI front on a little program is too often more work than its worth, especially
if that isn't a "massive, commercial, sellable" program. Sure, if you're making a commercial product
to sell to millions, you want a fancy user interface; but if it's just an one-off tool
for weekly departmental TPS reports, then it's good enough to have the kind of simple
GUI a web browser front-end makes easy.

slfsrv is a general-purpose localhost app maker, and as such it could be used to make
"sophisticated tools". For example, it could wrap around git to make a nice GUI UI (if there
weren't already so many of them). But I expect it more often to be used for situations
like these:

* you need a front end to an existing command-line tool (e.g. collect information, run
the command-line tool, the nicely format the output in the browser window)

* you want to create an application to run on your user's computers, but don't want to spend
a lot of time making it

* for creating your one all-powerful home screen, organizing your daily tasks by combining web activities and local-program activites
in the same menu UI (see the menu example under [getting started with demos](#getting-started))

* you want a first-order prototype of a new product idea

* Creation and wrapping of installation tools. Suppose you need to install a new executable
on everyone's system. You can embed that other executable and other scripts within a slfsrv executable, put a bunch of
checks on inputs in front of that with slfsrv, and the result would be a quick and slick and better
installation tool.

* Help people through your organization's procedures. Instead of publishing a list of steps people
must perform (e.g. "step 1, map to drive with this command and your name and password") you can create a GUI
front end around those, write out command in the background, embedding something like [Expect](http://en.wikipedia.org/wiki/Expect) if they're
complicated, and walk users through the procedures in a much friendlier manor).

------------------------------------------------------------------------------
<a name="usage-helpscreen"></a>
### Usage, helpscreen version:

From `#slfsrv --help`

    slfsrv - Launch web browser and serve local HTML/CSS/JS/ASSET files to
             it; or compile everything into a standalone executable, for
             multiple platforms, as distributable localhost apps;
             or wrap a GUI around a command-line program.
             For more info, see https://github.com/BrentNoorda/slfsrv/

    USAGE: slfsrv [flags...] [filePath [initialUrlPath]
    WHERE: filePath - path to initial file or directory to serve; if file (e.g.
               index.html) then start with that file; if directory then
               look for index.html (or index.htm) in that directory; if not
               specified then look for index.html (or index.htm) in the
               current directory
               path may also be a *.slfsrv bundle created with -bundle
           initialUrlPath - when filePath is simply a directory, then initialUrlPath
               is the path to the first page to load within that directory
    flags:
      -help - this help text
      -bundle <out.slfsrv> - zip all contents of path into self-contained file
                             to be executed later with slfsrv - the .slfsrv
                             extension is required
      -compile <output-name> - compile slfsrv and all contents of directory
                               tree into a standalone executable
      -compile-from <ss-source> - used with the '-compile' flag, this specifies a
                                  slfsrv executable to use (else will use
                                  the current slfsrv executable); this flag is
                                  useful for making cross-platform distributals (e.g.
                                  creating a windows executable from OSX)
      -replace - used with the '-compile' or '-bundle' flags, this specifies
                 whether it is OK to replace the <output-name> file
      -port <port#> - specify port to run server on, else will look for a random port
                      that is not in use
      -config-file <json file> - path to json config file name
      -verbose - write out lots of message about what's going on (else silent)
    EXAMPLES:
      slfsrv
      slfsrv /myapp
      slfsrv ../htmlgame/index.html
      slfsrv -compile /myapp
      slfsrv -compile /myapp -compile-from /tools/ss/win/slfsrv.exe
      slfsrv c:\user\me photo/tools/photoview.html
      slfsrv c:\user\me photo/tools/photoview.html -config-file ./config.json
    VERSION: 0.0.1

------------------------------------------------------------------------------
<a name="jsapi"></a>
### JavaScript API:

[&#x25BA; Javascript API](docs/JSAPI.md) - calls made available to javascript running in the browser

------------------------------------------------------------------------------
<a name="configfile"></a>
### Configuration file

If you supply a configuration file, it should be json with the following optional fields:

keepAliveSeconds: (int64) number of seconds to wait for a keepalive message before quitting.
Supply a large number to allow the app to stay alive even if page is pinging the server.

secretKey: (string) client has this key for rudimentary security. If none supplied, the app
will generate one.

port: (int) port to run server on.  If none supplied it will find an available port.
Note that the port can also be supplied on the command line.


------------------------------------------------------------------------------
<a name="getting-files"></a>
### Getting the source/executable files

#### Downloading pre-built executables

If you don't want to compile from original source, then download and unzip the "examples"
directory (you'll really want to run the examples), and then the executable appropriate to
your OS.

* [slfsrv-examples.zip](https://dl.dropboxusercontent.com/u/41075/slfsrv-downloads/slfsrv-examples.zip) - the examples you'll want
to start with for any of the following executables

* [slfsrv-darwin.zip](https://dl.dropboxusercontent.com/u/41075/slfsrv-downloads/slfsrv-darwin.zip) - for Macintosh/OSX users

* [slfsrv-windows.zip](https://dl.dropboxusercontent.com/u/41075/slfsrv-downloads/slfsrv-windows.zip) - for Windows users

* [slfsrv-linux.zip](https://dl.dropboxusercontent.com/u/41075/slfsrv-downloads/slfsrv-linux.zip) - for Ubuntu Linux - *NOTE:
I haven't tested this a lot, and I don't know how portable ubuntu linux portables are...*

#### Building from source code

`slfsrv` source is available at [github.com/BrentNoorda/slfsrv](https://github.com/BrentNoorda/slfsrv). It is written
in the "Go" language. If you don't already have Go installed, then follow these [Go Installation Instructions](http://golang.org/doc/install).

If you're familiar with Go then you already know how you like to install and build new Go applications, and so
you'll ignore the following OSX, Linux, or Windows instructions. If you are not familiar with Go,
then you may want to build `slfsrv` in it's own directory.

##### Building on OSX or Linux

To build in the `~/slfsrv` directory:

    $ mkdir ~/slfsrv
    $ cd ~/slfsrv
    $ export GOPATH=~/slfsrv
    $ go get github.com/BrentNoorda/slfsrv
    $ go install github.com/BrentNoorda/slfsrv

If those commands to download and build worked then the `slfsrv` executable will have been created in the `~/slfsrv/bin`
directory, and this command:

    $ ~/slfsrv/bin/slfsrv --help

will output the help screen seen [here](https://github.com/BrentNoorda/slfsrv#usage-helpscreen).

##### Building on Windows

To build in the `\slfsrv` directory:

    $ mkdir \slfsrv
    $ cd \slfsrv
    $ set GOPATH=\slfsrv
    $ go get github.com/BrentNoorda/slfsrv
    $ go install github.com/BrentNoorda/slfsrv

If those commands to download and build worked then the `slfsrv` executable will have been created in the `\slfsrv\bin`
directory, and this command:

    $ \slfsrv\bin\slfsrv --help

will output the help screen seen [here](https://github.com/BrentNoorda/slfsrv#usage-helpscreen).

------------------------------------------------------------------------------
<a name="getting-started"></a>
### Getting started with demos

When you have the executable and examples (and assuming the `slfsrv` executable is in your path),
then from the `src/github.com/BrentNoorda/slfsrv` directory run:

    $ slfsrv -verbose EXAMPLES.html

This should launch a browser that will allow you to execute some basic examples, along with viewing their source code.

My favorite example is "menu", which I now use (instead of my macintosh DOC) to launch just about everything, using the additional
`slfsrv-menu.js` script (documented [here](docs/SLFSRV-MENU.md)).

------------------------------------------------------------------------------
<a name="todo"></a>
### ISSUES / TODO / TOANSWER

[&#x25BA; TODO](docs/TODO.md) - what's still to be done

---------------------------------------------------------------------

random info to be better documented later:

For mac, to associate .slfsrv with an app so you can execute the bundle simply by double-clicking it
in the Finder, or by using the "open..." command in Terminal, follow these instructions for
[How to associate a file name extension with a darwin /linux/unix command-line executable on osx/macintosh](http://unlessimwrong.blogspot.com/2014/07/how-to-associate-file-name-extension.html)

------------------------------------------------------------------------------
<a name="history"></a>
### History

* 2014/10/07 - version 0.0.1 - first release

------------------------------------------------------------------------------
<a name="contact"></a>
### Who to contact

[&#x25BA; Brent Noorda](http://www.brent-noorda.com/) - all about the author and how to contact

------------------------------------------------------------------------------
<a name="home"></a>
### Project Homepage

[&#x25BA; github.com/BrentNoorda/slfsrv](https://github.com/BrentNoorda/slfsrv) - slfsrv project homepage

