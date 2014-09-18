/*jslint white:false plusplus:false browser:true nomen:false */
/*globals $, SLFSRV, alert, window*/

var gMyPath = SLFSRV.ROOTPATH + ("/" + SLFSRV.INITFILE).split('/index.html')[0].replace(/\//g,SLFSRV.dir.SLASH);

function output(msg) { // write message on console
    var html = '<p>';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    $('#output').append(html);
    window.scrollTo(0,document.body.scrollHeight);
}

function error(msg) {
    var html = '<p style="color:red;">';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    $('#output').append(html);
}

function hr() {
    $('#output').append("<hr/>");
}

var text_sample_binary = "English: Text sample" + "\n" +
                         "French: Il n\x27y a pas de fum\xC3\xA9e sans feu." + "\n" +
                         "Japanese: \xE3\x81\x93\xE3\x82\x8C\xE3\x81\xAF\xE3\x83\x86\xE3\x82\xB9\xE3\x83\x88\xE3\x81\xA7\xE3\x81\x99\xE3\x80\x82" + "\n" +
                         "Korean: \xEC\x9D\xB4\xEA\xB2\x83\xEC\x9D\x80 \xED\x85\x8C\xEC\x8A\xA4\xED\x8A\xB8\xEC\x9E\x85\xEB\x8B\x88\xEB\x8B\xA4." + "\n" +
                         "the end";

var text_sample_text = "English: Text sample" + "\n" +
                       "French: Il n'y a pas de fumée sans feu." + "\n" +
                       "Japanese: これはテストです。" + "\n" +
                       "Korean: 이것은 테스트입니다." + "\n" +
                       "the end";

var text_sample_encoded = "English: Text sample" + "\n" +
                          "French: Il n\'y a pas de fum\u00E9e sans feu." + "\n" +
                          "Japanese: \u3053\u308C\u306F\u30C6\u30B9\u30C8\u3067\u3059\u3002" + "\n" +
                          "Korean: \uC774\uAC83\uC740 \uD14C\uC2A4\uD2B8\uC785\uB2C8\uB2E4." + "\n" +
                          "the end";

function escStringToJSSource(str) {
    var i, c, outc, ret = "";

    for ( i = 0; i < str.length; i++ ) {
        c = str.charCodeAt(i);
        outc = c;
        if ( 256 <= c ) {
            outc = c.toString(16);
            while ( outc.length < 4 ) {
                outc = "0" + outc;
            }
            outc = "\\u" + outc.toUpperCase();
        } else if ( ( 32 === c ) ||
                    ( 48 <= outc && outc <= 57 ) ||
                    ( 65 <= outc && outc <= 90 ) ||
                    ( 97 <= outc && outc <= 122 ) ||
                    ( 58 <= outc && outc <= 64 ) ||
                    ( 36 <= outc && outc <= 46 ) ) {
            outc = String.fromCharCode(c);
        } else if ( c === 13 ) {
            outc = "\\r";
        } else if ( c === 10 ) {
            outc = "\\n\n";
        } else if ( c === 9 ) {
            outc = "\\t";
        } else {
            outc = c.toString(16);
            while ( outc.length < 2 ) {
                outc = "0" + outc;
            }
            outc = "\\x" + outc.toUpperCase();
        }
        ret += outc;
    }
    return ret;
}


function string_arrays_match(expected,got) { // call error() and return false, or return true
    var i, same = false, msg;
    expected.sort();
    got.sort();

    if ( expected.length === got.length ) {
        same = true;
        for ( i = 0; i < expected.length; i++ ) {
            if ( expected[i] !== got[i] ) {
                same = false;
            }
        }
    }
    if ( !same ) {
        msg = "results do not match\nEXPECTED:";
        for ( i = 0; i < expected.length; i++ ) {
            msg += "\n. " + expected[i];
        }
        msg += "\nGOT:";
        for ( i = 0; i < got.length; i++ ) {
            msg += "\n- " + got[i];
        }
        error(msg);
    }
    return same;
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
        name: 'test assumptions',
        test: function() {
            if ( text_sample_text === text_sample_encoded ) {
                output("done");
                next_test();
            } else {
                error("text_sample_text !== text_sample_encoded");
            }
        }
    },{
        name: 'save the current working directory',
        test: function() {
            SLFSRV.dir.getcwd(function(cwd){
                starting_cwd = cwd;
                output(starting_cwd);
                output("done");
                next_test();
            },function(err){
                error(err);
            });
        }
    },{
        name: 'change to directory of this test, which is: ' + gMyPath,
        test: function() {
            output("CWD = " + gMyPath);
            SLFSRV.dir.setcwd(gMyPath,function(){
                SLFSRV.dir.getcwd(function(cwd){
                    if ( cwd !== gMyPath ) {
                        error("changed to unexpected directory: " + cwd);
                    } else {
                        output("done");
                    }
                    next_test();
                },function(err){
                    error(err);
                });
            },function(err){
                error(err);
            });
        }
    },{
        name: 'get directory listing of cwd:',
        test: function() {
            var ignore_files = " .DS_Store text-sample-output.txt ", foundDirList = [];

            function dirList(rootname,onComplete) {
                SLFSRV.dir.list(rootname,function(list){
                    function next_from_list() {
                        var f, ignore;
                        if ( list.length === 0 ) {
                            onComplete();
                        } else {
                            f = list.shift();
                            ignore = ( -1 !== ignore_files.indexOf(' '+f.name+' ') );
                            output("name: " + rootname + f.name + "\n" +
                                   "size: " + f.size + "\n" +
                                   "dir: " + f.dir + "\n" +
                                   "perm: 0" + f.perm.toString(8) + "\n" +
                                   "mod: " + f.mod  +
                                    ( ignore ? "\nIGNORE" : "" ));
                            if ( !ignore ) {
                                foundDirList.push(rootname + f.name);
                            }
                            if ( f.dir ) {
                                dirList(f.name + SLFSRV.dir.SLASH,next_from_list);
                            } else {
                                next_from_list();
                            }
                        }
                    }

                    next_from_list();
                },function(err){
                    error(err);
                });
            }

            function compare_found_with_expected() {
                var expectedDirList = [
                    "index.html",
                    "javascript",
                    "javascript" + SLFSRV.dir.SLASH + "file-and-dir-test.js",
                    "datafiles",
                    "datafiles" + SLFSRV.dir.SLASH + "text-sample.txt",
                    "datafiles" + SLFSRV.dir.SLASH + "உரை.txt" ];

                if ( string_arrays_match(expectedDirList,foundDirList) ) {
                    output("done");
                    next_test();
                }
            }

            dirList("",compare_found_with_expected);
        }
    },{
        name: 'test reading file in text and binary modes:',
        test: function() {
            var srcfile, modes, mode;
            srcfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample.txt";
            modes = [ "text", undefined/*default text*/, "binary" ];

            function goodRead(data) {
                var expect;

                if ( mode === "binary" ) {
                    expect = text_sample_binary;
                    if ( expect !== data ) {
                        error('EXPECT binary:\n\n' + escStringToJSSource(expect) + "\n\nBUT GOT:\n\n" + escStringToJSSource(data));
                    } else {
                        output("OK");
                        setTimeout(test,0);
                    }
                } else {
                    // text mode
                    expect = text_sample_text;

                    if ( expect !== data ) {
                        error('EXPECT1:\n\n' + expect + "\n\nBUT GOT:\n\n" + data);
                    } else {
                        expect = text_sample_encoded;

                        if ( expect !== data ) {
                            error('EXPECT2:\n\n' + expect + "\n\nBUT GOT:\n\n" + data);
                        } else {
                            output("OK");
                            setTimeout(test,0);
                        }
                    }
                }
            }

            function badRead(err) {
                error(err);
            }

            function test() {
                if ( modes.length === 0 ) {
                    output("done");
                    next_test();
                } else {
                    mode = modes.shift();
                    output('test reading "' + srcfile + '" in ' + mode + ' mode...');
                    if ( mode === undefined ) {
                        SLFSRV.file.read(srcfile,goodRead,badRead);
                    } else {
                        SLFSRV.file.read(srcfile,mode,goodRead,badRead);
                    }
                }
            }

            test();
        }
    },{
        name: 'test read failure if file does not exist:',
        test: function() {
            var srcfile;
            srcfile = "datafiles" + SLFSRV.dir.SLASH + "no-such-text-sample.txt";

            function goodRead(text) {
                error("dang. should have failed to read");
            }

            function badRead(err) {
                output('Ignore this error, because it is expected "' + String(err) + '"');
                output("done");
                next_test();
            }

            SLFSRV.file.read(srcfile,goodRead,badRead);
        }
    },{
        name: 'test reading file in binary mode:',
        test: function() {
            output("done");
            next_test();
        }
    },{
        name: 'test file.exists() for file that exists:',
        test: function() {
            var srcfile;
            srcfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample.txt";

            function good(fileExists) {
                if ( fileExists ) {
                    output("done");
                    next_test();
                } else {
                    error("Error: " + srcfile + " does not exist");
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.exists(srcfile,good,bad);
        }
    },{
        name: 'test file.exists() for file that does not exist:',
        test: function() {
            var srcfile;
            srcfile = "datafiles" + SLFSRV.dir.SLASH + "no-such-text-sample.txt";

            function good(fileExists) {
                if ( fileExists ) {
                    error("Error: " + srcfile + " exists but should not");
                } else {
                    output("done");
                    next_test();
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.exists(srcfile,good,bad);
        }
    },{
        name: 'test file.exists() for directory that should fail as a file',
        test: function() {
            var srcfile;
            srcfile = "datafiles";

            function good(fileExists) {
                if ( fileExists ) {
                    error("Error: " + srcfile + " exists as a file but should not");
                } else {
                    output("done");
                    next_test();
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.exists(srcfile,good,bad);
        }
    },{
        name: 'test dir.exists() for dir that exists:',
        test: function() {
            var srcdir;
            srcdir = "datafiles";

            function good(dirExists) {
                if ( dirExists ) {
                    output("done");
                    next_test();
                } else {
                    error("Error: " + srcdir + " does not exist");
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.dir.exists(srcdir,good,bad);
        }
    },{
        name: 'test dir.exists() for dir that does not exist:',
        test: function() {
            var srcdir;
            srcdir = "no-such-datafiles";

            function good(dirExists) {
                if ( dirExists ) {
                    error("Error: " + srcdir + " exists but should not");
                } else {
                    output("done");
                    next_test();
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.dir.exists(srcdir,good,bad);
        }
    },{
        name: 'test dir.exists() for file that should fail as a dir',
        test: function() {
            var srcdir;
            srcdir = "datafiles" + SLFSRV.dir.SLASH + "text-sample.txt";

            function good(dirExists) {
                if ( dirExists ) {
                    error("Error: " + srcdir + " exists as a directory but should not");
                } else {
                    output("done");
                    next_test();
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.dir.exists(srcdir,good,bad);
        }
    },{
        name: 'test writing file in combinations of text/binary and create/append modes',
        test: function() {
            var dstfile, modes, mode, isBinary, isAppend, firstHalfOfAppend;
            dstfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample-output.txt";
            modes = [ "text", "text,create", "text,append", "create", "append", undefined/*default text,create*/,
                      "binary", "binary,create", "binary,append" ];

            function badWrite(err) {
                error(err);
            }

            function goodWrite(data) {
                var expect;

                if ( firstHalfOfAppend ) {
                    // for append mode do it all twice
                    firstHalfOfAppend = false;
                    SLFSRV.file.write(dstfile,mode,isBinary ? text_sample_binary : text_sample_text,
                                      goodWrite,badWrite);
                } else {

                    // read in what we just wrote to make sure it's what we expect
                    SLFSRV.file.read(dstfile,isBinary ? "binary" : "text").then( function(data){

                        if ( isBinary ) {
                            expect = text_sample_binary;
                            if ( isAppend ) {
                                expect += text_sample_binary;
                            }
                            if ( expect !== data ) {
                                error('EXPECT binary:\n\n' + escStringToJSSource(expect) + "\n\nBUT GOT:\n\n" + escStringToJSSource(data));
                            } else {
                                output("OK");
                                setTimeout(test,0);
                            }
                        } else {
                            // text mode
                            expect = text_sample_text;
                            if ( isAppend ) {
                                expect += text_sample_text;
                            }

                            if ( expect !== data ) {
                                error('EXPECT1:\n\n' + expect + "\n\nBUT GOT:\n\n" + data);
                            } else {
                                output("OK");
                                setTimeout(test,0);
                            }
                        }

                    },badWrite);
                }
            }

            function test() {
                if ( modes.length === 0 ) {
                    output("done");
                    next_test();
                } else {
                    mode = modes.shift();
                    output('test writing "' + dstfile + '" in ' + mode + ' mode...');
                    if ( mode === undefined ) {
                        isBinary = false;
                        isAppend = false;
                        firstHalfOfAppend = false;
                        SLFSRV.file.write(dstfile,text_sample_text,goodWrite,badWrite);
                    } else {
                        isBinary = ( -1 !== mode.indexOf("binary") );
                        isAppend = ( -1 !== mode.indexOf("append") );
                        firstHalfOfAppend = isAppend;
                        SLFSRV.file.write(dstfile,
                                          mode.replace("append","xxx"),
                                          isBinary ? text_sample_binary : text_sample_text,
                                          goodWrite,badWrite);
                    }
                }
            }

            test();
        }
    },{
        name: 'test file.exists() for file we just created in file.write test',
        test: function() {
            var dstfile;
            dstfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample-output.txt";

            function good(fileExists) {
                if ( fileExists ) {
                    output("done");
                    next_test();
                } else {
                    error("Error: " + dstfile + " does not exist but should");
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.exists(dstfile,good,bad);
        }
    },{
        name: 'test file.remove() for file we just created in file.write test',
        test: function() {
            var dstfile;
            dstfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample-output.txt";

            function good() {
                output("done");
                next_test();
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.remove(dstfile,good,bad);
        }
    },{
        name: 'test that that file is now gone',
        test: function() {
            var dstfile;
            dstfile = "datafiles" + SLFSRV.dir.SLASH + "text-sample-output.txt";

            function good(fileExists) {
                if ( fileExists ) {
                    error("Error: " + dstfile + " was note removed");
                } else {
                    output("done");
                    next_test();
                }
            }

            function bad(err) {
                error(err);
            }

            SLFSRV.file.exists(dstfile,good,bad);
        }
    },{
        name: 'test that cannot use file.remove() to remove a directory',
        test: function() {
            var dstfile;
            dstfile = "datafiles";

            function good() {
                error("should not have been able to use file.remove() to remove directory " + dstfile);
            }

            function bad(err) {
                output('Ignore this error, because it is expected "' + String(err) + '"');
                output("done");
                next_test();
            }

            SLFSRV.file.remove(dstfile,good,bad);
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




$( document ).ready(function() {
    test();
});