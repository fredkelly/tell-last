// Ionic Starter App

// angular.module is a global place for creating, registering and retrieving Angular modules
// 'starter' is the name of this angular module example (also set in a <body> attribute in index.html)
// the 2nd parameter is an array of 'requires'
angular.module('tellLast', ['ionic'])

.run(function($ionicPlatform) {
  $ionicPlatform.ready(function() {
    // Hide the accessory bar by default (remove this to show the accessory bar above the keyboard
    // for form inputs)
    if(window.cordova && window.cordova.plugins.Keyboard) {
      cordova.plugins.Keyboard.hideKeyboardAccessoryBar(true);
    }
    if(window.StatusBar) {
      StatusBar.styleDefault();
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
