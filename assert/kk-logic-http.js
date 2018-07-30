
kk.Logic.Http = function () {
    kk.Logic.apply(this, arguments);
};

kk.Logic.Http.prototype = kk.extend(kk.Logic.prototype, {

    exec: function (ctx, app) {
        kk.Logic.prototype.exec.apply(this, arguments);

        var options = {
            method: this.get(ctx, app, ["method"]) || 'GET',
            url: this.get(ctx, app, ["url"]),
            data: this.get(ctx, app, ["data"]),
            type: this.get(ctx, app, ["type"]),
            dataType: this.get(ctx, app, ["dataType"]),
        };

        var data = http.send(options);

        if (data === false) {
            ctx.set(["result"], undefined);
            ctx.set(["error"], http.errmsg);
            this.done("error", ctx, app)
            return;
        }

        ctx.set(["result"], data);

        this.done("done", ctx, app);
    }
});
