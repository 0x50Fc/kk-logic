
kk.Logic.Var = function () {
    kk.Logic.apply(this, arguments);
};

kk.Logic.Var.prototype = kk.extend(kk.Logic.prototype, {

    exec: function (ctx, app) {
        kk.Logic.prototype.exec.apply(this, arguments);

        var key = this.get(ctx, app, ["key"]);
        var value = this.get(ctx, app, ["value"]);

        if (key !== undefined) {
            if (typeof key == 'string' && key.endsWith("[]")) {
                key = key.substr(0, key.length - 2);
                var vs = ctx.get(key);
                if (!vs || !(vs instanceof Array)) {
                    vs = [value];
                } else {
                    vs.push(value);
                }
                ctx.set(key, vs);
            } else {
                ctx.set(key, value);
            }
        }

        this.done("done", ctx, app);
    }
});
