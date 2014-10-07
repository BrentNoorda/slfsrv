# TODO

Get this stuff done, somehow someday

[&#x25C0; back to README.cmd](../README.md)

------------------------------------------------------------

### add the ability to fetch files from other domains (avoiding the cross origin issue)

From [Hacker News comments](https://news.ycombinator.com/item?id=8424418)

    One thing that seems missing is the ability to fetch files from other domains
    (avoiding the cross origin issue). Can it somehow magically do this just
    using XMLHttpRequest?

    reply

        Good idea. I'll add "fetch files from other domains" to the TODO list, as either
        its own call, or maybe do something like allowing "http://" protocols in the
        SLFSRV.read() method.

        cross-origin issues would probably prevent some xmlhttprequest calls directly.
        Meanwhile, here is one current hack you could use to read files from other domains
        (although it assumes that "curl" is installed):

            SLFSRV.exec( { program:'curl', args:['curl','https://news.ycombinator.com/'] },
                         function(ret) {
                             console.log(ret.stdout);
                         }
                      );

### finish the menu example

it is too much a copy of my own, and doesn't have good windows examples

### a way for javascript to exit

SLFSRV.exit() ?

### improve golang packages layout

some of the ways to to import from subdirectories of github.com/BrentNoorda/slfsrv seem
silly, and would screw up anyone wanting to fork this. Am I doing it wrong?

### don't shut down while goroutines are still running

go probably has a convenient and efficitent way to add wait for a bunch of goroutines to be finished?

### catch ctrl-c shutdown

shut down niceley

### file permisssions

File permissions are mostly being ignored. Should probably include file permissions in
SLFSRV.write() and in unbundling files (made from -bundle and -compile options).

### answer some questions

Is GO the right tool? Maybe a minimal C web server?

compare with node-webkit?

compare with https://github.com/chriskiehl/Gooey

How can this be improved as a Chrome extension or Firefox plugin?

------------------------------------------------------------------------------

[&#x25C0; back to README.cmd](../README.md)
