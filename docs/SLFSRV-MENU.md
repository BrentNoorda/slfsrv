#SLFSRV-MENU API

`slfsrv-menu.js` provides a few functions for creating menus within a browser to access
common activities both on your local computer and on the web. It provides a good example
of the kind of thing done with the `slfsrv` localhost app maker.

[&#x25C0; back to README.cmd](../README.md)

------------------------------------------------------------

To use `slfsrv-menu.js`, first create a short HTML page looking something like this:

    <html>
    <head>
      <script src="slfsrv-core.js" type="text/javascript"></script>
      <script src="jfuncytag_and_style.0.0.4.min.js" type="text/javascript"></script>
      <script src="mymenu.js" type="text/javascript"></script>
    </head>
    <body>
      <div id="put-menus-here"></div>
      <script src="main_menu.js" type="text/javascript"></script>
    </body>
    </html>

Most of the work will be done in `mymenu.js` using the following functions (for a couple of
examples, see `main_menu.js` and `secondary_menu.js` in the examples/menu directory).

------------------------------------------------------------

slfsrv-menu functions:

* [SLFSRV_MENU_GLOBALS](#slfsrv-menu-globals) - global settings
* [addMenuCategory()](#slfsrv-menu-addmenucategory) - add new category of menu options
    * [menuOption()](#slfsrv-menu-menuoption) - define line-item selection within a menu category
* [pageTitle()](#slfsrv-menu-pagetitle) - give title to menu page, and add "edit this" link
* [launch_app()](#slfsrv-menu-launch_app) - launch a local application
* [shell_command()](#slfsrv-menu-shell_command) - run a shell command
* [browse_local_file()](#slfsrv-menu-browse_local_file) - open a local file in the default browser
* [record_menu_time()](#slfsrv-menu-record_menu_time) - record when this menu item was last selected
* [set_visible_if_not_selected_recently()](#slfsrv-menu-set_visible_if_not_selected_recently) - display option if it hasn't been selected in a while
* [set_visible_if_not_selected_today()](#slfsrv-menu-set_visible_if_not_selected_today) - display option it has not yet been selected today
* [set_visible_if_not_already_running()](#slfsrv-menu-set_visible_if_not_already_running) - display option if the program is not already executing

------------------------------------------------------------

<a name="slfsrv-menu-globals"></a>
### SLFSRV_MENU_GLOBALS - global settings on behavior

    SLFSRV_MENU_GLOBALS.columnCount // If defined, this is the number of columns used by menus.
                                    // If not defined then the number of columsn will expand
                                    // to fit the browser width.

------------------------------------------------------------

<a name="slfsrv-menu-addmenucategory"></a>
### addMenuCategory() - add new category of menu options

add new category of menu options, and insert the results into the current web page within the div with id="put-menus-here"

    addMenuCategory( title<string>,img<optional string>,elements<menuOption[]>)
    // title: string to represent this category (e.g. "Social Networks")
    // img: option URL for icon to represnet this category
    // elements: array of menuOption
    // return: undefined
    // ex: /* add a music category */
    //     addMenuCategory( "Music", "icons/playrecent-music-dark.png", [
    //         menuOption({...}),
    //         menuOption({...}),
    //         menuOption({...})
    //     ] );

------------------------------------------------------------

<a name="slfsrv-menu-menuOption"></a>
### menuOption() - define line-item selection within a menu category

define line-item selection within a menu category; an array of these will make up what is displayed within a category

    menuOption( options{properties...} )
    // .title: optional (but usually supplied) string to represent this selection, and to
    //         click on to select this option
    // .img: optional url to icon to represent this selection (and to select this option if
    //       clicked on) - best if small and square (size will be adjusted to match title)
    // .description: optional string to put beside title; not clicked on
    // .html: optional raw html string to use as selection (usually only if no title)
    // .visibleTest: optional <bool>function(setVisible<function>){...})
    //      if visibleTest is supplied, then the option is initially no displayed, but at
    //      menu initialization this will be called - if the item should be displayed then
    //      this function must call setVisible(), else the item will remain hidden - this
    //      is useful when it cannot be determined until run-time whether a menu item is
    //      applicable
    // .action: optional function to call, or url to open, if this menu item is selected
    // return: undefined
    // ex: /* web sites and apps for music */
    //     addMenuCategory( "Music", "icons/playrecent-music-dark.png", [
    //     menuOption({
    //         title: "Pandora",
    //         action: "http://www.pandora.com/",
    //         img: "http://www.pandora.com/favicon.ico"
    //     }),
    //     menuOption({
    //         title: "iTunes",
    //         action: function(){ launch_app('iTunes.app'); },
    //         img: "icons/iTunes.png"
    //     }),
    //     menuOption({
    //         title: "rdio",
    //         action: function(){ launch_app('rdio.app'); },
    //         img: "icons/rdio.png",
    //         html: ' - <a href="http://www.rdio.com/">web version</a>'
    //     })
    // ] );

------------------------------------------------------------

<a name="slfsrv-menu-pagetitle"></a>
### pageTitle() - give title to menu page, and add "edit this" link

    pageTitle( title<string> )
    // title: name of page
    // ex: pageTitle("my menu");

------------------------------------------------------------

<a name="slfsrv-menu-pagetitle"></a>
### pageTitle() - give title to menu page, and add "edit this" link

    pageTitle( title<string> )
    // title: name of page
    // ex: pageTitle("my menu");

------------------------------------------------------------

<a name="slfsrv-menu-launch_app"></a>
### launch_app() - launch a local application

Note: this currently only launch OSX applications

    launch_app( app<string> )
    // app: name of application to launch
    // ex: /* this option within will launch itunes when selected */
    //     menuOption({
    //         title: "iTunes",
    //         action: function(){ launch_app('iTunes.app'); },
    //         img: "icons/iTunes.png"
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-shell_cmd"></a>
### shell_cmd() - run a shell command

    shell_cmd( command<string> )
    // command: command to run in a local shell
    // ex: /* this option run the osx chrome browser at google.com */
    //     menuOption({
    //         title: "Chrome",
    //         action: function(){ shell_command('open -a "Google Chrome" http://google.com/'); },
    //         img: "icons/chrome.png"
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-browse_local_file"></a>
### browse_local_file() - open a local file in the default browser

    browse_local_file( filespec<string> )
    // filespec: local file to launch in browser
    // ex: /* this option run the osx chrome browser at google.com */
    //     menuOption({
    //         title: "FuncyTag Pages",
    //         action: function(){ browse_local_file("~/my_projects/FuncyTag/index.html"); }
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-record_menu_time"></a>
### record_menu_time() - record when this menu item was last selected

this call is used in conjunction with
[set_visible_if_not_selected_recently()](#slfsrv-menu-set_visible_if_not_selected_recently) or
[set_visible_if_not_selected_today()](#slfsrv-menu-set_visible_if_not_selected_today)
for making menu options that do not show up constantly

    record_menu_time( label<string>, callback<function>{} )
    // label: a unique string under which to store this time stamp
    // callback: function that will be called after time is recorded
    // ex: /* show 7-minute excercises only once per day */
    //     menuOption({
    //         title: "7 Minute Exercise",
    //         action: function() {
    //             record_menu_time("7 Minute Exercise",function(){
    //                 window.location = "http://www.7-min.com/";
    //             });
    //         },
    //         visibleTest: function(setVisible) { set_visible_if_not_selected_today('7 Minute Exercise',setVisible); }
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-set_visible_if_not_selected_recently"></a>
### set_visible_if_not_selected_recently() - display option if it hasn't been selected in a while

this call is used in conjunction with
[record_menu_time()](#slfsrv-menu-record_menu_time)
for making menu options that do not show up constantly

    set_visible_if_not_selected_recently( label<string> )
    // label: a unique string under which a time stamp was stored
    // ex: /* don't be distracted too frequently by Hacker News */
    //     menuOption({
    //         title: "Hacker News",
    //         action: function(){
    //             record_menu_time('Hacker News',function(){
    //                 location.href = "http://news.ycombinator.com";
    //             });
    //         },
    //         img: "https://news.ycombinator.com/favicon.ico",
    //         description: "3+ hours",
    //         visibleTest: function(setVisible) { set_visible_if_not_selected_recently('Hacker News',3*60,setVisible); }
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-set_visible_if_not_selected_today"></a>
### set_visible_if_not_selected_today() - display option it has not yet been selected today

this call is used in conjunction with
[record_menu_time()](#slfsrv-menu-record_menu_time)
for making menu options that do not show up constantly

    set_visible_if_not_selected_today( label<string> )
    // label: a unique string under which a time stamp was stored
    // ex: /* show 7-minute excercises only once per day */
    //     menuOption({
    //         title: "7 Minute Exercise",
    //         action: function() {
    //             record_menu_time("7 Minute Exercise",function(){
    //                 window.location = "http://www.7-min.com/";
    //             });
    //         },
    //         visibleTest: function(setVisible) { set_visible_if_not_selected_today('7 Minute Exercise',setVisible); }
    //     })

------------------------------------------------------------

<a name="slfsrv-menu-set_visible_if_not_already_running"></a>
### set_visible_if_not_already_running() - display option if the program is not already executing

Note: this does not currently work on Windows

    set_visible_if_not_selected_today( label<string>, setVisibleFunc<function> )
    // label: a unique string under which a time stamp was stored
    // ex: /* if Adium is not currently running, then show an option to launch it */
    //     menuOption({
    //         title: "Adium",
    //         action: function(){
    //             launch_app('Adium.app');
    //             location.reload();
    //         },
    //         img: "https://adium.im/favicon.ico",
    //         description: 'chat',
    //         visibleTest: function(setVisible) { set_visible_if_not_already_running('Adium.app',setVisible); }
    //     })

------------------------------------------------------------

[&#x25C0; back to README.cmd](../README.md)
