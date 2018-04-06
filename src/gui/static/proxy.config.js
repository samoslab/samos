const PROXY_CONFIG = {
  "/api/*": {
    "target": "http://127.0.0.1:8640",
    "pathRewrite": {
      "^/api" : ""
    },
    "secure": false,
    "logLevel": "debug",
    "bypass": function (req) {
      req.headers["host"] = '127.0.0.1:8640';
      req.headers["referer"] = 'http://127.0.0.1:8640';
      req.headers["origin"] = 'http://127.0.0.1:8640';
    }
},
  "/teller/*": {
    "target": "http://127.0.0.1:7071",
    "pathRewrite": {
      "^/teller" : "api/"
    },
    "secure": true,
    "logLevel": "debug"
  }
};

module.exports = PROXY_CONFIG;
