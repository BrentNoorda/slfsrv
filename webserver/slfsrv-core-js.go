package webserver

func selfServingJsSrc() string {
    return `


// http://javascript.nwbox.com/asyncAlert/
// http://javascript.about.com/library/blmodald2.htm
// http://blog.raventools.com/create-a-modal-dialog-using-css-and-javascript/

var SLFSRV = {
    SECRET_KEY : "<<<SECRET_KEY>>>",
    OS : "<<<OS>>>",
    PORT : 0/*<<<PORT>>>*/,
    ROOTPATH : "<<<ROOTPATH>>>",
    INITFILE : "<<<INITFILE>>>",
    SELF : "<<<SELF>>>"
};

(function () {
    'use strict';

    var sendRequest, XMLHttpFactories, createXMLHTTPObject, alertStack;

    function log(msg) {
        try {
            console.log(msg);
        } catch(e) {

        }
    }

    SLFSRV.esc4html = function(str)
    {
        return String(str).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
    };

    function return_callbacks_or_promise(onSuccess,onError,promise) {
        if ( onSuccess === undefined ) {
            return promise;
        } else {
            promise.then(onSuccess,function(err){
                if ( onError ) {
                    onError(err);
                } else {
                    throw err;
                }
            });
            return;
        }
    }

    SLFSRV.dir = {
        /* SLFSRV.dir.SLASH */
        SLASH : (SLFSRV.OS === "windows") ? "\\" : "/",

        /* SLFSRV.dir.getcwd() */
        getcwd : function(onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("cwd.get",{},function(result){
                    resolve(result.cwd);
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.dir.setcwd() */
        setcwd : function(dir,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("cwd.set",{ dir: dir },resolve,reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.dir.list() */
        list : function( optionalDir,onSuccess,onError ) {
            if ( (typeof optionalDir) !== "string" ) {
                onError = onSuccess;
                onSuccess = optionalDir;
                optionalDir = ".";
            }
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("dir.list",{ dir:optionalDir },function(ret){
                    var i, list;
                    list = ret.list;
                    for ( i = 0; i < list.length; i++ ) {
                        list[i].mod = new Date(list[i].mod);
                    }
                    resolve(list);
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.dir.exists() */
        exists : function(dirname,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("dir.exists",{ dir:dirname },function(result){
                    resolve( result.exists );
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        }
    };

    SLFSRV.env = {
        /* SLFSRV.env.get() */
        get : function(key,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("env.get",{ key:key },function(result){
                    resolve(result.value);
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.env.set() */
        set : function(key,value,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("env.set",{ key:key, value:value },resolve,reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        }
    };

    SLFSRV.file = {
        /* SLFSRV.file.read() */
        read : function(filename,optionalMode,onSuccess,onError) {
            if ( typeof(optionalMode) === "function" ) {
                onError = onSuccess;
                onSuccess = optionalMode;
                optionalMode = "text";
            }
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("file.read",{ file:filename, mode:optionalMode },function(result){
                    resolve( result.contents );
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.file.write() */
        write : function(filename,optionalMode,contents,onSuccess,onError) {
            if ( (typeof(contents) === "function") || (contents === undefined) ) {
                onError = onSuccess;
                onSuccess = contents;
                contents = optionalMode;
                optionalMode = ""; // same as "text,create"
            }
            var promise = new Promise(function(resolve,reject){
                // i suspect there's a bug in the JSON decoding on the golang server, so if binary
                // mode then send this whole thing as an array of bytes
                if ( -1 !== optionalMode.indexOf("binary") ) {
                    var i, bytes;
                    bytes = [];
                    for ( i = 0; i < contents.length; i++ ) {
                        bytes.push(contents.charCodeAt(i));
                    }
                    SLFSRV.callServer("file.write",{ file:filename, mode:optionalMode, bytes:bytes },resolve,reject);
                } else {
                    SLFSRV.callServer("file.write",{ file:filename, mode:optionalMode, contents:contents },resolve,reject);
                }
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.file.exists() */
        exists : function(filename,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("file.exists",{ file:filename },function(result){
                    resolve( result.exists );
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.file.remove() */
        remove : function(filename,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("file.remove",{ file:filename },resolve,reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        }
    };

    SLFSRV.store = {
        /* SLFSRV.store.get() */
        get : function(key,optionalDefaultValue,onSuccess,onError) {
            if ( typeof(optionalDefaultValue) === "function" ) {
                onError = onSuccess;
                onSuccess = optionalDefaultValue;
                optionalDefaultValue = undefined;
            }
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("store.get",{ key:key },function(result){
                    resolve( result.hasOwnProperty('value') ? result.value : optionalDefaultValue );
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.store.set() */
        set : function(key,value,onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                var obj = ( value === undefined ) ? { key:key } : { key:key, value:value };
                SLFSRV.callServer("store.set",obj,resolve,reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        },

        /* SLFSRV.store.list() */
        list : function(onSuccess,onError) {
            var promise = new Promise(function(resolve,reject){
                SLFSRV.callServer("store.list",{},function(obj){
                    resolve(obj.keys);
                },reject);
            });
            return return_callbacks_or_promise(onSuccess,onError,promise);
        }
    };

    /* SLFSRV.exec() */
    SLFSRV.exec = function(options/*program,args,input,timeout*/,onSuccess,onError) {
        var promise = new Promise(function(resolve,reject){
            var obj = { program: options.program,
                        args: (options.args === undefined) ? {} : options.args,
                        input: (options.input === undefined) ? "" : options.input };
            SLFSRV.callServer("exec",obj,resolve,reject,options.timeout);
        });
        return return_callbacks_or_promise(onSuccess,onError,promise);
    };

    /* SLFSRV.tempdir() */
    SLFSRV.tempdir = function(optionalUnpackDir,onSuccess,onError) {
        if ( typeof(optionalUnpackDir) === "function" ) {
            onError = onSuccess;
            onSuccess = optionalUnpackDir;
            optionalUnpackDir = "";
        }
        var promise = new Promise(function(resolve,reject){
            SLFSRV.callServer("tempdir",{dir:optionalUnpackDir},function(result){
                resolve(result.dir);
            },reject);
        });
        return return_callbacks_or_promise(onSuccess,onError,promise);
    };

    /* SLFSRV.alert() */
    function showTopAlert() {
        var html, el;
        html = '<div id="SLFSRV-alert" style="background-color:rgba(0,0,0,0.4);visibility:visible;position:absolute;left:0px;top: 0px;width:100%;height:100%;xxtext-align:center;z-index:1000;">';
        html += '<div style="width:' + (window.innerWidth * 3 / 4) + 'px;margin: 100px auto;background-color: #fff;border:1px solid #000;padding:15px;xxtext-align:center;">';
        html += '<tt><pre>' + SLFSRV.esc4html(alertStack[0]._msg) + '</pre></tt>';
        html += '<div style="text-align:center">';
        html += '<button id="SLFSRV-alert-ok" onclick="SLFSRV._alert_remove()">OK</button>';
        html += '</div>';
        html += '</div>';
        html += '</div>';
        el = document.createElement('div');
        el.innerHTML = html;
        document.body.appendChild(el);
        document.getElementById("SLFSRV-alert-ok").focus();

    }
    SLFSRV._alert_remove = function() {
        var el, parent, top;
        el = document.getElementById('SLFSRV-alert').parentNode;
        parent = el.parentNode;
        parent.removeChild(el);
        top = alertStack.shift();
        if ( alertStack.length !== 0 ) {
            setTimeout(showTopAlert,100);
        }
        if ( top._callback ) { top._callback(); }
    };
    alertStack = []; // will keep ._msg and ._callback of alerts to do

    SLFSRV.alert = function(message,optionalCallbackOnOK) {
        // based on http://blog.raventools.com/create-a-modal-dialog-using-css-and-javascript/
        var obj, promise;
        obj = {_msg:String(message),_callback:optionalCallbackOnOK};
        promise = new Promise(function(resolve,reject){
            obj._callback = resolve;
            alertStack.push(obj);
            if ( alertStack.length === 1 ) {
                showTopAlert();
            }
        });

        if ( optionalCallbackOnOK === undefined ) {
            obj._id = promise;
            return promise;
        } else {
            obj._id = obj;
            promise.then(optionalCallbackOnOK,function(err){
                throw err;
            });
            return obj;
        }
    };

    SLFSRV.alertClose = function(obj) {
        var i;
        for ( i = 0; i < alertStack.length; i++ ) {
            if ( alertStack[i]._id === obj ) {
                if ( i === 0 ) {
                    SLFSRV._alert_remove();
                } else {
                    SLFSRV.splice(i,1);
                }
                break;
            }
        }
    };

    SLFSRV.callServer = function(funcName,args,onSuccess,onError,timeout) {
        var jsonArgs, handleRequest;
        jsonArgs = JSON.stringify(args);
        if ( timeout === undefined ) { timeout = 60 * 60 * 24; }

        handleRequest = function(req) {

            //log("handleRequest req =\n" + JSON.stringify(req));
            var obj = JSON.parse(req.responseText);

            if ( obj.state === 'done' ) {
                if ( onSuccess !== undefined ) {
                    onSuccess(obj);
                }
            } else if ( obj.state === 'wait' ) {
                setTimeout(function(){
                        SLFSRV.callServer("check_wait_status",{key:obj.wait_key},onSuccess,onError);
                    },0);
            } else {
                // failed
                if ( onError !== undefined ) {
                    onError(new Error(obj.message));
                } else {
                    SLFSRV.alert("Error from call to \"" + funcName + "\": " + obj.message);
                }
            }
        };

        sendRequest("/call/" + SLFSRV.SECRET_KEY + "/" + funcName + "/" + timeout,handleRequest,jsonArgs);
    };


    /* the following minimal XMLHttpRequest is taken from http://www.quirksmode.org/js/xmlhttp.html */
    sendRequest = function(url,callback,postData) {
        var req, method;
        req = createXMLHTTPObject();
        if (!req) { return; }
        method = (postData) ? "POST" : "GET";
        req.open(method,url,true);
        //req.setRequestHeader('User-Agent','XMLHTTP/1.0');
        if (postData) {
            //req.setRequestHeader('Content-type','application/x-www-form-urlencoded');
            req.setRequestHeader('Content-type','application/json');
        }
        req.onreadystatechange = function () {
            if (req.readyState !== 4) { return; }
            if (req.status !== 200 && req.status !== 304) {
                if ( -1 !== url.indexOf("/keepalive/") ) {
                    return;
                }
                req = {responseText:JSON.stringify({message:"Unspecified error " + req.status + " on call to " + url})};
                // alert('HTTP error ' + req.status);
                //return;
            }
            callback(req);
        };
        if (req.readyState === 4) { return; }
        req.send(postData);
    };

    XMLHttpFactories = [
        function () {return new XMLHttpRequest();},
        function () {return new ActiveXObject("Msxml2.XMLHTTP");},
        function () {return new ActiveXObject("Msxml3.XMLHTTP");},
        function () {return new ActiveXObject("Microsoft.XMLHTTP");}
    ];

    createXMLHTTPObject = function createXMLHTTPObject() {
        var i, xmlhttp;
        xmlhttp = false;
        for (i=0;i<XMLHttpFactories.length;i++) {
            try {
                xmlhttp = XMLHttpFactories[i]();
            }
            catch (e) {
                continue;
            }
            break;
        }
        return xmlhttp;
    };

    function keep_alive_server_forever() {

        SLFSRV.callServer("keepalive",{animal:"dog",favoriteNumber:42},function(result) {
            log("result.message = " + result.message);
        });

        setTimeout(keep_alive_server_forever,500);
    }

    function init() {
        setTimeout(keep_alive_server_forever,500);
    }

    init();

}());

/* INSERT JSON2.JS HERE */

/* INSERT NATIVE_PROMISE_ONLY.JS HERE */
`
}
