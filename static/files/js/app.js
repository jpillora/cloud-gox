var app = angular.module("cloudgox", []);

app.controller("AppController", function($scope, $http) {
  //TODO platform select
  // var OSName="Unknown OS";
  // if (navigator.appVersion.indexOf("Win")!=-1) OSName="Windows";
  // if (navigator.appVersion.indexOf("Mac")!=-1) OSName="MacOS";
  // if (navigator.appVersion.indexOf("X11")!=-1) OSName="UNIX";
  // if (navigator.appVersion.indexOf("Linux")!=-1) OSName="Linux";
  window.cloudgox = $scope;
  $scope.package = {
    name: "github.com/jpillora/serve",
    version: "",
    versionVar: "main.VERSION",
    commitVar: "main.COMMIT",
    platforms: null,
    commitish: "",
    cgo: true,
    shrink: true,
    goGet: true
  };

  //pull server config
  $http.get("/config").success(function(config) {
    $scope.config = config;
    $scope.package.platforms = angular.copy(config.Platforms);
  });

  //link up to angular
  $scope.state = {};
  var v = velox("/sync", $scope.state);
  v.onupdate = function() {
    computeLog();
    $scope.$applyAsync();
  };
  v.onchange = function(connected) {
    $scope.$apply(function() {
      $scope.Connected = connected;
    });
  };

  $scope.compile = function(foo) {
    var data = angular.copy($scope.package);
    if (!data.version) data.version = "1.0.0";
    $scope.loading = true;
    $http
      .post("/compile", data)
      .then(
        function(resp) {
          // console.info("success", resp);
        },
        function(resp) {
          console.warn("failed", resp);
        }
      )
      .finally(function() {
        $scope.loading = false;
      });
  };

  $scope.ago = function(t) {
    return moment(t).fromNow();
  };

  $scope.ready = function() {
    return $scope.Connected && $scope.state.Ready;
  };

  $scope.osarch = function(file) {
    return /([a-z]+)_([a-z0-9]+)(\.exe)?\.gz$/.test(file)
      ? RegExp.$1 + "/" + RegExp.$2
      : file;
  };

  $scope.removeExt = function(file) {
    return file.replace(/\.gz/, "");
  };

  $scope.compilationsEmpty = true;
  $scope.compilations = function() {
    var c = [];
    if ($scope.state.Current) c.push($scope.state.Current);
    if ($scope.state.Done) c = c.concat($scope.state.Done);
    c.sort(function(a, b) {
      if (a.completedAt) return -1;
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
    var first = state.LogOffset;
    var last = first + state.LogCount;
    var i = first - 1;
    while ((elem = document.querySelector("#log" + i--))) elem.remove();
    i = last + 1;
    while ((elem = document.querySelector("#log" + i++))) elem.remove();

    //render new logs
    for (var i = state.LogOffset; i <= last; i++) {
      var l = state.Log[i];
      if (!l || l.$rendered) {
        continue;
      }
      var group = null;
      var p = state.Log[i - 1];
      if (p && p.src === l.src) {
        group = p.$group;
      } else {
        var div = document.createElement("div");
        var user = l.src !== "cloud-gox";
        div.className = "group " + (user ? "user shortened" : "cloudgox");
        div.setAttribute("src", l.src);
        if (user)
          angular.element(div).on("click", function() {
            angular.element(this).toggleClass("shortened");
          });
        group = div;
        logElem.prepend(div);
      }

      l.$group = group;

      var span = document.createElement("span");
      span.id = "log" + i;
      span.className = "msg " + l.type;
      l.$span = span;

      var html = l.msg
        .split("\n")
        .filter(function(l) {
          return !!l;
        })
        .reverse()
        .join("</br>");
      span.innerHTML = html + "</br>";

      var timestamp = document.createElement("span");
      timestamp.className = "timestamp";
      timestamp.innerHTML =
        moment(l.t).format("YYYY/MM/DD hh:mm:ss") + "&nbsp;";
      angular.element(span).prepend(timestamp);

      angular.element(group).prepend(span);
      l.$rendered = true;
    }
  };
});

app.controller("PkgURLController", function($scope) {
  var compilation = $scope.c;
  var urls = ($scope.urls = {});
  if (!/^([^\/]+\/[^\/]+\/[^\/]+)(\/(.+))?/.test(compilation.name)) {
    return;
  }
  var pkg = RegExp.$1;
  var target = RegExp.$3;
  var gh = /^github\.com/.test(pkg);
  var ghtree = compilation.commitish ? "/tree/" + compilation.commitish : "";
  urls.repo = "http://" + pkg;
  urls.repoName = pkg;
  if (gh) {
    urls.repo += ghtree;
  }
  if (target) {
    if (gh) {
      urls.target = "http://" + pkg + (ghtree || "/tree/master") + "/" + target;
    }
    urls.targetName = target;
  }
});
