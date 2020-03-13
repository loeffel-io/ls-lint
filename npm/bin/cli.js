#!/usr/bin/env node

var spawn = require('child_process').spawn;
var path = require('path');

var command_args = process.argv.slice(2);

function spawnCommand(binaryExecutable) {
  var child = spawn(
    path.join(__dirname, binaryExecutable),
    command_args,
    { stdio: [process.stdin, process.stdout, process.stderr] }
  );

  child.on('close', function(code) {
    if (code !== 0) {
      process.exit(1);
    }
  });
}

function getPlatformPath() {
  switch (process.platform) {
    case 'darwin':
      return 'ls-lint-darwin';
    case 'linux':
      return 'ls-lint-linux';
    case 'win32':
      return 'ls-lint-windows.exe';
    default:
      console.log('ls-lint builds are not available on platform: ' + process.platform)
      process.exit(1);
  }
}

spawnCommand(getPlatformPath());
