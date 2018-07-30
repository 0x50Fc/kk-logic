
kk.date = function (d, fmt) {

    if (!(d instanceof Date)) {
        d = new Date(parseInt(d));
    }

    var o = {
        "M+": d.getMonth() + 1, //月份 
        "d+": d.getDate(), //日 
        "h+": d.getHours(), //小时 
        "m+": d.getMinutes(), //分 
        "s+": d.getSeconds(), //秒 
        "q+": Math.floor((d.getMonth() + 3) / 3), //季度 
        "S": d.getMilliseconds() //毫秒 
    };

    fmt = fmt.replace(/(y+)/, d.getFullYear() + "");

    for (var k in o) {
        var v = o[k] + '';
        if (v.length < 2) {
            v = '0' + v;
        }
        fmt = fmt.replace(new RegExp("(" + k + ")"), v);
    }

    return fmt;
};
