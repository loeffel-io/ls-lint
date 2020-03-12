const { spawn } = require("child_process");
const { join } = require("path");

const isCI = process.env.CI;

if (!isCI) {
  if (process.platform === 'darwin') {
    binary = 'ls-lint-darwin';
  } else if (process.platform === 'linux') {
    binary = 'ls-lint-linux';
  } else if (process.platform === 'win32') {
    binary = 'ls-lint-windows.exe';
  } else {
    console.log("Unsupported OS");
    process.exit(0);
  }

  binpath = join(__dirname, 'bin', binary);

  result = spawn(binpath, ["install", "-f"], {
      cwd: process.env.INIT_CWD,
      stdio: [process.stdin, process.stdout, process.stderr]
  });
}
