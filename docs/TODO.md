# TODO

Get this stuff done, somehow someday

[&#x25C0; back to README.cmd](../README.md)

------------------------------------------------------------

### finish the menu example

it is too much a copy of my own, and doesn't have good windows examples

### way for javascript to exit

SLFSRV.exit() ?

### make the executable downloads available somewhere

for osx/darwin, windows, and linux

### figure out golang packages layout

some of the ways to to import from subdirectories of github.com/BrentNoorda/slfsrv seem
silly, and would screw up anyone wanting to fork this. Am I doing it wrong?

### don't shut down while goroutines are still running

go probably has a way to add an array waiting for the number of things that must shut down

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
