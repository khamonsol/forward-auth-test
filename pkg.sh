rm -rf fn.zip
mkdir pkg-build
cp -R src .logconfig go.mod go.sum main.go pkg-build
cd pkg-build
zip -r fn.zip .
cp fn.zip ../
cd ..
rm -rf pkg-build