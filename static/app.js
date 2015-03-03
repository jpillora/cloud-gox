
var config = document.querySelector("textarea");
var log = document.querySelector("#log");

var insert = function(cls, s) {
	s = s.trim();
	if(!s) return;
	var span = document.createElement("span");
	span.className = cls;
	span.innerText = s + "\n";
	log.appendChild(span);
}

var info = function(s) {
	insert("info", s);
};

var print = function(s) {
	insert("print", s);
};

var compile = function() {
	var body = config.value;
	try {
		JSON.parse(body);
	} catch(err) {
		info(err);
		return;
	}

	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/compile");
	xhr.onload = function() {
		if(xhr.status !== 200) {
			info("compilation queue error: " + xhr.responseText);
		}
	};
	xhr.send(body);
};

//====================

var t;
var ws;
(function reconnect() {
	info("connecting...");
	ws = new WebSocket(location.origin.replace("http","ws") + "/log");

	ws.onopen = function() {
		info("connected");
		t = 100;
	};

	ws.onclose = function() {
		info("disconnected (retry in "+t+"ms)");
		setTimeout(reconnect, t);
		t *= 2;
	};

	ws.onmessage = function(e) {
		console.log("'%s'", e.data);
		print(e.data);
	};
}());


