/*jslint white:false plusplus:false browser:true nomen:false */
/*globals SLFSRV, window*/

function show_file_source(fspec) {
    var li, html, fullspec;

    li = document.createElement("li");
    html = '<a href="#' + fspec + '">' + fspec.replace(/\//g,SLFSRV.dir.SLASH) + '</a>';
    li.innerHTML = html;
    document.getElementById("insert-links-here").appendChild(li);

    fullspec = SLFSRV.ROOTPATH + ("/" + fspec).replace(/\//g,SLFSRV.dir.SLASH);
    SLFSRV.file.read({filename:fullspec},function(contents){
        var div, html;
        html = '<a name="' + fspec +'"></a><h3>' + SLFSRV.esc4html(fspec.replace(/\//g,SLFSRV.dir.SLASH)) + '</h3>';
        html += '<code><pre>';
        html += SLFSRV.esc4html(contents);
        html += '</pre></code>';
        div = document.createElement("div");
        div.innerHTML = html;
        document.getElementById("insert-source-here").appendChild(div);
    },function(err){
        SLFSRV.alert("err " + err);
    });
}


function init() {
    var i, filenames, title, parts, query;
    query = window.location.href.split("?")[1].split("#")[0];
    parts = query.split("--");
    title = parts[0];
    filenames = parts[1].split(",");

    document.getElementById('title-here').innerHTML = SLFSRV.esc4html(decodeURI(title));

    for ( i = 0; i < filenames.length; i++ ) {
        show_file_source(filenames[i]);
    }
}

init();
