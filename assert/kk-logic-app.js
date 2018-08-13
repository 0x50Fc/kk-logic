
kk.Logic.App = function () {
    kk.Logic.apply(this, arguments);
};

kk.Logic.App.prototype = kk.extend(kk.Logic.prototype, {

    exec: function (ctx, app) {
        kk.Logic.prototype.exec.apply(this, arguments);

        if (this.app === undefined) {
            var path = this.get(ctx, app, ["path"]);
            this.app = require(path);
            if (this.app instanceof kk.App) {
                for (var name in this.app.logics) {

                    var logic = this.app.logics[name];
                    if (logic instanceof kk.Logic.Outlet) {

                        (function (name, logic, target) {

                            logic.ondone = function (ctx, app) {
                                target.outlet = name;
                            };

                        })(name, logic, this);

                    }
                }
            } else {
                this.app = false;
            }
        }

        if (this.app instanceof kk.App) {
            var params = this.get(ctx, app, ["params"]);
            var output = ctx.get(["output"]);
            ctx.begin();
            ctx.set(["params"], params);
            ctx.set(["result"], null);
            ctx.set(["output"], output);
            this.app.exec(ctx);
            var v = ctx.get(["result"]);
            ctx.end();
            if (v !== null) {
                ctx.set(["result"], v);
            }
            this.done(this.outlet || 'done', ctx, app);
        }
    }
});
