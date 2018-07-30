
(function(){
	var modules = {};
	require = function(path) {
		var m = modules[path];
		if(m === undefined) {
			m = { exports : {} };
			try {
				var fn = compile(path,'(function(module,exports){','})');
                if(typeof fn == 'function') {
                    fn = fn();
                    if(typeof fn == 'function') {
                        fn(m,m.exports);
                    } else {
                        echo("[REQUITE] [ERROR] " + (typeof fn) + "\n");
                    }
                } else {
                    echo("[REQUITE] [ERROR] Not Found " + path+ "\n");
                }
			} catch(e) {
				echo(e.toString()+ "\n");
			}
			modules[path] = m;
		}
		return m.exports;
	};
})();
