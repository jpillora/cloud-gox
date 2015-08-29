/* global window,WebSocket */
// Realtime - v0.1.0 - https://github.com/jpillora/go-realtime
// Jaime Pillora <dev@jpillora.com> - MIT Copyright 2015
(function(window, document) {

  if(!window.WebSocket)
    return alert("This browser does not support WebSockets");

  //realtime protocol version
  var proto = "v1";

  //public method
  var realtime = function(url) {
    var rt = new Realtime(url);
    rts.push(rt);
    return rt;
  };
  realtime.proto = proto;
  realtime.online = true;

  //special merge - ignore $properties
  // x <- y
  var merge = function(x, y) {
    if (!x || typeof x !== "object" ||
      !y || typeof y !== "object")
      return y;
    var k;
    if (x instanceof Array && y instanceof Array)
      while (x.length > y.length)
        x.pop();
    else
      for (k in x)
        if (k[0] !== "$" && !(k in y))
          delete x[k];
    for (k in y)
      x[k] = merge(x[k], y[k]);
    return x;
  };

  var rts = [];
  //global status change handler
  function onstatus(event) {
    realtime.online = navigator.onLine;
    for(var i = 0; i < rts.length; i++) {
      if(realtime.online && rts[i].autoretry)
        rts[i].retry();
    }
  }
  window.addEventListener('online',  onstatus);
  window.addEventListener('offline', onstatus);

  //helpers
  var events = ["message","error","open","close"];
  var loc = window.location;
  //Realtime class - represents a single websocket (User on the server-side)
  function Realtime(url) {
    if(!url)
      url = "/realtime";
    if(!(/^https?:/.test(url)))
      url = loc.protocol + "//" + loc.host + url;
    if(!(/^http(s?:\/\/.+)$/.test(url)))
      throw "Invalid URL: " + url;
    this.url = "ws" + RegExp.$1;
    this.connect();
    this.objs = {};
    this.subs = {};
    this.onupdates = {};
    this.connected = false;
  }
  Realtime.prototype = {
    add: function(key, obj, onupdate) {
      if(typeof key !== "string")
        throw "Invalid key - must be string";
      if(!obj || typeof obj !== "object")
        throw "Invalid object - must be an object";
      if(this.objs[key])
        throw "Duplicate key - already added";
      this.objs[key] = obj;
      this.subs[key] = 0;
      this.onupdates[key] = onupdate;
      this.subscribe();
    },
    connect: function() {
      this.autoretry = true;
      this.retry();
    },
    retry: function() {
      clearTimeout(this.retry.t);
      if(this.ws)
        this.cleanup();
      if(!this.delay)
        this.delay = 100;
      this.ws = new WebSocket(this.url, "rt-"+proto);
      var _this = this;
      events.forEach(function(e) {
        e = "on"+e;
        _this.ws[e] = _this[e].bind(_this);
      });
      this.ping.t = setInterval(this.ping.bind(this), 30 * 1000);
    },
    disconnect: function() {
      this.autoretry = false;
      this.cleanup();
    },
    cleanup: function(){
      if(!this.ws)
        return;
      var _this = this;
      events.forEach(function(e) {
        _this.ws["on"+e] = null;
      });
      if(this.ws.readyState !== WebSocket.CLOSED)
        this.ws.close();
      this.ws = null;
      clearInterval(this.ping.t);
    },
    send: function(data) {
      if(this.ws.readyState === WebSocket.OPEN)
        return this.ws.send(data);
    },
    ping: function() {
      this.send("ping");
    },
    subscribe: function() {
      this.send(JSON.stringify(this.subs));
    },
    onmessage: function(event) {
      var str = event.data;
      if (str === "ping") return;

      var updates;
      try {
        updates = JSON.parse(str);
      } catch(err) {
        return console.warn(err, str);
      }

      for(var i = 0; i < updates.length; i++) {
        var u = updates[i];
        var key = u.Key;
        var dst = this.objs[key];
        var src = u.Data;
        if(!src || !dst)
          continue;

        if(u.Delta)
          jsonpatch.apply(dst, src);
        else
          merge(dst, src);

        if(typeof dst.$apply === "function")
          dst.$apply();

        var onupdate = this.onupdates[key];
        if(typeof onupdate === "function")
          onupdate();

        this.subs[key] = u.Version;
      }
      //successful msg resets retry counter
      this.delay = 100;
    },
    onopen: function() {
      this.connected = true;
      if(this.onstatus) this.onstatus(true);
      this.subscribe();
    },
    onclose: function() {
      this.connected = false;
      if(this.onstatus) this.onstatus(false);
      this.delay *= 2;
      if(this.autoretry) {
        this.retry.t = setTimeout(this.connect.bind(this), this.delay);
      }
    },
    onerror: function(err) {
      // console.error("websocket error: %s", err);
    }
  };
  //publicise
  window.realtime = realtime;
}(window, document, undefined));

/*!
* https://github.com/Starcounter-Jack/JSON-Patch
* json-patch-duplex.js version: 0.5.2
* (c) 2013 Joachim Wester
* MIT license
*/
var __extends=this.__extends||function(a,b){function c(){this.constructor=a}for(var d in b)b.hasOwnProperty(d)&&(a[d]=b[d]);c.prototype=b.prototype,a.prototype=new c},OriginalError=Error,jsonpatch;!function(a){function b(a,c){switch(typeof a){case"undefined":case"boolean":case"string":case"number":return a===c;case"object":if(null===a)return null===c;if(g(a)){if(!g(c)||a.length!==c.length)return!1;for(var d=0,e=a.length;e>d;d++)if(!b(a[d],c[d]))return!1;return!0}var f=h(c),i=f.length;if(h(a).length!==i)return!1;for(var d=0;i>d;d++)if(!b(a[d],c[d]))return!1;return!0;default:return!1}}function c(a){for(var b,c=0,d=a.length;d>c;){b=a.charCodeAt(c);{if(!(b>=48&&57>=b))return!1;c++}}return!0}function d(a,b,d){for(var e,f,h=!1,m=0,n=b.length;n>m;){e=b[m],m++;for(var o=e.path.split("/"),p=a,q=1,r=o.length,s=void 0;;){if(f=o[q],d&&void 0===s&&(void 0===p[f]?s=o.slice(0,q).join("/"):q==r-1&&(s=e.path),void 0!==s&&this.validator(e,m-1,a,s)),q++,void 0===f&&q>=r){h=k[e.op].call(e,p,f,a);break}if(g(p)){if("-"===f)f=p.length;else{if(d&&!c(f))throw new l("Expected an unsigned base-10 integer value, making the new referenced value the array element with the zero-based index","OPERATION_PATH_ILLEGAL_ARRAY_INDEX",m-1,e.path,e);f=parseInt(f,10)}if(q>=r){if(d&&"add"===e.op&&f>p.length)throw new l("The specified index MUST NOT be greater than the number of elements in the array","OPERATION_VALUE_OUT_OF_BOUNDS",m-1,e.path,e);h=j[e.op].call(e,p,f,a);break}}else if(f&&-1!=f.indexOf("~")&&(f=f.replace(/~1/g,"/").replace(/~0/g,"~")),q>=r){h=i[e.op].call(e,p,f,a);break}p=p[f]}}return h}function e(b,c,d,e){if("object"!=typeof b||null===b||g(b))throw new l("Operation is not an object","OPERATION_NOT_AN_OBJECT",c,b,d);if(!i[b.op])throw new l("Operation `op` property is not one of operations defined in RFC-6902","OPERATION_OP_INVALID",c,b,d);if("string"!=typeof b.path)throw new l("Operation `path` property is not a string","OPERATION_PATH_INVALID",c,b,d);if(("move"===b.op||"copy"===b.op)&&"string"!=typeof b.from)throw new l("Operation `from` property is not present (applicable in `move` and `copy` operations)","OPERATION_FROM_REQUIRED",c,b,d);if(("add"===b.op||"replace"===b.op||"test"===b.op)&&void 0===b.value)throw new l("Operation `value` property is not present (applicable in `add`, `replace` and `test` operations)","OPERATION_VALUE_REQUIRED",c,b,d);if(d)if("add"==b.op){var f=b.path.split("/").length,h=e.split("/").length;if(f!==h+1&&f!==h)throw new l("Cannot perform an `add` operation at the desired path","OPERATION_PATH_CANNOT_ADD",c,b,d)}else if("replace"===b.op||"remove"===b.op||"_get"===b.op){if(b.path!==e)throw new l("Cannot perform the operation at a path that does not exist","OPERATION_PATH_UNRESOLVABLE",c,b,d)}else if("move"===b.op||"copy"===b.op){var j={op:"_get",path:b.from,value:void 0},k=a.validate([j],d);if(k&&"OPERATION_PATH_UNRESOLVABLE"===k.name)throw new l("Cannot perform the operation from a path that does not exist","OPERATION_FROM_UNRESOLVABLE",c,b,d)}}function f(a,b){try{if(!g(a))throw new l("Patch sequence must be an array","SEQUENCE_NOT_AN_ARRAY");if(b)b=JSON.parse(JSON.stringify(b)),d.call(this,b,a,!0);else for(var c=0;c<a.length;c++)this.validator(a[c],c)}catch(e){if(e instanceof l)return e;throw e}}if(!a.apply){var g,h=function(){return Object.keys?Object.keys:function(a){var b=[];for(var c in a)a.hasOwnProperty(c)&&b.push(c);return b}}(),i={add:function(a,b){return a[b]=this.value,!0},remove:function(a,b){return delete a[b],!0},replace:function(a,b){return a[b]=this.value,!0},move:function(a,b,c){var e={op:"_get",path:this.from};return d(c,[e]),d(c,[{op:"remove",path:this.from}]),d(c,[{op:"add",path:this.path,value:e.value}]),!0},copy:function(a,b,c){var e={op:"_get",path:this.from};return d(c,[e]),d(c,[{op:"add",path:this.path,value:e.value}]),!0},test:function(a,c){return b(a[c],this.value)},_get:function(a,b){this.value=a[b]}},j={add:function(a,b){return a.splice(b,0,this.value),!0},remove:function(a,b){return a.splice(b,1),!0},replace:function(a,b){return a[b]=this.value,!0},move:i.move,copy:i.copy,test:i.test,_get:i._get},k={add:function(a){k.remove.call(this,a);for(var b in this.value)this.value.hasOwnProperty(b)&&(a[b]=this.value[b]);return!0},remove:function(a){for(var b in a)a.hasOwnProperty(b)&&i.remove.call(this,a,b);return!0},replace:function(a){return d(a,[{op:"remove",path:this.path}]),d(a,[{op:"add",path:this.path,value:this.value}]),!0},move:i.move,copy:i.copy,test:function(a){return JSON.stringify(a)===JSON.stringify(this.value)},_get:function(a){this.value=a}};g=Array.isArray?Array.isArray:function(a){return a.push&&"number"==typeof a.length},a.apply=d;var l=function(a){function b(b,c,d,e,f){a.call(this,b),this.message=b,this.name=c,this.index=d,this.operation=e,this.tree=f}return __extends(b,a),b}(OriginalError);a.JsonPatchError=l,a.Error=l,a.validator=e,a.validate=f}}(jsonpatch||(jsonpatch={})),"undefined"!=typeof exports&&(exports.apply=jsonpatch.apply,exports.validate=jsonpatch.validate,exports.validator=jsonpatch.validator,exports.JsonPatchError=jsonpatch.JsonPatchError,exports.Error=jsonpatch.Error);
