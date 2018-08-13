
kk.Logic.Each = function () {
    kk.Logic.apply(this, arguments);
};

kk.Logic.Each.prototype = kk.extend(kk.Logic.prototype, {

    item: function (ctx, app, object, fields) {
        var v = {};
        for (var key in fields) {
            var vv = fields[key];
            if (typeof vv == 'string' ) {
                if(vv.startsWith("=")) {
                    ctx.begin();
                    ctx.set(["object"], object);
                    vv = this.evaluateValue(ctx, app, vv, object);
                    ctx.end();
                } else {
                    vv = kk.get(object, [vv]);
                }
            } else {
                vv = this.evaluateValue(ctx, app, vv, object);
            }
            if(vv !== undefined) {
                v.push(vv);
            }
        }
        return v;
    },

    exec: function (ctx, app) {
        kk.Logic.prototype.exec.apply(this, arguments);

        var type = this.get(ctx, app, ["type"]) || 'auto';
        var value = this.get(ctx, app, ["value"]);
        var fields = kk.get(this._object, ["fields"]) || {};

        if (type == 'auto') {
            if (value instanceof Array) {
                type = 'array';
            } else if (typeof value == 'object') {
                type = 'object';
            }
        }

        if (type == 'array') {
            var a = [];

            if (value instanceof Array) {
                for (var i = 0; i < value.length; i++) {
                    var v = value[i];
                    a.push(this.item(ctx, app, v, fields));
                }
            } else if (typeof value == 'object') {
                a.push(this.item(ctx, app, value, fields));
            }

            ctx.set(["result"], a);

            this.done("done", ctx, app);

        } else if (type == 'object') {

            var a = {};

            if (value instanceof Array) {
                if (value.length > 0) {
                    a = this.item(ctx, app, value[0], fields);
                }
            } else if (typeof value == 'object') {
                a = this.item(ctx, app, value, fields);
            }

            ctx.set(["result"], a);
            this.done("done", ctx, app);
        }

    }
});
