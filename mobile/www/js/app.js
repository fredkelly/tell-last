// Ionic Starter App

// angular.module is a global place for creating, registering and retrieving Angular modules
// 'starter' is the name of this angular module example (also set in a <body> attribute in index.html)
// the 2nd parameter is an array of 'requires'
angular.module('tellLast', ['ionic', 'ngResource', 'ngFacebook'])

.constant('API_BASE', 'http://localhost:3000')

.config(function($stateProvider, $urlRouterProvider, $facebookProvider) {
  $facebookProvider
    .setAppId('1390390824592299')
    .setCustomInit({ version: 'v1.0' })
    .setPermissions("email, user_friends")
    //.setVersion('v2.1');

  $urlRouterProvider.otherwise('/')

  /*
  $stateProvider.state('login', {
    url: '/login',
    views: {
      login: {
        templateUrl: 'login.html',
        //controller: 'LoginCtrl'
      }
    },
  })
  */


  $stateProvider.state('inbox', {
    url: '/',
    views: {
      inbox: {
        templateUrl: 'inbox.html',
        controller: 'InboxCtrl'
      }
    },
    requireLogin: true,
  })

  $stateProvider.state('tell', {
    url: '/tell',
    views: {
      tell: {
        templateUrl: 'tell.html',
        controller: 'TellCtrl'
      }
    },
    requireLogin: true,
  })
})

.factory('OAuthHttpInterceptor', function($rootScope) {
  return {
    request: function(config) {
      if ($rootScope.authResponse && !config.headers.Authorization) {
        config.headers.Authorization = $rootScope.authResponse.accessToken;
      }
      return config;
    }
  };
})

.config(function($httpProvider) {
  $httpProvider.interceptors.push('OAuthHttpInterceptor');
})

.run(function($ionicPlatform, $rootScope, $facebook, $state, $http) {
  $ionicPlatform.ready(function() {
    // Hide the accessory bar by default (remove this to show the accessory bar above the keyboard
    // for form inputs)
    if (window.cordova && window.cordova.plugins.Keyboard) {
      cordova.plugins.Keyboard.hideKeyboardAccessoryBar(true);
    }
    if (window.StatusBar) {
      StatusBar.styleDefault();
    }
  });

  // Facebook SDK
  (function(d, s, id){
     var js, fjs = d.getElementsByTagName(s)[0];
     if (d.getElementById(id)) {return;}
     js = d.createElement(s); js.id = id;
     js.src = "//connect.facebook.net/en_US/sdk.js";
     fjs.parentNode.insertBefore(js, fjs);
   }(document, 'script', 'facebook-jssdk'));

  $rootScope.$on('$stateChangeStart', function (event, toState, toParams) {
    if (toState.requireLogin && !$rootScope.authResponse) {
      event.preventDefault();
      $facebook.getLoginStatus().then(function(response) {
        if (response.authResponse) {
          $rootScope.authResponse = response.authResponse;
          return $state.go(toState.name, toParams);
        } else {
          // show login prompt
          $facebook.login();
        }
      },
      function(error) {
        console.log(error);
      });
    }
  });
})

.factory('Tell', function($resource, API_BASE) {
  return $resource(API_BASE + '/tells');
})

.controller('InboxCtrl', function($scope, Tell) {
  $scope.tells = Tell.query();
})

.controller('TellCtrl', function($scope, $facebook, Tell) {
  $scope.friends = [];
  $scope.tell = new Tell();

  $facebook.api('/me/friends').then(function(response) {
    $scope.friends = response.data;
  }, function(error) {
    // TODO
    console.log(error);
  });

  $scope.submit = function() {
    $scope.tell.$save();
  };
})
