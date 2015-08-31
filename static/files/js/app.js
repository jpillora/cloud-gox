
var app = angular.module('cloudgox', []);

app.controller("AppController", function($scope, $http) {
	window.cloudgox = $scope;
	$scope.package = {
		name: "github.com/jpillora/serve",
		version: "",
		versionVar: "VERSION",
		platforms: null,
		commitish: "",
	}

	//pull server config
	$http.get("/config").success(function(config) {
		$scope.config = config;
		$scope.package.platforms = angular.copy(config.Platforms);
	});

	//link up to angular
	var rt = realtime("/realtime");
	$scope.state = {};
	rt.add("state", $scope.state, function() {
		computeLog();
		$scope.$apply();
	});
	rt.onstatus = function(online) {
		$scope.$apply(function() {
			$scope.Connected = online;
		});
	};

	$scope.compile = function(foo) {
		var data = angular.copy($scope.package);
		if(!data.version) data.version = "1.0.0";
		$http.post("/compile", data).then(function(resp) {
			// console.info("success", resp);
		}, function(resp) {
			console.warn("failed", resp);
		});
	};

	$scope.ago = function(t) {
		return moment(t).fromNow();
	};

	$scope.ready = function() {
		return $scope.Connected && $scope.state.Ready;
	};

	$scope.osarch = function(file) {
		return /([a-z]+)_([a-z0-9]+)(\.exe)?\.gz$/.test(file) ? (RegExp.$1+"/"+RegExp.$2) : file;
	};

	$scope.compilationsEmpty = true;
	$scope.compilations = function() {
		var c = [];
		if($scope.state.Current)
			c.push($scope.state.Current);
		if($scope.state.Done)
			c = c.concat($scope.state.Done);
		c.sort(function(a,b) {
			if(a.completedAt) return -1;
			return a.completedAt > b.completedAt ? -1 : 1;
		});
		$scope.compilationsEmpty = c.length === 0;
		return c;
	};

	var logElem = angular.element(document.querySelector(".log"));
	var state = $scope.state;

	$scope.groups = [];

	var computeLog = function() {
		//remove the logs off the front and back
		var elem;
		var i = state.LogOffset-1;
		while(elem = document.querySelector("#log" + (i--)))
			elem.remove();
		i = state.LogCount+1;
		while(elem = document.querySelector("#log" + (i++)))
			elem.remove();

		//render new logs
		for (var i = state.LogOffset; i <= state.LogCount; i++) {
			var l = state.Log[i];
			if(l.$rendered) {
				continue;
			}
			var group = null;
			var p = state.Log[i-1];
			if(p && p.src === l.src) {
				group = p.$group;
			} else {
				var div = document.createElement("div");
				var user = l.src !== "cloud-gox";
				div.className = "group " + (user ? "user shortened" : "cloudgox");
				div.setAttribute("src", l.src);
				if(user) angular.element(div).on("click", function() {
					angular.element(this).toggleClass("shortened");
				});
				group = div;
				logElem.prepend(div);
			}

			l.$group = group;

			var span = document.createElement("span");
			span.id = "log"+i;
			span.className = "msg " + l.type;
			l.$span = span;

			var html = l.msg.split("\n").filter(function(l){return !!l;}).reverse().join("</br>");
			span.innerHTML = html + "</br>";

			var timestamp = document.createElement("span");
			timestamp.className = "timestamp";
			timestamp.innerHTML = moment(l.t).format("YYYY/MM/DD hh:mm:ss") + "&nbsp;";
			angular.element(span).prepend(timestamp);

			angular.element(group).prepend(span);
			l.$rendered = true;
		}
	};
});
