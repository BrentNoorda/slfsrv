/*jslint white:false plusplus:false browser:true nomen:false */
/*globals pageTitle, addMenuCategory, menuOption, launch_app, SLFSRV, SLFSRV_MENU_GLOBALS, window*/
/*globals set_visible_if_not_already_running*/
/*globals set_visible_if_not_selected_today, set_visible_if_not_selected_recently, record_menu_time*/
/*globals shell_command*/

SLFSRV_MENU_GLOBALS.columnCount = 3;

pageTitle("Main Menu (slfsrv-menu demo)");

addMenuCategory( "Common", [
    menuOption({
        title: "GMail",
        action: "http://gmail.com/",
        img: "https://ssl.gstatic.com/ui/v1/icons/mail/images/favicon2.ico"
    }),
    menuOption({
        title: "GCalendar",
        action: function(){
            record_menu_time('Daily View Calendar',function(){
                location.href = "http://calendar.google.com/";
            });
        },
        img: "https://calendar.google.com/googlecalendar/images/favicon_v2013_" + (new Date()).getDate() + ".ico"
    }),
    menuOption({
        title: "Workflowy",
        action: "https://workflowy.com/",
        img: "https://workflowy.com/media/i/favicon.ico"
    })
] );

addMenuCategory( "Don't look too often", [
    menuOption({
        title: "GCalendar",
        action: function(){
            record_menu_time('Daily View Calendar',function(){
                location.href = "http://calendar.google.com/";
            });
        },
        img: "https://calendar.google.com/googlecalendar/images/favicon_v2013_" + (new Date()).getDate() + ".ico",
        description: "daily",
        visibleTest: function(setVisible) { set_visible_if_not_selected_today('Daily View Calendar',setVisible); }
    }),
    menuOption({
        title: "7MIN",
        action: function() {
            record_menu_time('7MIN',function(){
                window.location = "http://www.7-min.com/";
            });
        },
        img: "http://www.7-min.com/favicon.ico",
        description: "daily exercise!!!",
        visibleTest: function(setVisible) { set_visible_if_not_selected_today('7MIN',setVisible); }
    }),
    menuOption({
        title: "Hacker News",
        action: function(){
            record_menu_time('Hacker News',function(){
                location.href = "http://news.ycombinator.com";
            });
        },
        img: "https://news.ycombinator.com/favicon.ico",
        description: "3+ hours",
        visibleTest: function(setVisible) { set_visible_if_not_selected_recently('Hacker News',3*60,setVisible); }
    }),
    menuOption({
        title: "BazQux Reader",
        action: function(){
            record_menu_time('BazQux',function(){
                location.href = "https://bazqux.com/";
            });
        },
        img: "https://bazqux.com/favicon.ico",
        description: "2.5+ hours",
        visibleTest: function(setVisible) { set_visible_if_not_selected_recently('BazQux',2.5*60,setVisible); }
    }),
    menuOption({
        title: "TechMeme",
        action: function(){
            record_menu_time('TechMeme',function(){
                location.href = "http://techmeme.com";
            });
        },
        img: "http://techmeme.com/img/favicon.ico",
        description: "8+ hours",
        visibleTest: function(setVisible) { set_visible_if_not_selected_recently('TechMeme',8*60,setVisible); }
    }),
    menuOption({
        title: "Facebook",
        action: function(){
            record_menu_time('Facebook',function(){
                location.href = "https://www.facebook.com";
            });
        },
        img: "http://facebook.com/favicon.ico",
        description: "5+ hours",
        visibleTest: function(setVisible) { set_visible_if_not_selected_recently('Facebook',5*60,setVisible); }
    }),
    menuOption({
        title: "Twitter",
        action: function(){
            record_menu_time('Twitter',function(){
                location.href = "https://twitter.com";
            });
        },
        img: "http://twitter.com/favicon.ico",
        description: "4+ hours",
        visibleTest: function(setVisible) { set_visible_if_not_selected_recently('Twitter',4*60,setVisible); }
    })
] );

addMenuCategory( "Daily Utils", [
    menuOption({
        title: "Terminal",
        action: function(){
            if ( SLFSRV.OS === "windows" ) {
                SLFSRV.exec({program:'cmd.exe',args:["cmd.exe","/c","start","cmd.exe"]});
            } else {
                launch_app('Terminal.app');
            }
        },
        img: "icons/Terminal.png"
    }),
    menuOption({
        title: "Contacts",
        action: function(){ launch_app('/Applications/Contacts.app'); },
        img: "icons/Contacts.png",
        description: "names & addresses",
        visibleTest: function(setVisible) {
            if (SLFSRV.OS === "darwin") {
                setVisible();
            }
        }
    }),
    menuOption({
        title: "Notes",
        action: function(){
            if ( SLFSRV.OS === "windows" ) {
                SLFSRV.exec({program:'notepad.exe',args:["notepad.exe"]});
            } else {
                launch_app('Notes.app');
            }
        },
        img: "icons/Notes.png",
        description: "quick post-its"
    })
] );

addMenuCategory( "Music", "icons/playrecent-music-dark.png", [
    menuOption({
        title: "rdio",
        action: function(){ launch_app('rdio.app'); },
        img: "icons/rdio.png",
        html: ' - <a href="http://www.rdio.com/">web version</a>'
    }),
    menuOption({
        title: "iTunes",
        action: function(){ launch_app('iTunes.app'); },
        img: "icons/iTunes.png"
    })
] );

addMenuCategory( "Video", "icons/playrecent-tv-shows-dark.png", [
    menuOption({
        title: "TV Worth Watching",
        action: "http://tvworthwatching.com/",
        img: "http://tvworthwatching.com/favicon.ico"
    }),
    menuOption({
        title: "Netflix",
        action: "http://netflix.com/",
        img: "http://cdn-0.nflximg.com/us/dvd/icons/favicon-16x16.ico"
    }),
    menuOption({
        title: "Rotten Tomatoes",
        action: "http://www.rottentomatoes.com/movie_times/94612/",
        img: "http://www.rottentomatoes.com/static/images/icons/favicon.ico"
    }),
    menuOption({
        title: "imdb",
        action: "http://imdb.com/",
        img: "http://www.imdb.com/favicon.ico"
    })
] );

addMenuCategory( "News", "https://ssl.gstatic.com/news-static/img/favicon.ico", [
    menuOption({
        title: "BazQux Reader",
        action: function(){
            record_menu_time('BazQux',function(){
                location.href = "https://bazqux.com/";
            });
        },
        img: "https://bazqux.com/favicon.ico"
    }),
    menuOption({
        title: "sfgate",
        action: "http://www.sfgate.com/",
        img: "http://www.sfgate.com/favicon.ico"
    }),
    menuOption({
        title: "bbc",
        action: "http://www.bbc.co.uk/",
        img: "http://www.bbc.co.uk/favicon.ico"
    }),
    menuOption({
        title: "cnn",
        action: "http://www.cnn.com/",
        img: "http://www.cnn.com/favicon.ie9.ico"
    })
] );

addMenuCategory( "Nerdy", "http://www.kia-forums.com/images/Kia-Forums/smilies/tango_face_glasses.png", [
    menuOption({
        title: "Hacker News",
        action: function(){
            record_menu_time('Hacker News',function(){
                location.href = "http://news.ycombinator.com";
            });
        },
        img: "https://news.ycombinator.com/favicon.ico",
        html: '<a href="https://news.ycombinator.com/shownew" style="font-size:75%;">Show HN New</a>'
    }),
    menuOption({
        title: "Stack Overflow",
        action: "http://stackoverflow.com/",
        img: "http://cdn.sstatic.net/stackoverflow/img/favicon.ico"
    }),
    menuOption({
        title: "QA/Testing Exchange",
        action: "http://sqa.stackexchange.com/",
        img: "http://cdn.sstatic.net/sqa/img/favicon.ico"
    }),
    menuOption({
        title: "TechMeme",
        action: function(){
            record_menu_time('TechMeme',function(){
                location.href = "http://techmeme.com";
            });
        },
        img: "http://techmeme.com/img/favicon.ico"
    }),
    menuOption({
        title: "Slashdot",
        action: "http://slashdot.com",
        img: "http://slashdot.org/favicon.ico"
    })
] );

addMenuCategory( "Search", "icons/websearch.png", [
    menuOption({
        title: "Google",
        action: "https://google.com/",
        img: "http://www.google.com/favicon.ico"
    }),
    menuOption({
        title: "DuckDuckGo",
        action: "https://duckduckgo.com/",
        img: "http://duckduckgo.com/favicon.ico"
    }),
    menuOption({
        title: "Yahoo",
        action: "http://yahoo.com/",
        img: "http://yahoo.com/favicon.ico"
    }),
    menuOption({
        title: "bing",
        action: "http://bing.com",
        img: "http://bing.com/favicon.ico"
    }),
    menuOption({
        title: "wikipedia",
        action: "http://wikipedia.org",
        img: "http://wikipedia.org/favicon.ico"
    })
] );

addMenuCategory( "Social", "https://cdn2.iconfinder.com/data/icons/seo-web-optomization-ultimate-set/512/team_building-512.png", [
    menuOption({
        title: "Adium",
        action: function(){
            launch_app('Adium.app');
            location.reload();
        },
        img: "https://adium.im/favicon.ico",
        description: 'chat',
        visibleTest: function(setVisible) { set_visible_if_not_already_running('Adium.app',setVisible); }
    }),
    menuOption({
        title: "Facebook",
        action: function(){
            record_menu_time('Facebook',function(){
                location.href = "https://www.facebook.com";
            });
        },
        img: "http://facebook.com/favicon.ico"
    }),
    menuOption({
        title: "Twitter",
        action: function(){
            record_menu_time('Twitter',function(){
                location.href = "https://twitter.com";
            });
        },
        img: "http://twitter.com/favicon.ico"
    }),
    menuOption({
        title: "Hootsuite",
        action: "https://hootsuite.com",
        img: "http://hootsuite.com/favicon.ico"
    }),
    menuOption({
        title: "groupme",
        action: "http://groupme.com/groups/2032623",
        img: "http://groupme.com/favicon.ico"
    }),
    menuOption({
        title: "Google Voice",
        action: "https://www.google.com/voice",
        img: "https://www.google.com/voice/favicon.ico"
    })
] );

addMenuCategory( "Browsers", "http://stevet3ch.com/img/world.png", [
    menuOption({
        title: "Chrome",
        action: function(){ shell_command('open -a "Google Chrome" ' + window.location); },
        img: "icons/chrome.png"
    }),
    menuOption({
        title: "Firefox",
        action: function(){ shell_command('open -a firefox ' + window.location); },
        img: "icons/firefox.png"
    }),
    menuOption({
        title: "Safari",
        action: function(){ shell_command('open -a Safari ' + window.location); },
        img: "icons/safari.png"
    })
] );

addMenuCategory( "Images", "https://cdn3.iconfinder.com/data/icons/mini-glyph/48/mini_glyph_gallery-256.png", [
    menuOption({
        title: "Picasa",
        description: "photos",
        action: function(){ launch_app('Picasa.app'); },
        img: "icons/Picasa.png"
    }),
    menuOption({
        title: "Image Capture",
        description: "scanner",
        action: function(){ launch_app('Image Capture.app'); },
        img: "http://icons.iconarchive.com/icons/tpdkdesign.net/refresh-cl/32/Hardware-Scanner-2-icon.png"
    })
] );

addMenuCategory( "etc...", [
    menuOption({
        title: "...another menu",
        action: "secondary_menu.html",
        img: "http://alumni.rice.edu/media/microsites/images/menu-icon.png"
    })
] );
