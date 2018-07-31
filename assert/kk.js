

kk = {
    extend: function (proto, object) {
        var v = {};
        if (object) {
            for (var key in object) {
                v[key] = { value: object[key] };
            }
        }
        return Object.create(proto, v)
    }
};

kk.get = function (object, keys, index) {

    if (keys === undefined) {
        return object;
    }

    if (typeof keys == 'string') {
        keys = keys.split(".");
    }

    if (keys instanceof Array) {

        if (index === undefined) {
            index = 0;
        }

        if (index < keys.length) {

            var key = keys[index];

            if (typeof object == 'object') {
                return kk.get(object[key], keys, index + 1);
            }

        } else {
            return object;
        }

    }

};

kk.set = function (object, keys, value, index) {

    if (keys === undefined) {
        return;
    }

    if (typeof object != 'object') {
        return;
    }

    if (typeof keys == 'string') {
        keys = keys.split(".");
    }

    if (keys instanceof Array) {

        if (index === undefined) {
            index = 0;
        }

        if (index + 1 < keys.length) {

            var key = keys[index];

            var v = object[key];

            if (v === undefined) {
                v = {};
                object[key] = v;
            }

            kk.set(v, keys, value, index + 1);

        } else if (index < keys.length) {
            var key = keys[index];
            object[key] = value;
        }
    }

};

kk.Context = function () {
    this._objects = [{}];
};

kk.Context.prototype = kk.extend(Object.prototype, {

    get: function (keys) {

        if (!(keys instanceof Array)) {
            keys = (keys || '').split('.');
        }

        if (keys.length == 0) {
            return this._objects[this._objects.length - 1];
        }

        var key = keys[0];
        var i = this._objects.length - 1;
        var object;

        while (i >= 0) {

            object = this._objects[i];

            if (object[key] !== undefined) {
                break;
            }

            i--;
        }

        return kk.get(object, keys);
    },

    set: function (keys, value) {
        var object = this._objects[this._objects.length - 1];
        kk.set(object, keys, value);
        return this;
    },

    evaluate: function (evaluate) {
        var _G;
        (function (object) {
            with (object) {
                _G = eval('(' + evaluate + ')');
            }
        })(this._objects[this._objects.length - 1]);
        return _G;
    },

    begin: function () {
        var object = this._objects[this._objects.length - 1];
        var v = {};
        for (var key in object) {
            v[key] = object[key];
        }
        this._objects.push(v);
    },

    end: function () {
        this._objects.pop();
    }
});

kk.Logic = function (object,app) {
    this._object = {};
    this._on = {};
    if (object) {
        for (var key in object) {
            if (key.startsWith("on")) {
                this.on(key.substr(2), object[key]);
            } else {
                this._object[key] = object[key];
            }
        }
    }
};

kk.Logic.prototype = kk.extend(Object.prototype, {

    evaluateValue: function (ctx, app, value, object) {

        if (typeof value == 'string') {

            if (value.startsWith("=")) {
                value = ctx.evaluate(value.substr(1));
            }

        } else if (typeof value == 'function') {
            ctx.begin();
            ctx.set(["object"], object);
            ctx.set(["result"], null);
            var vv = value(ctx, app);
            var v = ctx.get(["result"]);
            ctx.end();
            if (vv !== undefined) {
                return vv;
            }
            if (v !== null) {
                return v;
            }
            return undefined;
        } else if (typeof value == 'object') {

            if (value instanceof kk.Logic) {
                ctx.begin();
                ctx.set(["object"], object);
                ctx.set(["result"], null);
                value.exec(ctx, app);
                var v = ctx.get(["result"]);
                ctx.end();
                if (v !== null) {
                    return v;
                }
                return undefined;
            } else if (value instanceof Array) {
                var a = [];
                for (var i = 0; i < value.length; i++) {
                    var v = value[i];
                    a.push(this.evaluateValue(ctx, app, v, object));
                }
                return a;
            } else {
                var a = {};
                for (var key in value) {
                    a[key] = this.evaluateValue(ctx, app, value[key], object);
                }
                return a;
            }

        }

        return value;
    },

    get: function (ctx, app, keys, index, object) {

        if (index === undefined) {
            index = 0;
        }

        if (object === undefined) {
            object = this._object;
        }

        if (!(keys instanceof Array)) {
            keys = (keys || '').split(".");
        }

        if (index < keys.length) {

            var key = keys[index];

            var v = this.evaluateValue(ctx, app, kk.get(object, [key]), object);

            if (v === undefined) {
                return v;
            }

            return this.get(ctx, app, keys, index + 1, v);

        }

        return object;
    },

    exec: function (ctx, app) {

    },

    on: function (name, fn) {
        this._on[name] = fn;
        return this;
    },

    done: function (name, ctx, app) {
        var fn = this._on[name];
        if (typeof fn == 'function') {
            fn(ctx, app, this)
        } else if(typeof fn == 'string') {
            var v = app.get(fn);
            if(v && v instanceof kk.Logic) {
                v.exec(ctx,app);
            }
        } else if (fn instanceof kk.Logic) {
            fn.exec(ctx, app);
        }
    }

});

kk.App = function () {
    this.logics = {};
};

kk.App.prototype = kk.extend(Object.prototype, {

    exec: function (ctx, name) {
        if (name === undefined) {
            name = 'in';
        }
        var v = this.logics[name];
        if (v instanceof kk.Logic) {
            v.exec(ctx, this);
        }
    },

    set: function (name, logic) {
        this.logics[name] = logic;
        return this;
    },

    get : function(name) {
        return this.logics[name];
    },

    log: function (text) {

    }

});

kk.run = function(object) {

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
    ctx.set(["headers"], _HEADER);
    ctx.set(["url"], _REQUEST['url']);
    ctx.set(["path"], _REQUEST['path']);
    ctx.set(["hostname"], _REQUEST['hostname']);
    ctx.set(["protocol"], _REQUEST['protocol']);

    var app = new kk.App();

    for(var key in object) {
        app.set(key,object[key]);
    }
    
    app.exec(ctx);

    var v = ctx.get(["view"]);

    if(v !== undefined) {

        if(v.headers) {
            for(var key in v.headers) {
                header(key, v.headers[key]);
            }
        }

        echo(v.content);

    } else {
        header("Content-Type", "application/json; charset=utf-8");
        echo(JSON.stringify(ctx.get(["output"])));
    }
    
};

