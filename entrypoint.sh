#!/bin/bash
set -e

command -v gobenchdata
gobenchdata version

cd "${GITHUB_WORKSPACE}"
mkdir -p /tmp/{gobenchdata,build}

git config --global user.email "${GITHUB_ACTOR}@users.noreply.github.com"
git config --global user.name "${GITHUB_ACTOR}"

echo
echo '📊 Running benchmarks...'
RUN_OUTPUT="/tmp/gobenchdata/benchmarks.json"
go test \
  -bench "${INPUT_GO_BENCHMARKS:-"."}" \
  -benchmem \
  ${INPUT_GO_TEST_FLAGS} \
  ${INPUT_GO_TEST_PKGS:-"./..."} \
  | gobenchdata --json "${RUN_OUTPUT}" -v "${GITHUB_SHA}" -t "ref=${GITHUB_REF}"

echo
echo '📚 Checking out gh-pages...'
cd /tmp/build
git clone https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git .
git checkout gh-pages

echo
echo '☝️ Updating results...'
FINAL_OUTPUT="${INPUT_BENCHMARKS_OUT:-"benchmarks.json"}"
if [[ -f "${FINAL_OUTPUT}" ]]; then
  echo '📈 Existing report found - merging...'
  gobenchdata merge "${RUN_OUTPUT}" "${FINAL_OUTPUT}" \
    --flat \
    --prune "${INPUT_PRUNE_COUNT:-"0"}" \
    --json "${FINAL_OUTPUT}"
else
  cp "${RUN_OUTPUT}" "${FINAL_OUTPUT}"
fi

echo
echo '📷 Committing and pushing new benchmark data...'
git add .
git commit -m "${INPUT_GIT_COMMIT_MESSAGE:-"add benchmark run for ${GITHUB_SHA}"}"
git push -f origin gh-pages
cd ../

echo
echo '🚀 Done!'
