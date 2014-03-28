# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

module.exports = function(config) {
  config.set({
    basePath: '..',
    frameworks: ['jasmine'],

    files: [
      'js/angular.min.js',
      'tests/angular-mocks.js',
      'js/angular-animate.min.js',
      'js/angular-cookies.min.js',
      'js/angular-translate.min.js',
      'js/angular-translate-loader-static-files.min.js',
      'js/angular-translate-storage-cookie.min.js',
      'js/angular-translate-storage-local.min.js',
      'js/moment-with-langs.min.js',
      'js/shelter.js',
      'tests/statics.js',
      'tests/domain.js',
      'tests/domainform.js',
      'tests/domain-controller.js',
      'tests/domains-controller.js',
      'tests/scan.js',
      'tests/scan-controller.js',
      'tests/languages-controller.js',
      'tests/filters.js',
      'directives/*.html'
    ],

    exclude: [],

    preprocessors: {
      'directives/*.html': ['ng-html2js']
    },

    ngHtml2JsPreprocessor: {
      stripPrefix: '.*/shelter/templates/client',
      prependPrefix: '/',
      moduleName: 'directives'
    },

    reporters: ['progress'],
    port: 9876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: true,
    browsers: ['Firefox'],
    singleRun: false
  });
};
