/*jslint white:false plusplus:false browser:true nomen:false */
/*globals SLFSRV, window, $*/

function example(path,initfile) {
    // build and execute command line to run this executable
    var pathsep, fullpath, msg, program, args;

    pathsep = ( SLFSRV.OS==="windows" ) ? "\\" : "/";
    path = "examples" + pathsep + path.replace(/\//g,pathsep);

    msg = 'Press "OK" to launch the following command on this computer:\n\n   ' +
          (( SLFSRV.OS==="windows" ) ? ">" : "#") +
          " slfsrv -verbose " + path;
    if ( initfile ) {
        msg += ( " " + initfile );
    }
    SLFSRV.alert(msg,function() {

        var alertBox = SLFSRV.alert("Waiting for you to close the new example "+
            "window (which causes the example to exit).");

        program = SLFSRV.SELF;
        args = [ program, "-verbose", SLFSRV.ROOTPATH + pathsep + path ];
        if ( initfile ) {
            args.push(initfile);
        }

        function onError(msg) {
            SLFSRV.alertClose(alertBox);
            SLFSRV.alert("ERROR LAUNCHING APPLICATION:\n\n"+msg);
        }

        function onSuccess(ret) {
            if ( ret.stderr.length !== 0 ) {
                onError(ret.stderr);
            } else {
                SLFSRV.alertClose(alertBox);
                // don't show too many lines
                var results = ret.stdout.split("\n");
                if ( 45 < results.length ) {
                    results = results.slice(0,20).join("\n") +
                        "\n\n...too much data to show it all...\n\n" +
                        results.slice(results.length-25,results.length).join("\n");
                } else {
                    results = results.join("\n");
                }
                SLFSRV.alert("PROGRAM RESULTS:\n\n"+results);
            }
        }

        SLFSRV.exec({program:program,args:args},onSuccess,onError);
    });
}

function help() {
    var program, args;
    program = SLFSRV.SELF;
    args = [ program, "-help" ];

    SLFSRV.exec({program:program,args:args}).then(function(ret){
        if ( ret.stderr.length !== 0 ) {
            throw ret.stderr;
        } else {
            SLFSRV.alert("slfsrv -help:\n"+ret.stdout);
        }
    })["catch"](function(msg){
        SLFSRV.alert("ERROR LAUNCHING APPLICATION:\n\n"+msg);
    });
}

var gSources = [
    { id: 'example-helloworld',        name: "hello-world",
      src: ['helloworld/index.html'] },
    { id: 'example-globals',           name: "globals",
      src: ['globals/index.html','globals/javascript/globals.js'] },
    { id: 'example-cmd-shell',         name: "cmd-shell",
      src: ['cmd-shell/index.html','cmd-shell/javascript/shell.js'] },
    { id: 'example-store-test',        name: "store-test",
      src: ['store-test/index.html','store-test/javascript/store-test.js'] },
    { id: 'example-file-and-dir-test', name: "file-and-dir-test",
      src: ['file-and-dir-test/index.html','file-and-dir-test/javascript/file-and-dir-test.js'] },
    { id: 'example-tempdir-test', name: "tempdir-test",
      src: ['tempdir-test/index.html','tempdir-test/javascript/tempdir-test.js',
            'tempdir-test/shell-scripts/very-simple-bash.sh','tempdir-test/shell-scripts/very-simple-batch.bat'] },
    { id: 'example-callback-version',  name: "directory listing callback-version",
      src: ['recursive-file-list/callback-version.html','recursive-file-list/javascript/callback-version.js'] },
    { id: 'example-promises-version',  name: "directory listing promises-version",
      src: ['recursive-file-list/promises-version.html','recursive-file-list/javascript/promises-version.js'] },
    { id: 'example-menus',             name: "menus to web and local items",
      src: ['menu/main_menu.html','menu/main_menu.js',
            'menu/secondary_menu.html','menu/secondary_menu.js',
            'menu/javascript/slfsrv-menu.js'] },
    { id: 'example-video-test',        name: "video-test",
      src: ['video-test/index.html'] }
];

function initialize(count) {
    var idx, src, html, el;

    $(".command-prompt").text(( SLFSRV.OS==="windows" ) ? ">" : "#");

    // add options to view the source code, either in browser
    for ( idx = 0; idx < gSources.length; idx++ ) {
        src = gSources[idx];
        html = '<span style="font-size:80%;">(source: ';
        html += '<a href="javascript:void 0" onclick="showHtmlSource(' + idx + ')">html</a>';
        html += ' - ';
        html += '<a href="javascript:void 0" onclick="showEditorSource(' + idx + ')">editor</a>';
        html += ')</span>';
        el = document.getElementById(src.id);
        el.innerHTML += html;
    }
}

function showHtmlSource(idx) {
    var program, args, files, i, src;

    files = [];
    src = gSources[idx];
    for ( i = 0; i < src.src.length; i++ ) {
        files.push(src.src[i]);
    }

    program = SLFSRV.SELF;
    args = [ program, "-verbose",
             SLFSRV.ROOTPATH + SLFSRV.dir.SLASH + "examples",
             "showsourceinhtml/index.html?" + src.name.replace(/ /g,'%20') + "--" + files.join(',') ];
    SLFSRV.exec({program:program,args:args});
}

function showEditorSource(idx) {
    var i, src, filepath, program, args;
    src = gSources[idx];
    for ( i = 0; i < src.src.length; i++ ) {
        filepath = SLFSRV.ROOTPATH + ("/examples/" + src.src[i]).replace(/\//g,SLFSRV.dir.SLASH);
        if ( SLFSRV.OS==="windows" ) {
            program = 'notepad.exe';
            args = [program,filepath];
        } else {
            program = 'open';
            args = [program,'-t',filepath];
        }
        SLFSRV.exec({program:program,args:args});
    }

}

window.onload = function(e) {
    // i'm being kind of anal, but I really really need EXAMPLES.html to work correctly so
    // jquery and slfsrv-core.js both must be loaded for this to work correctly
    function look_for_jquery_and_self_serving_core(count) {
        if ( !("SLFSRV" in window) || !("$" in window)) {
            if ( count < 60 ) {
                setTimeout(function(){look_for_jquery_and_self_serving_core(count+1);},100);
            } else {
                if ( !("$" in window) ) {
                    document.getElementById("not-run-with-jquery").style.display = "block";
                } else {
                    $("#not-run-from-slfsrv").show();
                }
                document.getElementById("loading").style.display = "none";
            }
        } else {
            $( document ).ready(function() {
                $('#loading').hide();
                $('#post-loading').show();
                initialize();
            });
        }
    }
    look_for_jquery_and_self_serving_core(0);
};