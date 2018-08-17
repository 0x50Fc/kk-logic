
kk.extend = function (proto, object) {
    var p = {};
    if (typeof object == 'object') {
        for (var key in object) {
            p[key] = { value: object[key], configurable: true };
        }
    }
    return Object.create(proto, p);
};
