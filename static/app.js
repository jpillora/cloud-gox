
var config = document.querySelector("textarea");
var log = document.querySelector("#log");

var seq = 0;
var spans = [];

var compile = function() {
	var body = config.value;
	try {
		JSON.parse(body);
	} catch(err) {
		alert(""+err);
		return;
	}

	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/compile");
	xhr.onload = function() {
		if(xhr.status !== 200) {
			alert(xhr.responseText);
		}
	};
	xhr.send(body);
};


//====================

function status(type, msg, color) {
	var elem = document.querySelector("."+type+".status");
	if(!elem) return console.warn("no ", type);
	elem.innerText = msg;
	if(color)
		elem.setAttribute("color", color);
};

//=========================

var ongoxevent = function(e) {
	if(e.msg) {
		onmessage(e.msg);
	} else if(e.sts) {
		onstatus(e.sts);
	}
}

var onmessage = function(msg) {
	seq = Math.max(seq, msg.id);
	print(msg.txt);
};

var onstatus = function(s) {
	if(s.current)
		status("current", s.current.package, "green");
	else
		status("current", "not compiling", "grey");

	status("queue", s.numQueued + " items queued",
		s.numQueued ? "blue" : "grey");

	status("completed", s.numDone + " items completed",
		s.numDone ? "blue" : "grey");
};

var info = function(str) {
	insert("info", str);
};

var print = function(str) {
	insert("print", str);
};

var insert = function(cls, str) {
	var lines = str.split("\n").reverse();
	insertLines(cls, lines);
};

var insertLines = function(cls, lines) {
	for(var i = 0; i < lines.length; i++) {
		var l = lines[i];

		var top = spans[0];
		if(top && i === 0) {
			top.innerText += l;
			continue
		}

		var span = document.createElement("div");
		span.className = cls;
		span.innerText = l;

		log.insertBefore(span, top);
		spans.unshift(span);
	}
	while(spans.length > 1000) {
		log.removeChild(span.pop());
	}
};

//====================

var t;
var ws;
var seq = 0;
//auto ping
(setInterval(function ping() {
	if(ws && ws.readyState === 1)
		ws.send("ping!");
}, 30*1000));
//auto reconnect
(function reconnect() {
	status("connection", "connecting", "blue");
	ws = new WebSocket(location.origin.replace("http","ws") + "/log", ""+seq);

	ws.onopen = function() {
		status("connection", "connected", "green");
		t = 100;
	};

	ws.onclose = function() {
		var s = Math.round(t/100)/10;
		status("connection", "disconnected (retry in "+s+"s)", "red");
		setTimeout(reconnect, t);
		t *= 2;
	};

	ws.onmessage = function(msg) {
		var e;
		try {
			e = JSON.parse(msg.data);
		} catch(e) {
			return;
		}
		if(!e)
			return;
		if(e instanceof Array)
			e.forEach(ongoxevent);
		else
			ongoxevent(e);
	};
}());


