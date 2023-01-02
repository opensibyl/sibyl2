# this script was designed for some benchmark things
# /bin/zsh upload_batch.sh ~/github_workspace/jvm-sandbox 2 ~/github_workspace/sibyl2/scripts/sibyl
set -e

GIT_DIR=$1
BATCH=$2
SIBYL_EXE=$3

if [[ -z "${GIT_DIR}" ]]
then
  echo "git dir is required"
  exit 1
fi

if [[ -z "${BATCH}" ]]
then
  echo "batch is required"
  exit 1
fi

if [[ -z "${SIBYL_EXE}" ]]
then
  SIBYL_EXE = "./sibyl"
fi

echo "ready to upload dir ${GIT_DIR}, batch=${BATCH}"
CURRENT_BRANCH=`git symbolic-ref --short HEAD`
echo "current branch = ${CURRENT_BRANCH}"

cd ${GIT_DIR}

for i in {1..${BATCH}}
do
  echo "start uploading batch ${i} ..."
  ${SIBYL_EXE} upload --src "${GIT_DIR}"
  git checkout HEAD~1
done
git checkout ${CURRENT_BRANCH}
