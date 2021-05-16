#!/bin/bash

set -eou pipefail

PREFIX=${PREFIX:=/}

cp bin/mountmond "$PREFIX"/usr/bin/mountmond

cat >"$PREFIX"/etc/mountmond.yaml <<EOF
mounts:
EOF
