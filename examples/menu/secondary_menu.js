/*jslint white:false plusplus:false browser:true nomen:false */
/*globals pageTitle, addMenuCategory, menuOption, shell_command, browse_local_file, launch_app, SLFSRV_MENU_GLOBALS*/

SLFSRV_MENU_GLOBALS.columnCount = 3;

pageTitle("Secondary Menu (slfsrv-menu demo)");

addMenuCategory( "FuncyTag/FuncyStyle", "http://www.developergeekresources.com/images/javascript64.png", [
    menuOption({
        html: 'Github <a href="https://github.com/BrentNoorda/FuncyTag">project</a> - <a href="http://brentnoorda.github.io/FuncyTag/">pages</a>'
    }),
    menuOption({
        html: 'Local <a href="#" onclick="shell_command(\'mdview ~/my_projects/FuncyTag/README.md\')">readme</a> - ' +
              '<a href="#" onclick="browse_local_file(\'/Users/brent/my_projects/FuncyTag/index.html\')">pages</a>'
    }),
    menuOption({
        html: 'Examples <a href="#" onclick="browse_local_file(\'/Users/brent/my_projects/FuncyTag/examples_js.html\')">local</a> - ' +
              '<a href="http://brentnoorda.github.io/FuncyTag/examples_js.html">online</a>'
    })
] );

addMenuCategory( "slfsrv", "icons/slfsrv.png", [
    menuOption({
        title: 'Github',
        action: "https://github.com/BrentNoorda/slfsrv"
    }),
    menuOption({
        title: "Local README.md",
        action: function() { shell_command('mdview ~/my_projects/slfsrv-dev/src/github.com/BrentNoorda/slfsrv/README.md'); }
    }),
    menuOption({
        title: "Examples",
        action: function() { shell_command('/Users/brent/my_projects/slfsrv-dev/bin/slfsrv -verbose /Users/brent/my_projects/slfsrv-dev/src/github.com/BrentNoorda/slfsrv/EXAMPLES.html'); }
    }),
    menuOption({
        title: "Documentation",
        action: "https://drive.google.com/?tab=mo&authuser=0#folders/0B4rxOB63nnDMazljY0dNbUFFVDA"
    }),
    menuOption({
        title: "Screen Sharing",
        description: "for linux testing",
        action: function(){ launch_app('Screen Sharing.app'); }
    })
] );

addMenuCategory( "NoMorePasswordsJustEmail", "http://thebesthack.net/wp-content/uploads/2013/12/no-password.png", [
    menuOption({
        html: 'Run <a href="http://localhost:3000/">local</a> - <a href="http://NoMorePasswordsJustEmail.meteor.com/">live</a>'
    }),
    menuOption({
        title: "Github project",
        action: "https://github.com/BrentNoorda/NoMorePasswordsJustEmail"
    }),
    menuOption({
        title: "Local README.md",
        action: function() { shell_command('mdview ~/my_projects/NoMorePasswordsJustEmail/README.md'); }
    }),
    menuOption({
        title: "Documentation",
        action: "https://drive.google.com/?tab=mo&authuser=0#folders/0B4rxOB63nnDMQ25ub1QzNEFyRG8"
    })
] );

addMenuCategory( "Crappynet", "http://image.spreadshirt.com/image-server/v1/designs/12567870,width=178,height=178/HAPPY-POOP-EMOTICON.png", [
    menuOption({
        title: "Administration",
        action: "http://127.0.0.1:9090/"
    }),
    menuOption({
        title: "gateway",
        action: "http://127.0.0.1:8080/"
    }),
    menuOption({
        title: "Online Docs",
        action: "https://docs.google.com/?tab=mo&authuser=0#folders/0B7HKNGn35xLQMDJmYzY5MWEtY2JjNy00ZmZiLWE5ZWYtNGQyZmZkMWMyMmYx"
    }),
    menuOption({
        title: "CrappyNet - public view",
        action: "http://tinyurl.com/CrappyNet"
    })
] );

addMenuCategory( "toy-piano", "http://png-4.findicons.com/files/icons/2155/social_media_bookmark/32/piano.png", [
    menuOption({
        html: 'Run <a href="http://localhost:3000/">local</a> - <a href="http://toy-piano.meteor.com/">live</a>'
    }),
    menuOption({
        title: "Github project",
        action: "https://github.com/BrentNoorda/toy-piano"
    }),
    menuOption({
        title: "Local readme",
        action: function() { shell_command('mdview ~/my_projects/toy-piano/README.md'); }
    }),
    menuOption({
        title: "Documentation",
        action: "https://docs.google.com/document/d/1T713e2_-DPVGCc9079-DsQ64-d6b_QFKFx0kI2fVPpI/view"
    })
] );

addMenuCategory( "Airbnb Confessions", "https://www.airbnb.com/favicon.ico", [
    menuOption({
        html: 'Run <a href="http://localhost:3000/">local</a> - <a href="http://www.airbnbconfessions.com/">live</a>'
    }),
    menuOption({
        title: "Documentation",
        action: "https://drive.google.com/?tab=mo&authuser=0#folders/0B4rxOB63nnDMWTlaTE5LVFJtd3M"
    })
] );

addMenuCategory( "Django Unusual", "https://www.djangoproject.com/favicon.ico", [
    menuOption({
        title: "Github project",
        action: "https://github.com/BrentNoorda/django_unusual"
    }),
    menuOption({
        title: "Presentation online",
        action: "https://github.com/BrentNoorda/django_unusual/blob/master/SLIDESHOW/01.md"
    })
] );

addMenuCategory( "Problems Solving Problems", "http://www.problemssolvingproblems.com/_/rsrc/1333837178455/favicon.ico", [
    menuOption({
        title: "Mail",
        action: "https://mail.google.com/a/problemssolvingproblems.com",
        img: "https://ssl.gstatic.com/ui/v1/icons/mail/images/favicon2.ico"
    }),
    menuOption({
        title: "Apps",
        action: "https://www.google.com/a/cpanel/problemssolvingproblems.com/Dashboard",
        description: "Manage Apps",
        img: "https://admin.google.com/problemssolvingproblems.com/images/cpanelfavicon.ico"
    }),
    menuOption({
        title: "Web Site",
        action: "http://problemssolvingproblems.com/",
        img: "http://www.problemssolvingproblems.com/_/rsrc/1333837178455/favicon.ico"
    }),
    menuOption({
        title: "Blog",
        action: "http://blog.problemssolvingproblems.com/",
        img: "http://www.blogger.com/favicon.ico"
    })
] );

addMenuCategory( "etc...", [
    menuOption({
        html: '<a href="main_menu.html">MAIN MENU</a>',
        img: "http://keepass.info/help/images/b16x16_kfm_home.png"
    }),
    menuOption({
        title: "business ideas",
        action: function() { shell_command('open /Users/brent/stuff/personal/business_ideas'); },
        img: "http://static.iconsplace.com/icons/preview/black/idea-256.png"
    })
] );
