var app = angular.module("workflow", ["ngRoute"]);

app.config(function($routeProvider) {
  $routeProvider
  .when("/getdm", {
    templateUrl : "tmpl/dm.tmpl",
	controller : "dm"
  })
  .when("/getparams", {
    templateUrl : "tmpl/params.tmpl",
    controller : "params"
  })
  .when("/getinstances", {
    templateUrl : "tmpl/instances.tmpl",
    controller : "instances"
  })
});

app.controller("dm", function ($scope, $http) {
  console.log("Inside of getdm controller")
  $http({ method : "GET", url : "getdmjson"})
    .then(function mySuccess(response) {
      $scope.Objs = response.data;
    }, function myError(response) {
	  $scope.msg = response.statusText;
  });
});

app.controller("params", function ($scope, $http) {
  console.log("Inside of getparams controller");
  $http({ method : "GET", url : "getparamsjson"}).then(function mySuccess(response) {
		  $scope.Params = response.data;
    }, function myError(response) {
	  $scope.msg = response.statusText;
  });
  $http({ method : "GET", url : "geteps"}).then(function mySuccess(response) {
	  $scope.Names = ["Emil", "Tobias", "Linus"];
  }, function myError(response) {
	  $scope.msg = response.statusText;
  });
});

app.controller("instances", function ($scope, $http) {
  console.log("Inside of instances controller")
  $http({ method : "GET", url : "getinstancesjson"})
    .then(function mySuccess(response) {
      $scope.Instances = response.data;
    }, function myError(response) {
	  $scope.msg = response.statusText;
  });
});
