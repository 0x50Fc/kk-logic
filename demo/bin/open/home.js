// kk-cli logic


var app = {
	"in":new kk.Logic.App({
		"path":"../auth.js",
		"ondone":new kk.Logic.Var({
			"key":"output.version",
			"value":"1.0",
			"ondone":new kk.Logic.Http({
				"method":"GET",
				"ondone":new kk.Logic.Var({
					"key":"body",
					"value":"=result",
				}),
				"url":"http://www.baidu.com",
				"dataType":"text",
			}),
		}),
	}),
};

kk.run(app);

