#JAVASCRIPT API

All of these methods are available to your javascript files, assuming that the virtual file
`slfsrv-core.js` is included in your primary HTML file; e.g. with
this line is in the `<head>` section:

    <script src="slfsrv-core.js" type="text/javascript"></script>

[&#x25C0; back to README.cmd](../README.md)

------------------------------------------------------------

Javascript API:

* [Callbacks versus Promises](#callbacks-vs-promises) - most of these calls are available using either callback or Promise approches...
* [GLOBALS](#slfsrv-globals) - global variables
* [exec()](#slfsrv-exec) - execute a local computer program or command
* [tempdir()](#slfsrv-tempdir) - get a temporary directory, and unpack files to that directory
* [store](#slfsrv-store) - persistent data storage on local computer
    * [store.get()](#slfsrv-store-get) - get value from persistent store
    * [store.set()](#slfsrv-store-set) - set value in persistent store
    * [store.list()](#slfsrv-store-list) - list all stored keys
* [file](#slfsrv-file) - working with data in local files
    * [file.read()](#slfsrv-file-read) - read data from a file
    * [file.write()](#slfsrv-file-write) - write or append data to a file
    * [file.remove()](#slfsrv-file-remove) - delete a file
    * [file.exists()](#slfsrv-file-exists) - does this file exist
* [dir](#slfsrv-dir) - working with directories (i.e. folders) in local file system
    * [dir.list()](#slfsrv-dir-list) - list all files and directories within this directory
    * [dir.create()] - create a directory
    * [dir.remove()] - delete a directory
    * [dir.exists()](#slfsrv-dir-exists) - does this directory exist
    * [dir.getcwd()](#slfsrv-dir-getcwd) - get current workding directory
    * [dir.setcwd()](#slfsrv-dir-setcwd) - set current workding directory
* [alert](#slfsrv-alert) - non-blocking alert message
    * [alert()](#slfsrv-alert-alert) - display non-blocking alert message
    * [alertClose()](#slfsrv-alert-close) - close and alert message early
* [esc4html()](#slfsrv-esc4html) - escape string to be html-friendly
* [env](#slfsrv-env) - environment variables
    * [env.get()](#slfsrv-env-get) - get an environment variable
    * [env.set()](#slfsrv-env-set) - set an environment variable

------------------------------------------------------------

<a name="callbacks-vs-promises"></a>
### Callbacks versus Promises

Most of the functions provided by the `SLFSRV` object involve asynchronous calls; i.e., they
do not return immediately to your calling code. For all of these `SLFSRV` calls there are two
approaches for managing this: a callback version and a promises version.

Which version you use depends on your preference (some people love Promise, others aren't so sure).

##### callback version:

    SLFSRV.func( ...args...,
                 onsuccess<function(outcome){}>,
                 onerror<optional function(error<Error>){}> )

##### promise version:

    SLFSRV.func( ...args...)
          .then( resolve<function(outcome<object>){}>,
                 reject<optional function(error<Error>){}> )

To demonstrate the different ways these two methods are used, the following two code samples
each save the current directory, change to the directory "/foo", delete a file name "blah.txt",
then return to the original directory.

##### callback version:

    function cd_foo_del_blah_return() {
        SLFSRV.dir.getcwd(function(original_cwd){
            SLFSRV.dir.setcwd("/foo",function(){
                SLFSRV.file.remove("blah.txt",function(){
                    SLFSRV.dir.setcwd(original_cwd,function(){
                        alert("all went well!");
                    },function(err){
                        throw err;
                    });
                },function(err){
                    throw err;
                });
            },function(err){
                throw err;
            });
        },function(err){
            throw err;
        });
    }

##### promise version:

    function cd_foo_del_blah_return() {
        var original_cwd;
        SLFSRV.dir.getcwd().then(function(cwd){
            original_cwd = cwd;
            return SLFSRV.dir.setcwd("/foo");
        }).then(function(){
            return SLFSRV.file.remove("blah.txt");
        }).then(function(){
            return SLFSRV.dir.setcwd(original_cwd);
        }).then(function(){
            alert("all went well!");
        }).catch(function(err){
            throw err;
        });

------------------------------------------------------------

<a name="slfsrv-globals"></a>
### GLOBALS - global variables available to your javascript

    SLFSRV.OS // string for underlying OS (e.g. "darwin", "windows", "linux")

    SLFSRV.dir.SLASH // character used between directories. "\" for windows, else "/"

    SLFSRV.SECRET_KEY // random secret string generated when launching this local
                      // server instance - used to generate URLs that cannot be
                      // guessed by nefarious other systems on your network

    SLFSRV.PORT // what port the server is running on

    SLFSRV.ROOTPATH // root directory for the source html/js/css/etc... files

    SLFSRV.INITFILE // file to return is no filename specified in URL (e.g. "index.html")

    SLFSRV.SELF // full path to the executable that launched this program


------------------------------------------------------------

<a name="slfsrv-exec"></a>
### SLFSRV.exec() - execute a program on the system, and wait for result

##### callback version:

    SLFSRV.exec( options<object { program<string>, args<optional array>,
                                  input<optional string>, timeout<optional int> },
                 onsuccess<function(outcome<object>){}>,
                 onerror<optional function(error<Error>){}> )
    // options: object setting input values for exec
    //     .program: name of the command to execute
    //     .args: optional array of commands to pass to the program
    //     .input: optional string, if not "" then this will be passed to the program via stdin
    //     .timeout: optional int, how many seconds (approximately) to wait for call
    // onsuccess: called after program exits
    //     outcome: information about the program after it executes, containing these fields
    //         .exitStatus<int>: exit code returned by the program
    //         .stdout<string>: text program wrote to stdout
    //         .stderr<string>: text program wrote to stderr
    // onerror: called if there is an error launching program (or a timeout)
    //     error: Error object
    // return: undefined
    // ex: /* use shell to display full listing for the current directory */
    //     SLFSRV.dir.getcwd(function(cwd){
    //         var msg, program, args;
    //         msg = "Directory listing for: " + cwd + "\n\n";
    //
    //         function onsuccess(out) {
    //             msg += out.stdout;
    //             SLFSRV.alert(msg);
    //         }
    //
    //         function onerror(e) {
    //             msg += "ERROR GETTING DIRECTORY: " + e;
    //             SLFSRV.alert(msg);
    //         }
    //
    //         if ( SLFSRV.OS === "windows" ) {
    //             program = "cmd.exe";
    //             args = [ program, "/c", "dir" ];
    //         } else {
    //             program = "ls";
    //             args = [ program, "-l" ];
    //         }
    //
    //         SLFSRV.exec({program:program,args:args},onsuccess,onerror);
    //     });

##### promise version:

    SLFSRV.exec( options<object { program<string>, args<optional array>,
                                  input<optional string>, timeout<optional int> } )
          .then( resolve<function(outcome<object>){}>,
                 reject<optional function(error<Error>){}> )
    // options: object setting input values for exec
    //     .program: name of the command to execute
    //     .args: optional array of commands to pass to the program
    //     .input: optional string, if not "" then this will be passed to the program via stdin
    //     .timeout: optional int, how many seconds (approximately) to wait for call
    // resolve: called after program exits
    //     outcome: information about the program after it executes, containing these fields
    //         .exitStatus<int>: exit code returned by the program
    //         .stdout<string>: text program wrote to stdout
    //         .stderr<string>: text program wrote to stderr
    // reject: called if there is an error launching program (or a timeout)
    //     error: Error object
    // return: Promise
    // ex: /* use shell to display full listing for the current directory */
    //     var msg = "";
    //     SLFSRV.dir.getcwd().next(function(cwd){
    //         var program, args;
    //         msg += "Directory listing for: " + cwd + "\n\n";
    //
    //         if ( SLFSRV.OS === "windows" ) {
    //             program = "cmd.exe";
    //             args = [ program, "/c", "dir" ];
    //         } else {
    //             program = "ls";
    //             args = [ program, "-l" ];
    //         }
    //
    //         return SLFSRV.exec({program:program,args:args});
    //     }).then({function(out){
    //          msg += out.stdout;
    //     })["catch"](function(e){
    //         msg += "ERROR GETTING DIRECTORY: " + e;
    //     }).then(function(){
    //         SLFSRV.alert(msg);
    //     });

------------------------------------------------------------

<a name="slfsrv-tempdir"></a>
### SLFSRV.tempdir() - get a temporary directory, and unpack files to that directory

This functions creates a working directory where you can read/write/execute files. When the
program exits, that working directory will be deleted. If the "unpackDir" value is
provided, this will also copy all of the files from the specified `unpackDir` local directory into that
new working directory (subsequent calls with the same unpackDir will simply return the
name of the working directory).

This function is particularly useful in slfsrv files created with the "-bundle" or
"-compile-from" options, as it will copy files in the bundle to a working directory.

##### callback version:

    SLFSRV.tempdir( unpackDir<optional string>,
                    onsuccess<function(dirPath<string>){}>,
                    onerror<optional function(error<Error>){}> )
    // unpackDir - a label to five the temporary directory (this will not be the
    //             name of the tempdir). If there are any top-level directories within
    //             SLFSRV.ROOTPATH with this label, then on the first call everything
    //             within that directory will be copied to the new tempdir.
    //             If this is called again with the same unpackDir name then the same
    //             tempdir will be returned, with no changes. If this is not specified or
    //             is "" then a tempdir is returned with no copying.
    // onsuccess:
    //     dirPath: path to the temporary directory
    // onerror: called if there is an error
    //     error: Error object
    // return: undefined
    // ex: /* example to unpack Mac or Windows executable, then execute it */
    //     SLFSRV.tempdir( 'executables', function(dir) {
    //         var program, args;
    //         if ( SLFSRV.os === "windows" ) {
    //             program = "cmd.exe";
    //             args = [ program, "/c", dir + "\\myWindowsApp.exe" ];
    //         } else {
    //             program = "sh";
    //             args = [ program, "-c", dir + "/myMacApp" ];
    //         }
    //         SLFSRV.exec({program:program,args:args},function(ret) {
    //             if ( ret.stderr.length != 0 ) {
    //                 errFunc(ret.stderr);
    //             } else {
    //                 SLFSRV.alert("YEA! It worked. Output:\n\n" + ret.stdout);
    //             }
    //         },errFunc);
    //     }, errFunc );
    //     function errFunc(err) {
    //         SLFSRV.alert("ERROR: " + err);
    //     }

##### promise version:

    SLFSRV.tempdir( unpackDir<optional string> )
          .then( resolve<function(dirPath<string>){}>,
                 reject<optional function(error<Error>){}> )
    // unpackDir - a label to five the temporary directory (this will not be the
    //             name of the tempdir). If there are any top-level directories within
    //             SLFSRV.ROOTPATH with this label, then on the first call everything
    //             within that directory will be copied to the new tempdir.
    //             If this is called again with the same unpackDir name then the same
    //             tempdir will be returned, with no changes. If this is not specified or
    //             is "" then a tempdir is returned with no copying.
    // resolve: called after program exits
    //     dirPath: path to the temporary directory
    // reject: called if there is an error
    //     error: Error object
    // return: Promise
    // ex: /* example to unpack Mac or Windows executable, then execute it */
    //     SLFSRV.tempdir( 'executables' ).then( function(dir) {
    //         var program, args;
    //         if ( SLFSRV.os === "windows" ) {
    //             program = "cmd.exe";
    //             args = [ program, "/c", dir + "\\myWindowsApp.exe" ];
    //         } else {
    //             program = "sh";
    //             args = [ program, "-c", dir + "/myMacApp" ];
    //         }
    //         return SLFSRV.exec({program:program,args:args});
    //     }).then( function(ret) {
    //         if ( ret.stderr.length != 0 ) {
    //             errFunc(ret.stderr);
    //         } else {
    //             SLFSRV.alert("YEA! It worked. Output:\n\n" + ret.stdout);
    //         }
    //     })["catch"](function(err){
    //         SLFSRV.alert("ERROR: " + err)
    //     });

------------------------------------------------------------

<a name="slfsrv-store"></a>
### store - persistent data storage

<a name="slfsrv-store-get"></a>
### SLFSRV.store.get() - get value from persistent store

##### callback version:

    SLFSRV.store.get( key<string>, default<optional value if key does not exist>,
                      onsuccess<function(value){}>,
                      onerror<optional function(error<Error>){}> )
    // key: unique name for stored value
    // default: optional value if key does not exist, else would return undefined
    // onsuccess: return the stored value (or undefined if no stored value for that key)
    // onerror: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: undefined
    // note: if there is no value store for key, onsuccess will return undefined (and not onerror)
    // ex: SLFSRV.store.get( "favoriteColor", function(value) {
    //         if ( value === undefined ) {
    //             SLFSRV.alert( "no favorite color" );
    //         } else {
    //             SLFSRV.alert( "favorite color is " + value );
    //         }
    //     });
    //     SLFSRV.store.get( "counter", 0, function(value) {
    //         SLFSRV.alert( "counter called " + counter + " times" );
    //     });

##### promise version:

    SLFSRV.store.get( key<string>, default<optional value if key does not exist> )
                .then( resolve<function(value){}>,
                       reject<optional function(error<Error>){}> )
    // key: unique name for stored value
    // default: optional value if key does not exist, else would return undefined
    // resolve: return the stored value (or undefined if no stored value for that key)
    // reject: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: Promise
    // note: if there is no value store for key, resolve will return undefined (and not onerror)
    // ex: SLFSRV.store.get( "favoriteColor" ).next(function(value) {
    //         if ( value === undefined ) {
    //             SLFSRV.alert( "no favorite color" );
    //         } else {
    //             SLFSRV.alert( "favorite color is " + value );
    //         }
    //     });
    //     SLFSRV.store.get( "counter", 0 ).next(function(value) {
    //         SLFSRV.alert( "counter called " + counter + " times" );
    //     });

<a name="slfsrv-store-set"></a>
### SLFSRV.store.set() - set value in persistent store

##### callback version:

    SLFSRV.store.set( key<string>, value, onsuccess<function(){}>,
                      onerror<optional function(error<Error>){}> )
    // key: unique name for stored value
    // value: any javascript value that can be stored as json (or undefined to remove key)
    // onsuccess: called when value has been stored (or removed for undefined)
    // onerror: called if there is an error storing value
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.store.set( "favoritColor", "green", function(){} );

##### promise version:

    SLFSRV.store.set( key<string>, value )
                .then( resolve<function(){}>,
                       reject<optional function(error<Error>){}> )
    // key: unique name for stored value
    // value: any javascript value that can be stored as json (or undefined to remove key)
    // resolve: called when value has been stored (or removed for undefined)
    // reject: called if there is an error storing value
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.store.set( "favoritColor", "green" ).next( function(){} );

<a name="slfsrv-store-list"></a>
### SLFSRV.store.list() - list all stored keys

##### callback version:

    SLFSRV.store.list( onsuccess<function(keys<string array>){}>,
                       onerror<optional function(error<Error>){}> )
    // onsuccess: function called with an array of stored key names
    //     keys: array of stored key names
    // onerror: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.store.list( function(keys) {
    //         SLFSRV.alert( "stored keys are:\n" + keys.join("\n") );
    //     });

##### promise version:

    SLFSRV.store.list( )
                .then( resolve<function(keys<string array>){}>,
                       reject<optional function(error<Error>){}> )
    // resolve: function called with an array of stored key names
    //     keys: array of stored key names
    // reject: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.store.list().then( function(keys) {
    //         SLFSRV.alert( "stored keys are:\n" + keys.join("\n") );
    //     });

------------------------------------------------------------

<a name="slfsrv-file"></a>
### file - read/write access to the local file system

<a name="slfsrv-file-read"></a>
### SLFSRV.file.read() - read data from a file

##### callback version:

    SLFSRV.file.read( filename<string>, mode<optional string>,
                      onsuccess<function(contents<string>){}>,
                      onerror<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // mode - options "binary" or "text" - binary returns full file in javascript
    //        string of bytes - text returns unicode strings - defaults to "text"
    // contents - contents of the file (binary or text)
    // onerror: called if error reading file
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.file.read( "config.txt", "binary", function(contents) {
    //         SLFSRV.alert( "config.txt binary contents are:\n" + contents );
    //     }
    //     SLFSRV.file.read( "config.txt", function(contents) {
    //         SLFSRV.alert( "config.txt text contents are:\n" + contents );
    //     }

##### promise version:

    SLFSRV.file.read( filename<string>, mode<optional string> )
               .then( resolve<function (contents<string>){}>,
                      reject<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // mode - options "binary" or "text" - binary returns full file in javascript
    //        string of bytes - text returns unicode strings - defaults to "text"
    // resolve: function called with the contents of the file read
    //     contents - contents of the file (binary or text)
    // reject: called if error getting directory list
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.file.read( "config.txt", "binary" ).then( function(contents) {
    //         SLFSRV.alert( "config.txt binary contents are:\n" + contents );
    //     } );
    //     SLFSRV.file.read( "config.txt" ).then(function(contents) {
    //         SLFSRV.alert( "config.txt text contents are:\n" + contents );
    //     })["catch"](function(e){
    //         msg += "ERROR READING FILE: " + e;
    //     });

<a name="slfsrv-file-write"></a>
### SLFSRV.file.write() - write or append data to a file

##### callback version:

    SLFSRV.file.write( filename<string>, mode<optional string>, contents<string>,
                       onsuccess<optional function(){}>,
                       onerror<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // mode - options "binary" or "text", and "create" or "append"
    //        "binary" writes bytes - "text" writes unicode (utf-8) strings
    //        "create" creates (or overwrites) file - "append" appends to file (or creates)
    //        default: "text,create"
    // contents - contents of the file (binary or text)
    // onsuccess: called after file has been written
    // onerror: called if error writing file
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.file.write( "data.bin", "binary,append", "owl\x89\x43cat", function() {
    //         SLFSRV.alert( "binary data has been appended to data.bin" );
    //     }
    //     SLFSRV.file.write( "data1.txt", "555.333.2222", function() {
    //         SLFSRV.alert( "data1.txt was created, and data written in text mode" );
    //     }
    //     SLFSRV.file.write( "data2.txt", "text,create", "555.333.2222", function() {
    //         SLFSRV.alert( "data2.txt was created, and data written in text mode" );
    //     }

##### promise version:

    SLFSRV.file.write( filename<string>, mode<optional string>, contents<string> )
               .then( resolve<function (){}>,
                      reject<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // mode - options "binary" or "text", and "create" or "append"
    //        "binary" writes bytes - "text" writes unicode (utf-8) strings
    //        "create" creates (or overwrites) file - "append" appends to file (or creates)
    //        default: "text,create"
    // contents - contents of the file (binary or text)
    // resolve: function called after file has been written
    // reject: called if error getting directory list
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.file.write( "data.bin", "binary,append", "owl\x89\x43cat"  ).then( function() {
    //         SLFSRV.alert( "binary data has been appended to data.bin" );
    //     } );
    //     SLFSRV.file.write( "data1.txt", "555.333.2222" ).then(function() {
    //         SLFSRV.alert( "data1.txt was created, and data written in text mode" );
    //     });
    //     SLFSRV.file.write( "data2.txt", "text,create", "555.333.2222" ).then(function() {
    //         SLFSRV.alert( "data2.txt was created, and data written in text mode" );
    //     })["catch"](function(e){
    //         msg += "ERROR WRITING FILE: " + e;
    //     });

<a name="slfsrv-file-exists"></a>
### SLFSRV.file.exists() - does this file exist

##### callback version:

    SLFSRV.file.exists( filename<string>,
                        onsuccess<function(exists<bool>){}>,
                        onerror<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // exists: true if file exists (and not a directory), else false
    // onerror: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.file.exists( "config.txt", function(exists) {
    //         SLFSRV.alert( "file " + (exists ? "DOES" : "DOES NOT") + " exist");
    //     }

##### promise version:

    SLFSRV.file.exists( filename<string> )
               .then( resolve<function (exists<bool>){}>,
                      reject<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // resolve: function called with the determination if file exists
    //     exists: true if file exists (and not a directory), else false
    // reject: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.file.exists( "config.txt" ).then( function(exists) {
    //         SLFSRV.alert( "file " + (exists ? "DOES" : "DOES NOT") + " exist");
    //     } );

<a name="slfsrv-file-remove"></a>
### SLFSRV.file.remove() - delete a file

##### callback version:

    SLFSRV.file.remove( filename<string>,
                        onsuccess<optional function(){}>,
                        onerror<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // onsuccess: called if all is well
    // onerror: called if there is an error removing the file
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.file.remove( "tempdata.txt", function(){} );

##### promise version:

    SLFSRV.file.remove( filename<string> )
               .then( resolve<optional function(){}>,
                      reject<optional function(error<Error>){}> )
    // filename - path to file - if not a full file path then relative to cwd
    // resolve: called if all is well
    // reject: called if there is an error removing the file
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.file.remove( "tempdata.txt" ).next( function(){} );

------------------------------------------------------------

<a name="slfsrv-alert"></a>
### alert - non-blocking alert message

<a name="slfsrv-alert-alert"></a>
### SLFSRV.alert() - display non-blocking alert message

##### callback version:

    SLFSRV.alert( message<string>, onsuccess<function(){}> )
    // message: message to display in popup text box
    // onsuccess: function called after the user presses "OK" to dismiss the alert
    // return: an object that can be used with SLFSRV.alertClose()
    // note: this should be used instead of standard browser alert() because it is non-
    //       blocking and so allows the client to continue to ping the server with
    //       keep-alive calls
    // ex: SLFSRV.alert( "hello, squirrels", function(){} );
    //     SLFSRV.alert( "message 1", function(){
    //         SLFSRV.alert( "message 2: OK was pressed on message 1", function(){
    //             SLFSRV.alert( "message 3: OK was pressed on message 2", function(){
    //                 SLFSRV.alert( "message 4: OK was pressed on message 3", function(){} )
    //             });
    //         });
    //     });

##### promise version:

    SLFSRV.alert( message<string> )
          .then( resolve<function(){}> )
    // message: message to display in popup text box
    // resolve: function called after the user presses "OK" to dismiss the alert
    // return: Promise (which can be used with SLFSRV.alertClose())
    // note: this should be used instead of standard browser alert() because it is non-
    //       blocking and so allows the client to continue to ping the server with
    //       keep-alive calls
    // ex: SLFSRV.alert( "hello, squirrels" );
    //     SLFSRV.alert( "message 1" ).then( function(){
    //         return SLFSRV.alert( "message 2: OK was pressed on message 1" );
    //     }).then( function() {
    //         return SLFSRV.alert( "message 3: OK was pressed on message 2" );
    //     }).then( function() {
    //         SLFSRV.alert( "message 4: OK was pressed on message 3" );
    //     });

<a name="slfsrv-alert-close"></a>
### SLFSRV.alertClose() - close an alert message early

    SLFSRV.closeAlert( alert<object> )
    // alert: object that was returned by SLFSRV.alert()
    // return: undefined
    // note: if the alert window has not yet shown, this will prevent it from showing. if the
    //       alert has already been dismissed then this function will do nothing
    // ex: var a = SLFSRV.alert( "waiting..." );
    //     ... stuff is happening ...
    //     SLFSRV.closeAlert(a);

------------------------------------------------------------

<a name="slfsrv-esc4html"></a>
### SLFSRV.esc4html() - replace html-encoding characters in strings

    SLFSRV.esc4html( str<string> )
    // str: relative path to change current directory to
    // return: string that is safe to insert into html
    // ex: SLFSRV.esc4html( "14 < 15" ); // returns "14 &lt; 15"

------------------------------------------------------------

<a name="slfsrv-dir"></a>
### dir - working with directories (i.e. folders) in local file system

<a name="slfsrv-dir-list"></a>
### SLFSRV.dir.list() - list all files and directories within this directory

##### callback version:

    SLFSRV.dir.list( dirname<optional string>, callback<function (list<array>){}>,
                     onerror<optional function(error<Error>){}> )
    // dirname: name of directory to list, relaive to current working directory if root not
    //          specified; if this parameter not supplied then will return from current working directory
    // callback: called to return directory listing
    //     list: array of file listing, where each entry has these values
    //         .name<string> - name of the file or subdirectory
    //         .size<int> - size of the file
    //         .dir<bool> - true if this is a subdirectory name
    //         .perm<int> - permission bits for *nix systems
    //         .mod<Date> - file modification time
    // onerror: called if error getting directory list
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.dir.list( function(list) {
    //         for ( var i=0; i < list.length; i++ ) {
    //             if ( !list[file].dir ) {
    //                 SLFSRV.alert( "file " + list[file].name );
    //             }
    //         }
    //     });

##### promise version:

    SLFSRV.dir.list( dirname<optional string> )
              .then( resolve<function (list<array>){}>,
                     reject<optional function(error<Error>){}> )
    // dirname: name of directory to list, relaive to current working directory if root not
    //          specified; if this parameter not supplied then will return from current working directory
    // reject: called to return directory listing
    //     list: array of file listing, where each entry has these values
    //         .name<string> - name of the file or subdirectory
    //         .size<int> - size of the file
    //         .dir<bool> - true if this is a subdirectory name
    //         .perm<int> - permission bits for *nix systems
    //         .mod<Date> - file modification time
    // reject: called if error getting directory list
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.dir.list().next( function(list) {
    //         for ( var i=0; i < list.length; i++ ) {
    //             if ( !list[file].dir ) {
    //                 SLFSRV.alert( "file " + list[file].name );
    //             }
    //         }
    //     });

<a name="slfsrv-dir-getcwd"></a>
### SLFSRV.dir.getcwd() - get current working directory

##### callback version:

    SLFSRV.dir.getcwd( onsuccess<function(cwd<string>){}>,
                       onerror<optional function(error<Error>){}>  )
    // onsuccess: function called with the name of the current directory
    //     cwd: current working directory
    // onerror: called if there is an error, which is very much not expected
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.dir.getcwd( function(cwd) {
    //         SLFSRV.alert( "current directory is " + cwd );
    //     });

##### promise version:

    SLFSRV.dir.getcwd()
              .then( resolve<function(cwd<string>){}>,
                     reject<optional function(error<Error>){}> )
    // resolve: function called with the name of the current directory
    //     cwd: current working directory
    // reject: called if there is an error, which is very much not expected
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.dir.getcwd().then(function(cwd) {
    //         SLFSRV.alert( "current directory is " + cwd );
    //     });

<a name="slfsrv-dir-setcwd"></a>
### SLFSRV.dir.setcwd() - set current working directory

##### callback version:

    SLFSRV.dir.setcwd( cwd<string>, onsuccess<function(){}>,
                       onerror<optional function(error<Error>){}> )
    // cwd: relative path to change current directory to
    // onsuccess: called when directory has changed
    // onerror: called if there is an error changing path
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.dir.setcwd( "../parent-dir-name", function(){} );

##### promise version:

    SLFSRV.dir.setcwd( cwd<string> )
              .then( resolve<optional function(){}>,
                     reject<optional function(error<Error>){}> )
    // cwd: relative path to change current directory to
    // resolve: function called after current working directory has changed
    // reject: called if there is an error changing path
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.dir.setcwd( "../parent-dir-name" ).next(function(){});

<a name="slfsrv-dir-exists"></a>
### SLFSRV.dir.exists() - does this directory exist

##### callback version:

    SLFSRV.dir.exists( dirname<string>,
                       onsuccess<function(exists<bool>){}>,
                       onerror<optional function(error<Error>){}> )
    // dirname - path to direcotry - if not a full path then relative to cwd
    // exists: true if directory exists, else false
    // onerror: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.dir.exists( "pebbled", function(exists) {
    //         SLFSRV.alert( "directory " + (exists ? "DOES" : "DOES NOT") + " exist");
    //     }

##### promise version:

    SLFSRV.dir.exists( dirname<string> )
              .then( resolve<function (exists<bool>){}>,
                     reject<optional function(error<Error>){}> )
    // dirname - path to directory - if not a full path then relative to cwd
    // resolve: function called with the determination if directory exists
    //     exists: true if directory exists, else false
    // reject: called if there is an error (rare that this would happen)
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.dir.exists( "pebbles" ).then( function(exists) {
    //         SLFSRV.alert( "directory " + (exists ? "DOES" : "DOES NOT") + " exist");
    //     } );

------------------------------------------------------------

<a name="slfsrv-end"></a>
### env - environment variables

<a name="slfsrv-env-get"></a>
### SLFSRV.env.get() - return an environment variable

##### callback version:

    SLFSRV.env.get( key<string>, onsuccess<function (value<string>){}>,
                                 onerror<optional function(error<Error>){}> )
    // key: name of the environment variable to retrieve
    // onsuccess: called to return environment value
    //     value: value of the variable, may be "" if no such variable
    // onerror: called if there is an error, which is very much not expected
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.env.get( "PATH", function(path) {
    //         SLFSRV.alert( "current path is " + path );
    //     });

##### promise version:

    SLFSRV.env.get( key<string> )
              .then( resolve<function (value<string>){}>,
                     reject<optional function(error<Error>){}> )
    // key: name of the environment variable to retrieve
    // resolve: called to return environment value
    //     value: value of the variable, may be "" if no such variable
    // reject: called if there is an error, which is very much not expected
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.env.get( "PATH" ).next( function(path) {
    //         SLFSRV.alert( "current path is " + path );
    //     });

<a name="slfsrv-env-set"></a>
### SLFSRV.env.set() - set an environment variable

##### callback version:

    SLFSRV.env.set( key<string>, value<string>,
                    onsuccess<optional function(){}>,
                    onerror<optional function(error<Error>){}> )
    // key: environment variable to set
    // value: value of that variable
    // onsuccess: called if all is well
    // onerror: called if there is an error changing path
    //     error: Error object
    // return: undefined
    // ex: SLFSRV.env.set( "FEELING", "happy", function(){} );

##### promise version:

    SLFSRV.env.set( key<string>, value<string> )
              .then( resolve<optional function(){}>,
                     reject<optional function(error<Error>){}> )
    // key: environment variable to set
    // value: value of that variable
    // resolve: called if all is well
    // reject: called if there is an error changing path
    //     error: Error object
    // return: Promise
    // ex: SLFSRV.env.set( "FEELING", "happy" ).next( function(){} );

------------------------------------------------------------

[&#x25C0; back to README.cmd](../README.md)
