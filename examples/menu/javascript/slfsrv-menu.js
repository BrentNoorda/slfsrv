/*jslint white:false plusplus:false browser:true nomen:false */
/*globals SLFSRV, funcyTag, funcyStyle, window, alert*/

var SLFSRV_MENU_GLOBALS = {
    // override any of the following to change behavior
    columnCount: undefined  // define as integer if you want to set column count (else will wrap with browser)
};

//----------------------------------------------------------------

var gColumnCount = 0;

function close_myself() {
    //window.open("close.html","_self");
}

function launch_app(app_name) {
    var program, args;
    program = "open";
    args = [program,"-a",app_name];
    SLFSRV.exec({program:program,args:args}).then(function(){close_myself();});
}

function edit_this_menu() {
    var filename, program, args;
    filename = window.location.pathname.substring(1+SLFSRV.SECRET_KEY.length+1).replace('.html','.js');
    program = "open";
    //var args = [program,"-a","TextEdit",filename];
    args = [program,"-t",filename];
    SLFSRV.exec({program:program,args:args}).then(function(){close_myself();});

    //var program = "/bin/sh";
    //var args = [program,"-c","/Users/brent/bin/e " + SLFSRV.ROOTPATH + SLFSRV.dir.SLASH + filename];
    //SLFSRV.exec({program:program,args:args},function(){});
}

function shell_command(cmd) {
    var program, args;
    program = "/bin/sh";
    args = [program,"-c",cmd];
    SLFSRV.exec({program:program,args:args},function(){});
}

function browse_local_file(filespec) {
    var program, args;
    program = "open";
    args = [program,filespec];
    SLFSRV.exec({program:program,args:args}).then(close_myself);
}

function record_menu_time(dataKey,callback) {
    SLFSRV.store.set( dataKey, Math.floor((new Date()).getTime() / 1000), callback );
}

function set_visible_if_not_selected_recently(dataKey,minutesSincePreviousSelection,setVisible) {
    SLFSRV.store.get(dataKey,0,function(thenEpochTime){
        var nowEpochTime, elapsedMinutes;
        nowEpochTime = (new Date()).getTime() / 1000;
        elapsedMinutes = (nowEpochTime - thenEpochTime) / 60;
        if ( minutesSincePreviousSelection <= elapsedMinutes ) {
            setVisible();
        }
    });
}

function set_visible_if_not_selected_today(dataKey,setVisible) {
    SLFSRV.store.get(dataKey,0,function(epochTime){
        var now = new Date(), then = new Date( epochTime * 1000 );
        if ( (now.getDate() !== then.getDate()) || (now.getMonth() !== then.getMonth()) || (now.getFullYear() !== then.getFullYear()) ) {
            setVisible();
        }
    });
}

function set_visible_if_not_already_running(appSearchtring,setVisible) {
    SLFSRV.exec( {program:'ps',args:['ps','-ax']} ).then(function(out){
        if ( -1 === out.stdout.indexOf(appSearchtring) ) {
            setVisible();
        }
    });
}

window.onerror = function (msg, url, line, col, error) {
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

//////////////////////////////////////////////////////////////////////////////////

// initialize choice image to be 1x1 pixel, but that will be adjusted later
funcyStyle( '.choice-img', {
                width_px: 1,
                height_px: 1,
                verticalAlign: 'bottom',
                marginRight: 0
            } );
funcyStyle( '.category-img', {
                width_px: 1,
                height_px: 1,
                verticalAlign: 'bottom'
            } );
funcyStyle( 'h2', {
                margin_px: 0,
                marginBottom_px: -10
            } );
funcyStyle( 'h3', {
                paddingLeft_em: 0.5,
                marginBottom_px: 0,
                marginTop_em: 1
            } );
funcyStyle( 'ul', {
                paddingLeft_em: 1.5,
                marginTop_px: 4,
                marginBottom_px: 0,
                marginLeft_px: 4,
                listStyleType: 'none'
            } );
funcyStyle( 'li', {
                marginTop_px: 4
            } );
funcyStyle( 'a', {
                color: funcyStyle.color(0,0,238),
                textDecoration: 'none',
                fontWeight: 'bold'
            } );
funcyStyle( 'a:hover', {
                backgroundColor: '#000066',
                fontStyle: 'normal',
                color: '#FFFFFF'
            } );



var h2 = funcyTag('h2'), h3 = funcyTag('h3'), div = funcyTag('div'), span = funcyTag('span'),
     a = funcyTag('a'), ul = funcyTag('ul'), li = funcyTag('li');
var choiceImg = funcyTag( 'img', { class_:'choice-img' });
var categoryImg = funcyTag( 'img', { class_:'category-img' });
var choiceA = funcyTag( 'a', { href:'#', _nobrin:true } ); //, cssFontWeight:'bold' } );
var firstChoiceInitialized = false, firstCategoryInitialized = false;

function pageTitle(title) {
    document.title = title;
    var el = h2( title,
                 a( { href: '#',
                      onclick: edit_this_menu,
                      cssMarginLeft_em: 2,
                      cssFontSize_px: 13,
                      cssFontWeight: 'normal'
                    },
                    'edit this' ) );

    document.getElementById("put-menus-here").appendChild( el.createElement() );
}

function menuOption(opt/*{ required properties:title  optional properties:action,html,img,description,visibleTest}*/) {
    // return undefined if conditional says not to display this action
    var els, list, actnTag;
    els = [];

    function doAction() {
        if ( opt.action ) {
            if ( typeof(opt.action) === "function" ) {
                opt.action();
            } else {
                //location.href = opt.action;
                window.open(opt.action,"_self");
            }
        }
    }

    actnTag = ( opt.action ) ? choiceA : span ;

    if ( opt.img ) {
        els.push( actnTag( { onclick: doAction },
                           choiceImg({src:opt.img}) ) );
    } else {
        els.push( actnTag( { onclick: doAction, class_:'choice-img', cssDisplay:'inline-block',
                             cssTextAlign:'center', cssVerticalAlign:'middle', cssColor:'#999999' },
                                "&#8226;" ) );
    }

    if ( opt.title ) {
        els.push( actnTag( {
                    onclick: doAction,
                    oninit: function(elem) {
                        // if this is the first choice ever initialized, use it's height
                        // to adjust the height of all choiceImg elements
                        if ( !firstChoiceInitialized && !opt.visibleTest ) {
                            firstChoiceInitialized = true;
                            funcyStyle.injectStyle(
                                funcyStyle( '.choice-img', {
                                                width_px: Math.floor(elem.offsetHeight/16) * 16,
                                                height_px: Math.floor(elem.offsetHeight/16) * 16,
                                                lineHeight_px: elem.offsetHeight - 1,
                                                marginBottom_px:1,
                                                verticalAlign: 'bottom'
                                            }
                                          )
                            );
                        }
                    }
                  }, SLFSRV.esc4html(opt.title) ) );

        if ( opt.description ) {
            els.push("-");
            els.push(SLFSRV.esc4html(opt.description));
        }
    }
    if ( opt.html ) {
        //els.push( span( { oninit: function(elem) { elem.innerHTML = opt.html; } } ) );
        els.push( opt.html );
    }

    if ( opt.visibleTest ) {
        list = li( { cssDisplay:'none',
                     oninit: function(elem) {
                         opt.visibleTest( function() {
                             elem.style.display = 'block';
                         });
                     } },
                   els );
    } else {
        list = li(els);
    }
    return list;
}

function addMenuCategory(title,img,elements) {
    var attrs, mc, parentDiv, addNewline = false, parentElem;

    parentElem = document.getElementById("put-menus-here");
    attrs = { cssVerticalAlign:'top', cssMarginRight_em:2, cssWidth_em:17,
              cssDisplay:'inline-block', cssWhiteSpace:'normal' };

    if ( SLFSRV_MENU_GLOBALS.columnCount ) {
        if ( 0 === (gColumnCount % SLFSRV_MENU_GLOBALS.columnCount) ) {
            // need new div to prevent unwanted newlines between divs
            parentDiv = div( { cssWhiteSpace:'nowrap' } );
            parentElem.appendChild( parentDiv.createElement() );
        }
        parentElem = parentElem.lastChild;
    }

    mc =  div( attrs,
               h3( { oninit: function(elem) {
                        if ( !firstCategoryInitialized ) {
                            firstCategoryInitialized = true;
                            funcyStyle.injectStyle(
                                funcyStyle( '.category-img', {
                                                width_px: elem.offsetHeight,
                                                height_px: elem.offsetHeight
                                            }
                                          )
                               );
                        }
                   } },
                   img ? categoryImg({src:img}) : undefined,
                   SLFSRV.esc4html(title) ),
               ul( elements )
          );

    //SLFSRV.alert(String(mc));
    parentElem.appendChild( mc.createElement() );
    gColumnCount++;
}
