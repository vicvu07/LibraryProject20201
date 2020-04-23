
my_dir=`dirname $0`
cd $my_dir

# minify js
uglifyjs --compress -o src/js/bundle.min.js src/js/bootstrap-notify.min.js src/js/credUtil.js src/js/domUtil.js src/js/storage.js src/js/main.js src/js/now-ui-kit.js src/js/router.js src/js/axiosUtil.js src/js/date.format.js

# remove now-ui-kit img
rm -rf src/node_modules/now-ui-kit/assets/img
rm -rf src/node_modules/now-ui-kit/documentation

# build dist
rm -rf dist && mkdir dist && cp -R src/* dist/
