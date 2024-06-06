

# init handles the setup of the library
# 
# use the library by adding the following line to a script:
# source $(dirname $0)/../common/libbackend.sh
#
# If you want to limit what targets a backend can be used on, set the variable LIMIT_TARGETS to a
# space separated list of valid targets BEFORE sourcing the library, for example to only allow a backend
# to be used on CUDA and CPU backends:
#
# LIMIT_TARGETS="cublas cpu"
# source $(dirname $0)/../common/libbackend.sh
#
# You can use any valid BUILD_TYPE or BUILD_PROFILE, if you need to limit a backend to CUDA 12 only:
#
# LIMIT_TARGETS="cublas12"
# source $(dirname $0)/../common/libbackend.sh
#
function init() {
    BACKEND_NAME=${PWD##*/}
    MY_DIR=$(realpath `dirname $0`)
    BUILD_PROFILE=$(getBuildProfile)

    # If a backend has defined a list of valid build profiles...
    if [ ! -z "${LIMIT_TARGETS}" ]; then
        isValidTarget=$(checkTargets ${LIMIT_TARGETS})
        if [ ${isValidTarget} != true ]; then
            echo "${BACKEND_NAME} can only be used on the following targets: ${LIMIT_TARGETS}"
            exit 0
        fi
    fi

    echo "Initializing libbackend for ${BACKEND_NAME}"
}

# getBuildProfile will inspect the system to determine which build profile is appropriate:
# returns one of the following:
# - cublas11
# - cublas12
# - hipblas
# - intel
function getBuildProfile() {
    # First check if we are a cublas build, and if so report the correct build profile
    if [ x"${BUILD_TYPE}" == "xcublas" ]; then
        if [ ! -z ${CUDA_MAJOR_VERSION} ]; then
            # If we have been given a CUDA version, we trust it
            echo ${BUILD_TYPE}${CUDA_MAJOR_VERSION}
        else
            # We don't know what version of cuda we are, so we report ourselves as a generic cublas
            echo ${BUILD_TYPE}
        fi
        return 0
    fi

    # If /opt/intel exists, then we are doing an intel/ARC build
    if [ -d "/opt/intel" ]; then
        echo "intel"
        return 0
    fi

    # If for any other values of BUILD_TYPE, we don't need any special handling/discovery
    if [ ! -z ${BUILD_TYPE} ]; then
        echo ${BUILD_TYPE}
        return 0
    fi

    # If there is no BUILD_TYPE set at all, set a build-profile value of CPU, we aren't building for any GPU targets
    echo "cpu"
}

# ensureVenv makes sure that the venv for the backend both exists, and is activated.
#
# This function is idempotent, so you can call it as many times as you want and it will
# always result in an activated virtual environment
function ensureVenv() {
    if [ ! -d "${MY_DIR}/venv" ]; then
        uv venv ${MY_DIR}/venv
        echo "virtualenv created"
    fi
    
    if [ "x${VIRTUAL_ENV}" != "x${MY_DIR}/venv" ]; then
        source ${MY_DIR}/venv/bin/activate
        echo "virtualenv activated"
    fi

    echo "activated virtualenv has been ensured"
}

# installRequirements looks for several requirements files and if they exist runs the install for them in order
#
#  - requirements-install.txt
#  - requirements.txt
#  - requirements-${BUILD_TYPE}.txt
#  - requirements-${BUILD_PROFILE}.txt
#
# BUILD_PROFILE is a pore specific version of BUILD_TYPE, ex: cuda11 or cuda12
# it can also include some options that we do not have BUILD_TYPES for, ex: intel
#
# NOTE: for BUILD_PROFILE==intel, this function does NOT automatically use the Intel python package index.
# you may want to add the following line to a requirements-intel.txt if you use one:
#
# --index-url https://pytorch-extension.intel.com/release-whl/stable/xpu/us/
#
# If you need to add extra flags into the pip install command you can do so by setting the variable EXTRA_PIP_INSTALL_FLAGS
# before calling installRequirements.  For example:
#
# source $(dirname $0)/../common/libbackend.sh
# EXTRA_PIP_INSTALL_FLAGS="--no-build-isolation"
# installRequirements
function installRequirements() {
    ensureVenv

    # These are the requirements files we will attempt to install, in order
    declare -a requirementFiles=(
        "${MY_DIR}/requirements-install.txt"
        "${MY_DIR}/requirements.txt"
        "${MY_DIR}/requirements-${BUILD_TYPE}.txt"
    )

    if [ "x${BUILD_TYPE}" != "x${BUILD_PROFILE}" ]; then
        requirementFiles+=("${MY_DIR}/requirements-${BUILD_PROFILE}.txt")
    fi

    for reqFile in ${requirementFiles[@]}; do
        if [ -f ${reqFile} ]; then
            echo "starting requirements install for ${reqFile}"
            uv pip install ${EXTRA_PIP_INSTALL_FLAGS} --requirement ${reqFile}
            echo "finished requirements install for ${reqFile}"
        fi
    done
}

# startBackend discovers and runs the backend GRPC server
#
# You can specify a specific backend file to execute by setting BACKEND_FILE before calling startBackend.
# example:
#
# source ../common/libbackend.sh
# BACKEND_FILE="${MY_DIR}/source/backend.py"
# startBackend $@
#
# valid filenames for autodiscovered backend servers are:
#  - server.py
#  - backend.py
#  - ${BACKEND_NAME}.py
function startBackend() {
    ensureVenv

    if [ ! -z ${BACKEND_FILE} ]; then
        python ${BACKEND_FILE} $@
    elif [ -e "${MY_DIR}/server.py" ]; then
        python ${MY_DIR}/server.py $@
    elif [ -e "${MY_DIR}/backend.py" ]; then
        python ${MY_DIR}/backend.py $@
    elif [ -e "${MY_DIR}/${BACKEND_NAME}.py" ]; then
        python ${MY_DIR}/${BACKEND_NAME}.py $@
    fi
}

# runUnittests discovers and runs python unittests
#
# You can specify a specific test file to use by setting TEST_FILE before calling runUnittests.
# example:
#
# source ../common/libbackend.sh
# TEST_FILE="${MY_DIR}/source/test.py"
# runUnittests $@
#
# be default a file named test.py in the backends directory will be used
function runUnittests() {
    ensureVenv

    if [ ! -z ${TEST_FILE} ]; then
        testDir=$(dirname `realpath ${TEST_FILE}`)
        testFile=$(basename ${TEST_FILE})
        pushd ${testDir}
        python -m unittest ${testFile}
        popd
    elif [ -f "${MY_DIR}/test.py" ]; then
        pushd ${MY_DIR}
        python -m unittest test.py
        popd
    else
        echo "no tests defined for ${BACKEND_NAME}"
    fi
}

##################################################################################
# Below here are helper functions not intended to be used outside of the library #
##################################################################################

# checkTargets determines if the current BUILD_TYPE or BUILD_PROFILE is in a list of valid targets
function checkTargets() {
    # Collect all provided targets into a variable and...
    targets=$@
    # ...convert it into an array
    declare -a targets=($targets)

    for target in ${targets[@]}; do
        if [ "x${BUILD_TYPE}" == "x${target}" ]; then
            echo true
            return 0
        fi
        if [ "x${BUILD_PROFILE}" == "x${target}" ]; then
            echo true
            return 0
        fi
    done
    echo false
}

init