/*jslint white:false plusplus:false browser:true nomen:false */
/*globals $, SLFSRV, window*/

function output(msg) { // write message on console
    var html = '<p>';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    $('#output').append(html);
    window.scrollTo(0,document.body.scrollHeight);
}

function error(msg) {
    var html = '<p style="color:red;">';
    html += SLFSRV.esc4html(msg).replace(/\n/g,"<br/>");
    html += '</p>';
    $('#output').append(html);
    window.scrollTo(0,document.body.scrollHeight);
}

function hr() {
    $('#output').append("<hr/>");
}

function test() {
    var tests, test_idx;

    function next_test() {
        test_idx++;
        hr();
        if ( tests.length <= test_idx ) {
            output("FINISHED!");
        } else {
            output("TEST " + test_idx + ": " + tests[test_idx].name);
            setTimeout(function(){tests[test_idx].test();},0);
        }
    }

    tests = [{
        name: 'delete store in case any old data is there',
        test: function() {
            SLFSRV.store.list({}).then(function(keys){
                function delete_first_key() {
                    if ( keys.length === 0 ) {
                        output("done");
                        next_test();
                    } else {
                        var key = keys.shift();
                        output("deleting key: "+key);
                        SLFSRV.store.set({key:key,val:undefined},function(){
                            delete_first_key();
                        },function(err){
                            error(err);
                        });
                    }
                }
                delete_first_key();
            })["catch"](function(err){
                error(err);
            });
        }
    },{
        name: 'store must now be empty',
        test: function() {
            SLFSRV.store.list({},function(keys){
                if ( keys.length !== 0 ) {
                    error("wasn't empty, instead had these keys: " + JSON.stringify(keys));
                } else {
                    output("done");
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    },{
        name: 'read key "foo" should be undefined',
        test: function() {
            SLFSRV.store.get({key:"foo"},function(value){
                output("value = " + value);
                if ( value !== undefined ) {
                    error("value should have been undefined");
                } else {
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    },{
        name: 'write key "foo" as 42',
        test: function() {
            SLFSRV.store.set({key:"foo",val:42},function(){
                output("done");
                next_test();
            },function(err){
                error(err);
            });
        }
    },{
        name: 'read key "foo" should be 42',
        test: function() {
            SLFSRV.store.get({key:"foo"},function(value){
                output("value = " + value);
                if ( value !== 42 ) {
                    error("value should have been 42");
                } else {
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    },{
        name: 'add a more complex object',
        test: function() {
            var obj = {
                poop: "doggy doo doo",
                person: {
                    fname: "Elmer",
                    "middle initial": "J",
                    "last-name": "Fudd",
                    occupation: "millionaire"
                },
                wealth: 1000000.0
            };
            SLFSRV.store.set({key:"complex-object",val:obj},function(){
                output("done");
                next_test();
            },function(err){
                error(err);
            });
        }
    },{
        name: 'the only objects should now be "foo" and "complex-object"',
        test: function() {
            SLFSRV.store.list({},function(keys){
                if ( keys.length !== 2 ) {
                    error("keys.length = " + keys.length + "but expected 2");
                } else {
                    for ( var i = 0; i < keys.length; i++ ) {
                        if ( (keys[i] !== "foo") && (keys[i] !== "complex-object") ) {
                            error('found key "' + keys[i] + '" but only expect "foo" or "complex-object"');
                            return;
                        }
                    }
                    output("done");
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    },{
        name: 'verify contents of "foo"',
        test: function() {
            SLFSRV.store.get({key:"foo"},function(value){
                if ( value !== 42 ) {
                    error("value was " + value + " but should have been 42");
                } else {
                    output("done");
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    },{
        name: 'verify contents of "complex-object"',
        test: function() {
            SLFSRV.store.get({key:"complex-object"},function(obj){
                if ( !('wealth' in obj) ) {
                    error('"wealth" field is missing');
                } else if ( obj.wealth !== 1000000 ) {
                    error('"wealth" was ' + obj.wealth + " but should have been " + 1000000);
                } else if ( !("person" in obj) ) {
                    error('"person" field is missing');
                } else if ( typeof(obj.person) !== "object" ) {
                    error('"person" should be an object, but is type ' + typeof(obj.person));
                } else if ( !("fname" in obj.person) ) {
                    error('"person" does not contain property "fname"');
                } else if ( obj.person.fname !== "Elmer" ) {
                    error('"person.fname" = ' + obj.person.fname + ' but should === "Elmer"');
                } else if ( !("occupation" in obj.person) ) {
                    error('"person" does not contain property "occupation"');
                } else if ( obj.person.occupation !== "millionaire" ) {
                    error('"person.occupation" = ' + obj.person.occupation + ' but should === "millionaire"');
                } else if ( !("last-name" in obj.person) ) {
                    error('"person" does not contain property "middle initial"');
                } else if ( obj.person["middle initial"] !== "J" ) {
                    error("'person[\"middle initial\"]' = " + obj.person["middle initial"] + ' but should === "J"');
                } else if ( !("last-name" in obj.person) ) {
                    error('"person" does not contain property "last-name"');
                } else if ( obj.person["last-name"] !== "Fudd" ) {
                    error("'person[\"last-name\"]' = " + obj.person["last-name"] + ' but should === "Fudd"');
                } else {
                    output("done");
                    next_test();
                }
            },function(err){
                error(err);
            });
        }
    }];


    test_idx = -1;
    next_test();
}


$( document ).ready(function() {
    test();
});