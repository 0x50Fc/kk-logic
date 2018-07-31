// kk-cli logic


kk.run({
	"in":new kk.Logic.Var({
			"key":"output.version",
			"value":"1.0",
			"ondone":new kk.Logic.Http({
				"url":"http://www.baidu.com",
				"dataType":"text",
				"method":"GET",
				"ondone":new kk.Logic.Var({
					"value":"=result",
					"key":"output.body",
				}),
			}),
		}),
});

