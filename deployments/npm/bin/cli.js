#!/usr/bin/env node

const spawn = require('child_process').spawn;
const path = require('path');
const command_args = process.argv.slice(2);

function spawnCommand(binaryExecutable) {
    const child = spawn(
        path.join(__dirname, binaryExecutable),
        command_args,
        {stdio: [process.stdin, process.stdout, process.stderr]}
    );

    child.on('close', function (code) {
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
                    return 'ls-lint-darwin-amd64';
                case 'arm64':
                    return 'ls-lint-darwin-arm64';
                default:
                    console.log('ls-lint builds are not available on platform: ' + process.platform + ' arch: ' + process.arch);
                    process.exit(1);
            }
            break
        case 'linux':
            switch (process.arch) {
                case 'x64':
                    return 'ls-lint-linux-amd64';
                case 'arm64':
                    return 'ls-lint-linux-arm64';
                case 's390x':
                    return 'ls-lint-linux-s390x';
                case 'ppc64le':
                    return 'ls-lint-linux-ppc64le';
                default:
                    console.log('ls-lint builds are not available on platform: ' + process.platform + ' arch: ' + process.arch);
                    process.exit(1);
            }
            break
        case 'win32':
            return 'ls-lint-windows-amd64.exe';
        default:
            console.log('ls-lint builds are not available on platform: ' + process.platform)
            process.exit(1);
    }
}

spawnCommand(getPlatformPath());
