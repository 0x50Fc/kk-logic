
kk.Logic.Outlet = function () {
    kk.Logic.apply(this, arguments);
};

kk.Logic.Outlet.prototype = kk.extend(kk.Logic.prototype, {

    exec: function (ctx, app) {
        kk.Logic.prototype.exec.apply(this, arguments);

        if (typeof this.ondone == 'function') {
            this.ondone(ctx, app);
        }

    }
});
