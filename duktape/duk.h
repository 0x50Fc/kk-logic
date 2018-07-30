
struct Context {
	duk_context * ctx;
	void * scope;
};

struct Context * NewContext();

void RecycleContext(struct Context * v);

duk_ret_t Throw(duk_context * ctx, const char * errmsg);