

var ctx = new kk.Context();

var input = {};

if (typeof _GET == 'object') {
    for (var key in _GET) {
        input[key] = _GET[key];
    }
}

if (typeof _POST == 'object') {
    for (var key in _POST) {
        input[key] = _POST[key];
    }
}

ctx.set(["input"], input);
ctx.set(["output"], {});
ctx.set(["userAgent"], _HEADER['User-Agent']);

var app = new kk.App();

app.set("in", new kk.Logic.Var(
    {
        key: 'output.version',
        value: '1.0',
        ondone: new kk.Logic.Var({
            key: 'output.body',
            value: new kk.Logic.Http({
                url: 'http://www.baidu.com',
                dataType: 'text',
                headers : {
                    'User-Agent' : '=userAgent'
                }
            })
        })
    }
));

app.exec(ctx);

header("Content-Type", "application/json; charset=utf-8");

echo(JSON.stringify(ctx.get(["output"])));


