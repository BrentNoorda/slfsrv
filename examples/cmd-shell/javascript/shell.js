/*jslint white:false plusplus:false browser:true nomen:false */
/*globals $, SLFSRV*/

function shell() {

    var prompt, runCommand, history;

    history = [];

    prompt = function(cwd) {

        function gotCwd(cwd) {
            var curDir, html, historyIdx = history.length;

            // get the current directory name (not the full path)
            curDir = cwd.split((SLFSRV.OS === "windows") ? "\\" : "/");
            curDir = curDir[curDir.length-1];

            // show a new prompt
            html = '<div class="prompt">';
            html +=  '<table>';
            html +=    '<tr>';
            html +=      '<td width="1%" id="curdir">';
            html +=        '<tt>' + SLFSRV.esc4html(curDir+">") + '</tt>';
            html +=      '</td>';
            html +=      '<td width="99%" id="command">';
            html +=        '<input id="cmd-input">';
            html +=      '</td>';
            html +=    '</tr>';
            html +=  '</table>';
            html +='</div>';

            // add this new div at the end of the console
            $('#console').append(html);

            $('#cmd-input').focus();

            // command-line completion
            function command_line_completion(el) {

                var cmdline, splits, lastword, lastslash, lastdir, lastfile;

                cmdline = el.val();
                splits = cmdline.split(" ");
                lastword = splits.pop();

                lastslash = lastword.lastIndexOf(SLFSRV.dir.SLASH);
                if ( lastslash === -1 ) {
                    lastdir = "";
                    lastfile = lastword;
                } else {
                    lastdir = lastword.substring(0,lastslash+1);
                    lastfile = lastword.substring(lastslash+1);
                }

                SLFSRV.dir.list((lastdir==="")?".":lastdir).then(function(list){
                    var i, len, candidate, candidates, lastfileLower, a, isDir;

                    if ( el.size() === 0 ) {
                        return; // this element no longer valid
                    }

                    len = lastfile.length;
                    lastfileLower = lastfile.toLowerCase();
                    candidates = [];
                    for ( i = 0; i < list.length; i++ ) {
                        candidate = list[i].name;
                        if ( candidate.toLowerCase().substring(0,len) === lastfileLower ) {
                            candidates.push(candidate);
                            isDir = list[i].dir;
                        }
                    }
                    if ( candidates.length === 1 ) {
                        // found a match
                        cmdline = splits.join(" ") + " " + lastdir + candidates[0];
                        if ( isDir && (lastfile === candidates[0]) ) {
                            cmdline += SLFSRV.dir.SLASH;
                        }
                        el.val(cmdline);
                    } else if ( candidates.length !== 0 ) {
                        // found many matches
                        a = SLFSRV.alert(candidates.join("\n"),function(){
                            el.focus();
                        });
                        setTimeout(function(){
                            SLFSRV.alertClose(a);
                        },1000);
                    }
                });
            }

            // from http://stackoverflow.com/questions/1314450/jquery-how-to-capture-the-tab-keypress-within-a-textbox
            $('#cmd-input').on('keydown',function(e){
                var keyCode = e.keyCode || e.which;

                if (keyCode === 9) {
                    e.preventDefault();
                    command_line_completion($(e.target));
                }
            });

            // run until enter is pressed in that field
            $("#cmd-input").keyup(function (e) {
                if ( e.keyCode === 13 ) {

                    // field is finished, so get the command, replace prompt with a non-input version
                    // then move on to doing the command
                    var cmd = $("#cmd-input").val();

                    $('#cmd-input').parent().html('<tt>' + SLFSRV.esc4html(cmd) + '</tt>');

                    runCommand($.trim(cmd));
                } else if ( e.keyCode === 38 ) {
                    // up
                    historyIdx--;
                    if ( 0 <= historyIdx ) {
                        $("#cmd-input").val(history[historyIdx]);
                    } else {
                        historyIdx = history.length;
                        $("#cmd-input").val('');
                    }
                } else if ( e.keyCode === 40 ) {
                    // down
                    if ( historyIdx === history.length ) {
                        historyIdx = 0;
                    } else {
                        historyIdx++;
                    }
                    if ( historyIdx < history.length ) {
                        $("#cmd-input").val(history[historyIdx]);
                    } else {
                        $("#cmd-input").val('');
                    }
                }
            });

        }

        runCommand = function(cmd) {
            var cmdParts, program;

            if ( cmd.length === 0 ) {
                prompt();
                return;
            }

            function onSuccess(ret) {
                var html, el;

                if ( ret.stderr.length !== 0 ) {
                    html = '<div><pre class="stderr"></pre></div>';
                    $('#console').append(html);
                    el = $('#console .stderr').last();
                    el.text(ret.stderr);
                }

                if ( ret.stdout.length !== 0 ) {
                    html = '<div><pre class="stdout"></pre></div>';
                    $('#console').append(html);
                    el = $('#console .stdout').last();
                    el.text(ret.stdout);
                }

                prompt();
            }

            function onError(errMsg) {
                onSuccess({exitStatus:-1,stdout:"",stderr:errMsg});
            }

            history.push(cmd);

            cmdParts = cmd.split(" ");

            if ( cmdParts[0].toLowerCase() === "cd" || cmdParts[0].toLowerCase() === "chdir" ) {
                if ( cmdParts.length === 1 ) {
                    SLFSRV.dir.getcwd(function(cwd){
                        onSuccess({exitStatus:0,stdout:cwd,stderr:""});
                    });
                } else {
                    cmdParts.shift();
                    SLFSRV.dir.setcwd(cmdParts.join(' '),function(ret){
                        SLFSRV.dir.getcwd(function(cwd){
                            onSuccess({exitStatus:0,stdout:cwd,stderr:""});
                        });
                    },onError);
                }
            } else if ( cmdParts.length === 1 && cmdParts[0].toLowerCase() === "pwd" ) {
                SLFSRV.dir.getcwd(function(cwd){
                    onSuccess({exitStatus:0,stdout:cwd,stderr:""});
                });
            } else {

                if ( SLFSRV.OS === "windows" ) {
					program = "cmd.exe";
					cmdParts = [program,"/c",cmd];
				} else {
					program = "sh";
					cmdParts = [program,"-c",cmd];
				}
                SLFSRV.exec({program:program,args:cmdParts},onSuccess,onError);
            }
        };

        SLFSRV.dir.getcwd(gotCwd);
    };

    prompt();
}



$( document ).ready(function() {
    shell();
});