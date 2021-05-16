#!/bin/bash

set -eou pipefail

tar -czf "${PACKAGE_NAME}_${VERSION}.orig.tar.gz" /src/

pushd /src/
make clean
make
debuild -us -uc

popd
mv "${PACKAGE_NAME}_${VERSION}-1_amd64.deb" /src/
