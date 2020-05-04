#!/bin/bash
set -e

# core configuration
INPUT_SUBDIRECTORY="${INPUT_SUBDIRECTORY:-"."}"
INPUT_PRUNE_COUNT="${INPUT_PRUNE_COUNT:-"0"}"
INPUT_BENCHMARKS_OUT="${INPUT_BENCHMARKS_OUT:-"benchmarks.json"}"
INPUT_GO_TEST_PKGS="${INPUT_GO_TEST_PKGS:-"./..."}"
INPUT_GO_BENCHMARKS="${INPUT_GO_BENCHMARKS:-"."}"
INPUT_GIT_COMMIT_MESSAGE="${INPUT_GIT_COMMIT_MESSAGE:-"add benchmark run for ${GITHUB_SHA}"}"

# publishing configuration
INPUT_PUBLISH_REPO="${INPUT_PUBLSIH_REPO:-${GITHUB_REPOSITORY}}"
INPUT_PUBLISH_BRANCH="${INPUT_PUBLISH_BRANCH:-"gh-pages"}"

# output build data
echo '========================'
command -v gobenchdata
gobenchdata version
env | grep 'INPUT_'
echo "GITHUB_ACTOR=${GITHUB_ACTOR}"
echo "GITHUB_WORKSPACE=${GITHUB_WORKSPACE}"
echo "GITHUB_REPOSITORY=${GITHUB_REPOSITORY}"
echo "GITHUB_SHA=${GITHUB_SHA}"
echo "GITHUB_REF=${GITHUB_REF}"
echo '========================'

# setup
mkdir -p /tmp/{gobenchdata,build}
git config --global user.email "${GITHUB_ACTOR}@users.noreply.github.com"
git config --global user.name "${GITHUB_ACTOR}"

# run benchmarks from configured directory
echo
echo '📊 Running benchmarks...'
RUN_OUTPUT="/tmp/gobenchdata/benchmarks.json"
cd "${GITHUB_WORKSPACE}"
cd "${INPUT_SUBDIRECTORY}"
go test \
  -bench "${INPUT_GO_BENCHMARKS}" \
  -benchmem \
  ${INPUT_GO_TEST_FLAGS} \
  ${INPUT_GO_TEST_PKGS} \
  | gobenchdata --json "${RUN_OUTPUT}" -v "${GITHUB_SHA}" -t "ref=${GITHUB_REF}"
cd "${GITHUB_WORKSPACE}"

# fetch publish destination
echo
echo "📚 Checking out ${INPUT_PUBLISH_REPO}@${INPUT_PUBLISH_BRANCH}..."
cd /tmp/build
git clone https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${INPUT_PUBLISH_REPO}.git .
git checkout ${INPUT_PUBLISH_BRANCH}

# generate output
echo
echo '☝️ Updating results...'
if [[ -f "${INPUT_BENCHMARKS_OUT}" ]]; then
  echo '📈 Existing report found - merging...'
  gobenchdata merge "${RUN_OUTPUT}" "${INPUT_BENCHMARKS_OUT}" \
    --flat \
    --prune "${INPUT_PRUNE_COUNT}" \
    --json "${INPUT_BENCHMARKS_OUT}"
else
  cp "${RUN_OUTPUT}" "${INPUT_BENCHMARKS_OUT}"
fi

# publish results
echo
echo '📷 Committing and pushing new benchmark data...'
git add .
git commit -m "${INPUT_GIT_COMMIT_MESSAGE}"
git push -f origin ${INPUT_PUBLISH_BRANCH}

echo
echo '🚀 Done!'
