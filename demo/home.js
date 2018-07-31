
kk.run({
    in : new kk.Logic.Var(
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
    )
});
