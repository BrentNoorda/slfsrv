/*jslint white:false plusplus:false browser:true nomen:false */
/*globals SLFSRV, alert, window*/

var gMyPath = SLFSRV.ROOTPATH + ("/" + SLFSRV.INITFILE).split('/index.html')[0].replace(/\//g,SLFSRV.dir.SLASH);

function output(msg) { // write message on console
    var html = '<p>';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    document.getElementById('output').innerHTML += html;
    window.scrollTo(0,document.body.scrollHeight);
}

function error(msg) {
    var html = '<p style="color:red;">';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    document.getElementById('output').innerHTML += html;
}

function hr() {
    document.getElementById('output').innerHTML += "<hr/>";
}

function test() {
    var tests, test_idx, starting_cwd;

    function next_test() {
        test_idx++;
        hr();
        if ( tests.length <= test_idx ) {
            output("FINISHED!");
        } else {
            output("TEST " + test_idx + ": " + tests[test_idx].name);
            setTimeout(function(){tests[test_idx].test();},0);
        }
    }

    tests = [{
        name: 'simply create a blank, unnamed tempdir, and it better be empty',
        test: function() {
            var tempDir;
            SLFSRV.tempdir().then( function(dirPath) {
                output("tempdir() dirpath = " + dirPath);

                // verify that the directory exists, and is not empty
                tempDir = dirPath;
                return SLFSRV.dir.exists( dirPath );
            }).then( function( exists ) {
                output("confirmed that " + tempDir + " exists");
                return SLFSRV.dir.list( tempDir );
            }).then( function(list) {
                if ( list.length === 0 ) {
                    output("confirmed that " + tempDir + " contains no files");
                    // if called again it better return the same dirname
                    return SLFSRV.tempdir();
                } else {
                    throw new Error("tempdir \"" + tempDir + "\" should be empty but contains " + list.length + " files");
                }
            }).then( function(dirPath) {
                if ( dirPath === tempDir ) {
                    output("confirmed that same tempDir call returns same dir the second time");
                    output("done");
                    next_test();
                } else {
                    throw new Error("tempdir \"" + tempDir + "\" does not match second call \"" + dirPath + "\"");
                }
            })["catch"]( function(err) {
                error(err);
            });
        }
    },{
        name: 'simply create a labeled blank, but it should still be empty',
        test: function() {
            var tempDir;
            SLFSRV.tempdir("unmatched-name").then( function(dirPath) {
                output("tempdir() dirpath = " + dirPath);

                // verify that the directory exists, and is not empty
                tempDir = dirPath;
                return SLFSRV.dir.exists( dirPath );
            }).then( function( exists ) {
                output("confirmed that " + tempDir + " exists");
                return SLFSRV.dir.list( tempDir );
            }).then( function(list) {
                if ( list.length === 0 ) {
                    output("confirmed that " + tempDir + " contains no files");
                    // if called again it better return the same dirname
                    return SLFSRV.tempdir("unmatched-name");
                } else {
                    throw new Error("tempdir \"" + tempDir + "\" should be empty but contains " + list.length + " files");
                }
            }).then( function(dirPath) {
                if ( dirPath === tempDir ) {
                    output("confirmed that same tempDir call returns same dir the second time");
                    output("done");
                    next_test();
                } else {
                    throw new Error("tempdir \"" + tempDir + "\" does not match second call \"" + dirPath + "\"");
                }
            })["catch"]( function(err) {
                error(err);
            });
        }
    },{
        name: 'create tempdir from the "shell-scripts" and make sure the correct files are there',
        test: function() {
            var tempDir;
            SLFSRV.tempdir("shell-scripts").then( function(dirPath) {
                output("tempdir() dirpath = " + dirPath);

                // verify that the directory exists, and is not empty
                tempDir = dirPath;
                return SLFSRV.dir.exists( dirPath );
            }).then( function( exists ) {
                output("confirmed that " + tempDir + " exists");
                return SLFSRV.dir.list( tempDir );
            }).then( function(list) {
                // verify that this contains only the files "very-simple-bash.sh" and "very-simple-batch.bat"
                var i, name, fileCount = 0;
                for ( i = 0; i < list.length; i++ ) {
                    name = list[i].name;
                    if ( (name === "very-simple-bash.sh") || (name === "very-simple-batch.bat") || (name === "testdir") ) {
                        fileCount++;
                    } else if ( name === ".DS_Store" ) {
                        // ignore these files
                    } else {
                        throw new Error("Found file \"" + name + "\" but expect only \"very-simple-bash.sh\" or \"very-simple-batch.bat\" or \"testdir\"");
                    }
                }
                if ( fileCount !== 3 ) {
                    throw new Error("Expected to find 3 matching files but found " + fileCount );
                } else {
                    output("confirmed that " + tempDir + " contains only the files we expect");
                    // if called again it better return the same dirname
                    return SLFSRV.dir.list( tempDir + SLFSRV.dir.SLASH + "testdir" );
                }
            }).then( function(list) {
                // verify that this contains only the file "unusedfile.txt"
                var i, name, fileCount = 0;
                for ( i = 0; i < list.length; i++ ) {
                    name = list[i].name;
                    if ( name === "unusedfile.txt" ) {
                        fileCount++;
                    } else if ( name === ".DS_Store" ) {
                        // ignore these files
                    } else {
                        throw new Error("Found file \"" + name + "\" but expect only \"unusedfile.txt\"");
                    }
                }
                if ( fileCount !== 1 ) {
                    throw new Error("Expected to find 1 matching files but found " + fileCount );
                } else {
                    output("confirmed that " + (tempDir + SLFSRV.dir.SLASH + "testdir") + " contains only the files we expect");
                    // if called again it better return the same dirname
                    return SLFSRV.tempdir("shell-scripts");
                }
            }).then( function(dirPath) {
                if ( dirPath === tempDir ) {
                    output("confirmed that same tempDir call returns same dir the second time");
                    output("done");
                    next_test();
                } else {
                    throw new Error("tempdir \"" + tempDir + "\" does not match second call \"" + dirPath + "\"");
                }
            })["catch"]( function(err) {
                error(err);
            });
        }
    },{
        name: 'verify that we can execute an unpacked shell script correctly',
        test: function() {
            SLFSRV.tempdir("shell-scripts").then( function(dir) {
                output("tempdir() dirpath = " + dir);

                // execute the os-appropriate shell script
                var program, args;
                if ( SLFSRV.os === "windows" ) {
                    program = "cmd.exe";
                    args = [ program, "/c", dir + "\\very-simple-batch.bat" ];
                } else {
                    program = "sh";
                    args = [ program, "-c", dir + "/very-simple-bash.sh" ];
                }
                return SLFSRV.exec({program:program,args:args});
           }).then( function(ret) {
                if ( ret.stderr.length !== 0 ) {
                    throw new Error("Error running shell-script: " + ret.stderr);
                } else {
                    var stdout, expect;
                    stdout = ret.stdout.replace(/^\s+|\s+$/g,''); // remove newline and spaces from beginning and end
                    expect = "if you see this, then tempdir worked";
                    if ( stdout !== expect ) {
                        throw new Error("shell script returned\n\n\"" + stdout + "\"\n\nbut we expected\n\n\"" + expect + "\"");
                    } else {
                        output("confirmed that shell script returned \"" + expect + "\"");
                        output("done");
                        next_test();
                    }
                }
            })["catch"]( function(err) {
                error(err);
            });
        }
    }];

    test_idx = -1;
    next_test();
}

window.onerror = function(msg, url, line, col, error) {
   // Note that col & error are new to the HTML 5 spec and may not be
   // supported in every browser.  It worked for me in Chrome.
   var extra, suppressErrorAlert;
   extra = !col ? '' : '\ncolumn: ' + col;
   extra += !error ? '' : '\nerror: ' + error;

   // You can view the information in an alert to see things working like this:
   alert("Error: " + msg + "\nurl: " + url + "\nline: " + line + extra);

   // TODO: Report this error via ajax so you can keep track
   //       of what pages have JS issues

   suppressErrorAlert = true;
   // If you return true, then error alerts (like in older versions of
   // Internet Explorer) will be suppressed.
   return suppressErrorAlert;
};

test();
