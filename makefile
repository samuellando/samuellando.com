SRC = $(wildcard frontend/jsx/*.jsx)
LIB = $(SRC:frontend/jsx/%.jsx=frontend/js/%.js)

lib: $(LIB)
frontend/js/%.js: frontend/jsx/%.jsx .babelrc
	mkdir -p $(@D)
	./node_modules/\@babel/cli/bin/babel.js $< -o $@
