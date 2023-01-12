#!/bin/bash
set -e pipefail

# core configuration
export GO_BINARY="${GO_BINARY:-"go"}"
export GOBENCHDATA_BINARY="${GOBENCHDATA_BINARY:-"gobenchdata"}"
export INPUT_SUBDIRECTORY="${INPUT_SUBDIRECTORY:-"."}"
export INPUT_PRUNE_COUNT="${INPUT_PRUNE_COUNT:-"0"}"
export INPUT_BENCHMARKS_OUT="${INPUT_BENCHMARKS_OUT:-"benchmarks.json"}"
export INPUT_GO_TEST_PKGS="${INPUT_GO_TEST_PKGS:-"./..."}"
export INPUT_GO_BENCHMARKS="${INPUT_GO_BENCHMARKS:-"."}"
export INPUT_GIT_COMMIT_MESSAGE="${INPUT_GIT_COMMIT_MESSAGE:-"add benchmark run for ${GITHUB_SHA}"}"
export INPUT_GOBENCHDATA_PARSE_FLAGS="${INPUT_GOBENCHDATA_PARSE_FLAGS:-""}"

# publishing configuration
export INPUT_PUBLISH_REPO="${INPUT_PUBLISH_REPO:-${GITHUB_REPOSITORY}}"
export INPUT_PUBLISH_BRANCH="${INPUT_PUBLISH_BRANCH:-"gh-pages"}"

# pull request checks
export INPUT_CHECKS="${INPUT_CHECKS:-"false"}"
export INPUT_CHECKS_CONFIG="${INPUT_CHECKS_CONFIG:-"gobenchdata-checks.yml"}"

# output build data
echo '========================'
echo "👨‍⚕️ Checking configuration..."
echo "GO_BINARY=${GO_BINARY}"
${GO_BINARY} version
echo "GOBENCHDATA_BINARY=${GOBENCHDATA_BINARY}"
${GOBENCHDATA_BINARY} version
env | grep 'INPUT_'
echo "GITHUB_ACTOR=${GITHUB_ACTOR}"
echo "GITHUB_WORKSPACE=${GITHUB_WORKSPACE}"
echo "GITHUB_REPOSITORY=${GITHUB_REPOSITORY}"
echo "GITHUB_SHA=${GITHUB_SHA}"
echo "GITHUB_REF=${GITHUB_REF}"
echo '========================'

# setup
mkdir -p /tmp/{gobenchdata,build}

# run benchmarks from configured directory
echo
echo '📊 Running benchmarks...'
RUN_OUTPUT="/tmp/gobenchdata/benchmarks.json"
cd "${GITHUB_WORKSPACE}"
cd "${INPUT_SUBDIRECTORY}"
${GO_BINARY} test \
  -bench "${INPUT_GO_BENCHMARKS}" \
  -benchmem \
  ${INPUT_GO_TEST_FLAGS} \
  ${INPUT_GO_TEST_PKGS} |
  ${GOBENCHDATA_BINARY} ${INPUT_GOBENCHDATA_PARSE_FLAGS} --json "${RUN_OUTPUT}" -v "${GITHUB_SHA}" -t "ref=${GITHUB_REF}"
cd "${GITHUB_WORKSPACE}"

# fetch published data
if [[ "${INPUT_PUBLISH}" == "true" || "${INPUT_CHECKS}" == "true" ]]; then
  echo
  echo "📚 Checking out ${INPUT_PUBLISH_REPO}@${INPUT_PUBLISH_BRANCH}..."
  cd /tmp/build
  git clone https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${INPUT_PUBLISH_REPO}.git .
  git checkout ${INPUT_PUBLISH_BRANCH}
  echo

  if [[ "${INPUT_CHECKS}" == "true" ]]; then

    # check results against published
    echo '🔎 Evaluating results against base runs...'
    CHECKS_OUTPUT="/tmp/gobenchdata/checks-results.json"
    ${GOBENCHDATA_BINARY} checks eval "${INPUT_BENCHMARKS_OUT}" "${RUN_OUTPUT}" \
      --checks.config "${GITHUB_WORKSPACE}/${INPUT_CHECKS_CONFIG}" \
      --json ${CHECKS_OUTPUT} \
      --flat
    RESULTS=$(cat ${CHECKS_OUTPUT})
    echo "checks-results=$RESULTS" >>${GITHUB_OUTPUT}

    # output results
    echo
    echo '📝 Generating checks report...'
    ${GOBENCHDATA_BINARY} checks report ${CHECKS_OUTPUT}

  fi

  if [[ "${INPUT_PUBLISH}" == "true" ]]; then

    # merge results with published
    echo '☝️ Updating results...'
    if [[ -f "${INPUT_BENCHMARKS_OUT}" ]]; then
      echo '📈 Existing report found - merging...'
      ${GOBENCHDATA_BINARY} merge "${RUN_OUTPUT}" "${INPUT_BENCHMARKS_OUT}" \
        --prune "${INPUT_PRUNE_COUNT}" \
        --json "${INPUT_BENCHMARKS_OUT}" \
        --flat
    else
      cp "${RUN_OUTPUT}" "${INPUT_BENCHMARKS_OUT}"
    fi

    # publish results
    echo
    echo '📷 Committing and pushing new benchmark data...'
    if [[ "${SET_GIT_USER}" != "false" ]]; then
      git config --local user.email "${GITHUB_ACTOR}@users.noreply.github.com"
      git config --local user.name "${GITHUB_ACTOR}"
    fi
    git add .
    git commit -m "${INPUT_GIT_COMMIT_MESSAGE}"
    git push -f origin ${INPUT_PUBLISH_BRANCH}

  fi
fi

echo
echo '🚀 Done!'
