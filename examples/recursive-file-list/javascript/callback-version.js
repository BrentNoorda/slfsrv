/*jslint white:false plusplus:false browser:true nomen:false */
/*globals $, SLFSRV, window*/

function output(msg) { // write message on console
    var html = '<span>';//'<p>';
    html += SLFSRV.esc4html(msg).replace(/\t/g,"&nbsp;&nbsp;&nbsp;&nbsp;").replace(/\n/g,"<br/>");
    html += '</span>';//'</p>';
    $('#output').append(html);
    window.scrollTo(0,document.body.scrollHeight);
}

function error(msg) {
    var html = '<p style="color:red;">';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    $('#output').append(html);
}

function change_to_self_serving_root_dir(cwd_callback) {
    // full path will be two directories above this one
    SLFSRV.dir.getcwd({},function(original_cwd){
        SLFSRV.dir.setcwd({dirname:SLFSRV.ROOTPATH},function(){
            SLFSRV.dir.setcwd({dirname:".."},function(){
                SLFSRV.dir.setcwd({dirname:".."},function(){
                    cwd_callback(original_cwd);
                });
            });
        });
    });
}

function restore_original_cwd(original_cwd,callback) {
    SLFSRV.dir.setcwd({dirname:original_cwd},function(){
        callback(callback);
    });
}

function list_all(callback) {

    function list_current_dir(indent,callback) {
        SLFSRV.dir.list({},function(list){

            function display_next_list_item() {
                if ( list.length === 0 ) {
                    callback();
                } else {
                    var li;
                    li = list.shift();
                    output(indent + li.name + '\n');
                    if ( li.dir ) {
                        SLFSRV.dir.setcwd({dirname:li.name},function(){
                            list_current_dir(indent+"\t",function(){
                                SLFSRV.dir.setcwd({dirname:".."},display_next_list_item);
                            });
                        });
                    } else {
                        setTimeout(display_next_list_item,0);
                    }
                }
            }

            display_next_list_item();
        });
    }

    output("Directory listing for " + SLFSRV.ROOTPATH + ':\n\n');

    list_current_dir("",callback);
}

$( document ).ready(function() {
    change_to_self_serving_root_dir(function(original_cwd) {
        list_all(function(){
            restore_original_cwd(original_cwd,function(){
                output("\nALL DONE!");
            });
        });
    });
});