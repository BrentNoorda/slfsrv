/*jslint white:false plusplus:false browser:true nomen:false */
/*globals $, SLFSRV, Promise, window*/

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

function change_to_self_serving_root_dir() {

    // full path will be two directories above this one
    return new Promise(function(resolve,reject){
        var original_cwd;
        SLFSRV.dir.getcwd({}).then(function(cwd){
            original_cwd = cwd;
            return SLFSRV.dir.setcwd({dirname:SLFSRV.ROOTPATH});
        }).then(function(){
            return SLFSRV.dir.setcwd({dirname:".."});
        }).then(function(){
            return SLFSRV.dir.setcwd({dirname:".."});
        }).then(function(){
            resolve(original_cwd);
        });
    });
}

function restore_original_cwd(original_cwd) {
    return new Promise(function(resolve,reject){
        SLFSRV.dir.setcwd({dirname:original_cwd}).then(resolve);
    });
}

function list_all() {

    return new Promise(function(resolve,reject){

        function list_current_dir(indent) {
            return new Promise(function(resolve,reject){
                SLFSRV.dir.list({}).then(function(list){

                    function display_next_list_item() {
                        if ( list.length === 0 ) {
                            resolve();
                        } else {
                            var li;
                            li = list.shift();
                            output(indent + li.name + '\n');
                            if ( li.dir ) {
                                SLFSRV.dir.setcwd({dirname:li.name}).then(function(){
                                    return list_current_dir(indent+"\t");
                                }).then(function(){
                                    return SLFSRV.dir.setcwd({dirname:".."});
                                }).then(
                                    display_next_list_item
                                );
                            } else {
                                setTimeout(display_next_list_item,0);
                            }
                        }
                    }

                    display_next_list_item();
                });
            });
        }

        output("Directory listing for " + SLFSRV.ROOTPATH + ':\n\n');
        list_current_dir("").then(resolve);
    });
}

$( document ).ready(function() {
    var original_cwd;
    change_to_self_serving_root_dir().then(function(cwd) {
        original_cwd = cwd;
        return list_all();
    }).then(function(){
        return restore_original_cwd(original_cwd);
    }).then(function(){
        output("\nALL DONE!");
    });
});