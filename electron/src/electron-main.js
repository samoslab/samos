'use strict'

const { app, Menu, BrowserWindow, dialog } = require('electron');

var log = require('electron-log');

const path = require('path');

const childProcess = require('child_process');

const cwd = require('process').cwd();

// This adds refresh and devtools console keybindings
// Page can refresh with cmd+r, ctrl+r, F5
// Devtools can be toggled with cmd+alt+i, ctrl+shift+i, F12
require('electron-debug')({enabled: true, showDevTools: false});
require('electron-context-menu')({});


global.eval = function() { throw new Error('bad!!'); }

const defaultURL = 'http://127.0.0.1:8640/';
let currentURL;

// Force everything localhost, in case of a leak
app.commandLine.appendSwitch('host-rules', 'MAP * 127.0.0.1, EXCLUDE api.github.com');
app.commandLine.appendSwitch('ssl-version-fallback-min', 'tls1.2');
app.commandLine.appendSwitch('--no-proxy-server');
app.setAsDefaultProtocolClient('samos');



// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win;

var samos = null;

function startSamos() {
  console.log('Starting samos from electron');

  if (samos) {
    console.log('Samos already running');
    app.emit('samos-ready');
    return
  }

  var reset = () => {
    samos = null;
  }

  // Resolve samos binary location
  var appPath = app.getPath('exe');
  var exe = (() => {
        switch (process.platform) {
  case 'darwin':
    return path.join(appPath, '../../Resources/app/samos');
  case 'win32':
    // Use only the relative path on windows due to short path length
    // limits
    return './resources/app/samos.exe';
  case 'linux':
    return path.join(path.dirname(appPath), './resources/app/samos');
  default:
    return './resources/app/samos';
  }
})()

  var args = [
    '-launch-browser=false',
    '-gui-dir=' + path.dirname(exe),
    '-color-log=false', // must be disabled for web interface detection
    '-logtofile=true',
    '-download-peerlist=true',
    '-enable-seed-api=true',
    '-enable-wallet-api=true',
    '-rpc-interface=false',
    "-disable-csrf=false"
    // will break
    // broken (automatically generated certs do not work):
    // '-web-interface-https=true',
  ]
  samos = childProcess.spawn(exe, args);

  samos.on('error', (e) => {
    dialog.showErrorBox('Failed to start samos', e.toString());
    app.quit();
  });

  samos.stdout.on('data', (data) => {
    console.log(data.toString());
    // Scan for the web URL string
    if (currentURL) {
      return
    }
    const marker = 'Starting web interface on ';
    var i = data.indexOf(marker);
    if (i === -1) {
      return
    }
    currentURL = defaultURL;
    app.emit('samos-ready', { url: currentURL });
  });

  samos.stderr.on('data', (data) => {
    console.log(data.toString());
  });

  samos.on('close', (code) => {
    // log.info('Samos closed');
    console.log('Samos closed');
    reset();
  });

  samos.on('exit', (code) => {
    // log.info('Samos exited');
    console.log('Samos exited');
    reset();
  });
}

function createWindow(url) {
  if (!url) {
    url = defaultURL;
  }

  // To fix appImage doesn't show icon in dock issue.
  var appPath = app.getPath('exe');
  var iconPath = (() => {
    switch (process.platform) {
      case 'linux':
        return path.join(path.dirname(appPath), './resources/icon512x512.png');
    }
  })()

  // Create the browser window.
  win = new BrowserWindow({
    width: 1200,
    height: 900,
    title: 'Samos',
    icon: iconPath,
    nodeIntegration: false,
    webPreferences: {
      webgl: false,
      webaudio: false,
    },
  });

  // patch out eval
  win.eval = global.eval;

  const ses = win.webContents.session
  ses.clearCache(function () {
    console.log('Cleared the caching of the samos wallet.');
  });

  ses.clearStorageData([],function(){
    console.log('Cleared the stored cached data');
  });

  win.loadURL(url);

  // Open the DevTools.
  // win.webContents.openDevTools();

  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;
  });

  win.webContents.on('will-navigate', function(e, url) {
    e.preventDefault();
    require('electron').shell.openExternal(url);
  });

  // create application's main menu
  var template = [{
    label: "Samos",
    submenu: [
      { label: "About Samos", selector: "orderFrontStandardAboutPanel:" },
      { type: "separator" },
      { label: "Quit", accelerator: "Command+Q", click: function() { app.quit(); } }
    ]
  }, {
    label: "Edit",
    submenu: [
      { label: "Undo", accelerator: "CmdOrCtrl+Z", selector: "undo:" },
      { label: "Redo", accelerator: "Shift+CmdOrCtrl+Z", selector: "redo:" },
      { type: "separator" },
      { label: "Cut", accelerator: "CmdOrCtrl+X", selector: "cut:" },
      { label: "Copy", accelerator: "CmdOrCtrl+C", selector: "copy:" },
      { label: "Paste", accelerator: "CmdOrCtrl+V", selector: "paste:" },
      { label: "Select All", accelerator: "CmdOrCtrl+A", selector: "selectAll:" }
    ]
  }];

  Menu.setApplicationMenu(Menu.buildFromTemplate(template));
}

// Enforce single instance
const alreadyRunning = app.makeSingleInstance((commandLine, workingDirectory) => {
      // Someone tried to run a second instance, we should focus our window.
      if (win) {
        if (win.isMinimized()) {
          win.restore();
        }
        win.focus();
      } else {
        createWindow(currentURL || defaultURL);
}
});

if (alreadyRunning) {
  app.quit();
  return;
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', startSamos);

app.on('samos-ready', (e) => {
  createWindow(e.url);
});

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
  app.quit();
}
});

app.on('activate', () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
  createWindow();
}
});

app.on('will-quit', () => {
  if (samos) {
    samos.kill('SIGINT');
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
