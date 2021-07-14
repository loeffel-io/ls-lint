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
      switch (process.arch) {
        case 'x64':
          return 'ls-lint-darwin';
        case 'arm64':
          return 'ls-lint-darwin-arm64';
        default:
          console.log('ls-lint builds are not available on platform: ' + process.platform + ' arch: ' + process.arch);
          process.exit(1);
      }
      return 'ls-lint-darwin';
    case 'linux':
      switch (process.arch) {
        case 'x64':
          return 'ls-lint-linux';
        case 'arm64':
          return 'ls-lint-linux-arm64';
        default:
          console.log('ls-lint builds are not available on platform: ' + process.platform + ' arch: ' + process.arch);
          process.exit(1);
      }
    case 'win32':
      return 'ls-lint-windows.exe';
    default:
      console.log('ls-lint builds are not available on platform: ' + process.platform)
      process.exit(1);
  }
}

spawnCommand(getPlatformPath());
