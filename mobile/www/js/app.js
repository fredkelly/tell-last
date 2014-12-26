// Ionic Starter App

// angular.module is a global place for creating, registering and retrieving Angular modules
// 'starter' is the name of this angular module example (also set in a <body> attribute in index.html)
// the 2nd parameter is an array of 'requires'
angular.module('tellLast', ['ionic', 'ngFacebook'])

.config(function($stateProvider, $urlRouterProvider, $facebookProvider) {
  $facebookProvider
    .setAppId('1390390824592299')
    .setCustomInit({ version: 'v1.0' })
    //.setVersion('v2.1');

  $urlRouterProvider.otherwise('/')

  $stateProvider.state('inbox', {
    url: '/',
    views: {
      inbox: { templateUrl: 'inbox.html' }
    }
  })

  $stateProvider.state('new', {
    url: '/new',
    views: {
      new: { templateUrl: 'new.html' }
    },
    requireLogin: true
  })
})

.run(function($ionicPlatform, $rootScope, $facebook) {
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
    if (toState.requireLogin && typeof $rootScope.currentUser === 'undefined') {
      event.preventDefault();
      var response = $facebook.getAuthResponse();
      console.log(response);
      //$facebook.api('/me').then(function(user) {
      //  console.log(user);
      //},
      //function(error) {
      //  console.log(error);
      //});
    }
  });
})

.controller('TellCtrl', function($scope) {
  $scope.friends = [
    { id: 1, name: 'Joe Bloggs' },
    { id: 2, name: 'Dave Smith' }
  ];

  $scope.body = '';
  $scope.from = $scope.friends[0];
  $scope.to = $scope.friends[1];
})

.controller('InboxCtrl', function($scope) {
  $scope.items = [
    { created: new Date(), title: 'Test item', body: 'Some lovely thoughts about said person.', uuid: '656691226' },
    { created: new Date(), title: 'Secret Complement', body: 'More lovely thoughts about said person.', uuid: '590344115'}
  ];
})
