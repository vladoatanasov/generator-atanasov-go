'use strict';
var yeoman = require('yeoman-generator');
var chalk = require('chalk');
var yosay = require('yosay');
var path = require("path");


module.exports = yeoman.generators.Base.extend({
  prompting: function () {
    var done = this.async();

    // Have Yeoman greet the user.
    this.log(yosay(
      'Welcome to the riveting ' + chalk.red('generator-atanasovgo') + ' generator!'
    ));

    var prompts = [{
      type: 'input',
      name: 'package',
      message: 'Please enter package name',
      default: process.cwd().split(path.sep).pop()
    }];

    this.prompt(prompts, function (props) {
      this.props = props;
      // To access props later use this.props.someOption;

      done();
    }.bind(this));
  },

  writing: function () {
    this.fs.copyTpl(
      this.templatePath('main.go'),
      this.destinationPath('main.go'), {
        package: this.props.package
      }
    );

    this.fs.copy(
      this.templatePath('Makefile'),
      this.destinationPath('Makefile')
    );

    this.fs.copy(
      this.templatePath('config.json.template'),
      this.destinationPath('config.json')
    );

    this.fs.copy(
      this.templatePath('rest/handler.go'),
      this.destinationPath('rest/handler.go')
    );

    this.fs.copy(
      this.templatePath('rest/httpserv.go'),
      this.destinationPath('rest/httpserv.go')
    );

    this.fs.copy(
      this.templatePath('rest/router.go'),
      this.destinationPath('rest/router.go')
    );

    this.fs.copy(
      this.templatePath('rest/server.go'),
      this.destinationPath('rest/server.go')
    );
  }
});
