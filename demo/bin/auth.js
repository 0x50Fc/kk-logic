// kk-cli logic


var app = {
	"in":new kk.Logic.Var({
		"key":"output.input",
		"value":"=input",
		"ondone":"done",
	}),
	"done":new kk.Logic.Outlet({
		"title":"验证成功",
	}),
};

module.exports = new kk.App(app);

