
var config = document.querySelector("textarea");
var log = document.querySelector("#log");


var spans = [];

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
}

var info = function(str) {
	insert("info", str);
};

var print = function(str) {
	insert("print", str);
};

var compile = function() {
	var body = config.value;
	try {
		JSON.parse(body);
	} catch(err) {
		info(""+err);
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
//auto ping
(setInterval(function ping() {
	if(ws && ws.readyState === 1)
		ws.send("ping!");
}, 30*1000));
//auto reconnect
(function reconnect() {
	info("connecting...\n");
	ws = new WebSocket(location.origin.replace("http","ws") + "/log");

	ws.onopen = function() {
		info("connected\n");
		t = 100;
	};

	ws.onclose = function() {
		info("disconnected (retry in "+t+"ms)\n");
		setTimeout(reconnect, t);
		t *= 2;
	};

	ws.onmessage = function(e) {
		print(e.data);
	};
}());


