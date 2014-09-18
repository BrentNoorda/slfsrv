/*jslint white:false plusplus:false browser:true nomen:false */
/*globals SLFSRV*/

function display(id,str) {
    var el;
    el = document.getElementById(id);
    el.innerHTML = SLFSRV.esc4html(str);
}

display('SLFSRV.OS',SLFSRV.OS);
display('SLFSRV.SECRET_KEY',SLFSRV.SECRET_KEY);
display('SLFSRV.PORT',String(SLFSRV.PORT));
display('SLFSRV.ROOTPATH',SLFSRV.ROOTPATH);
display('SLFSRV.INITFILE',SLFSRV.INITFILE);
display('SLFSRV.SELF',SLFSRV.SELF);
display('SLFSRV.dir.SLASH',SLFSRV.dir.SLASH);
